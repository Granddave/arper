package main

import (
	"net"
	"testing"
)

func TestHostCollection(t *testing.T) {
	// Create a new empty host collection
	collection := HostCollection{}

	// Create a host to add to the collection
	host := Host{
		MAC:      net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
		IP:       net.ParseIP("192.168.1.1"),
		Hostname: "example.com",
	}

	// Test AddHost method
	collection.AddHost(host)
	if len(collection.Hosts) != 1 {
		t.Errorf("AddHost failed, expected collection length to be 1, got %d", len(collection.Hosts))
	}

	if collection.Hosts[0].Hostname == "example.com" {
		t.Errorf("AddHost failed, expected Hostname to be populated, got empty string")
	}

	// Test HasHost method
	if !collection.HasHost(host) {
		t.Errorf("HasHost failed, expected host to exist in collection, got false")
	}

	// Test Len method
	if collection.Len() != 1 {
		t.Errorf("Len failed, expected collection length to be 1, got %d", collection.Len())
	}
}
