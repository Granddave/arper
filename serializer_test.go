package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"os"
	"testing"
)

func TestSerializeDeserialize(t *testing.T) {
	// Arrange
	tempFile, err := ioutil.TempFile("", "serializer")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	origHosts := []Host{
		{MAC: net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}, IP: net.ParseIP("192.168.1.1"), Hostname: "host1"},
		{MAC: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, IP: net.ParseIP("192.168.1.2"), Hostname: "host2"},
	}

	// Act
	err = Serialize(origHosts, tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	serializedHosts := []Host{}
	if err := Deserialize(&serializedHosts, tempFile.Name()); err != nil {
		t.Fatal(err)
	}

	// Assert
	if len(origHosts) != len(serializedHosts) {
		t.Errorf("Hosts lengths do not match: %d != %d", len(origHosts), len(serializedHosts))
	}

	for i := range origHosts {
		original := &origHosts[i]
		serialized := &serializedHosts[i]
		if !bytes.Equal(original.MAC, serialized.MAC) {
			t.Errorf("MAC addresses do not match: %s != %s", original.MAC.String(), serialized.MAC.String())
		}
		if !original.IP.Equal(serialized.IP) {
			t.Errorf("IP addresses do not match: %s != %s", original.IP.String(), serialized.IP.String())
		}
		if original.Hostname != serialized.Hostname {
			t.Errorf("Hostnames do not match: %s != %s", original.Hostname, serialized.Hostname)
		}
	}
}
