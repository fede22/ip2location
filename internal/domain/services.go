package domain

type ProxyService struct {
	r ProxyRepository
}

// NewService creates an adding service with the necessary dependencies
func NewService(r ProxyRepository) ProxyService {
	return ProxyService{r}
}

func (ps ProxyService) GetByIP(address string) (Proxy, error) {
	proxy, err := ps.r.GetByIP(address)
	if err != nil {
		return Proxy{}, err
	}
	return proxy, nil
}

func (ps ProxyService) GetISPs(countryCode string) ([]string, error) {
	ipsNames, err := ps.r.GetISPs(countryCode)
	if err != nil {
		return nil, err
	}
	return ipsNames, nil
}

func (ps ProxyService) GetIPs(countryCode string, limit int) ([]IP, error) {
	proxies, err := ps.r.GetByCountryCode(countryCode, limit)
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
			ips = append(ips, IP{Address: ip.String(), CountryName: p.CountryName, CityName: p.CityName})
		}
	}
	return ips[:min(len(ips),limit)], nil
}

func (ps ProxyService) GetIPCount(countryCode string) (int, error) {
	count, err := ps.r.GetIPCount(countryCode)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ps ProxyService) GetTopProxyTypes(limit int) ([]string, error) {
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
