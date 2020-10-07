package proxy

import (
	"net"
	"testing"
)

func TestIPv4ToDecimalAndBack(t *testing.T) {
	ip := NetIP{net.ParseIP("192.0.2.1")}
	ip2 := ip.ToBigIntIP().ToNetIP()
	if !ip.Equal(ip2.IP) {
		t.Errorf("expected ip %s, got instead %s", ip, ip2)
	}
}

func TestIPv6ToDecimalAndBack(t *testing.T) {
	ip := NetIP{net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")}
	ip2 := ip.ToBigIntIP().ToNetIP()
	if !ip.Equal(ip2.IP) {
		t.Errorf("expected ip %s, got instead %s", ip, ip2)
	}
}
