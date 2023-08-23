package command

import "github.com/rivo/tview"

type ArgType int

const (
	StringArg ArgType = iota
	BoolArg
)

type ArgOptions struct {
	Name            string
	AutoComplete    func(currentText string, currentArgs *[]string) []string
	OnAutoCompleted func(text string, index int, field *tview.InputField) bool
	Type            ArgType
}

func (a *ArgOptions) AutoCompleteWithArgs(currentArgs *[]string) func(currentText string) []string {
	return func(currentText string) []string {
		return a.AutoComplete(currentText, currentArgs)
	}
}

func (a *ArgOptions) OnAutoCompletedWithField(field *tview.InputField) func(text string, index int, source int) bool {
	return func(text string, index int, source int) bool {
		if source != tview.AutocompletedNavigate && source != tview.AutocompletedTab {
			a.OnAutoCompleted(text, index, field)
			return true
		}
		return false
	}
}
