package main

import "net"

type Host struct {
	MAC      net.HardwareAddr
	IP       net.IP
	Hostname string
}

type HostCollection struct {
	Hosts []Host
}

func (collection *HostCollection) Len() int {
	return len(collection.Hosts)
}

func (collection *HostCollection) AddHost(host Host) {
	host.Hostname = TryGetHostname(host.IP)
	collection.Hosts = append(collection.Hosts, host)
}

func (collection *HostCollection) HasHost(newHost Host) bool {
	hasHost := false
	for _, host := range collection.Hosts {
		if newHost.MAC.String() == host.MAC.String() {
			hasHost = true
		}
	}
	return hasHost
}
