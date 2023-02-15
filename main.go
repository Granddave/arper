package main

import (
	"flag"
	"log"
	"net"
	"sync"
	"syscall"
)

type NetworkHost struct {
	MAC      net.HardwareAddr
	IP       net.IP
	Hostname string
}

func collectArpPackets(ifaceName string, list *[]NetworkHost, mu *sync.Mutex, cond *sync.Cond) {
	socket := CreateSocket(ifaceName)
	defer syscall.Close(socket)

	// Listen for incoming ARP packets
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
		*list = append(*list, NetworkHost{IP: header.SenderIP[:], MAC: header.SenderMAC[:]})
		cond.Signal()
		mu.Unlock()
	}
}

func consumeArpPackets(list *[]NetworkHost, mu *sync.Mutex, cond *sync.Cond) {
	allHosts := make([]NetworkHost, 0)
	for {
		mu.Lock()
		if len(*list) == 0 {
			cond.Wait()
		}
		var isNewHost bool = true
		for _, host := range allHosts {
			if host.MAC.String() == (*list)[0].MAC.String() {
				isNewHost = false
			}
		}
		if isNewHost {
			var newHost = (*list)[0]
			newHost.Hostname = TryGetHostname(newHost.IP)
			allHosts = append(allHosts, newHost)
			log.Printf("New host: %v, total hosts: %v", newHost, len(allHosts))
		}
		*list = (*list)[1:]
		mu.Unlock()
	}
}

func main() {
	var ifaceName string
	flag.StringVar(&ifaceName, "iface", "eth0", "network interface to use")
	flag.Parse()

	newlyFoundHosts := make([]NetworkHost, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup
	cond := sync.NewCond(&mu)

	wg.Add(2)

	go func() {
		defer wg.Done()
		collectArpPackets(ifaceName, &newlyFoundHosts, &mu, cond)
	}()

	go func() {
		defer wg.Done()
		consumeArpPackets(&newlyFoundHosts, &mu, cond)
	}()

	wg.Wait()
}
