package mysql

import (
	"bufio"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fede22/ip2location/internal/domain"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
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
	columns := []string{"ip_from", "ip_to", "proxy_type", "country_code", "country_name", "region_name", "city_name",
		"isp", "domain", "usage_type", "asn", "as"}
	ip := domain.NetIP{IP: net.ParseIP("1.0.4.1")}

	file, err := os.Open("testdata/get_proxy.csv")
	assert.Nil(t, err, "couldn't open csv file")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mock.ExpectQuery(getProxyQuery).
			WithArgs(toDecimalIP(ip)).
			WillReturnRows(
				sqlmock.NewRows(columns).
					FromCSVString(scanner.Text()))
	}
	assert.Nil(t, scanner.Err(), "non EOF error encountered by scanner")

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

func TestClient_GetISPs(t *testing.T) {
	client, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}

	countryCode := "FR"
	p, err := client.GetISPs(countryCode)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(p), p)
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
