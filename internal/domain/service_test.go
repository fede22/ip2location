package domain_test

import (
	"encoding/json"
	"fmt"
	"github.com/fede22/ip2location/internal/domain"
	"github.com/fede22/ip2location/internal/domain/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

//TODO use test tables
//TODO failure cases

func TestMockRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	expected := fmt.Errorf("test error")
	m.EXPECT().GetProxies(gomock.Eq("AR"), gomock.Eq(50)).Return(nil, expected)

	s := domain.NewService(m)
	_, err := s.GetIPs("AR", 50)
	if err != expected {
		t.Errorf("expected error '%v', got instead '%v'", expected, err)
	}
}

func TestService_GetProxy(t *testing.T) {
	b, err := ioutil.ReadFile("mock/testdata/proxy_ar.json")
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
	m := mock.NewMockRepository(ctrl)
	m.EXPECT().GetProxy(gomock.Eq(expected.AddressFrom)).Return(expected, nil)

	s := domain.NewService(m)
	p, err := s.GetProxy(expected.AddressFrom.String())
	assert.Nil(t, err, "error in GetProxy(%v)", expected.AddressFrom.String())
	assert.Equal(t, expected, p)
}

func TestService_GetIPs(t *testing.T) {
	b, err := ioutil.ReadFile("mock/testdata/proxies_ar.json")
	if err != nil {
		t.Fatalf("error loading golden file: %s", err)
	}
	var proxies []domain.Proxy
	err = json.Unmarshal(b, &proxies)
	if err != nil {
		t.Fatalf("error unmarshaling proxies: %s", err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	countryCode, limit := "AR", 50
	m.EXPECT().GetProxies(gomock.Eq(countryCode), gomock.Eq(limit)).Return(proxies, nil)

	expected := make([]domain.IP, 0)
	for _, p := range proxies {
		expected = append(expected, domain.IP{
			Address:     p.AddressFrom,
			CountryName: p.CountryName,
			CityName:    p.CityName,
		})
	}

	s := domain.NewService(m)
	ips, err := s.GetIPs(countryCode, limit)
	assert.Nil(t, err, "error in GetIPs(%v, %v)", countryCode, limit)
	assert.Equal(t, expected, ips)
}

func TestService_GetISPs(t *testing.T) {
	expected := []string{"FDCServers.net", "Telecom Argentina S.A."}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	countryCode := "AR"
	m.EXPECT().GetISPs(gomock.Eq(countryCode)).Return(expected, nil)

	s := domain.NewService(m)
	ispNames, err := s.GetISPs(countryCode)
	assert.Nil(t, err, "error in GetISPs(%v)", countryCode)
	assert.Equal(t, expected, ispNames)
}

func TestService_GetIPCount(t *testing.T) {
	expected := 4

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	countryCode := "AR"
	m.EXPECT().GetIPCount(gomock.Eq(countryCode)).Return(expected, nil)

	s := domain.NewService(m)
	count, err := s.GetIPCount(countryCode)
	assert.Nil(t, err, "error in GetIPCount(%v)", countryCode)
	assert.Equal(t, expected, count)
}

func TestService_GetTopProxyTypes(t *testing.T) {
	expected := []domain.ProxyType{{"PUB", 12345}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	limit := 3
	m.EXPECT().TopProxyTypes(gomock.Eq(limit)).Return(expected, nil)

	s := domain.NewService(m)
	proxyTypes, err := s.GetTopProxyTypes(limit)
	assert.Nil(t, err, "error in TopProxyTypes(%v)", limit)
	assert.Equal(t, expected, proxyTypes)
}
