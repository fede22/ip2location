package domain

import (
	"fmt"
	"math/big"
	"net"
)

type Proxy struct {
	AddressFrom string
	AddressTo   string
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
	x := new(big.Int)
	x, ok := x.SetString(p.AddressFrom, 10)
	if !ok {
		return nil, fmt.Errorf("error converting address %s to big int", p.AddressFrom)
	}
	y := new(big.Int)
	y, ok = y.SetString(p.AddressTo, 10)
	if !ok {
		return nil, fmt.Errorf("error converting address %s to big int", p.AddressTo)
	}
	ips := make([]net.IP, 0)
	for x.Cmp(y) != 1 {
		ips = append(ips, x.Bytes())
		x.Add(x, big.NewInt(1))
	}
	return ips, nil
}
