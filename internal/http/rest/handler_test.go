package rest

import (
	"encoding/json"
	"fmt"
	"github.com/fede22/ip2location/internal/domain"
	"github.com/fede22/ip2location/internal/http/rest/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := mock.NewMockService(ctrl)

	r := NewRouter(s)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestRouter_GetIPs(t *testing.T) {
	b, err := ioutil.ReadFile("mock/testdata/ip_ar.json")
	assert.Nil(t, err, "error loading golden file")
	var expected []domain.IP
	err = json.Unmarshal(b, &expected)
	assert.Nil(t, err, "error unmarshalling golden file IPs")
	assert.NotEqual(t, 0, len(expected))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := mock.NewMockService(ctrl)
	countryCode, limit := "ar", 50
	s.EXPECT().GetIPs(gomock.Eq(countryCode), gomock.Eq(limit)).Return(expected, nil)

	r := NewRouter(s)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/country/%s/ip", countryCode), nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var current []domain.IP
	err = json.Unmarshal(w.Body.Bytes(), &current)
	assert.Nil(t, err, "error unmarshalling response IPs")
	assert.Equal(t, expected, current)
}

func TestRouter_GetISPs(t *testing.T) {
	expected := []string{"FDCServers.net", "Telecom Argentina S.A."}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := mock.NewMockService(ctrl)
	countryCode := "ar"
	s.EXPECT().GetISPs(gomock.Eq(countryCode)).Return(expected, nil)

	r := NewRouter(s)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/country/%s/isp", countryCode), nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var current []string
	err = json.Unmarshal(w.Body.Bytes(), &current)
	assert.Nil(t, err, "error unmarshalling response ISPs")
	assert.Equal(t, expected, current)
}

func TestRouter_GetIPCount(t *testing.T) {
	expected := 4
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := mock.NewMockService(ctrl)
	countryCode := "ar"
	s.EXPECT().GetIPCount(gomock.Eq(countryCode)).Return(expected, nil)

	r := NewRouter(s)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/country/%s/ip_count", countryCode), nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var current struct {
		Count int `json:"ip_count"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &current)
	assert.Nil(t, err, "error unmarshalling response IP count")
	assert.Equal(t, expected, current.Count)
}

func TestRouter_GetProxy(t *testing.T) {
	b, err := ioutil.ReadFile("mock/testdata/proxy_ar.json")
	assert.Nil(t, err, "error loading golden file")
	var expected domain.Proxy
	err = json.Unmarshal(b, &expected)
	assert.Nil(t, err, "error unmarshalling golden file proxy")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := mock.NewMockService(ctrl)
	address := expected.AddressFrom.String()
	s.EXPECT().GetProxy(gomock.Eq(address)).Return(expected, nil)

	r := NewRouter(s)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/ip/%s", address), nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var current domain.Proxy
	err = json.Unmarshal(w.Body.Bytes(), &current)
	assert.Nil(t, err, "error unmarshalling response Proxy")
	assert.Equal(t, expected, current)
}

func TestRouter_GetTopProxyTypes(t *testing.T) {
	expected := []domain.ProxyType{{ProxyType: "PUB", Count: 12134}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := mock.NewMockService(ctrl)
	limit := 3
	s.EXPECT().GetTopProxyTypes(gomock.Eq(limit)).Return(expected, nil)

	r := NewRouter(s)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/top_proxy_types", nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var current []domain.ProxyType
	err = json.Unmarshal(w.Body.Bytes(), &current)
	assert.Nil(t, err, "error unmarshalling response proxy types")
	assert.Equal(t, expected, current)
}
