package main

import (
	"log"
	"net"
	"syscall"
)

func TryGetHostname(ip net.IP) string {
	names, err := net.LookupAddr(ip.String())

	if err != nil {
		return ""
	}

	return names[0]
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func CreateSocket(ifaceName string) int {
	// Open a raw socket for ARP packets
	rawSocket, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		log.Fatal(err)
	}

	// Bind the raw socket to the desired interface
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatal(err)
	}

	llAddr := syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ARP),
		Ifindex:  iface.Index,
	}

	err = syscall.Bind(rawSocket, &llAddr)
	if err != nil {
		log.Fatal(err)
	}

	return rawSocket
}
