package config

import (
	"github.com/andrewbenington/go-ledger/ledger"
)

type VenmoConfig struct {
	IgnoreTransfers bool `yaml:"ignore_transfers"`
}

func (c *VenmoConfig) IgnoreEntry(e *ledger.Entry) bool {
	if c.IgnoreTransfers && e.Type == "Standard Transfer" {
		return true
	}
	return false
}
