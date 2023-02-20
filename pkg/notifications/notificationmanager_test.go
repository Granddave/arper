package notifications

import (
	"net"
	"reflect"
	"testing"

	"github.com/granddave/arper/pkg/arp"
)

type mockNotifier struct {
	notifiedHosts []*arp.Host
}

func (mn *mockNotifier) NotifyNewHost(host *arp.Host) error {
	mn.notifiedHosts = append(mn.notifiedHosts, host)
	return nil
}

func (mn *mockNotifier) Name() string {
	return "mock"
}

func TestNotificationManager_NotifyAllNewHosts(t *testing.T) {
	// Arrange
	mock := &mockNotifier{}
	nm := &NotificationManager{
		notifiers: []Notifier{mock},
	}

	hosts := []*arp.Host{
		{
			IP:       net.ParseIP("192.168.1.1"),
			Hostname: "my-host-1",
			MAC:      net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		},
		{
			IP:       net.ParseIP("192.168.1.2"),
			Hostname: "my-host-2",
			MAC:      net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x56},
		},
	}

	// Act
	for _, host := range hosts {
		nm.NotifyNewHost(host)
	}

	// Assert
	for i, host := range hosts {
		if len(mock.notifiedHosts) <= i || !reflect.DeepEqual(mock.notifiedHosts[i], host) {
			t.Errorf("mockNotifier.NotifyNewHost not called with expected host at index %d", i)
		}
	}
}
