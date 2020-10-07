package db

import (
	"fmt"
	"github.com/ip2location/ip2proxy-go"
	"net"
	"testing"
)

const path = "../../ignore"

func TestIp2proxy_sampleDB(t *testing.T) {
	t.Log("running")

	db, err := ip2proxy.OpenDB(path + "/sample.bin.px7/IP2PROXY-IP-PROXYTYPE-COUNTRY-REGION-CITY-ISP-DOMAIN-USAGETYPE-ASN.BIN")

	if err != nil {
		return
	}
	ip := "199.83.103.79"
	all, err := db.GetAll(ip)

	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("ModuleVersion: %s\n", ip2proxy.ModuleVersion())
	fmt.Printf("PackageVersion: %s\n", db.PackageVersion())
	fmt.Printf("DatabaseVersion: %s\n", db.DatabaseVersion())

	fmt.Printf("isProxy: %s\n", all["isProxy"])
	fmt.Printf("ProxyType: %s\n", all["ProxyType"])
	fmt.Printf("CountryShort: %s\n", all["CountryShort"])
	fmt.Printf("CountryLong: %s\n", all["CountryLong"])
	fmt.Printf("RegionName: %s\n", all["RegionName"])
	fmt.Printf("CityName: %s\n", all["CityName"])
	fmt.Printf("ISP: %s\n", all["ISP"])
	fmt.Printf("Domain: %s\n", all["Domain"])
	fmt.Printf("UsageType: %s\n", all["UsageType"])
	fmt.Printf("ASN: %s\n", all["ASN"])
	fmt.Printf("AS: %s\n", all["AS"])
	fmt.Printf("LastSeen: %s\n", all["LastSeen"])
	fmt.Printf("Threat: %s\n", all["Threat"])

	db.Close()

}

func TestIp2proxy_realDB(t *testing.T) {
	t.Log("running")

	db, err := ip2proxy.OpenDB(path + "/IP2PROXY-LITE-PX7.BIN/IP2PROXY-LITE-PX7.BIN")

	if err != nil {
		return
	}
	ip := "1.0.132.50"
	all, err := db.GetAll(ip)

	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("ModuleVersion: %s\n", ip2proxy.ModuleVersion())
	fmt.Printf("PackageVersion: %s\n", db.PackageVersion())
	fmt.Printf("DatabaseVersion: %s\n", db.DatabaseVersion())

	fmt.Printf("isProxy: %s\n", all["isProxy"])
	fmt.Printf("ProxyType: %s\n", all["ProxyType"])
	fmt.Printf("CountryShort: %s\n", all["CountryShort"])
	fmt.Printf("CountryLong: %s\n", all["CountryLong"])
	fmt.Printf("RegionName: %s\n", all["RegionName"])
	fmt.Printf("CityName: %s\n", all["CityName"])
	fmt.Printf("ISP: %s\n", all["ISP"])
	fmt.Printf("Domain: %s\n", all["Domain"])
	fmt.Printf("UsageType: %s\n", all["UsageType"])
	fmt.Printf("ASN: %s\n", all["ASN"])
	fmt.Printf("AS: %s\n", all["AS"])
	fmt.Printf("LastSeen: %s\n", all["LastSeen"])
	fmt.Printf("Threat: %s\n", all["Threat"])

	db.Close()

}

func TestMySQL_localDB(t *testing.T) {
	client, err := NewClient()
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

func TestMySQL_GetIP(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	ip := net.ParseIP("1.0.4.1")
	p, err := client.GetProxy(ip)
	if err != nil {
		t.Fatal(err)
	}
	if !p.AddressFrom.Equal(ip) {
		t.Errorf("expected address_from %s, got instead %s", ip, p.AddressFrom)
	}
	t.Log(p)
}

func TestMySQL_GetProxies(t *testing.T) {
	client, err := NewClient()
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

func TestMySQL_GetISPs(t *testing.T) {
	client, err := NewClient()
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

func TestMySQL_GetIPCount(t *testing.T) {
	client, err := NewClient()
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

func TestMySQL_TopProxyTypes(t *testing.T) {
	client, err := NewClient()
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

func TestMySQL_IPv4ToDecimalAndBack(t *testing.T) {
	ip := net.ParseIP("192.0.2.1")
	ip2, err := decimalToIP(ipToDecimal(ip))
	if err != nil {
		t.Fatal(err)
	}
	if !ip.Equal(ip2) {
		t.Errorf("expected ip %s, got instead %s", ip, ip2)
	}
}

func TestMySQL_IPv6ToDecimalAndBack(t *testing.T) {
	ip := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	ip2, err := decimalToIP(ipToDecimal(ip))
	if err != nil {
		t.Fatal(err)
	}
	if !ip.Equal(ip2) {
		t.Errorf("expected ip %s, got instead %s", ip, ip2)
	}
}