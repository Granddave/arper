package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"os"
	"testing"
)

func TestSerializeDeserializeHostCollection(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "host-collection-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Create a HostCollection instance for testing
	originalHosts := HostCollection{
		Hosts: []Host{
			{MAC: net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}, IP: net.ParseIP("192.168.1.1"), Hostname: "host1"},
			{MAC: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, IP: net.ParseIP("192.168.1.2"), Hostname: "host2"},
		},
	}

	// Serialize the HostCollection to the temporary file
	err = Serialize(originalHosts, tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Deserialize the HostCollection from the temporary file
	deserializedHosts := HostCollection{}
	if err := Deserialize(&deserializedHosts, tempFile.Name()); err != nil {
		t.Fatal(err)
	}

	// Check that the original and deserialized HostCollections are equal
	if len(originalHosts.Hosts) != len(deserializedHosts.Hosts) {
		t.Errorf("Host collection lengths do not match: %d != %d", len(originalHosts.Hosts), len(deserializedHosts.Hosts))
	}

	for i := range originalHosts.Hosts {
		original := &originalHosts.Hosts[i]
		deserialized := &deserializedHosts.Hosts[i]
		if !bytes.Equal(original.MAC, deserialized.MAC) {
			t.Errorf("MAC addresses do not match: %s != %s", original.MAC.String(), deserialized.MAC.String())
		}
		if !original.IP.Equal(deserialized.IP) {
			t.Errorf("IP addresses do not match: %s != %s", original.IP.String(), deserialized.IP.String())
		}
		if original.Hostname != deserialized.Hostname {
			t.Errorf("Hostnames do not match: %s != %s", original.Hostname, deserialized.Hostname)
		}
	}
}
