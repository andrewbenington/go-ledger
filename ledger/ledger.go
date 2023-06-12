package ledger

import "fmt"

type Ledger struct {
	Entries []LedgerEntry
}

func (l *Ledger) InsertEntries(entries []LedgerEntry) {
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

func Assemble(sourceEntries [][]LedgerEntry) []LedgerEntry {
	fmt.Printf("%d entry slices\n", len(sourceEntries))
	if len(sourceEntries) == 0 {
		return []LedgerEntry{}
	}
	allEntries := sourceEntries[0]
	for i := 1; i < len(sourceEntries); i++ {
		allEntries = mergeSortedEntrySlices(allEntries, sourceEntries[i])
	}
	return allEntries
}

func mergeSortedEntrySlices(slice1 []LedgerEntry, slice2 []LedgerEntry) []LedgerEntry {
	fmt.Println("Merging entries")
	sortedEntries := make([]LedgerEntry, len(slice1)+len(slice2))
	for i1, i2 := 0, 0; i1 < len(slice1) && i2 < len(slice2); {
		// fmt.Printf("%d, %d\n", i1, i2)
		if i2 == len(slice2) || slice1[i1].Date.Before(slice2[i2].Date) {
			if i2 < len(slice2) {
				fmt.Printf("%s < %s\n", slice1[i1].Date, slice2[i2].Date)
			}
			slice1[i1].Notes = fmt.Sprintf("i1: %d", i1)
			sortedEntries[i1+i2] = slice1[i1]
			i1++
		} else {
			slice2[i2].Notes = fmt.Sprintf("i2: %d", i2)
			sortedEntries[i1+i2] = slice2[i2]
			i2++
		}
	}
	fmt.Println(sortedEntries[0].Notes)
	return sortedEntries
}
