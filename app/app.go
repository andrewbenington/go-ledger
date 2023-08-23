package app

import (
	"log"

	"github.com/andrewbenington/go-ledger/cmd"
	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/config"
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	view  *tview.Frame
	stack []*command.Command
)

// Start creates a new tview application and runs it, allowing the user to
// run commands within the tview interface until the application exits
func Start() {
	if err := config.ReadConfig(); err != nil {
		log.Fatal(err)
	}
	app = tview.NewApplication()
	view = tview.NewFrame(nil)
	stack = []*command.Command{cmd.GetCommand()}
	doCommand(cmd.GetCommand(), nil)
	app.SetRoot(view, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// listFromSubcommands returns a tview.List including the subcommands for
// the given command.Command
func listFromSubcommands(c *command.Command) *tview.List {
	list := tview.NewList()
	for i := range c.SubCommands {
		sc := c.SubCommands[i]
		list.AddItem(sc.Name, sc.Short, []rune(sc.Name)[0], func() {
			LogStack()
			stack = append(stack, sc)
			doCommand(sc, nil)
		})
	}
	if len(stack) > 1 {
		return list.AddItem("back", "return to previous", 'b', popStack).SetDoneFunc(popStack)
	}
	return list.AddItem("exit", "exit go-ledger", 'x', func() { app.Stop() }).SetDoneFunc(func() { app.Stop() })
}

// popStack removes the most recent command from the application stack
// and executes the command now at the top
func popStack() {
	stack = stack[:len(stack)-1]
	prevCommand := stack[len(stack)-1]
	doCommand(prevCommand, nil)
}

// doCommand executes the given command, with input, output or
// neither depending on the command's usage
func doCommand(c *command.Command, args []string) {
	LogInterface(c)
	LogInterface(args)
	view.SetBorder(true).SetTitle(c.Name)
	if len(c.SubCommands) > 0 {
		view.SetPrimitive(listFromSubcommands(c))
		return
	}
	if len(c.ExpectedArgs) > 0 && len(args) == 0 {
		runCommandWithInput(c)
		return
	}
	if c.ShowLogs {
		runCommandWithLogs(c, args)
		return
	}
	runCommand(c, args)
}
