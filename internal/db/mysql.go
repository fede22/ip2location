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

func Client() (*sql.DB, error) {
	return client(10, 10, time.Minute * 3)
}


func client(maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	if connMaxLifetime.Minutes() > 5 {
		return nil, fmt.Errorf("connMaxLifetime should not be greater than 5 min")
	}
	db.SetConnMaxLifetime(connMaxLifetime)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

