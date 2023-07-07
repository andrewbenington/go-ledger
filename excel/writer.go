package excel

import (
	"fmt"
	"reflect"

	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/xuri/excelize/v2"
)

func WriteLedger(l ledger.Ledger) error {
	file := excelize.NewFile()
	year := l.Entries()[0].Date.Year()
	month := l.Entries()[0].Date.Month()
	addMonthSheet(file, month.String())
	row := 0
	for _, entry := range l.Entries() {
		if entry.Date.Year() > year {
			fmt.Printf("setting year to %d\n", entry.Date.Year())
			err := saveAndCloseFile(file, fmt.Sprintf("%d", year))
			if err != nil {
				return fmt.Errorf("error saving %d ledger: %w", year, err)
			}
			year = entry.Date.Year()
			month = entry.Date.Month()
			file = excelize.NewFile()
			addMonthSheet(file, month.String())
		}
		if entry.Date.Month() > month {
			if row > 0 {
				err := addPivotTable(file, month.String(), row)
				if err != nil {
					return fmt.Errorf("error adding %s pivot table: %w", month.String(), err)
				}
			}
			month = entry.Date.Month()
			addMonthSheet(file, month.String())
			row = 0
		}
		// fmt.Printf("setting sheet %s row at %s for %s\n", entry.Date.Month().String(), cell, entry.Memo)
		// err = file.SetSheetRow(entry.Date.Month().String(), cell, &[]interface{}{entry.ID, entry.Date, entry.Source, entry.Person, entry.Memo, entry.Value, entry.Type, entry.Balance, entry.Label, entry.Notes})
		err := writeEntryRow(file, row, entry)
		if err != nil {
			_ = saveAndCloseFile(file, fmt.Sprintf("%d", year))
			return fmt.Errorf("error writing entry row in %s %d: %w", entry.Date.Month().String(), year, err)
		}
		row++
	}
	err := addPivotTable(file, month.String(), row)
	if err != nil {
		return fmt.Errorf("error adding %s pivot table: %w", month.String(), err)
	}
	err = saveAndCloseFile(file, fmt.Sprintf("%d", year))
	if err != nil {
		return fmt.Errorf("error saving %d ledger: %w", year, err)
	}
	return nil
}

func addMonthSheet(file *excelize.File, month string) error {
	_, err := file.NewSheet(month)
	if err != nil {
		return fmt.Errorf("error adding sheet for sheet %s: %w", month, err)
	}
	if i, _ := file.GetSheetIndex("Sheet1"); i != -1 {
		// delete default sheet
		err = file.DeleteSheet("Sheet1")
		if err != nil {
			return fmt.Errorf("error deleting placeholder Sheet1: %w", err)
		}
	}
	err = writeHeaderRow(file, month)
	if err != nil {
		return fmt.Errorf("error writing headers for sheet %s: %w", month, err)
	}
	err = setColumnWidths(file, month)
	if err != nil {
		return fmt.Errorf("error setting column widths for sheet %s: %w", month, err)
	}
	return nil
}

func writeHeaderRow(file *excelize.File, sheet string) error {
	for columnIndex, fieldName := range ledger.Columns {
		cell, err := excelize.CoordinatesToCellName(columnIndex+1, 1)
		if err == nil {
			err = file.SetCellValue(sheet, cell, fieldName)
			if err != nil {
				return fmt.Errorf("Error setting cell %s: %w", cell, err)
			}
		}
	}
	return nil
}

func writeEntryRow(file *excelize.File, row int, entry ledger.Entry) error {
	for columnIndex, fieldName := range ledger.Columns {
		cell, err := excelize.CoordinatesToCellName(columnIndex+1, row+2)
		if err != nil {
			return fmt.Errorf("Error getting cell name for %d,%d: %w", columnIndex+1, 1, err)
		}
		refVal := reflect.ValueOf(entry)
		field := refVal.FieldByName(fieldName)
		if !field.IsValid() {
			continue
		}
		// fmt.Printf("%+v, %+v, %s, %s, %+v\n", file, entry, entry.Date.Month().String(), cell, field)
		err = file.SetCellValue(entry.Date.Month().String(), cell, field.Interface())
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
	return f.Close()
}

func setColumnWidths(file *excelize.File, sheetName string) error {
	dateColumn, err := excelize.ColumnNumberToName(ledger.DATE_COLUMN)
	if err == nil {
		err = file.SetColWidth(sheetName, dateColumn, dateColumn, 15)
		if err != nil {
			return err
		}
	}
	memoColumn, err := excelize.ColumnNumberToName(ledger.MEMO_COLUMN)
	if err == nil {
		err = file.SetColWidth(sheetName, memoColumn, memoColumn, 70)
		if err != nil {
			return err
		}
	}
	return nil
}

func addPivotTable(file *excelize.File, sheet string, rowCount int) error {
	firstDataCell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		return fmt.Errorf("error getting first column: %w", err)
	}
	lastDataColumn := ledger.LABEL_COLUMN
	if ledger.VALUE_COLUMN > lastDataColumn {
		lastDataColumn = ledger.VALUE_COLUMN
	}
	lastDataCell, err := excelize.CoordinatesToCellName(lastDataColumn, rowCount+1)
	if err != nil {
		return fmt.Errorf("error getting last column: %w", err)
	}
	firstPivotCell, err := excelize.CoordinatesToCellName(ledger.SWAP_TABLE_START_COLUMN, 1)
	if err != nil {
		return fmt.Errorf("error getting first pivot column: %w", err)
	}
	// lastPivotCell, err := excelize.CoordinatesToCellName(ledger.SWAP_TABLE_START_COLUMN+2, rowCount)
	// if err != nil {
	// 	return fmt.Errorf("error getting last pivot column: %w", err)
	// }
	// err = file.AddPivotTable(&excelize.PivotTableOptions{
	// 	DataRange:       fmt.Sprintf("%s!%s:%s", sheet, firstDataCell, lastDataCell),
	// 	PivotTableRange: fmt.Sprintf("%s!%s:%s", sheet, firstPivotCell, lastPivotCell),
	// 	Rows:            []excelize.PivotTableField{{Data: "Label"}, {Data: "Source"}},
	// 	Columns:         []excelize.PivotTableField{{Data: "Value"}},
	// 	RowGrandTotals:  true,
	// 	ColGrandTotals:  true,
	// 	ShowDrill:       true,
	// 	ShowRowHeaders:  true,
	// 	ShowColHeaders:  true,
	// 	ShowLastColumn:  true,
	// })
	// lastDataCell = "J34"
	options := &excelize.PivotTableOptions{
		DataRange:       fmt.Sprintf("%s!$%s:$%s", sheet, firstDataCell, lastDataCell),
		PivotTableRange: fmt.Sprintf("%s!$%s:$Q34", sheet, firstPivotCell),
		Rows: []excelize.PivotTableField{
			{Data: "Label", Name: "Assigned Label", DefaultSubtotal: true}},
		Data: []excelize.PivotTableField{
			{Data: "Value", Subtotal: "Sum"}},
		ShowDrill:      true,
		ShowRowHeaders: true,
		ShowColHeaders: true,
		ShowLastColumn: true,
	}
	err = file.AddPivotTable(options)
	return err
}
