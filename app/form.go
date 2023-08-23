package app

import (
	"github.com/andrewbenington/go-ledger/command"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Form struct {
	fields       []tview.Primitive
	values       []string
	view         *tview.Flex
	doneBtn      *tview.Button
	cancelBtn    *tview.Button
	currentIndex int
}

type FormInput interface {
	SetInputCapture(key *tcell.EventKey) *tcell.EventKey
}

func formViewFromCommand(c *command.Command) *Form {
	form := &Form{
		view:   tview.NewFlex().SetDirection(tview.FlexRow),
		values: make([]string, len(c.ExpectedArgs)),
		fields: make([]tview.Primitive, len(c.ExpectedArgs)+2),
	}
	form.doneBtn = tview.NewButton("Done").SetSelectedFunc(func() {
		if c.ShowLogs {
			runCommandWithLogs(c, form.values)
			return
		}
		runCommand(c, form.values)
	})
	form.doneBtn.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyTAB, tcell.KeyDown, tcell.KeyESC:
			form.NavigateForwards()
		case tcell.KeyBacktab, tcell.KeyUp:
			form.NavigateBack()
		}
		return key
	})

	form.cancelBtn = tview.NewButton("Cancel").SetSelectedFunc(popStack)
	form.cancelBtn.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyBacktab, tcell.KeyUp:
			form.NavigateBack()
		}
		return key
	})

	for i := range c.ExpectedArgs {
		index := i
		arg := c.ExpectedArgs[index]
		// add field label
		form.view.AddItem(tview.NewTextView().SetText(arg.Name), 1, 0, false)
		var field tview.Primitive
		switch arg.Type {
		case command.StringArg:
			field = form.buildStringFieldAtIndex(arg, index)
		case command.BoolArg:
			field = form.buildCheckboxFieldAtIndex(arg, index)
		}
		form.view.AddItem(field, 1, 0, index == 0)
		form.view.AddItem(nil, 1, 0, false)
	}

	// create buttons and add to Form
	form.view.AddItem(form.doneBtn, 1, 0, false)
	form.fields[len(form.fields)-2] = form.doneBtn
	form.view.AddItem(form.cancelBtn, 1, 0, false)
	form.fields[len(form.fields)-1] = form.cancelBtn
	LogInterface(form)
	return form
}

func (f *Form) buildStringFieldAtIndex(arg command.ArgOptions, index int) *tview.InputField {
	// handle field completion
	field := tview.NewInputField()

	field.SetChangedFunc(func(text string) {
		f.values[index] = text
	})
	field.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTAB, tcell.KeyEnter:
			f.NavigateForwards()
		case tcell.KeyBacktab:
			f.NavigateBack()
		case tcell.KeyESC:
			app.SetFocus(f.cancelBtn)
		}
	})
	field.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyUp:
			f.NavigateBack()
		case tcell.KeyDown:
			f.NavigateForwards()
		}
		return key
	})

	if arg.AutoComplete != nil {
		field.SetAutocompleteFunc(arg.AutoCompleteWithArgs(&f.values))
	}
	if arg.OnAutoCompleted != nil {
		field.SetAutocompletedFunc(arg.OnAutoCompletedWithField(field))
	}
	f.fields[index] = field
	return field
}

func (f *Form) buildCheckboxFieldAtIndex(arg command.ArgOptions, index int) *tview.Checkbox {
	f.values[index] = "false"
	// handle field completion
	field := tview.NewCheckbox()
	field.SetChangedFunc(func(checked bool) {
		if checked {
			f.values[index] = "true"
		} else {
			f.values[index] = "false"
		}
	})
	field.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTAB, tcell.KeyEnter:
			f.NavigateForwards()
		case tcell.KeyBacktab:
			f.NavigateBack()
		case tcell.KeyESC:
			app.SetFocus(f.cancelBtn)
		}
	})
	field.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyUp:
			f.NavigateBack()
		case tcell.KeyDown:
			f.NavigateForwards()
		}
		return key
	})
	f.fields[index] = field

	return field
}

func (f *Form) NavigateBack() {
	if f.currentIndex > 0 {
		f.currentIndex -= 1
	}
	Log("currentIndex: %d", f.currentIndex)
	Log("setting focus to %+v (%p)", f.fields[f.currentIndex], f.fields[f.currentIndex])
	app.SetFocus(f.fields[f.currentIndex])
}

func (f *Form) NavigateForwards() {
	if f.currentIndex < len(f.fields)-1 {
		f.currentIndex += 1
	}
	Log("currentIndex: %d", f.currentIndex)
	Log("setting focus to %+v (%p)", f.fields[f.currentIndex], f.fields[f.currentIndex])
	app.SetFocus(f.fields[f.currentIndex])
}
