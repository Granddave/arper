package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/granddave/arper/pkg/arp"
	"github.com/granddave/arper/pkg/config"
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
	database := arp.NewDatabase(config.DatabaseFilepath)
	notificationManager := notifications.NewNotificationManager(config)

	for {
		mu.Lock()
		if len(*newHosts) == 0 {
			cond.Wait()
		}

		var host = (*newHosts)[0]
		*newHosts = (*newHosts)[1:]
		mu.Unlock()

		if !database.HasHost(host) {
			database.AddHost(host)
			database.Save()
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
