package command

import "github.com/rivo/tview"

type ArgType int

const (
	StringArg ArgType = iota
	BoolArg
	SelectArg
)

type ArgOption struct {
	Value string
	Label string
}

type Argument struct {
	Name            string
	IsConstant      bool
	AutoComplete    func(currentText string, currentArgs *[]string) []string
	OnAutoCompleted func(text string, index int, field *tview.InputField) bool
	Type            ArgType
	Options         []ArgOption
}

func (a *Argument) AutoCompleteWithArgs(currentArgs *[]string) func(currentText string) []string {
	return func(currentText string) []string {
		return a.AutoComplete(currentText, currentArgs)
	}
}

func (a *Argument) OnAutoCompletedWithField(field *tview.InputField) func(text string, index int, source int) bool {
	return func(text string, index int, source int) bool {
		if source != tview.AutocompletedNavigate && source != tview.AutocompletedTab {
			a.OnAutoCompleted(text, index, field)
			return true
		}
		return false
	}
}
