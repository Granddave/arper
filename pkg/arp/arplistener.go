package arp

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"syscall"
)

type ArpListener struct {
	socketFd  int
	ifaceName string
}

func NewArpListener(ifaceName string) *ArpListener {
	arpListener := ArpListener{ifaceName: ifaceName}
	arpListener.createSocket()
	return &arpListener
}

func (al *ArpListener) createSocket() {
	rawSocket, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		log.Fatalf("Failed to open socket: %v", err)
	}

	iface, err := net.InterfaceByName(al.ifaceName)
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

	al.socketFd = rawSocket
}

func (al *ArpListener) CloseSocket() {
	syscall.Close(al.socketFd)
}

func (al *ArpListener) AwaitArpResponse() *Host {
	var buffer [128]byte

	n, _, err := syscall.Recvfrom(al.socketFd, buffer[:], 0)
	if err != nil {
		log.Fatal(err)
	}

	if n < 14+28 {
		log.Fatalf("Not enough bytes read: %v", n)
	}

	// Parse the Arp packet by skip the Ethernet header
	header, err := parseARPPacket(buffer[14 : 14+28])
	if err != nil {
		log.Fatal(err)
	}

	return NewHost(header.SenderMAC[:], header.SenderIP[:])
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

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

func parseARPPacket(packet []byte) (*ARPHeader, error) {
	var header ARPHeader

	buf := bytes.NewReader(packet)
	err := binary.Read(buf, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	return &header, nil
}
