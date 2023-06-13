package ledger

type Ledger struct {
	Entries []Entry
}

func (l *Ledger) InsertEntries(entries []Entry) {
	l.Entries = append(l.Entries, entries...)
}

func (l Ledger) Len() int {
	return len(l.Entries)
}

func (l Ledger) Less(i, j int) bool {
	return l.Entries[i].Date.Before(l.Entries[j].Date)
}

func (l Ledger) Swap(i, j int) {
	l.Entries[i], l.Entries[j] = l.Entries[j], l.Entries[i]
}
