package notifications

import (
	"bytes"
	"net/http"

	"github.com/granddave/arper/pkg/arp"
)

type DiscordNotifier struct {
	WebhookURL string
}

func (n *DiscordNotifier) Name() string {
	return "discord"
}

func (n *DiscordNotifier) NotifyNewHost(host *arp.Host) error {
	jsonPayload := []byte(`{"content":"` + host.NotificationText() + `"}`)
	req, err := http.NewRequest("POST", n.WebhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
