package chase

import (
	"fmt"
	"regexp"
	"time"

	"github.com/andrewbenington/go-ledger/cmd/source/csv"
	"github.com/andrewbenington/go-ledger/cmd/source/file"
	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/andrewbenington/go-ledger/util"
)

const (
	DATE_FORMAT = "01/02/2006"
)

var (
	CreditColumns = csv.ColumnFormat{
		Date:  1,
		Memo:  3,
		Type:  5,
		Value: 6,
	}
	NonCreditColumns = csv.ColumnFormat{
		Date:    2,
		Memo:    3,
		Value:   4,
		Type:    5,
		Balance: 6,
	}
	FourDigitPattern = regexp.MustCompile(`[0-9]{4}`)
)

type Source struct {
	SourceName  string
	Directories []string
	LastDigits  string
	AccountType string
	csvSource   *csv.Source
}

func (s *Source) Name() string {
	return s.SourceName
}

func (s *Source) Validate() error {
	if len(s.LastDigits) != 4 || !FourDigitPattern.MatchString(s.LastDigits) {
		return fmt.Errorf("Chase CSV source validation: last_four_digits must be exactly four digits")
	}
	fmt.Println(s)
	if s.AccountType != "credit" && s.AccountType != "non-credit" {
		return fmt.Errorf("Chase CSV source validation: account_type must be 'credit' or 'non-credit'")
	}
	return nil
}

func (s *Source) GetLedgerEntries() ([]ledger.LedgerEntry, error) {
	return s.csvSource.GetLedgerEntries()
}

func (s *Source) UnmarshalYAML(unmarshal func(interface{}) error) error {
	fields := struct {
		SourceName  string   `yaml:"name"`
		Directories []string `yaml:"directories"`
		LastDigits  string   `yaml:"last_four_digits"`
		AccountType string   `yaml:"account_type"`
	}{}
	err := unmarshal(&fields)
	if err != nil {
		return err
	}
	s.SourceName = fields.SourceName
	s.Directories = fields.Directories
	s.LastDigits = fields.LastDigits
	s.AccountType = fields.AccountType
	s.Directories = fields.Directories
	s.csvSource = &csv.Source{
		SourceName:       s.SourceName,
		DateFormat:       DATE_FORMAT,
		HeaderRows:       1,
		OrderDescending:  true,
		PostProcessEntry: PostProcessEntry,
	}
	fileNamePattern, err := regexp.Compile(fmt.Sprintf("Chase%s_Activity.*", s.LastDigits))
	if err != nil {
		return fmt.Errorf("Error building pattern for source %s: %w", s.SourceName, err)
	}
	s.csvSource.FileSearchPattern = file.SearchPattern{
		Directories:      s.Directories,
		FileNamePatterns: []regexp.Regexp{*fileNamePattern},
	}
	if s.AccountType == "credit" {
		s.csvSource.Columns = CreditColumns
	} else if s.AccountType == "non-credit" {
		s.csvSource.Columns = NonCreditColumns
	}
	return nil
}

func PostProcessEntry(entry *ledger.LedgerEntry, row []string) error {
	month, day := util.ExtractDateFromTitle(entry.Memo)
	if month > 0 && day > 0 {
		entry.Date = time.Date(entry.Date.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	return nil
}
