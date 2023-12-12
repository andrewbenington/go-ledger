package excel

import (
	"fmt"
	"os"
	"time"

	"github.com/andrewbenington/go-ledger/config"
	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/andrewbenington/go-ledger/source"
	"github.com/andrewbenington/go-ledger/util"
	"github.com/xuri/excelize/v2"
)

// LedgerFromFile will return a pointer to a ledger.Ledger, including entries
// from the excel file named for the given year
func LedgerFromFile(year int) (*ledger.Ledger, error) {
	filename := fmt.Sprintf("%d.xlsx", year)
	l := &ledger.Ledger{Year: year}
	if _, err := os.Stat(filename); err != nil {
		// file doesn't exist yet
		return l, nil
	}
	file, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", filename, err)
	}
	err = addEntriesFromFile(l, file)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func addEntriesFromFile(l *ledger.Ledger, file *excelize.File) error {
	for _, sheet := range file.GetSheetList() {
		if sheet == "Overview" {
			continue
		}
		newEntries, err := ledgerEntriesFromSheet(file, sheet)
		if err != nil {
			return fmt.Errorf("read sheet %s: %w", sheet, err)
		}
		l.InsertEntries(newEntries)
	}
	return nil
}

func ledgerEntriesFromSheet(file *excelize.File, sheet string) ([]ledger.Entry, error) {
	cfg := config.GetConfig()
	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("error getting rows from sheet %s: %w", sheet, err)
	}
	entries := []ledger.Entry{}
	for _, row := range rows[1:] {
		entry, err := ledgerEntryFromRow(row)
		if err != nil {
			fmt.Printf("get entry from from %s row %v: %s\n", sheet, row, err)
			continue
		}
		if cfg.IgnoreEntry(entry) {
			continue
		}
		entries = append(entries, *entry)
	}
	return entries, nil
}

func ledgerEntryFromRow(row []string) (entry *ledger.Entry, err error) {
	entry = &ledger.Entry{}
	if len(row) > ledger.DateIndex {
		entry.Date, err = time.Parse("1/2/06 15:04", row[ledger.DateIndex])
		if err != nil {
			return nil, fmt.Errorf("Error parsing date: %s", err)
		}
	}
	if len(row) > ledger.MemoIndex {
		entry.Memo = util.NormalizeUnicode(row[ledger.MemoIndex])
	}
	if len(row) > ledger.ValueIndex {
		entry.Value, err = util.ParseMoneyAmount(row[ledger.ValueIndex])
		if err != nil {
			return nil, fmt.Errorf("Error parsing value amount: %s", err)
		}
	}
	if len(row) > ledger.BalanceIndex {
		entry.Balance, err = util.ParseMoneyAmount(row[ledger.BalanceIndex])
		if err != nil {
			return nil, fmt.Errorf("Error parsing balance amount: %s", err)
		}
	}
	if len(row) > ledger.TypeIndex {
		entry.Type = util.NormalizeUnicode(row[ledger.TypeIndex])
		if entry.Type == "" {
			entry.Type = "Other"
		}
	}
	if len(row) > ledger.SourceNameIndex {
		entry.SourceName = util.NormalizeUnicode(row[ledger.SourceNameIndex])
		if entry.SourceName == "" {
			entry.SourceName = "Other"
		}
	}
	if len(row) > ledger.SourceTypeIndex {
		entry.SourceType = util.NormalizeUnicode(row[ledger.SourceTypeIndex])
		if entry.SourceType == "" {
			entry.SourceType = "OTHER"
		}
	}
	if len(row) > ledger.PersonIndex {
		entry.Person = util.NormalizeUnicode(row[ledger.PersonIndex])
	}
	if len(row) > ledger.LabelIndex {
		entry.Label = util.NormalizeUnicode(row[ledger.LabelIndex])
	}
	if len(row) > ledger.NotesIndex {
		entry.Notes = util.NormalizeUnicode(row[ledger.NotesIndex])
	}
	if len(row) > ledger.IDIndex {
		entry.ID = util.NormalizeUnicode(row[ledger.IDIndex])
		if entry.ID == "" {
			entry.ID = source.GenerateSourceID(*entry)
		}
	}
	// entry.ID = source.GenerateSourceID(*entry)
	return entry, nil
}
