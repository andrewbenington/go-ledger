package excel

import (
	"fmt"
	"reflect"

	"github.com/andrewbenington/go-ledger/cmd/label"
	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/andrewbenington/go-ledger/util"
	"github.com/xuri/excelize/v2"
)

func WriteLedger(l ledger.Ledger) error {
	if len(l.Entries()) == 0 {
		return fmt.Errorf("no entries")
	}
	file := excelize.NewFile()
	err := addOverviewSheet(file)
	if err != nil {
		fmt.Println(err)
	}
	year := l.Entries()[0].Date.Year()
	month := l.Entries()[0].Date.Month()
	err = initMonth(file, month.String(), &l)
	if err != nil {
		return err
	}
	row := 0
	for _, entry := range l.Entries() {
		if entry.Date.Year() > year {
			fmt.Printf("setting year to %d\n", entry.Date.Year())
			err = saveAndCloseFile(file, fmt.Sprintf("%d", year))
			if err != nil {
				return fmt.Errorf("error saving %d ledger: %w", year, err)
			}
			year = entry.Date.Year()
			month = entry.Date.Month()
			file = excelize.NewFile()
			err = initMonth(file, month.String(), &l)
			if err != nil {
				return err
			}
		}
		if entry.Date.Month() > month {
			if row > 0 {
				err := addPivotTable(file, month.String(), row)
				if err != nil {
					return fmt.Errorf("error adding %s pivot table: %w", month.String(), err)
				}
			}
			month = entry.Date.Month()
			err = initMonth(file, month.String(), &l)
			if err != nil {
				return err
			}
			row = 0
		}
		// fmt.Printf("setting sheet %s row at %s for %s\n", entry.Date.Month().String(), cell, entry.Memo)
		// err = file.SetSheetRow(entry.Date.Month().String(), cell, &[]interface{}{entry.ID, entry.Date, entry.Source, entry.Person, entry.Memo, entry.Value, entry.Type, entry.Balance, entry.Label, entry.Notes})
		err = writeEntryRow(file, row, entry, &l)
		if err != nil {
			_ = saveAndCloseFile(file, fmt.Sprintf("%d", year))
			return fmt.Errorf("error writing entry row in %s %d: %w", entry.Date.Month().String(), year, err)
		}
		row++
	}
	err = addPivotTable(file, month.String(), row)
	if err != nil {
		return fmt.Errorf("error adding %s pivot table: %w", month.String(), err)
	}
	err = addOverviewTotal(file)
	if err != nil {
		fmt.Println(err)
	}
	err = saveAndCloseFile(file, fmt.Sprintf("%d", year))
	if err != nil {
		return fmt.Errorf("error saving %d ledger: %w", year, err)
	}
	return nil
}

