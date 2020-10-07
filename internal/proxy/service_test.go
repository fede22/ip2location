package proxy_test

import (
	"fmt"
	"github.com/fede22/ip2location/internal/proxy"
	"github.com/fede22/ip2location/internal/proxy/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestMockRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepository(ctrl)
	expected := fmt.Errorf("wololo")
	m.EXPECT().GetProxies(gomock.Eq("AR"), gomock.Eq(50)).Return(nil, expected)

	s := proxy.NewService(m)
	_, err := s.GetIPs("AR", 50)
	if err != expected {
		t.Errorf("expected error '%v', got instead '%v'", expected, err)
	}
}
