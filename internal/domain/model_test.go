package domain

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"testing"
)

func TestNetIP_ToBigIntIP_IPv4(t *testing.T) {
	ip := NetIP{net.ParseIP("192.0.2.1")}
	ip2 := ip.ToBigIntIP().ToNetIP()
	if !ip.Equal(ip2.IP) {
		t.Errorf("expected ip %s, got instead %s", ip, ip2)
	}
}

func TestNetIP_ToBigIntIP_IPv6(t *testing.T) {
	ip := NetIP{net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")}
	ip2 := ip.ToBigIntIP().ToNetIP()
	if !ip.Equal(ip2.IP) {
		t.Errorf("expected ip %s, got instead %s", ip, ip2)
	}
}

func TestProxy_Netblock(t *testing.T) {
	b, err := ioutil.ReadFile("mock/testdata/proxy_ar.json")
	if err != nil {
		t.Fatalf("error loading golden file: %s", err)
	}
	var p Proxy
	err = json.Unmarshal(b, &p)
	if err != nil {
		t.Fatalf("error unmarshaling proxy: %s", err)
	}
	nb, err := p.Netblock()
	if err != nil {
		t.Fatalf("error in NetBlock(): %s", err)
	}
	expected := []NetIP{
		{net.ParseIP("23.237.23.73")},
		{net.ParseIP("23.237.23.74")},
		{net.ParseIP("23.237.23.75")},
	}
	assert.Equal(t, expected, nb)
}
