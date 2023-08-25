package source

import (
	"regexp"

	"github.com/andrewbenington/go-ledger/csv"
	"github.com/andrewbenington/go-ledger/file"
	"github.com/andrewbenington/go-ledger/ledger"
)

const (
	VenmoDateFormat = "2006-01-02T15:04:05"
	VenmoSourceType = "VENMO"
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
		DateFormat: VenmoDateFormat,
		HeaderRows: 4,
		FooterRows: 1,
	}
)

type VenmoSource struct {
	SourceName        string   `yaml:"name"`
	AccountHolderName string   `yaml:"account_holder_name"`
	Directories       []string `yaml:"directories"`
	HideTransfers     bool     `yaml:"hide_transfers"`
	csvSource         csv.Source
}

func (s *VenmoSource) Name() string {
	return s.SourceName
}

func (s *VenmoSource) Validate() error {
	return nil
}

func (s *VenmoSource) GetLedgerEntries(year int) ([]ledger.Entry, error) {
	return s.csvSource.GetLedgerEntries(year)
}

func (s *VenmoSource) UnmarshalYAML(unmarshal func(interface{}) error) error {
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
		return postProcessVenmo(entry, row, s.AccountHolderName)
	}
	s.csvSource.FileSearchPattern = file.SearchPattern{
		Directories:      s.Directories,
		FileNamePatterns: []regexp.Regexp{*FilePattern},
	}
	return nil
}

func postProcessVenmo(entry *ledger.Entry, row []string, accountHolder string) error {
	entry.SourceType = VenmoSourceType
	// if payment, use "To" column
	if entry.Person == accountHolder {
		entry.Person = row[7]
	}
	return nil
}

func AddVenmoSource(vs VenmoSource) error {
	s, err := Get()
	if err != nil {
		return err
	}
	err = vs.Validate()
	if err != nil {
		return err
	}
	s.Venmo = append(s.Venmo, vs)
	return saveSources(s)
}
