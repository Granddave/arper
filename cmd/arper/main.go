package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/granddave/arper/pkg/arp"
	"github.com/granddave/arper/pkg/config"
	"github.com/granddave/arper/pkg/database"
	"github.com/granddave/arper/pkg/notifications"
	"github.com/granddave/arper/pkg/utils"
)

func collectArpPackets(config *config.Config, newHosts *[]arp.Host, mu *sync.Mutex, cond *sync.Cond) {
	arpListener := arp.NewArpListener(config.Iface)
	defer arpListener.CloseSocket()

	log.Printf("Listening for Arp responses on %v", config.Iface)

	for {
		newHost := arpListener.AwaitArpResponse()

		mu.Lock()
		*newHosts = append(*newHosts, *newHost)
		cond.Signal()
		mu.Unlock()
	}
}

func consumeDiscoveredHosts(config *config.Config, newHosts *[]arp.Host, mu *sync.Mutex, cond *sync.Cond) {
	db := database.NewDatabase(config.DatabaseFilepath)
	notificationManager := notifications.NewNotificationManager(config)

	for {
		mu.Lock()
		if len(*newHosts) == 0 {
			cond.Wait()
		}

		var host = (*newHosts)[0]
		*newHosts = (*newHosts)[1:]
		mu.Unlock()

		if !db.HasHost(host) {
			vendorPart := arp.GetVendorPart(host.MAC.String())
			if vendor := db.GetVendorIfExists(vendorPart); vendor != "" {
				log.Printf("Found cached vendor='%v'", vendor)
				host.Vendor = vendor
			} else {
				host.TryLookupVendor()
				if host.Vendor != "" {
					db.AddVendor(vendorPart, host.Vendor)
					log.Printf("Lookup vendor, found='%v'", host.Vendor)
				}
			}
			db.AddHost(&host)
			db.Save()
			notificationManager.NotifyNewHost(&host)
		}
	}
}

func main() {
	utils.InitLogging()

	config := config.NewConfig()
	newHosts := make([]arp.Host, 0)
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	go collectArpPackets(config, &newHosts, &mu, cond)
	go consumeDiscoveredHosts(config, &newHosts, &mu, cond)

	// Set up signal handler
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, syscall.SIGINT)
	select {
	case <-osChan:
		log.Println("Shutting down...")
	}
}
