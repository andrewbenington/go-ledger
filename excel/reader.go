package excel

import (
	"fmt"
	"os"
	"time"

	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/andrewbenington/go-ledger/util"
	"github.com/xuri/excelize/v2"
)

func ReadLedger(years []int) (*ledger.Ledger, error) {
	l := &ledger.Ledger{}
	for _, year := range years {
		filename := fmt.Sprintf("%d.xlsx", year)
		if _, err := os.Stat(filename); err != nil {
			// file doesn't exist yet
			continue
		}
		entries, err := ledgerEntriesFromFile(filename)
		if err != nil {
			return nil, err
		}
		l.InsertEntries(entries)
	}
	return l, nil
}

func ledgerEntriesFromFile(filename string) ([]ledger.Entry, error) {
	file, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %w", filename, err)
	}
	entries := []ledger.Entry{}
	for _, sheet := range file.GetSheetList() {
		newEntries, err := ledgerEntriesFromSheet(file, sheet)
		if err != nil {
			return nil, fmt.Errorf("error reading %s: %w", filename, err)
		}
		entries = append(entries, newEntries...)
	}
	return entries, nil
}

func ledgerEntriesFromSheet(file *excelize.File, sheet string) ([]ledger.Entry, error) {
	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("error getting rows from sheet %s: %w", sheet, err)
	}
	entries := []ledger.Entry{}
	for _, row := range rows[1:] {
		entry, err := ledgerEntryFromRow(row)
		if err != nil {
			fmt.Printf("error getting entry from from row %v: %s\n", row, err)
			continue
		}
		entries = append(entries, *entry)
	}
	return entries, nil
}

func ledgerEntryFromRow(row []string) (entry *ledger.Entry, err error) {
	entry = &ledger.Entry{}
	if len(row) > ledger.ID_COLUMN {
		entry.ID = util.NormalizeUnicode(row[ledger.ID_COLUMN-1])
	}
	if len(row) > ledger.DATE_COLUMN {
		entry.Date, err = time.Parse("1/2/06 03:04", row[ledger.DATE_COLUMN-1])
		if err != nil {
			return nil, fmt.Errorf("Error parsing date: %s", err)
		}
	}
	if len(row) > ledger.MEMO_COLUMN {
		entry.Memo = util.NormalizeUnicode(row[ledger.MEMO_COLUMN-1])
	}
	if len(row) > ledger.VALUE_COLUMN {
		entry.Value, err = util.ParseMoneyAmount(row[ledger.VALUE_COLUMN-1])
		if err != nil {
			return nil, fmt.Errorf("Error parsing value amount: %s", err)
		}
	}
	if len(row) > ledger.BALANCE_COLUMN {
		entry.Balance, err = util.ParseMoneyAmount(row[ledger.BALANCE_COLUMN-1])
		if err != nil {
			return nil, fmt.Errorf("Error parsing balance amount: %s", err)
		}
	}
	if len(row) > ledger.TYPE_COLUMN {
		entry.Type = util.NormalizeUnicode(row[ledger.TYPE_COLUMN-1])
	}
	if len(row) > ledger.SOURCE_COLUMN {
		entry.Source = util.NormalizeUnicode(row[ledger.SOURCE_COLUMN-1])
	}
	if len(row) > ledger.PERSON_COLUMN {
		entry.Person = util.NormalizeUnicode(row[ledger.PERSON_COLUMN-1])
	}
	if len(row) > ledger.LABEL_COLUMN {
		entry.Label = util.NormalizeUnicode(row[ledger.LABEL_COLUMN-1])
	}
	return entry, nil
}
