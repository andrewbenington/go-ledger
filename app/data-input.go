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

func formViewFromCommand(c *command.Command) *DataInput {
	dataInput := &DataInput{
		form:   tview.NewForm(),
		values: make([]string, len(c.ExpectedArgs)),
	}

	for i := range c.ExpectedArgs {
		index := i
		arg := c.ExpectedArgs[index]
		// add field label
		// form.view.AddItem(tview.NewTextView().SetText(arg.Name), 1, 0, false)
		// var field tview.Primitive
		switch arg.Type {
		case command.StringArg:
			dataInput.addStringArgField(arg, i)
		case command.BoolArg:
			dataInput.form.AddCheckbox(arg.Name, false, func(checked bool) {
				if checked {
					dataInput.values[index] = "true"
				} else {
					dataInput.values[index] = "false"
				}
			})
			// field = form.buildCheckboxFieldAtIndex(arg, index)
		}
		// form.view.AddItem(field, 1, 0, index == 0)
		// form.view.AddItem(nil, 1, 0, false)
	}
	dataInput.form.AddButton("Done", func() {
		if c.ShowLogs {
			runCommandWithLogs(c, dataInput.values)
			return
		}
		runCommand(c, dataInput.values)
	})
	dataInput.form.AddButton("Cancel", popStack)
	// create buttons and add to Form
	// form.view.AddItem(form.doneBtn, 1, 0, false)
	// form.fields[len(form.fields)-2] = form.doneBtn
	// form.view.AddItem(form.cancelBtn, 1, 0, false)
	// form.fields[len(form.fields)-1] = form.cancelBtn
	// LogInterface(form)
	return dataInput
}

func (d *DataInput) addStringArgField(arg command.ArgOptions, index int) {
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
