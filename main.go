package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type NetworkHost struct {
	MAC      net.HardwareAddr
	IP       net.IP
	Hostname string
}

type NetworkHostCollection struct {
	Hosts []NetworkHost
}

func collectArpPackets(ifaceName string, newHosts *[]NetworkHost, mu *sync.Mutex, cond *sync.Cond) {
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
		*newHosts = append(*newHosts, NetworkHost{IP: header.SenderIP[:], MAC: header.SenderMAC[:]})
		cond.Signal()
		mu.Unlock()
	}
}

func consumeDiscoveredHosts(newHosts *[]NetworkHost, mu *sync.Mutex, cond *sync.Cond) {
	hosts := make([]NetworkHost, 0)

	for {
		mu.Lock()
		if len(*newHosts) == 0 {
			cond.Wait()
		}
		isNewHost := true
		for _, host := range hosts {
			if host.MAC.String() == (*newHosts)[0].MAC.String() {
				isNewHost = false
			}
		}
		if isNewHost {
			var newHost = (*newHosts)[0]
			newHost.Hostname = TryGetHostname(newHost.IP)
			hosts = append(hosts, newHost)
			log.Printf("New host: %v, total hosts: %v", newHost, len(hosts))
		}
		*newHosts = (*newHosts)[1:]
		mu.Unlock()
	}
}

func main() {
	var ifaceName string
	flag.StringVar(&ifaceName, "iface", "eth0", "network interface to use")
	flag.Parse()

	newHosts := make([]NetworkHost, 1)
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
