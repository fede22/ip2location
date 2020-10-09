package mysql

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fede22/ip2location/internal/domain"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"strings"
	"testing"
)

func TestNewRepository(t *testing.T) {
	client, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}

	var (
		from string
		to   string
	)

	query := "select ip_from, ip_to from ip2proxy_px7 limit 1"
	rows, err := client.db.Query(query)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&from, &to)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(from, to)
}

func TestClient_GetProxy(t *testing.T) {
	client, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}

	ip := domain.NetIP{IP: net.ParseIP("1.0.4.1")}
	p, err := client.GetProxy(ip)
	if err != nil {
		t.Fatal(err)
	}
	if !p.AddressFrom.Equal(ip.IP) {
		t.Errorf("expected address_from %s, got instead %s", ip, p.AddressFrom)
	}
	t.Log(p)
}

func TestClient_GetProxy_Mock(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err, "error opening mock database connection")
	defer db.Close()
	defer func() {
		assert.Nil(t, mock.ExpectationsWereMet(), "there were unfulfilled expectations")
	}()
	b, err := ioutil.ReadFile("testdata/get_proxy.csv")
	assert.Nil(t, err, "couldn't read csv file")
	columns := []string{"ip_from", "ip_to", "proxy_type", "country_code", "country_name", "region_name", "city_name",
		"isp", "domain", "usage_type", "asn", "as"}
	ip := domain.NetIP{IP: net.ParseIP("1.0.4.1")}
	mock.ExpectQuery(getProxyQuery).
		WithArgs(toDecimalIP(ip)).
		WillReturnRows(
			sqlmock.NewRows(columns).
				FromCSVString(string(b)))

	c := client{db}
	p, err := c.GetProxy(ip)
	assert.Nil(t, err, "error in GetProxy(%v)", ip)
	assert.Equal(t, ip, p.AddressFrom)
}

func TestClient_GetProxies(t *testing.T) {
	client, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}

	countryCode, limit := "AR", 50
	p, err := client.GetProxies(countryCode, limit)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
}

func TestClient_GetProxies_Mock(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err, "error opening mock database connection")
	defer db.Close()
	defer func() {
		assert.Nil(t, mock.ExpectationsWereMet(), "there were unfulfilled expectations")
	}()
	b, err := ioutil.ReadFile("testdata/get_proxies.csv")
	assert.Nil(t, err, "couldn't read csv file")
	columns := []string{"ip_from", "ip_to", "proxy_type", "country_code", "country_name", "region_name", "city_name",
		"isp", "domain", "usage_type", "asn", "as"}
	countryCodeUpper, limit := "AR", 50
	mock.ExpectQuery(getProxiesQuery).
		WithArgs(countryCodeUpper, limit).
		WillReturnRows(
			sqlmock.NewRows(columns).
				FromCSVString(string(b)))
	c := client{db}

	countryCodeLower, limit := strings.ToLower(countryCodeUpper), 50
	proxies, err := c.GetProxies(countryCodeLower, limit)
	assert.Nil(t, err, "error in GetProxies(%v)", countryCodeLower, limit)
	assert.Equal(t, 2, len(proxies))
	assert.Equal(t, decimalIP("281471083157321"), toDecimalIP(proxies[0].AddressFrom))
	assert.Equal(t, "Telecom Argentina S.A.", *proxies[1].AS)
}

func TestClient_GetISPs(t *testing.T) {
	client, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}

	countryCode := "FR"
	ipsNames, err := client.GetISPs(countryCode)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(ipsNames), ipsNames)
}

func TestClient_GetISPs_Mock(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err, "error opening mock database connection")
	defer db.Close()
	defer func() {
		assert.Nil(t, mock.ExpectationsWereMet(), "there were unfulfilled expectations")
	}()
	b, err := ioutil.ReadFile("testdata/get_isp.csv")
	assert.Nil(t, err, "couldn't read csv file")
	columns := []string{"isp"}
	countryCode := "AR"
	mock.ExpectQuery(getISPsQuery).
		WithArgs(countryCode).
		WillReturnRows(
			sqlmock.NewRows(columns).
				FromCSVString(string(b)))
	c := client{db}

	ispNames, err := c.GetISPs(countryCode)
	assert.Nil(t, err, "error in GetISPs(%v)", countryCode)
	assert.Equal(t, 10, len(ispNames))
	assert.Equal(t, "Azul Networks S.R.L", ispNames[9])
}

func TestClient_GetIPCount_Mock(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err, "error opening mock database connection")
	defer db.Close()
	defer func() {
		assert.Nil(t, mock.ExpectationsWereMet(), "there were unfulfilled expectations")
	}()
	columns := []string{"ip_count"}
	countryCodeUpper := "AR"
	expected := 89780
	mock.ExpectQuery(getIPCountQuery).
		WithArgs(countryCodeUpper).
		WillReturnRows(
			sqlmock.NewRows(columns).
				FromCSVString(fmt.Sprintf("%d", expected)))
	c := client{db}

	countryCodeLower := strings.ToLower(countryCodeUpper)
	count, err := c.GetIPCount(countryCodeLower)
	assert.Nil(t, err, "error in GetIPCount(%v)", countryCodeLower)
	assert.Equal(t, expected, count)
}

func TestClient_GetIPCount(t *testing.T) {
	client, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}

	countryCode := "AR"
	count, err := client.GetIPCount(countryCode)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func TestClient_TopProxyTypes(t *testing.T) {
	client, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}

	limit := 3
	proxyTypes, err := client.TopProxyTypes(limit)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(proxyTypes), proxyTypes)
}

func TestClient_TopProxyTypes_Mock(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err, "error opening mock database connection")
	defer db.Close()
	defer func() {
		assert.Nil(t, mock.ExpectationsWereMet(), "there were unfulfilled expectations")
	}()
	b, err := ioutil.ReadFile("testdata/top_proxy_types.csv")
	assert.Nil(t, err, "couldn't read csv file")
	columns := []string{"proxy_type", "proxy_type_count"}
	limit := 3
	mock.ExpectQuery(topProxyTypesQuery).
		WithArgs(limit).
		WillReturnRows(
			sqlmock.NewRows(columns).
				FromCSVString(string(b)))
	c := client{db}

	proxyTypes, err := c.TopProxyTypes(limit)
	assert.Nil(t, err, "error in TopProxyTypes(%v)", limit)
	assert.Equal(t, []string{"PUB", "WEB"}, proxyTypes)
}
