package config

import (
	"strings"

	"github.com/andrewbenington/go-ledger/ledger"
)

type ChaseConfig struct {
	IgnoreTransfers bool `yaml:"ignore_transfers"`
	IgnoreVenmo     bool `yaml:"ignore_venmo"`
}

func (c *ChaseConfig) IgnoreEntry(e *ledger.Entry) bool {
	if c.IgnoreVenmo && strings.HasPrefix(e.Memo, "VENMO") {
		return true
	}
	if c.IgnoreTransfers && (e.Type == "LOAN_PMT" ||
		e.Type == "Payment" ||
		strings.HasPrefix(e.Memo, "CHASE CREDIT CRD") ||
		strings.HasPrefix(e.Memo, "Payment to Chase card")) {
		return true
	}
	return false
}
