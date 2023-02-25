package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DefaultDatabaseFilepath = "/var/lib/arper/database.json"

type Config struct {
	Iface             string
	DatabaseFilepath  string
	DiscordWebhookURL string
}

func NewConfig() *Config {
	cfg := Config{}

	rootCmd := &cobra.Command{
		Use: "arper",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Config:\n%s", cfg.String())
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfg.Iface, "iface", "eth0", "network interface to use")
	rootCmd.PersistentFlags().StringVar(&cfg.DatabaseFilepath, "database-filepath", DefaultDatabaseFilepath, "path to the database file")
	rootCmd.PersistentFlags().StringVar(&cfg.DiscordWebhookURL, "discord-webhook-url", "", "Discord webhook URL")

	viper.SetConfigName("arper")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/arper")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %s", err)
	}

	return &cfg
}

func (c *Config) String() string {
	v := reflect.ValueOf(*c)
	t := v.Type()
	s := ""

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		s += fmt.Sprintf("  %s: %v\n", fieldName, field.Interface())
	}

	return s
}
