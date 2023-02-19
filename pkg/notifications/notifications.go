package notifications

import (
	"bytes"
	"log"
	"net/http"

	"github.com/granddave/arper/pkg/arp"
)

type Notifier struct {
	WebhookURL string
}

func NewNotifier(webhookUrl string) *Notifier {
	return &Notifier{webhookUrl}
}

func (n *Notifier) NotifyNewHost(host *arp.Host) {
	if n.WebhookURL == "" {
		return
	}

	err := n.send(host.NotificationText())
	if err != nil {
		log.Printf("Failed to send Discord notification: %v", err)
	}
}

func (n *Notifier) send(message string) error {
	jsonPayload := []byte(`{"content":"` + message + `"}`)
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
