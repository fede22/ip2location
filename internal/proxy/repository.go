package proxy

type Repository interface {
	GetProxies(countryCode string, limit int) ([]Proxy, error)
	GetProxy(address NetIP) (Proxy, error)
	GetISPs(countryCode string) ([]string, error)
	GetIPCount(countryCode string) (int, error)
	TopProxyTypes(limit int) ([]string, error)
}