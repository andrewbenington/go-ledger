package ledger

import (
	"fmt"
	"regexp"
	"sort"
)

type Ledger struct {
	Year       int
	entries    []Entry
	entryMap   map[string]*Entry
	patternMap map[*regexp.Regexp][]*Entry
}

func (l *Ledger) InsertEntries(entries []Entry) {
	if l.entryMap == nil {
		l.entryMap = make(map[string]*Entry)
	}
	for _, e := range entries {
		if existingEntry, ok := l.entryMap[e.ID]; !ok {
			l.entries = append(l.entries, e)
			l.entryMap[e.ID] = &l.entries[len(l.entries)-1]
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
		entries, err := source.GetLedgerEntries(l.Year)
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

func (l *Ledger) TrackPattern(re *regexp.Regexp) {
	l.patternMap[re] = []*Entry{}
	for i := range l.entries {
		entry := &l.entries[i]
		if re.MatchString(entry.Memo) {
			l.patternMap[re] = append(l.patternMap[re], entry)
		}
	}
}

func (l *Ledger) EntriesWithPattern(re *regexp.Regexp) []*Entry {
	return l.patternMap[re]
}
