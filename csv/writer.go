package csv

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/andrewbenington/go-ledger/ledger"
)

func WriteLedger(l ledger.Ledger, filename string) error {
	fileOut, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	writer := csv.NewWriter(fileOut)
	err = writer.Write([]string{"Date", "Source", "Person", "Memo", "Value", "Type", "Balance", "Label", "Notes"})
	if err != nil {
		return fmt.Errorf("error writing ledger header: %w", err)
	}
	for _, entry := range l.Entries() {
		err = writer.Write([]string{entry.Date.String(), entry.Source, entry.Person, entry.Memo, fmt.Sprintf("%f", entry.Value), entry.Type, fmt.Sprintf("%f", entry.Balance), entry.Label, entry.Notes})
		if err != nil {
			return fmt.Errorf("error writing ledger entry row: %w", err)
		}
	}
	writer.Flush()
	return nil
}
