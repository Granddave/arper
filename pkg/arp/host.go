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
	Vendor    string
	Timestamp time.Time
}

func NewHost(MAC net.HardwareAddr, IP net.IP) *Host {
	return &Host{
		MAC:       MAC,
		IP:        IP,
		Timestamp: time.Now(),
	}
}

func (h *Host) TryLookupHostname() {
	names, err := net.LookupAddr(h.IP.String())

	if err != nil {
		return
	}

	h.Hostname = names[0]
}

func (h *Host) TryLookupVendor() {
	h.Vendor = LookupVendorName(h.MAC.String())
}

func (h *Host) String() string {
	return fmt.Sprintf("MAC=%v IP=%v Vendor='%s' Hostname='%v'", h.MAC, h.IP, h.Vendor, h.Hostname)
}

func (h *Host) NotificationText() string {
	hostname := h.Hostname
	if hostname == "" {
		hostname = "-"
	}
	return fmt.Sprintf("**New host:** Timestamp=`%v` MAC=`%v` IP=`%v` Vendor='%s' Hostname=`%v`",
		h.Timestamp.Format(time.RFC3339),
		h.MAC,
		h.IP,
		h.Vendor,
		hostname,
	)
}
