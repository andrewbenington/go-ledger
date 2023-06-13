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

func EmptyEntry() *Entry {
	return &Entry{
		Value:   -1,
		Balance: -1,
	}
}
