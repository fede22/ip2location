package domain

import "net"

type ProxyRepository interface {
	GetProxies(countryCode string, limit int) ([]Proxy, error)
	GetProxy(address net.IP) (Proxy, error)
	GetISPs(countryCode string) ([]string, error)
	GetIPCount(countryCode string) (int, error)
	TopProxyTypes(limit int) ([]string, error)
}