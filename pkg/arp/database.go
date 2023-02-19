package arp

import (
	"log"
	"os"
	"path"

	"github.com/granddave/arper/pkg/utils"
)

type Database struct {
	Hosts    []Host
	Filepath string
}

func (db *Database) Len() int {
	return len(db.Hosts)
}

func (db *Database) AddHost(host Host) {
	db.Hosts = append(db.Hosts, host)
}

func (db *Database) HasHost(other Host) bool {
	for _, host := range db.Hosts {
		if other.MAC.String() == host.MAC.String() {
			return true
		}
	}
	return false
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

	if err := utils.Deserialize(&db.Hosts, db.Filepath); err != nil {
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
	if err := utils.Serialize(db.Hosts, db.Filepath); err != nil {
		log.Printf("Failed to serialize database: %v", err)
	}
}
