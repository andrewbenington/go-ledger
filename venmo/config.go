package venmo

import "github.com/andrewbenington/go-ledger/ledger"

const (
	transferType = "Standard Transfer"
)

type Config struct {
	IgnoreTransfers bool `yaml:"ignore_transfers"`
}

func (c *Config) IgnoreEntry(e *ledger.Entry) bool {
	if e.SourceType != SourceType {
		return false
	}
	if c.IgnoreTransfers && e.Type == transferType {
		return true
	}
	return false
}
