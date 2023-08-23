package config

import (
	"fmt"
	"os"
	"path/filepath"

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

const (
	configFolder = "config"
	configFile   = "config.yaml"
)

func (g *Config) read() error {
	_, err := os.Stat("config/config.yaml")
	if err != nil {
		err = initialize()
		if err != nil {
			return fmt.Errorf("initialize config: %w", err)
		}
	}
	yamlFile, err := os.ReadFile(filepath.Join(configFolder, configFile))
	if err != nil {
		return fmt.Errorf("read %s: %v ", filepath.Join(configFolder, configFile), err)
	}
	err = yaml.Unmarshal(yamlFile, g)
	if err != nil {
		return fmt.Errorf("parse config: %s", err)
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

func initialize() error {
	err := os.MkdirAll("config", 0755)
	if err != nil {
		return fmt.Errorf("mkdir 'config': %w", err)
	}
	return os.WriteFile("config/config.yaml", []byte{}, 0644)
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
