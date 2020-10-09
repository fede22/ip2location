package mysql

import (
	"database/sql"
	"fmt"
	"github.com/fede22/ip2location/internal/domain"
	"math/big"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type client struct {
	db *sql.DB
}

type decimalIP string

const getProxyQuery = "select ip_from, ip_to, proxy_type, country_code, country_name, region_name, city_name, isp, domain, usage_type, asn, `as` from ip2proxy_px7 where ? between ip_from and ip_to;"

func NewRepository() (client, error) {
	dataSourceName := "root:rootroot@tcp(localhost:3306)/ip2proxy?charset=utf8"
	driverName := "mysql"
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return client{}, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)
	if err := db.Ping(); err != nil {
		return client{}, err
	}
	return client{db}, nil
}

func (c client) GetProxy(address domain.NetIP) (domain.Proxy, error) {
	var p domain.Proxy
	var addressFrom, addressTo decimalIP
	err := c.db.QueryRow(getProxyQuery, toDecimalIP(address)).Scan(&addressFrom, &addressTo, &p.ProxyType, &p.CountryCode,
		&p.CountryName, &p.RegionName, &p.CityName, &p.ISP, &p.Domain, &p.UsageType, &p.ASN, &p.AS)
	if err != nil {
		return domain.Proxy{}, err
	}
	p, err = setAddresses(p, addressFrom, addressTo)
	if err != nil {
		return domain.Proxy{}, err
	}
	return p, nil
}

func (c client) GetProxies(countryCode string, limit int) ([]domain.Proxy, error) {
	query := "select ip_from, ip_to, proxy_type, country_code, country_name, region_name, city_name, isp, domain, usage_type, asn," +
		" `as` from ip2proxy.ip2proxy_px7 where country_code = ? limit ?;"
	rows, err := c.db.Query(query, strings.ToUpper(countryCode), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	proxies := make([]domain.Proxy, 0)
	for rows.Next() {
		var p domain.Proxy
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
	rows, err := c.db.Query(query, strings.ToUpper(countryCode))
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
	err := c.db.QueryRow(query, strings.ToUpper(countryCode)).Scan(&count)
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

func setAddresses(p domain.Proxy, addressFrom, addressTo decimalIP) (domain.Proxy, error) {
	ip, err := addressFrom.toNetIP()
	if err != nil {
		return domain.Proxy{}, err
	}
	p.AddressFrom = ip
	ip, err = addressTo.toNetIP()
	if err != nil {
		return domain.Proxy{}, err
	}
	p.AddressTo = ip
	return p, nil
}

func (dec decimalIP) toNetIP() (domain.NetIP, error) {
	x, err := dec.toBigIntIP()
	if err != nil {
		return domain.NetIP{}, err
	}
	return x.ToNetIP(), nil
}

func (dec decimalIP) toBigIntIP() (domain.BigIntIP, error) {
	x, ok := big.NewInt(0).SetString(string(dec), 10)
	if !ok {
		return domain.BigIntIP{}, fmt.Errorf("error parsing address %s to big int", dec)
	}
	return domain.BigIntIP{Int: x}, nil
}

func toDecimalIP(ip domain.NetIP) decimalIP {
	return decimalIP(ip.ToBigIntIP().String())
}
