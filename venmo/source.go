package venmo

import (
	"regexp"

	"github.com/andrewbenington/go-ledger/csv"
	"github.com/andrewbenington/go-ledger/file"
	"github.com/andrewbenington/go-ledger/ledger"
)

const (
	DATE_FORMAT = "2006-01-02T15:04:05"
	SourceType  = "VENMO"
)

var (
	FilePattern = regexp.MustCompile("transaction_history.*")
	CSVSource   = csv.Source{
		SourceName: "Venmo",
		Columns: csv.ColumnFormat{
			ID:     2,
			Date:   3,
			Type:   4,
			Memo:   6,
			Person: 7,
			Value:  9,
		},
		DateFormat: DATE_FORMAT,
		HeaderRows: 4,
		FooterRows: 1,
	}
)

type Source struct {
	SourceName        string   `yaml:"name"`
	AccountHolderName string   `yaml:"account_holder_name"`
	Directories       []string `yaml:"directories"`
	HideTransfers     bool     `yaml:"hide_transfers"`
	csvSource         csv.Source
}

func (s *Source) Name() string {
	return s.SourceName
}

func (s *Source) Validate() error {
	return nil
}

func (s *Source) GetLedgerEntries(year int) ([]ledger.Entry, error) {
	return s.csvSource.GetLedgerEntries(year)
}

func (s *Source) UnmarshalYAML(unmarshal func(interface{}) error) error {
	fields := struct {
		SourceName        string   `yaml:"name"`
		AccountHolderName string   `yaml:"account_holder_name"`
		Directories       []string `yaml:"directories"`
	}{}
	err := unmarshal(&fields)
	if err != nil {
		return err
	}
	s.SourceName = fields.SourceName
	s.AccountHolderName = fields.AccountHolderName
	s.Directories = fields.Directories
	s.csvSource = CSVSource
	s.csvSource.PostProcessEntry = func(entry *ledger.Entry, row []string) error {
		return PostProcessEntry(entry, row, s.AccountHolderName)
	}
	s.csvSource.FileSearchPattern = file.SearchPattern{
		Directories:      s.Directories,
		FileNamePatterns: []regexp.Regexp{*FilePattern},
	}
	return nil
}

func PostProcessEntry(entry *ledger.Entry, row []string, accountHolder string) error {
	entry.SourceType = SourceType
	// if payment, use "To" column
	if entry.Person == accountHolder {
		entry.Person = row[7]
	}
	return nil
}
