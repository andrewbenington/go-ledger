package chase

import (
	"strings"

	"github.com/andrewbenington/go-ledger/ledger"
)

type Config struct {
	IgnoreVenmo bool `yaml:"ignore_venmo"`
}

func (c *Config) IgnoreEntry(e *ledger.Entry) bool {
	if e.SourceType != SourceType {
		return false
	}
	if c.IgnoreVenmo && strings.Contains(e.Memo, "VENMO") {
		return true
	}
	return false
}
