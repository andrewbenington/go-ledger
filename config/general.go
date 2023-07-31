package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/andrewbenington/go-ledger/chase"
	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/andrewbenington/go-ledger/venmo"
	"gopkg.in/yaml.v2"
)

type GeneralConfig struct {
	Chase chase.Config `yaml:"chase"`
	Venmo venmo.Config `yaml:"venmo"`
}

func (g *GeneralConfig) read() error {
	yamlFile, err := os.ReadFile("config/general.yaml")
	if err != nil {
		return fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, g)
	if err != nil {
		return fmt.Errorf("error parsing general config: %s", err)
	}
	return nil
}

func IgnoreEntry(e *ledger.Entry) bool {
	cfg := GetConfig().General
	if cfg.Chase.IgnoreVenmo && strings.Contains(e.Memo, "VENMO") {
		return true
	}
	if cfg.Venmo.IgnoreTransfers && strings.Contains(e.Memo, "Standard Transfer") {
		return true
	}
	return false
}
