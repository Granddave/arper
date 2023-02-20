package config

import (
	"flag"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		args     []string
		expected *Config
	}{
		{
			[]string{},
			&Config{
				Iface:             "eth0",
				DatabaseFilepath:  "/var/lib/arper/database.json",
				DiscordWebhookURL: "",
			},
		},
		{
			[]string{"-iface", "eth1", "-db", "/tmp/hosts.json", "https://discord.com/api/webhooks/891230"},
			&Config{
				Iface:             "eth1",
				DatabaseFilepath:  "/tmp/hosts.json",
				DiscordWebhookURL: "https://discord.com/api/webhooks/891230",
			},
		},
		{
			[]string{"-iface", "lo"},
			&Config{
				Iface:             "lo",
				DatabaseFilepath:  "/var/lib/arper/database.json",
				DiscordWebhookURL: "",
			},
		},
		{
			[]string{"-db", "/tmp/hosts2.json"},
			&Config{
				Iface:             "eth0",
				DatabaseFilepath:  "/tmp/hosts2.json",
				DiscordWebhookURL: "",
			},
		},
		{
			[]string{"-discord-webhook", "https://discord.com/api/webhooks/891230"},
			&Config{
				Iface:             "eth0",
				DatabaseFilepath:  "/var/lib/arper/database.json",
				DiscordWebhookURL: "https://discord.com/api/webhooks/891230",
			},
		},
	}

	for _, test := range tests {
		// Arrange
		origArgs := os.Args
		defer func() { os.Args = origArgs }()
		os.Args = append(os.Args[:1], test.args...)
		flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)

		// Act
		c := NewConfig()
		flag.CommandLine.Parse(test.args)

		// Assert
		if c.Iface != test.expected.Iface {
			t.Errorf("Expected iface to be %s, but got %s", test.expected.Iface, c.Iface)
		}
		if c.DatabaseFilepath != test.expected.DatabaseFilepath {
			t.Errorf("Expected db to be %s, but got %s", test.expected.DatabaseFilepath, c.DatabaseFilepath)
		}
	}
}
