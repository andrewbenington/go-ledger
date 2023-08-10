package config

import (
	"fmt"
	"os"

	"github.com/andrewbenington/go-ledger/ledger"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Chase ChaseConfig `yaml:"chase"`
	Venmo VenmoConfig `yaml:"venmo"`
}

var (
	config Config
)

func (g *Config) read() error {
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, g)
	if err != nil {
		return fmt.Errorf("error parsing general config: %s", err)
	}
	return nil
}

func GetConfig() Config {
	err := ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read config: %s", err)
		os.Exit(1)
	}
	return config
}

func ReadConfig() error {
	err := config.read()
	if err != nil {
		return fmt.Errorf("read config: %s", err)
	}
	return nil
}

func (c *Config) IgnoreEntry(e *ledger.Entry) bool {
	if e.SourceType == "CHASE" {
		return c.Chase.IgnoreEntry(e)
	}
	if e.SourceType == "VENMO" {
		return c.Venmo.IgnoreEntry(e)
	}
	return false
}
