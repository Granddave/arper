package arp

import (
	"log"
	"os"
	"path"

	"github.com/granddave/arper/pkg/utils"
)

type Database struct {
	Contents struct {
		Hosts   []Host
		Vendors map[string]string
	}
	Filepath string
}

func (db *Database) Len() int {
	return len(db.Contents.Hosts)
}

func (db *Database) AddHost(host *Host) {
	host.TryLookupHostname()
	db.setVendorForHost(host)

	db.Contents.Hosts = append(db.Contents.Hosts, *host)
}

func (db *Database) HasHost(other Host) bool {
	for _, host := range db.Contents.Hosts {
		if other.MAC.String() == host.MAC.String() {
			return true
		}
	}
	return false
}

func (db *Database) setVendorForHost(host *Host) {
	vendorPart := GetVendorPart(host.MAC.String())
	if vendor, exists := db.Contents.Vendors[vendorPart]; exists {
		host.Vendor = vendor
	} else {
		host.TryLookupVendor()
		if host.Vendor != "" {
			db.Contents.Vendors[vendorPart] = host.Vendor
		}
	}
}

func initializeDatabaseFile(filepath string) bool {
	dir := path.Dir(filepath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Println("Error creating directory:", err)
		return false
	}

	if _, err := os.Stat(filepath); err != nil {
		log.Printf("Database doesn't exist, creating '%v'", filepath)
		if err := os.WriteFile(filepath, []byte("[]"), 0644); err != nil {
			log.Printf("Error creating database file '%v': %v", filepath, err)
			return false
		}
	}

	return true
}

func NewDatabase(databaseFilepath string) *Database {
	initializeDatabaseFile(databaseFilepath)

	db := Database{Filepath: databaseFilepath}

	if err := utils.Deserialize(&db.Contents, db.Filepath); err != nil {
		log.Printf("Failed to deserialize database: %v", err)

		log.Printf("Removing and recreating")
		// TODO: Make backup and recreate
		initializeDatabaseFile(databaseFilepath)

		return nil
	}

	numHosts := db.Len()
	if numHosts == 1 {
		log.Printf("Parsed 1 host")
	} else if numHosts > 1 {
		log.Printf("Parsed %v hosts", numHosts)
	}

	return &db
}

func (db *Database) Save() {
	if err := utils.Serialize(db.Contents.Hosts, db.Filepath); err != nil {
		log.Printf("Failed to serialize database: %v", err)
	}
}
