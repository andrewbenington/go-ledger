package app

import (
	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/util"
	"github.com/rivo/tview"
)

// displayOutput displays either a single modal or a list of items
// depending on if there are multiple outputs
func displayOutput(outputs []command.Output) {
	if len(outputs) == 0 {
		popStack()
		return
	}
	if outputs[0].IsMessage {
		modal := tview.NewModal()
		modal.AddButtons([]string{"Done"})
		modal.SetText(outputs[0].String)
		view.SetPrimitive(modal)
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) { popStack() })
		return
	}
	list := tview.NewList()
	for _, output := range outputs {
		options := output.Options
		list.AddItem(output.String, "", ' ', func() {
			if len(output.Options) > 0 {
				opt := options[0]
				util.Log("Running %+v with %+v", opt.Name, opt.Args)
				view.SetBorder(true).SetTitle(opt.Name)
				stack = append(stack, opt.Select)
				doCommand(opt.Select, opt.Args)
			}
		})
	}
	list.AddItem("Back", "", 'b', popStack).SetDoneFunc(popStack)
	view.SetPrimitive(list)
}
