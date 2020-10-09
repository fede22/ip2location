package domain

import (
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

type IP struct {
	Address     NetIP
	CountryName *string
	CityName    *string
}

type BigIntIP struct {
	*big.Int
}

func (n NetIP) ToBigIntIP() BigIntIP {
	return BigIntIP{big.NewInt(0).SetBytes(n.IP)}
}

func (bi BigIntIP) ToNetIP() NetIP {
	b := bi.Bytes()
	arr := [net.IPv6len]byte{}
	copy(arr[net.IPv6len-len(b):], b)
	return NetIP{arr[:]}
}

func (p Proxy) Netblock() ([]NetIP, error) {
	x := p.AddressFrom.ToBigIntIP()
	y := p.AddressTo.ToBigIntIP()
	ips := make([]NetIP, 0)
	for x.Cmp(y.Int) != 1 {
		ips = append(ips, x.ToNetIP())
		x.Add(x.Int, big.NewInt(1))
	}
	return ips, nil
}
