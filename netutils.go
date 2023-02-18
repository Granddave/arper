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
	rawSocket, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		log.Fatalf("Failed to open socket: %v", err)
	}

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("Failed to get interface by name: %v", err)
	}

	llAddr := syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ARP),
		Ifindex:  iface.Index,
	}

	err = syscall.Bind(rawSocket, &llAddr)
	if err != nil {
		log.Fatalf("Failed to bind socket to interface: %v", err)
	}

	return rawSocket
}
