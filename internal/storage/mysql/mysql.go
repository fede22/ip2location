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

//TODO check if address should be compared to address_from, address_to or if it's between them.
//TODO test null value for any column
func (c Client) GetIP(address string) (Proxy, error) {
	var p Proxy
	query := "select ip_from, ip_to, proxy_type, country_code, country_name, region_name, city_name, isp, domain, usage_type, asn," +
		" `as` from ip2proxy.ip2proxy_px7 where ip_from=?;"
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