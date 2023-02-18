package main

import (
	"log"
	"os"
	"path"
)

func initializeDatabaseFile(filepath string) bool {
	dir := path.Dir(filepath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Println("Error creating directory:", err)
		return false
	}

	if _, err := os.Stat(filepath); err != nil {
		if err := os.WriteFile(filepath, []byte("{}"), 0644); err != nil {
			log.Printf("Error creating database file '%v': %v", filepath, err)
			return false
		}
	}

	return true
}

func NewHostCollection(databaseFilepath string) *HostCollection {
	initializeDatabaseFile(databaseFilepath)

	hostCollection := HostCollection{}

	if err := Deserialize(&hostCollection, databaseFilepath); err != nil {
		log.Printf("Failed to deserialize database: %v", err)

		log.Printf("Removing and recreating")
		// TODO: Make backup and recreate
		initializeDatabaseFile(databaseFilepath)

		return nil
	}

	numHosts := hostCollection.Len()
	if numHosts == 1 {
		log.Printf("Parsed 1 host")
	} else {
		log.Printf("Parsed %v hosts", numHosts)
	}

	return &hostCollection
}

func SaveHostCollection(hostCollection *HostCollection, databaseFilepath string) {
	if err := Serialize(hostCollection, databaseFilepath); err != nil {
		log.Printf("Failed to serialize database: %v", err)
	}
}
