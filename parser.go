package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
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
	return fmt.Sprintf("Sender MAC: %s Sender IP: %s Target MAC: %s Target IP: %s",
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
