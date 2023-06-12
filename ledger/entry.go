package ledger

import "time"

type LedgerEntry struct {
	ID      string
	Date    time.Time
	Source  string
	Person  string
	Memo    string
	Value   float64
	Type    string
	Balance float64
	Label   string
	Notes   string
}

func EmptyEntry() *LedgerEntry {
	return &LedgerEntry{
		Value:   -1,
		Balance: -1,
	}
}
