//TODO use the same convention for naming all methods
package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/ip2location/ip2proxy-go"
)

const (
	driverName     = "mysql"
	dataSourceName = "root:rootroot@tcp(localhost:3306)/ip2proxy?charset=utf8"
)

type Client struct  {
	db *sql.DB
}

func NewClient() (Client, error) {
	return client(10, 10, time.Minute * 3)
}


func client(maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (Client, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return Client{}, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	if connMaxLifetime.Minutes() > 5 {
		return Client{}, fmt.Errorf("connMaxLifetime should not be greater than 5 min")
	}
	db.SetConnMaxLifetime(connMaxLifetime)
	if err := db.Ping(); err != nil {
		return Client{}, err
	}
	return Client{db}, nil
}

func (c Client) GetByIP(address string) (Proxy, error) {
	var p Proxy
	query := "select ip_from, ip_to, proxy_type, country_code, country_name, region_name, city_name, isp, domain, usage_type, asn," +
		" `as` from ip2proxy.ip2proxy_px7 where ? between ip_from and ip_to;"
	rows, err := c.db.Query(query, address)
	if err != nil {
		return Proxy{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&p.AddressFrom, &p.AddressTo, &p.ProxyType, &p.CountryCode,
				&p.CountryName, &p.RegionName, &p.CityName, &p.ISP, &p.Domain, &p.UsageType, &p.ASN, &p.AS)
		if err != nil {
			return Proxy{}, err
		}
	}
	if err := rows.Err(); err != nil {
		return Proxy{}, err
	}
	return p, nil
}

func (c Client) GetByCountryCode(countryCode string, limit int) ([]Proxy, error) {
	query := "select ip_from, ip_to, proxy_type, country_code, country_name, region_name, city_name, isp, domain, usage_type, asn," +
		" `as` from ip2proxy.ip2proxy_px7 where country_code = ? limit ?;"
	rows, err := c.db.Query(query, countryCode, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	proxies := make([]Proxy, 0)
	for rows.Next() {
		var p Proxy
		err := rows.Scan(&p.AddressFrom, &p.AddressTo, &p.ProxyType, &p.CountryCode,
			&p.CountryName, &p.RegionName, &p.CityName, &p.ISP, &p.Domain, &p.UsageType, &p.ASN, &p.AS)
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

func (c Client) GetISPNames(countryCode string) ([]string, error) {
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

func (c Client) GetIPCount(countryCode string) (int, error) {
	query := "select sum((ip_to - ip_from) + 1) from ip2proxy.ip2proxy_px7 where country_code = ?;"
	var count int
	err := c.db.QueryRow(query, countryCode).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

//TODO rename
func (c Client) TopProxyTypes(limit int) ([]string, error) {
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

