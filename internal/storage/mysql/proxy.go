package db

//TODO use net.IP for addresses
type Proxy struct {
	AddressFrom string
	AddressTo   string
	CountryCode *string
	CountryName *string
	RegionName  *string
	CityName    *string
	ISP         *string
	ProxyType   *string
	Domain      *string
	UsageType   *string
	ASN         *int
	AS          *string
}
