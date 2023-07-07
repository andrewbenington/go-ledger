package ledger

import "time"

type Entry struct {
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
