package domain

type ProxyRepository interface {
	GetByCountryCode(countryCode string, limit int) ([]Proxy, error)
	GetByIP(address string) (Proxy, error)
	GetISPs(countryCode string) ([]string, error)
	GetIPCount(countryCode string) (int, error)
	TopProxyTypes(limit int) ([]string, error)
}