package domain

import (
	"fmt"
	"net"
)

type Service interface {
	GetProxy(address string) (Proxy, error)
	GetISPs(countryCode string) ([]string, error)
	GetIPs(countryCode string, limit int) ([]IP, error)
	GetIPCount(countryCode string) (int, error)
	GetTopProxyTypes(limit int) ([]string, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) service {
	return service{r}
}

func (ps service) GetProxy(address string) (Proxy, error) {
	ip := net.ParseIP(address)
	if ip == nil {
		return Proxy{}, fmt.Errorf("error parsing address %s as an IP address", address)
	}
	proxy, err := ps.r.GetProxy(NetIP{IP: ip})
	if err != nil {
		return Proxy{}, err
	}
	return proxy, nil
}

func (ps service) GetISPs(countryCode string) ([]string, error) {
	ispNames, err := ps.r.GetISPs(countryCode)
	if err != nil {
		return nil, err
	}
	return ispNames, nil
}

func (ps service) GetIPs(countryCode string, limit int) ([]IP, error) {
	proxies, err := ps.r.GetProxies(countryCode, limit)
	if err != nil {
		return nil, err
	}
	ips := make([]IP, 0)
	for _, p := range proxies {
		if len(ips) >= limit {
			break
		}
		nb, err := p.Netblock()
		if err != nil {
			return nil, err
		}
		for _, ip := range nb {
			ips = append(ips, IP{Address: ip, CountryName: p.CountryName, CityName: p.CityName})
		}
	}
	return ips[:min(len(ips), limit)], nil
}

func (ps service) GetIPCount(countryCode string) (int, error) {
	count, err := ps.r.GetIPCount(countryCode)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ps service) GetTopProxyTypes(limit int) ([]string, error) {
	proxyTypes, err := ps.r.TopProxyTypes(limit)
	if err != nil {
		return nil, err
	}
	return proxyTypes, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
