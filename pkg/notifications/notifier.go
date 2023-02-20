package notifications

import (
	"github.com/granddave/arper/pkg/arp"
)

type Notifier interface {
	NotifyNewHost(host *arp.Host) error
	Name() string
}
