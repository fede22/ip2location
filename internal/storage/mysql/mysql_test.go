package db

import (
	"fmt"
	"github.com/ip2location/ip2proxy-go"
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
		t.Log(from, to)
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMySQL_GetIP(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	ip := "281470698521601"
	p, err := client.GetByIP(ip)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
}

func TestMySQL_GetByCountryCode(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	countryCode, limit := "AR", 50
	p, err := client.GetByCountryCode(countryCode, limit)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
}

func TestMySQL_GetISPNames(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	countryCode := "FR"
	p, err := client.GetISPNames(countryCode)
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
