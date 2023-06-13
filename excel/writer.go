package excel

import (
	"fmt"

	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/xuri/excelize/v2"
)

func WriteLedger(l ledger.Ledger) error {
	var file *excelize.File
	year := -1
	month := -1
	row := 0
	for _, entry := range l.Entries {
		if entry.Date.Year() > year {
			fmt.Printf("setting year to %d\n", entry.Date.Year())
			err := saveAndCloseFile(file, fmt.Sprintf("%d", year))
			if err != nil {
				return fmt.Errorf("error saving %d ledger: %w", year, err)
			}
			year = entry.Date.Year()
			file = excelize.NewFile()
			month = -1
		}
		if int(entry.Date.Month()) > month {
			month = int(entry.Date.Month())
			_, err := file.NewSheet(entry.Date.Month().String())
			if err != nil {
				return fmt.Errorf("error adding sheet for %s %d: %w", entry.Date.Month().String(), year, err)
			}
			err = writeHeaderRow(file, entry.Date.Month().String())
			if err != nil {
				return fmt.Errorf("error writing headers for %s %d: %w", entry.Date.Month().String(), year, err)
			}
			dateColumn, err := excelize.ColumnNumberToName(ledger.DATE_COLUMN + 1)
			if err == nil {
				_ = file.SetColWidth(entry.Date.Month().String(), dateColumn, dateColumn, 15)
			}
			memoColumn, err := excelize.ColumnNumberToName(ledger.MEMO_COLUMN + 1)
			if err == nil {
				_ = file.SetColWidth(entry.Date.Month().String(), memoColumn, memoColumn, 70)
			}
			err = file.DeleteSheet("Sheet1")
			if err != nil {
				return fmt.Errorf("error deleting placeholder Sheet1: %w", err)
			}
			row = 0
		}
		// fmt.Printf("setting sheet %s row at %s for %s\n", entry.Date.Month().String(), cell, entry.Memo)
		// err = file.SetSheetRow(entry.Date.Month().String(), cell, &[]interface{}{entry.ID, entry.Date, entry.Source, entry.Person, entry.Memo, entry.Value, entry.Type, entry.Balance, entry.Label, entry.Notes})
		err := writeEntryRow(file, entry.Date.Month().String(), row, entry)
		if err != nil {
			_ = saveAndCloseFile(file, fmt.Sprintf("%d", year))
			return fmt.Errorf("error setting sheet row for %s %d: %w", entry.Date.Month().String(), year, err)
		}
		row++
	}
	err := saveAndCloseFile(file, fmt.Sprintf("%d", year))
	if err != nil {
		return fmt.Errorf("error saving %d ledger: %w", year, err)
	}
	return nil
}

func writeHeaderRow(file *excelize.File, sheet string) error {
	cell, err := excelize.CoordinatesToCellName(ledger.ID_COLUMN+1, 1)
	if err == nil {
		err = file.SetCellValue(sheet, cell, "ID")
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.DATE_COLUMN+1, 1)
	if err == nil {
		err = file.SetCellValue(sheet, cell, "Date")
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.SOURCE_COLUMN+1, 1)
	if err == nil {
		err = file.SetCellValue(sheet, cell, "Source")
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.PERSON_COLUMN+1, 1)
	if err == nil {
		err = file.SetCellValue(sheet, cell, "Person")
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.MEMO_COLUMN+1, 1)
	if err == nil {
		err = file.SetCellValue(sheet, cell, "Memo")
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.VALUE_COLUMN+1, 1)
	if err == nil {
		err = file.SetCellValue(sheet, cell, "Value")
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.TYPE_COLUMN+1, 1)
	if err == nil {
		err = file.SetCellValue(sheet, cell, "Type")
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.BALANCE_COLUMN+1, 1)
	if err == nil {
		err = file.SetCellValue(sheet, cell, "Balance")
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.LABEL_COLUMN+1, 1)
	if err == nil {
		err = file.SetCellValue(sheet, cell, "Label")
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.NOTES_COLUMN+1, 1)
	if err == nil {
		err = file.SetCellValue(sheet, cell, "Notes")
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	return nil
}

func writeEntryRow(file *excelize.File, sheet string, row int, entry ledger.Entry) error {
	cell, err := excelize.CoordinatesToCellName(ledger.ID_COLUMN+1, row+2)
	if err == nil {
		err = file.SetCellValue(sheet, cell, entry.ID)
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.DATE_COLUMN+1, row+2)
	if err == nil {
		err = file.SetCellValue(sheet, cell, entry.Date)
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.SOURCE_COLUMN+1, row+2)
	if err == nil {
		err = file.SetCellValue(sheet, cell, entry.Source)
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.PERSON_COLUMN+1, row+2)
	if err == nil {
		err = file.SetCellValue(sheet, cell, entry.Person)
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.MEMO_COLUMN+1, row+2)
	if err == nil {
		err = file.SetCellValue(sheet, cell, entry.Memo)
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.VALUE_COLUMN+1, row+2)
	if err == nil {
		err = file.SetCellValue(sheet, cell, entry.Value)
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.TYPE_COLUMN+1, row+2)
	if err == nil {
		err = file.SetCellValue(sheet, cell, entry.Type)
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.BALANCE_COLUMN+1, row+2)
	if err == nil && entry.Balance != -1 {
		err = file.SetCellValue(sheet, cell, entry.Balance)
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.LABEL_COLUMN+1, row+2)
	if err == nil {
		err = file.SetCellValue(sheet, cell, entry.Label)
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	cell, err = excelize.CoordinatesToCellName(ledger.NOTES_COLUMN+1, row+2)
	if err == nil {
		err = file.SetCellValue(sheet, cell, entry.Notes)
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
	}
	return nil
}

func saveAndCloseFile(f *excelize.File, filename string) error {
	if f == nil {
		return nil
	}
	if err := f.SaveAs(fmt.Sprintf("%s.xlsx", filename)); err != nil {
		return err
	}
	return nil
}
