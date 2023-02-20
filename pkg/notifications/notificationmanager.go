package notifications

import (
	"log"

	"github.com/granddave/arper/pkg/arp"
	"github.com/granddave/arper/pkg/config"
)

type NotificationManager struct {
	notifiers []Notifier
}

func NewNotificationManager(config *config.Config) *NotificationManager {
	nm := &NotificationManager{}

	// Log notifier
	nm.notifiers = append(nm.notifiers, &LogNotifier{})

	// Discord notifier
	if config.DiscordWebhookURL != "" {
		nm.notifiers = append(nm.notifiers, &DiscordNotifier{config.DiscordWebhookURL})
	}

	return nm
}

func (nm *NotificationManager) NotifyNewHost(host *arp.Host) {
	for _, notifier := range nm.notifiers {
		err := notifier.NotifyNewHost(host)
		if err != nil {
			log.Printf("Failed to send notification (%v): %v", notifier.Name(), err)
		}
	}
}
