package proxy_test

import (
	"encoding/json"
	"fmt"
	"github.com/fede22/ip2location/internal/proxy"
	"github.com/fede22/ip2location/internal/proxy/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"testing"
)

//TODO rename package to domain
//TODO use test tables

func TestMockRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepository(ctrl)
	expected := fmt.Errorf("test error")
	m.EXPECT().GetProxies(gomock.Eq("AR"), gomock.Eq(50)).Return(nil, expected)

	s := proxy.NewService(m)
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
	var expected proxy.Proxy
	err = json.Unmarshal(b, &expected)
	if err != nil {
		t.Fatalf("error unmarshaling proxy: %s", err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepository(ctrl)
	address := "23.237.23.73"
	m.EXPECT().GetProxy(
		gomock.Eq(proxy.NetIP{net.ParseIP(address)})).Return(expected, nil)

	s := proxy.NewService(m)
	p, err := s.GetProxy(address)
	assert.Nil(t, err, "error in GetProxy(%s)", address)
	assert.Equal(t, expected, p)
}