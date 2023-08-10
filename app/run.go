package app

import (
	"bufio"
	"os"

	"github.com/andrewbenington/go-ledger/command"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// runCommand executes the the command with the given arguments and
// calls DisplayOutput with the output
func runCommand(c *command.Command, args []string) {
	output, err := c.Run(args)
	if err != nil {
		output = []command.Output{{
			IsMessage: true,
			String:    err.Error(),
		}}
	}
	displayOutput(output)
}

// runCommandWithInput takes user input for the given command, runs the
// command with those inputs as arguments, and displays the output
func runCommandWithInput(c *command.Command) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	args := make([]string, len(c.ExpectedArgs))
	fields := make([]*tview.InputField, len(c.ExpectedArgs))

	doneBtn := tview.NewButton("Done").SetSelectedFunc(func() {
		if c.ShowLogs {
			runCommandWithLogs(c, args)
			return
		}
		runCommand(c, args)
	})
	cancelBtn := tview.NewButton("Cancel").SetSelectedFunc(popStack)

	doneBtn.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyTAB, tcell.KeyDown, tcell.KeyESC:
			app.SetFocus(cancelBtn)
		case tcell.KeyBacktab, tcell.KeyUp:
			app.SetFocus(fields[len(c.ExpectedArgs)-1])
		}
		return key
	})
	cancelBtn.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyBacktab, tcell.KeyUp:
			app.SetFocus(doneBtn)
		}
		return key
	})

	for i := range c.ExpectedArgs {
		index := i
		arg := c.ExpectedArgs[index]
		flex.AddItem(tview.NewTextView().SetText(arg.Name), 1, 0, false)
		// handle field completion
		fields[index] = tview.NewInputField().
			SetChangedFunc(func(text string) {
				args[index] = text
			}).
			SetDoneFunc(func(key tcell.Key) {
				switch key {
				case tcell.KeyTAB, tcell.KeyEnter:
					if index < len(c.ExpectedArgs)-1 {
						app.SetFocus(fields[index+1])
					} else {
						app.SetFocus(doneBtn)
					}
				case tcell.KeyBacktab:
					if index > 0 {
						app.SetFocus(fields[index-1])
					}
				case tcell.KeyESC:
					app.SetFocus(cancelBtn)
				}
			})
		// handle arrow navigation
		fields[index].SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
			if arg.AutoComplete != nil {
				return key
			}
			switch key.Key() {
			case tcell.KeyUp:
				if index > 0 {
					app.SetFocus(fields[index-1])
				}
			case tcell.KeyDown:
				if index < len(c.ExpectedArgs)-1 {
					app.SetFocus(fields[index+1])
				} else {
					app.SetFocus(doneBtn)
				}
			}
			return key
		})
		if arg.AutoComplete != nil {
			fields[index].SetAutocompleteFunc(arg.AutoCompleteWithArgs(&args))
		}
		if arg.OnAutoCompleted != nil {
			fields[index].SetAutocompletedFunc(arg.OnAutoCompletedWithField(fields[index]))
		}
		flex.AddItem(fields[index], 1, 0, index == 0)
		flex.AddItem(nil, 1, 0, false)
	}
	flex.AddItem(doneBtn, 1, 0, false)
	flex.AddItem(cancelBtn, 1, 0, false)
	view.SetPrimitive(flex)
}

// runCommandWithLogs executes the command with the given arguments, outputs
// its stdout to a scrollable text view, and calls DisplayOutput with the output
// after the user is done viewing the logs
func runCommandWithLogs(c *command.Command, args []string) {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	view.SetPrimitive(textView)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	go func() {
		defer func() {
		}()
		output, err := c.Run(args)
		if err != nil {
			output = []command.Output{{
				IsMessage: true,
				String:    err.Error(),
			}}
		}

		r.Close()
		w.Close()
		os.Stdout = old

		_, _ = textView.Write([]byte("(Enter or ESC to continue)\n"))

		// Queue because we aren't in main goroutine
		app.QueueUpdate(func() {
			textView.
				SetDoneFunc(func(key tcell.Key) {
					if key == tcell.KeyEnter || key == tcell.KeyESC {
						displayOutput(output)
					}
				})
		})
	}()

	// print stdout to the textview
	go func() {
		reader := bufio.NewReader(r)
		var err error = nil
		var line string
		for err == nil {
			line, err = reader.ReadString('\n')
			_, _ = textView.Write([]byte(line))
		}
	}()
}
