package chase

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"time"

	"github.com/andrewbenington/go-ledger/csv"
	"github.com/andrewbenington/go-ledger/file"
	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/andrewbenington/go-ledger/util"
)

const (
	DATE_FORMAT = "01/02/2006"
	SourceType  = "CHASE"
)

var (
	creditColumns = csv.ColumnFormat{
		Date:  1,
		Memo:  3,
		Type:  5,
		Value: 6,
	}
	nonCreditColumns = csv.ColumnFormat{
		Date:    2,
		Memo:    3,
		Value:   4,
		Type:    5,
		Balance: 6,
	}
	fourDigitPattern  = regexp.MustCompile(`[0-9]{4}`)
	maxDateAdjustment = 7 * 24 * time.Hour
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
	if len(s.LastDigits) != 4 || !fourDigitPattern.MatchString(s.LastDigits) {
		return fmt.Errorf("Chase CSV source validation: last_four_digits must be exactly four digits")
	}
	if s.AccountType != "credit" && s.AccountType != "non-credit" {
		return fmt.Errorf("Chase CSV source validation: account_type must be 'credit' or 'non-credit'")
	}
	return nil
}

func (s *Source) GetLedgerEntries(year int) ([]ledger.Entry, error) {
	return s.csvSource.GetLedgerEntries(year)
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
		s.csvSource.Columns = creditColumns
	} else if s.AccountType == "non-credit" {
		s.csvSource.Columns = nonCreditColumns
	}
	return nil
}

func PostProcessEntry(entry *ledger.Entry, row []string) error {
	entry.SourceType = SourceType
	month, day := util.ExtractDateFromTitle(entry.Memo)
	if month > 0 && day > 0 {
		newDate := time.Date(entry.Date.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
		if newDate.After(entry.Date) && newDate.Sub(entry.Date) < maxDateAdjustment {
			entry.Notes = fmt.Sprintf("original date: %s", entry.Date)
			entry.Date = newDate
		}
	}
	entry.ID = hashEntry(entry)
	return nil
}

func hashEntry(entry *ledger.Entry) string {
	str := fmt.Sprintf("%s_%f", entry.Date.Format(DATE_FORMAT), entry.Balance)
	algorithm := fnv.New32a()
	algorithm.Write([]byte(str))
	return fmt.Sprintf("%d", algorithm.Sum32())
}
