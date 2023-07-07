package ledger

import (
	"fmt"
	"sort"
)

type Ledger struct {
	entries  []Entry
	entryMap map[string]int
}

func (l *Ledger) InsertEntries(entries []Entry) {
	if l.entryMap == nil {
		l.entryMap = make(map[string]int)
	}
	for _, entry := range entries {
		if i, ok := l.entryMap[entry.ID]; !ok {
			l.entryMap[entry.ID] = len(l.entries)
			l.entries = append(l.entries, entry)
		} else if l.entries[i].Label == "" && entry.Label != "" {
			l.entries[i].Label = entry.Label
			fmt.Println(l.entries[i].Label)

		}
	}
}

func (l Ledger) Len() int {
	return len(l.entries)
}

func (l Ledger) Less(i, j int) bool {
	return l.entries[i].Date.Before(l.entries[j].Date)
}

func (l Ledger) Swap(i, j int) {
	l.entries[i], l.entries[j] = l.entries[j], l.entries[i]
}

func (l *Ledger) UpdateFromSources(allSources []Source) error {
	for _, source := range allSources {
		fmt.Printf("getting entries from %s...\n", source.Name())
		entries, err := source.GetLedgerEntries()
		if err != nil {
			fmt.Printf("Error getting entries from %s: %e", source.Name(), err)
			continue
		}
		l.InsertEntries(entries)
	}
	sort.Sort(l)
	return nil
}

func (l *Ledger) Entries() []Entry {
	return l.entries
}
