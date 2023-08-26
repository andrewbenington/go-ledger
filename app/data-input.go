package app

import (
	"strings"

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

func dataInputFromCommand(c *command.Command, initialValues []string) *DataInput {
	dataInput := &DataInput{
		form:   tview.NewForm(),
		values: make([]string, len(c.ExpectedArgs)),
	}

	for i := range c.ExpectedArgs {
		index := i
		arg := c.ExpectedArgs[index]
		if len(initialValues) > index {
			dataInput.values[index] = initialValues[index]
			if arg.IsConstant {
				dataInput.form.AddTextView(arg.Name, initialValues[index], 0, 1, false, false)
				continue
			}
		}
		switch arg.Type {
		case command.StringArg:
			dataInput.addStringArgField(arg, index, initialValues)
		case command.BoolArg:
			dataInput.addBoolArgField(arg, index, initialValues)
		case command.SelectArg:
			dataInput.addSelectArgField(arg, index, initialValues)
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

func (d *DataInput) addStringArgField(arg command.Argument, index int, initial []string) {
	initialValue := ""
	if len(initial) > index {
		initialValue = initial[index]
	}

	d.form.AddInputField(arg.Name, initialValue, 0, nil, func(text string) {
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

func (d *DataInput) addBoolArgField(arg command.Argument, index int, initial []string) {
	initialValue := len(initial) > index && strings.EqualFold(initial[index], "true")
	d.form.AddCheckbox(arg.Name, initialValue, func(checked bool) {
		if checked {
			d.values[index] = "true"
		} else {
			d.values[index] = "false"
		}
	})
}

func (d *DataInput) addSelectArgField(arg command.Argument, index int, initial []string) {
	initialValue := 0
	optionLabels := []string{}
	for i, option := range arg.Options {
		optionLabels = append(optionLabels, option.Label)
		if len(initial) > i && strings.EqualFold(option.Value, initial[i]) {
			initialValue = i
		}
	}
	d.form.AddDropDown(arg.Name, optionLabels, initialValue, func(_ string, i int) {
		d.values[index] = arg.Options[i].Value
	})
}
