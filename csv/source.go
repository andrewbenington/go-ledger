package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/andrewbenington/go-ledger/cmd/label"
	"github.com/andrewbenington/go-ledger/file"
	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/andrewbenington/go-ledger/util"
)

type Source struct {
	SourceName      string
	Columns         ColumnFormat
	DateFormat      string `yaml:"date_format"`
	HeaderRows      int    `yaml:"header_rows"`
	FooterRows      int    `yaml:"header_rows"`
	OrderDescending bool   `yaml:"order_descending"`
	// called after the entry has been populated with defined columns
	PostProcessEntry  func(*ledger.Entry, []string) error
	FileSearchPattern file.SearchPattern
}

// Columns are 1-indexed; 0 means column not present
type ColumnFormat struct {
	ID      int
	Date    int
	Person  int
	Memo    int
	Value   int
	Type    int
	Balance int
	Label   int
	Others  []int
}

func (s *Source) GetLedgerEntries() ([]ledger.Entry, error) {
	filePaths, err := s.FileSearchPattern.FindMatchingFiles()
	if err != nil {
		return nil, fmt.Errorf("error finding files for source %s: %w\n", s.SourceName, err)
	}
	ledgerEntries := []ledger.Entry{}
	for _, path := range filePaths {
		fmt.Printf("\tReading file %s...\n", path)
		fileEntries, err := s.LedgerEntriesFromFile(path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ledgerEntries = append(ledgerEntries, fileEntries...)
	}
	return ledgerEntries, nil
}

func (s *Source) LedgerEntriesFromFile(filename string) ([]ledger.Entry, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening CSV file: %s", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading CSV records: %s", err)
	}

	entries := []ledger.Entry{}
	for i, row := range records {
		if i >= len(records)-s.FooterRows {
			break
		}
		if i > s.HeaderRows {
			entry := &ledger.Entry{}
			entry.Source = s.SourceName
			err := s.fillDefinedColumns(entry, row)
			if err != nil {
				return nil, err
			}
			err = s.PostProcessEntry(entry, row)
			if err != nil {
				return nil, err
			}
			entry.Label = label.FindLabel(entry.Memo)
			if s.OrderDescending {
				entries = append([]ledger.Entry{*entry}, entries...)
			} else {
				entries = append(entries, *entry)
			}
		}
	}
	return entries, nil
}

func (s *Source) fillDefinedColumns(entry *ledger.Entry, row []string) (err error) {
	if s.Columns.ID > 0 {
		entry.ID = util.NormalizeUnicode(row[s.Columns.ID-1])
	}
	if s.Columns.Date > 0 {
		entry.Date, err = time.Parse(s.DateFormat, row[s.Columns.Date-1])
		if err != nil {
			return fmt.Errorf("Error parsing date: %s", err)
		}
	}
	if s.Columns.Memo > 0 {
		entry.Memo = util.NormalizeUnicode(row[s.Columns.Memo-1])
	}
	if s.Columns.Value > 0 {
		entry.Value, err = util.ParseMoneyAmount(row[s.Columns.Value-1])
		if err != nil {
			return fmt.Errorf("Error parsing value amount: %s", err)
		}
	}
	if s.Columns.Balance > 0 {
		entry.Balance, err = util.ParseMoneyAmount(row[s.Columns.Balance-1])
		if err != nil {
			return fmt.Errorf("Error parsing balance amount: %s", err)
		}
	}
	if s.Columns.Type > 0 {
		entry.Type = util.NormalizeUnicode(row[s.Columns.Type-1])
	}
	if s.Columns.Person > 0 {
		entry.Person = util.NormalizeUnicode(row[s.Columns.Person-1])
	}
	return nil
}
