package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func collectArpPackets(ifaceName string, newHosts *[]Host, mu *sync.Mutex, cond *sync.Cond) {
	socket := CreateSocket(ifaceName)
	defer syscall.Close(socket)

	log.Printf("Listening for Arp responses on %v", ifaceName)

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

func consumeDiscoveredHosts(newHosts *[]Host, mu *sync.Mutex, cond *sync.Cond) {
	hostCollection := HostCollection{}

	for {
		mu.Lock()
		if len(*newHosts) == 0 {
			cond.Wait()
		}

		var host = (*newHosts)[0]
		if !hostCollection.HasHost(host) {
			hostCollection.AddHost(host)
			log.Printf("New host: %v, total hosts: %v", host, hostCollection.Len())
		}

		*newHosts = (*newHosts)[1:]

		mu.Unlock()
	}
}

func main() {
	var ifaceName string
	flag.StringVar(&ifaceName, "iface", "eth0", "network interface to use")
	flag.Parse()

	newHosts := make([]Host, 0)
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	go collectArpPackets(ifaceName, &newHosts, &mu, cond)
	go consumeDiscoveredHosts(&newHosts, &mu, cond)

	// Set up signal handler
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, syscall.SIGINT)
	select {
	case <-osChan:
		log.Println("Shutting down...")
	}
}
