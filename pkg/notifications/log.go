package notifications

import (
	"log"

	"github.com/granddave/arper/pkg/arp"
)

type LogNotifier struct {
}

func (n *LogNotifier) Name() string {
	return "log"
}

func (n *LogNotifier) NotifyNewHost(host *arp.Host) error {
	log.Printf("New host: %v", host)
	return nil
}
