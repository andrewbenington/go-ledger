package app

import (
	"github.com/andrewbenington/go-ledger/cmd/command"
	"github.com/rivo/tview"
)

func DisplayText(text string) *tview.Modal {
	Log("DisplayText %s", text)
	modal := tview.NewModal()
	modal.AddButtons([]string{"Done"})
	modal.SetText(text)
	currentView := view.GetPrimitive()
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		view.SetPrimitive(currentView)
	})
	view.SetPrimitive(modal)
	return modal
}

func DisplayOutput(outputs []command.Output) {
	Log("DisplayOutput")
	LogInterface(outputs)
	if len(outputs) == 0 {
		return
	}
	if outputs[0].IsMessage {
		DisplayText(outputs[0].String).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) { PopStack() })
		return
	}
	list := tview.NewList()
	for _, output := range outputs {
		options := output.Options
		list.AddItem(output.String, "", ' ', func() {
			if len(output.Options) > 0 {
				opt := options[0]
				Log("Running %+v with %+v", opt.Name, opt.Args)
				view.SetBorder(true).SetTitle(opt.Name)
				stack = append(stack, opt.Select)
				RunCommandWithOutput(opt.Select, opt.Args)
			}
		})
	}
	list.AddItem("Back", "", 'b', PopStack).SetDoneFunc(PopStack)
	view.SetPrimitive(list)
}
