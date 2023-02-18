package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func collectArpPackets(config *Config, newHosts *[]Host, mu *sync.Mutex, cond *sync.Cond) {
	socket := CreateSocket(config.IfaceName)
	defer syscall.Close(socket)

	log.Printf("Listening for Arp responses on %v", config.IfaceName)

	for {
		var buffer [128]byte

		n, _, err := syscall.Recvfrom(socket, buffer[:], 0)
		if err != nil {
			log.Fatal(err)
		}

		if n < 14+28 {
			log.Fatalf("Not enough bytes read: %v", n)
		}

		header, err := ParseARPPacket(buffer[14 : 14+28])
		if err != nil {
			log.Fatal(err)
		}

		mu.Lock()
		*newHosts = append(*newHosts, Host{IP: header.SenderIP[:], MAC: header.SenderMAC[:]})
		cond.Signal()
		mu.Unlock()
	}
}

func consumeDiscoveredHosts(config *Config, newHosts *[]Host, mu *sync.Mutex, cond *sync.Cond) {
	hostCollection := NewHostCollection(config.DatabaseFilepath)

	for {
		mu.Lock()
		if len(*newHosts) == 0 {
			cond.Wait()
		}

		var host = (*newHosts)[0]
		*newHosts = (*newHosts)[1:]
		mu.Unlock()

		if !hostCollection.HasHost(host) {
			hostCollection.AddHost(host)
			log.Printf("New host: %v, total hosts: %v", host, hostCollection.Len())

			SaveHostCollection(hostCollection, config.DatabaseFilepath)
		}
	}
}

func main() {
	config := NewConfig()
	newHosts := make([]Host, 0)
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
