package arp

import (
	"fmt"
	"net"
	"time"
)

type Host struct {
	MAC       net.HardwareAddr
	IP        net.IP
	Hostname  string
	Timestamp time.Time
}

func NewHost(MAC net.HardwareAddr, IP net.IP) *Host {
	return &Host{MAC: MAC, IP: IP, Hostname: TryGetHostname(IP), Timestamp: time.Now()}
}

func (h Host) String() string {
	return fmt.Sprintf("MAC=%v IP=%v Hostname='%v'", h.MAC, h.IP, h.Hostname)
}

func (h Host) NotificationText() string {
	hostname := h.Hostname
	if hostname == "" {
		hostname = "-"
	}
	return fmt.Sprintf("**New host:** Timestamp=`%v`, MAC=`%v` IP=`%v` Hostname=`%v`", h.Timestamp.Format(time.RFC3339), h.MAC, h.IP, hostname)
}
