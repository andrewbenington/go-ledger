package ledger

import (
	"fmt"
	"sort"
)

type Ledger struct {
	entries  []Entry
	entryMap map[string]*Entry
}

func (l *Ledger) InsertEntries(entries []Entry) {
	if l.entryMap == nil {
		l.entryMap = make(map[string]*Entry)
	}
	for _, e := range entries {
		if existingEntry, ok := l.entryMap[e.ID]; !ok {
			newEntry := e
			l.entries = append(l.entries, newEntry)
			l.entryMap[newEntry.ID] = &newEntry
		} else if existingEntry.Label == "" && e.Label != "" {
			existingEntry.Label = e.Label
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
