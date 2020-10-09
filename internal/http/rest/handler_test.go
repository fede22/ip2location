package rest

import (
	"encoding/json"
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
	req, err := http.NewRequest("GET", "/country/ar/ip", nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var current []domain.IP
	err = json.Unmarshal(b, &current)
	assert.Nil(t, err, "error unmarshalling response IPs")
	assert.Equal(t, expected, current)
}


