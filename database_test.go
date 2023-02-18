package main

import (
	"os"
	"path"
	"testing"
)

func TestInitializeDatabaseFile(t *testing.T) {
	testDir := t.TempDir()
	testFile := path.Join(testDir, "test.db")

	// Call the function to initialize the test file
	result := initializeDatabaseFile(testFile)

	// Check that the file was created and has the expected permissions
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

	// Check that the function returned true
	if !result {
		t.Errorf("Test failed: initializeDatabaseFile returned false")
	}
}
