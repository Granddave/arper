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
	socket := arp.CreateSocket(config.Iface)
	defer syscall.Close(socket)

	log.Printf("Listening for Arp responses on %v", config.Iface)

	for {
		var buffer [128]byte

		n, _, err := syscall.Recvfrom(socket, buffer[:], 0)
		if err != nil {
			log.Fatal(err)
		}

		if n < 14+28 {
			log.Fatalf("Not enough bytes read: %v", n)
		}

		header, err := arp.ParseARPPacket(buffer[14 : 14+28])
		if err != nil {
			log.Fatal(err)
		}

		mu.Lock()
		*newHosts = append(*newHosts, *arp.NewHost(header.SenderMAC[:], header.SenderIP[:]))
		cond.Signal()
		mu.Unlock()
	}
}

func consumeDiscoveredHosts(config *config.Config, newHosts *[]arp.Host, mu *sync.Mutex, cond *sync.Cond) {
	database := arp.NewDatabase(config.DatabaseFilepath)
	notifier := notifications.NewNotifier(config.DiscordWebhookURL)

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
			log.Printf("New host (%v): %v", database.Len(), host)
			notifier.NotifyNewHost(&host)
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
