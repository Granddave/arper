package database

import (
	"log"
	"os"
	"path"

	"github.com/granddave/arper/pkg/arp"
	"github.com/granddave/arper/pkg/utils"
)

type DatabaseContents struct {
	Hosts   []arp.Host
	Vendors map[string]string
}

type Database struct {
	Contents DatabaseContents
	Filepath string
}

func (db *Database) NumHosts() int {
	return len(db.Contents.Hosts)
}

func (db *Database) NumVendors() int {
	return len(db.Contents.Vendors)
}

func (db *Database) AddHost(host *arp.Host) {
	db.Contents.Hosts = append(db.Contents.Hosts, *host)
}

func (db *Database) HasHost(other arp.Host) bool {
	for _, host := range db.Contents.Hosts {
		if other.MAC.String() == host.MAC.String() {
			return true
		}
	}
	return false
}

func (db *Database) GetVendorIfExists(mac string) string {
	vendorPart := arp.GetVendorPart(mac)

	if vendor, exists := db.Contents.Vendors[vendorPart]; exists {
		return vendor
	}

	return ""
}

func (db *Database) AddVendor(mac string, vendor string) {
	db.Contents.Vendors[mac] = vendor
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

	db := Database{
		Contents: DatabaseContents{
			Hosts:   make([]arp.Host, 0),
			Vendors: make(map[string]string),
		},
		Filepath: databaseFilepath,
	}

	if err := utils.Deserialize(&db.Contents, db.Filepath); err != nil {
		log.Printf("Failed to deserialize database: %v", err)

		log.Printf("Removing and recreating")
		// TODO: Make backup
		initializeDatabaseFile(databaseFilepath)
	}

	numHosts := db.NumHosts()
	if numHosts == 1 {
		log.Printf("Parsed 1 host")
	} else if numHosts > 1 {
		log.Printf("Parsed %v hosts", numHosts)
	}
	numVendors := db.NumVendors()
	if numVendors == 1 {
		log.Printf("Parsed 1 vendor")
	} else if numHosts > 1 {
		log.Printf("Parsed %v vendors", numVendors)
	}

	return &db
}

func (db *Database) Save() {
	if err := utils.Serialize(db.Contents, db.Filepath); err != nil {
		log.Printf("Failed to serialize database: %v", err)
	}
}
