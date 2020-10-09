package rest

import (
	"github.com/fede22/ip2location/internal/http/rest/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockService(ctrl)
	r := NewRouter(m)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

