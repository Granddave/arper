package main

import (
	"net"
	"testing"
	"time"
)

func TestNewHost(t *testing.T) {
	// Arrange
	mac, _ := net.ParseMAC("00:11:22:33:44:55")
	ip := net.ParseIP("192.168.1.100")

	// Act
	host := NewHost(mac, ip)
	host.Hostname = "example.com"

	// Assert
	if host.MAC.String() != "00:11:22:33:44:55" {
		t.Errorf("unexpected MAC: got %q, want %q", host.MAC.String(), "00:11:22:33:44:55")
	}

	if host.IP.String() != "192.168.1.100" {
		t.Errorf("unexpected IP: got %q, want %q", host.IP.String(), "192.168.1.100")
	}

	if time.Since(host.Timestamp) > 1*time.Second {
		t.Errorf("unexpected Timestamp: got %v, want a value within the last second", host.Timestamp)
	}
}
