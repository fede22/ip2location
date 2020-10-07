package domain

import (
	"math/big"
	"net"
)

type Proxy struct {
	AddressFrom net.IP
	AddressTo   net.IP
	ProxyType   *string
	CountryCode *string
	CountryName *string
	RegionName  *string
	CityName    *string
	ISP         *string
	Domain      *string
	UsageType   *string
	ASN         *int
	AS          *string
}

//TODO move to another package
type IP struct {
	Address     string
	CountryName *string
	CityName    *string
}

func (p Proxy) Netblock() ([]net.IP, error) {
	x := big.NewInt(0).SetBytes(p.AddressFrom)
	y := big.NewInt(0).SetBytes(p.AddressTo)
	ips := make([]net.IP, 0)
	for x.Cmp(y) != 1 {
		ips = append(ips, x.Bytes())
		x.Add(x, big.NewInt(1))
	}
	return ips, nil
}
