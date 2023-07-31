package app

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/andrewbenington/go-ledger/cmd"
	"github.com/andrewbenington/go-ledger/cmd/command"
	"github.com/andrewbenington/go-ledger/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	view  *tview.Frame
	stack []*command.Command
)

func Start() {
	if err := config.ReadConfig(); err != nil {
		log.Fatal(err)
	}
	app = tview.NewApplication()
	view = tview.NewFrame(nil)
	stack = []*command.Command{cmd.GetCommand()}
	DoCommand(cmd.GetCommand())
	app.SetRoot(view, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func ListFromSubcommands(c *command.Command) *tview.List {
	// Log("ListCommand", "go-ledger.log")
	// LogInterface(c)
	list := tview.NewList()
	for i := range c.SubCommands {
		sc := c.SubCommands[i]
		list.AddItem(sc.Name, sc.Short, []rune(sc.Name)[0], func() {
			LogStack()
			stack = append(stack, sc)
			DoCommand(sc)
		})
	}
	if len(stack) > 1 {
		return list.AddItem("back", "return to previous", 'b', PopStack).SetDoneFunc(PopStack)
	}
	return list.AddItem("exit", "exit go-ledger", 'x', func() { app.Stop() }).SetDoneFunc(func() { app.Stop() })
}

func PopStack() {
	Log("PopStack")
	stack = stack[:len(stack)-1]
	prevCommand := stack[len(stack)-1]
	LogStack()
	DoCommand(prevCommand)
}

func DoCommand(c *command.Command) {
	Log("DoCommand %s", c.Name)
	fmt.Fprintln(os.Stderr, "test error")
	LogStack()
	view.SetBorder(true).SetTitle(c.Name)
	if len(c.SubCommands) > 0 {
		view.SetPrimitive(ListFromSubcommands(c))
		return
	}
	if len(c.ExpectedArgs) > 0 {
		RunCommandWithInput(c)
		return
	}
	if c.ShowOutput {
		RunCommandWithOutput(c, nil)
		return
	}
	RunCommand(c, nil)
}

func RunCommandWithInput(c *command.Command) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	args := make([]string, len(c.ExpectedArgs))
	fields := make([]*tview.InputField, len(c.ExpectedArgs))

	doneBtn := tview.NewButton("Done").SetSelectedFunc(func() {
		if c.ShowOutput {
			RunCommandWithOutput(c, args)
			return
		}
		RunCommand(c, args)
	})
	cancelBtn := tview.NewButton("Cancel").SetSelectedFunc(PopStack)

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

func RunCommand(c *command.Command, args []string) {
	output, err := c.Run(args)
	if err != nil {
		output = []command.Output{{
			IsMessage: true,
			String:    err.Error(),
		}}
	}
	DisplayOutput(output)
}

func RunCommandWithOutput(c *command.Command, args []string) {
	Log("RunCommandWithOutput %s %+v", c.Name, args)
	LogStack()
	Log("stderr: %p", os.Stderr)
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
			r.Close()
			w.Close()
			os.Stdout = old
		}()
		output, err := c.Run(args)
		if err != nil {
			output = []command.Output{{
				IsMessage: true,
				String:    err.Error(),
			}}
		}
		_, err = w.WriteString("(Enter or ESC to continue)\n")
		if err != nil {
			LogErr(err.Error())
		}
		app.QueueUpdate(func() {
			textView.
				SetDoneFunc(func(key tcell.Key) {
					if key == tcell.KeyEnter || key == tcell.KeyESC {
						// Queue because we aren't in main goroutine
						DisplayOutput(output)
					}
				})
		})

	}()

	// display standard out to screen
	go func() {
		reader := bufio.NewReader(r)
		var err error = nil
		var line string
		for err == nil {
			line, err = reader.ReadString('\n')
			Log(line)
			_, _ = textView.Write([]byte(line))
		}
	}()
}
