package arp

import (
	"net"
)

func TryGetHostname(ip net.IP) string {
	names, err := net.LookupAddr(ip.String())

	if err != nil {
		return ""
	}

	return names[0]
}
