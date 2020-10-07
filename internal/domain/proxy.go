package domain

import (
	"fmt"
	"math/big"
	"net"
)

type NetIP struct {
	net.IP
}

type Proxy struct {
	AddressFrom NetIP
	AddressTo   NetIP
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

type DecimalIP string

//TODO move to another package
type IP struct {
	Address     string
	CountryName *string
	CityName    *string
}

type bigIntIP struct {
	*big.Int
}

func (dec DecimalIP) ToNetIP() (NetIP, error) {
	x, err := dec.toBigIntIP()
	if err != nil {
		return NetIP{}, err
	}
	return x.toNetIP(), nil
}

func (dec DecimalIP) toBigIntIP() (bigIntIP, error) {
	x, ok := big.NewInt(0).SetString(string(dec), 10)
	if !ok {
		return bigIntIP{}, fmt.Errorf("error parsing address %s to big int", dec)
	}
	return bigIntIP{x}, nil
}

func (n NetIP) ToDecimalIP() DecimalIP {
	return DecimalIP(n.toBigIntIP().String())
}

func (n NetIP) toBigIntIP() bigIntIP {
	return bigIntIP{big.NewInt(0).SetBytes(n.IP)}
}

func (bi bigIntIP) toNetIP() NetIP {
	b := bi.Bytes()
	arr := [net.IPv6len]byte{}
	copy(arr[net.IPv6len-len(b):], b)
	return NetIP{arr[:]}
}

func (p Proxy) Netblock() ([]net.IP, error) {
	x := p.AddressFrom.toBigIntIP()
	y := p.AddressTo.toBigIntIP()
	ips := make([]net.IP, 0)
	for x.Cmp(y.Int) != 1 {
		ips = append(ips, x.toNetIP().IP)
		x.Add(x.Int, big.NewInt(1))
	}
	return ips, nil
}