func initMonth(file *excelize.File, month string, l *ledger.Ledger) error {
	err := addMonthSheet(file, month, l)
	if err != nil {
		return err
	}
	err = addOverviewRow(file, month)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func addMonthSheet(file *excelize.File, month string, l *ledger.Ledger) error {
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
	err = writeHeaderRow(file, month, l)
	if err != nil {
		return fmt.Errorf("error writing headers for sheet %s: %w", month, err)
	}
	err = setColumnWidths(file, month, l)
	if err != nil {
		return fmt.Errorf("error setting column widths for sheet %s: %w", month, err)
	}
	return nil
}

func writeHeaderRow(file *excelize.File, sheet string, l *ledger.Ledger) error {
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

func writeEntryRow(file *excelize.File, row int, entry ledger.Entry, l *ledger.Ledger) error {
	for columnIndex, fieldName := range ledger.Columns {
		fname := fieldName
		cell, err := excelize.CoordinatesToCellName(columnIndex+1, row+2)
		if err != nil {
			return fmt.Errorf("Error getting cell name for %d,%d: %w", columnIndex+1, 1, err)
		}
		if fname == "Source Name" {
			err = file.SetCellValue(entry.Date.Month().String(), cell, entry.SourceName)
			if err != nil {
				return fmt.Errorf("Error setting cell %s: %w", cell, err)
			}
			continue
		}
		if fname == "Source Type" {
			err = file.SetCellValue(entry.Date.Month().String(), cell, entry.SourceType)
			if err != nil {
				return fmt.Errorf("Error setting cell %s: %w", cell, err)
			}
			continue
		}
		refVal := reflect.ValueOf(entry)
		field := refVal.FieldByName(fname)
		if !field.IsValid() {
			continue
		}
		// fmt.Printf("%+v, %+v, %s, %s, %+v\n", file, entry, entry.Date.Month().String(), cell, field)
		err = file.SetCellValue(entry.Date.Month().String(), cell, field.Interface())
		if err != nil {
			return fmt.Errorf("Error setting cell %s: %w", cell, err)
		}
		if fieldName == "Label" && entry.Label != "" {
			err = applyLabelTextColor(file, entry.Label, entry.Date.Month().String(), cell)
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

func setColumnWidths(file *excelize.File, sheetName string, l *ledger.Ledger) error {
	dateColumn, err := excelize.ColumnNumberToName(ledger.DateIndex + 1)
	if err == nil {
		err = file.SetColWidth(sheetName, dateColumn, dateColumn, 15)
		if err != nil {
			return err
		}
	}
	memoColumn, err := excelize.ColumnNumberToName(ledger.MemoIndex + 1)
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
	lastDataCell, err := excelize.CoordinatesToCellName(ledger.FieldCount-1, rowCount+1)
	if err != nil {
		return fmt.Errorf("error getting last column: %w", err)
	}
	firstPivotCell, err := excelize.CoordinatesToCellName(ledger.SwapTableStart, 1)
	if err != nil {
		return fmt.Errorf("error getting first pivot column: %w", err)
	}
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

func addOverviewSheet(file *excelize.File) error {
	_, err := file.NewSheet("Overview")
	if err != nil {
		return fmt.Errorf("add overview sheet: %w", err)
	}
	for i, l := range label.All() {
		cell, err := excelize.CoordinatesToCellName(i+2, 1)
		if err != nil {
			return fmt.Errorf("add overview sheet: %w", err)
		}
		file.SetCellValue("Overview", cell, l.Name)
		err = applyLabelTextColor(file, l.Name, "Overview", cell)
	}
	return nil
}

func addOverviewRow(file *excelize.File, month string) error {
	rows, err := file.GetRows("Overview")
	if err != nil {
		return fmt.Errorf("add %s overview row: %w", month, err)
	}
	maxRow := len(rows)
	if err != nil {
		return fmt.Errorf("add %s overview row: %w", month, err)
	}
	cell, err := excelize.CoordinatesToCellName(1, maxRow+1)
	if err != nil {
		return fmt.Errorf("add %s overview row: %w", month, err)
	}
	file.SetCellValue("Overview", cell, month)
	for i, l := range label.All() {
		cell, err := excelize.CoordinatesToCellName(i+2, maxRow+1)
		if err != nil {
			return fmt.Errorf("add %s overview row: %w", month, err)
		}
		labelCol, err := excelize.ColumnNumberToName(ledger.SwapTableStart)
		if err != nil {
			return fmt.Errorf("add %s overview row: %w", month, err)
		}
		sumCol, err := excelize.ColumnNumberToName(ledger.SwapTableStart + 1)
		if err != nil {
			return fmt.Errorf("add %s overview row: %w", month, err)
		}
		formula := fmt.Sprintf(
			"=IFERROR(VLOOKUP(\"%s\", %s!%s:%s, 2, FALSE), 0)",
			l.Name,
			month,
			labelCol,
			sumCol)
		file.SetCellFormula("Overview", cell, formula)
	}
	return nil
}

func addOverviewTotal(file *excelize.File) error {
	rows, err := file.GetRows("Overview")
	if err != nil {
		return fmt.Errorf("add overview total row: %w", err)
	}
	maxRow := len(rows)
	if err != nil {
		return fmt.Errorf("add overview total row: %w", err)
	}
	cell, err := excelize.CoordinatesToCellName(1, maxRow+1)
	if err != nil {
		return fmt.Errorf("add overview total row: %w", err)
	}
	file.SetCellValue("Overview", cell, "Total")
	for i := range label.All() {
		cell, err := excelize.CoordinatesToCellName(i+2, maxRow+1)
		if err != nil {
			return fmt.Errorf("add overview total row: %w", err)
		}
		startCell, err := excelize.CoordinatesToCellName(i+2, 1)
		if err != nil {
			return fmt.Errorf("add overview total row: %w", err)
		}
		endCell, err := excelize.CoordinatesToCellName(i+2, maxRow)
		if err != nil {
			return fmt.Errorf("add overview total row: %w", err)
		}
		formula := fmt.Sprintf("=SUM(%s:%s)", startCell, endCell)
		file.SetCellFormula("Overview", cell, formula)
	}
	return nil
}

func applyLabelTextColor(file *excelize.File, label string, sheet string, cell string) error {
	style := &excelize.Style{
		Font: &excelize.Font{
			Color: util.StringToRGB(label),
			Bold:  true,
		},
	}
	styleIndex, err := file.NewStyle(style)
	if err != nil {
		return err
	}
	return file.SetCellStyle(sheet, cell, cell, styleIndex)
}
