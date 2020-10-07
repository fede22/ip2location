package mysql

import (
	"database/sql"
	"fmt"
	"github.com/fede22/ip2location/internal/proxy"
	"math/big"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type client struct {
	db *sql.DB
}

type decimalIP string

func NewProxyRepository() (client, error) {
	dataSourceName := "root:rootroot@tcp(localhost:3306)/ip2proxy?charset=utf8"
	return newClient(dataSourceName, 10, 10, time.Minute*3)
}

func newClient(dataSourceName string, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (client, error) {
	driverName := "mysql"
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return client{}, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	if connMaxLifetime.Minutes() > 5 {
		return client{}, fmt.Errorf("connMaxLifetime should not be greater than 5 min")
	}
	db.SetConnMaxLifetime(connMaxLifetime)
	if err := db.Ping(); err != nil {
		return client{}, err
	}
	return client{db}, nil
}

func (c client) GetProxy(address proxy.NetIP) (proxy.Proxy, error) {
	var p proxy.Proxy
	query := "select ip_from, ip_to, proxy_type, country_code, country_name, region_name, city_name, isp, domain, usage_type, asn," +
		" `as` from ip2proxy.ip2proxy_px7 where ? between ip_from and ip_to;"
	rows, err := c.db.Query(query, toDecimalIP(address))
	if err != nil {
		return proxy.Proxy{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var addressFrom, addressTo decimalIP
		err := rows.Scan(&addressFrom, &addressTo, &p.ProxyType, &p.CountryCode,
			&p.CountryName, &p.RegionName, &p.CityName, &p.ISP, &p.Domain, &p.UsageType, &p.ASN, &p.AS)
		if err != nil {
			return proxy.Proxy{}, err
		}
		p, err = setAddresses(p, addressFrom, addressTo)
		if err != nil {
			return proxy.Proxy{}, err
		}
	}
	if err := rows.Err(); err != nil {
		return proxy.Proxy{}, err
	}
	return p, nil
}

func (c client) GetProxies(countryCode string, limit int) ([]proxy.Proxy, error) {
	query := "select ip_from, ip_to, proxy_type, country_code, country_name, region_name, city_name, isp, domain, usage_type, asn," +
		" `as` from ip2proxy.ip2proxy_px7 where country_code = ? limit ?;"
	rows, err := c.db.Query(query, countryCode, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	proxies := make([]proxy.Proxy, 0)
	for rows.Next() {
		var p proxy.Proxy
		var addressFrom, addressTo decimalIP
		err := rows.Scan(&addressFrom, &addressTo, &p.ProxyType, &p.CountryCode,
			&p.CountryName, &p.RegionName, &p.CityName, &p.ISP, &p.Domain, &p.UsageType, &p.ASN, &p.AS)
		if err != nil {
			return nil, err
		}
		p, err = setAddresses(p, addressFrom, addressTo)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return proxies, nil
}

func (c client) GetISPs(countryCode string) ([]string, error) {
	query := "select isp from ip2proxy.ip2proxy_px7 where country_code = ? group by isp;"
	rows, err := c.db.Query(query, countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ispNames := make([]string, 0)
	for rows.Next() {
		var isp sql.NullString
		err := rows.Scan(&isp)
		if err != nil {
			return nil, err
		}
		if isp.Valid {
			ispNames = append(ispNames, isp.String)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ispNames, nil
}

func (c client) GetIPCount(countryCode string) (int, error) {
	query := "select sum((ip_to - ip_from) + 1) from ip2proxy.ip2proxy_px7 where country_code = ?;"
	var count int
	err := c.db.QueryRow(query, countryCode).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c client) TopProxyTypes(limit int) ([]string, error) {
	query := "select proxy_type, count(*) from ip2proxy.ip2proxy_px7 group by proxy_type order by count(*) desc limit ?;"
	rows, err := c.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	proxyTypes := make([]string, 0)
	for rows.Next() {
		var pt sql.NullString
		var count int
		err := rows.Scan(&pt, &count)
		if err != nil {
			return nil, err
		}
		if pt.Valid {
			proxyTypes = append(proxyTypes, pt.String)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return proxyTypes, nil
}

func setAddresses(p proxy.Proxy, addressFrom, addressTo decimalIP) (proxy.Proxy, error) {
	ip, err := addressFrom.toNetIP()
	if err != nil {
		return proxy.Proxy{}, err
	}
	p.AddressFrom = ip
	ip, err = addressTo.toNetIP()
	if err != nil {
		return proxy.Proxy{}, err
	}
	p.AddressTo = ip
	return p, nil
}

func (dec decimalIP) toNetIP() (proxy.NetIP, error) {
	x, err := dec.toBigIntIP()
	if err != nil {
		return proxy.NetIP{}, err
	}
	return x.ToNetIP(), nil
}

func (dec decimalIP) toBigIntIP() (proxy.BigIntIP, error) {
	x, ok := big.NewInt(0).SetString(string(dec), 10)
	if !ok {
		return proxy.BigIntIP{}, fmt.Errorf("error parsing address %s to big int", dec)
	}
	return proxy.BigIntIP{Int: x}, nil
}

func toDecimalIP(ip proxy.NetIP) decimalIP {
	return decimalIP(ip.ToBigIntIP().String())
}
