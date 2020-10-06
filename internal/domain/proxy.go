package domain

//TODO use net.IP for addresses
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