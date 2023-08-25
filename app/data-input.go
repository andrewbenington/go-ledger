package app

import (
	"github.com/andrewbenington/go-ledger/command"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DataInput struct {
	values []string
	form   *tview.Form
}

type FormInput interface {
	SetInputCapture(key *tcell.EventKey) *tcell.EventKey
}

func dataInputFromCommand(c *command.Command) *DataInput {
	dataInput := &DataInput{
		form:   tview.NewForm(),
		values: make([]string, len(c.ExpectedArgs)),
	}

	for i := range c.ExpectedArgs {
		index := i
		arg := c.ExpectedArgs[index]
		switch arg.Type {
		case command.StringArg:
			dataInput.addStringArgField(arg, i)
		case command.BoolArg:
			dataInput.addBoolArgField(arg, i)
		case command.SelectArg:
			dataInput.addSelectArgField(arg, i)
		}
	}
	dataInput.form.AddButton("Done", func() {
		if c.ShowLogs {
			runCommandWithLogs(c, dataInput.values)
			return
		}
		runCommand(c, dataInput.values)
	})
	dataInput.form.AddButton("Cancel", popStack)
	return dataInput
}

func (d *DataInput) addStringArgField(arg command.Argument, index int) {
	d.form.AddInputField(arg.Name, "", 0, nil, func(text string) {
		d.values[index] = text
	})
	formItem := d.form.GetFormItem(index)
	field, ok := formItem.(*tview.InputField)
	if ok {
		if arg.AutoComplete != nil {
			field.SetAutocompleteFunc(arg.AutoCompleteWithArgs(&d.values))
		}
		if arg.OnAutoCompleted != nil {
			field.SetAutocompletedFunc(arg.OnAutoCompletedWithField(field))
		}
	}
}

func (d *DataInput) addBoolArgField(arg command.Argument, index int) {
	d.form.AddCheckbox(arg.Name, false, func(checked bool) {
		if checked {
			d.values[index] = "true"
		} else {
			d.values[index] = "false"
		}
	})
}

func (d *DataInput) addSelectArgField(arg command.Argument, index int) {
	optionLabels := []string{}
	for _, option := range arg.Options {
		optionLabels = append(optionLabels, option.Label)
	}
	d.form.AddDropDown(arg.Name, optionLabels, 0, func(_ string, i int) {
		d.values[index] = arg.Options[i].Value
	})
}
