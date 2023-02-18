package main

import (
	"flag"
	"log"
	"reflect"
)

const DefaultConfigFilepath = "/etc/arper/config.json"

type Config struct {
	Iface            string
	DatabaseFilepath string
}

func DefaultConfig() *Config {
	return &Config{
		Iface:            "eth0",
		DatabaseFilepath: "/var/lib/arper/hosts.json",
	}
}

func NewConfig() *Config {
	config := DefaultConfig()
	config.readConfigFromFile(DefaultConfigFilepath)
	config.parseFlags()

	config.logCurrentConfigs()

	return config
}

func (c *Config) readConfigFromFile(configFilename string) {
	// TODO: Implement
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.Iface, "iface", c.Iface, "network interface to use")
	flag.StringVar(&c.DatabaseFilepath, "db", c.DatabaseFilepath, "filepath to database")
	flag.Parse()
}

func (c *Config) logCurrentConfigs() {
	v := reflect.ValueOf(*c)
	t := v.Type()

	log.Println("Active configuration:")
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		log.Printf("  %s: %v\n", fieldName, field.Interface())
	}
}
