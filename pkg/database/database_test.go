package database

import (
	"net"
	"os"
	"path"
	"testing"
	"time"

	"github.com/granddave/arper/pkg/arp"
)

func TestInitializeDatabaseFile(t *testing.T) {
	// Arrange
	testDir := t.TempDir()
	testFile := path.Join(testDir, "test.db")

	// Act
	result := initializeDatabaseFile(testFile)

	// Assert
	if !result {
		t.Errorf("Test failed: initializeDatabaseFile returned false")
	}

	fileInfo, err := os.Stat(testFile)
	if err != nil {
		t.Errorf("Test failed: %v", err)
	}
	if !fileInfo.Mode().IsRegular() {
		t.Errorf("Test failed: file is not a regular file")
	}
	if fileInfo.Mode().Perm() != 0644 {
		t.Errorf("Test failed: file has incorrect permissions")
	}
}

func TestDatabase(t *testing.T) {
	// Arrange
	host := arp.Host{
		MAC:       net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
		IP:        net.ParseIP("192.0.2.1"),
		Hostname:  "example.com",
		Timestamp: time.Now(),
	}

	// Act
	// TODO: Test NewDatabase
	db := Database{}
	db.AddHost(&host)
	// TODO: Test Save()
	// TODO: Test NewDatabase() by reading serialized database

	// Assert
	if len(db.Contents.Hosts) != 1 {
		t.Errorf("AddHost failed, expected length to be 1, got %d", len(db.Contents.Hosts))
	}

	if !db.HasHost(host) {
		t.Errorf("HasHost failed, expected host to exist, got false")
	}

	if db.NumHosts() != 1 {
		t.Errorf("Len failed, expected length to be 1, got %d", db.NumHosts())
	}
}
