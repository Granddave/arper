package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"syscall"
)

// ARPHeader represents the header of an ARP packet
type ARPHeader struct {
	HardwareType uint16
	ProtocolType uint16
	HardwareSize uint8
	ProtocolSize uint8
	OpCode       uint16
	SenderMAC    [6]byte
	SenderIP     [4]byte
	TargetMAC    [6]byte
	TargetIP     [4]byte
}

func (h *ARPHeader) String() string {
	return fmt.Sprintf("Sender MAC: %s Sender IP: %s\nTarget MAC: %s Target IP: %s",
		net.HardwareAddr(h.SenderMAC[:]).String(),
		net.IP(h.SenderIP[:]).String(),
		net.HardwareAddr(h.TargetMAC[:]).String(),
		net.IP(h.TargetIP[:]).String(),
	)
}

// ParseARPPacket parses an ARP packet into an ARPHeader struct
func ParseARPPacket(packet []byte) (*ARPHeader, error) {
	var header ARPHeader

	// Read the header from the packet
	buf := bytes.NewReader(packet)
	err := binary.Read(buf, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	return &header, nil
}

func main() {
	// Open a raw socket for ARP packets
	rawSocket, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(rawSocket)

	// Bind the raw socket to the desired interface
	iface, err := net.InterfaceByName("enp0s31f6")
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

	// Listen for incoming ARP packets
	for {
		var buffer [1500]byte

		n, _, err := syscall.Recvfrom(rawSocket, buffer[:], 0)
		if err != nil {
			log.Fatal(err)
		}

		header, err := ParseARPPacket(buffer[:n])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received ARP packet: %v\n", header)
	}
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}
