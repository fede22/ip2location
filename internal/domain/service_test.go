package domain_test

import (
	"encoding/json"
	"fmt"
	"github.com/fede22/ip2location/internal/domain"
	"github.com/fede22/ip2location/internal/domain/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"testing"
)

//TODO use test tables
//TODO failure cases

func TestMockRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepository(ctrl)
	expected := fmt.Errorf("test error")
	m.EXPECT().GetProxies(gomock.Eq("AR"), gomock.Eq(50)).Return(nil, expected)

	s := domain.NewService(m)
	_, err := s.GetIPs("AR", 50)
	if err != expected {
		t.Errorf("expected error '%v', got instead '%v'", expected, err)
	}
}

func TestGetProxy(t *testing.T) {
	b, err := ioutil.ReadFile("mocks/testdata/proxy_ar.json")
	if err != nil {
		t.Fatalf("error loading golden file: %s", err)
	}
	var expected domain.Proxy
	err = json.Unmarshal(b, &expected)
	if err != nil {
		t.Fatalf("error unmarshaling proxy: %s", err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepository(ctrl)
	address := "23.237.23.73"
	m.EXPECT().GetProxy(
		gomock.Eq(domain.NetIP{net.ParseIP(address)})).Return(expected, nil)

	s := domain.NewService(m)
	p, err := s.GetProxy(address)
	assert.Nil(t, err, "error in GetProxy(%s)", address)
	assert.Equal(t, expected, p)
}

func TestGetIPs(t *testing.T) {
	b, err := ioutil.ReadFile("mocks/testdata/proxies_ar.json")
	if err != nil {
		t.Fatalf("error loading golden file: %s", err)
	}
	var proxies []domain.Proxy
	err = json.Unmarshal(b, &proxies)
	if err != nil {
		t.Fatalf("error unmarshaling proxy: %s", err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepository(ctrl)
	countryCode, limit := "AR", 50
	m.EXPECT().GetProxies(gomock.Eq(countryCode), gomock.Eq(limit)).Return(proxies, nil)

	expected := make([]domain.IP, 0)
	for _, p := range proxies {
		expected = append(expected, domain.IP{
			Address: p.AddressFrom,
			CountryName: p.CountryName,
			CityName: p.CityName,
		})
	}

	s := domain.NewService(m)
	ips, err := s.GetIPs(countryCode, limit)
	assert.Nil(t, err, "error in GetProxies(%v, %v)", countryCode, limit)
	assert.Equal(t, expected, ips)
}

