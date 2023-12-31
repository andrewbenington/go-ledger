package app

import (
	"bufio"
	"fmt"
	"os"

	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// runCommand executes the the command with the given arguments and
// calls DisplayOutput with the output
func runCommand(c *command.Command, args []string) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// send stdout to the log file
	go func() {
		reader := bufio.NewReader(r)
		var err error = nil
		var line string
		for err == nil {
			line, err = reader.ReadString('\n')
			util.Log(line)
		}
		util.Log("Done logging")
	}()

	output, err := c.Run(args)
	if err != nil {
		output = []command.Output{{
			IsMessage: true,
			String:    fmt.Sprintf("ERROR:\n%s", err),
		}}
	}

	r.Close()
	w.Close()
	os.Stdout = old

	displayOutput(output)
}

// runCommandWithInput takes user input for the given command, runs the
// command with those inputs as arguments, and displays the output
func runCommandWithInput(c *command.Command, args []string) {
	dataInput := dataInputFromCommand(c, args)
	view.SetPrimitive(dataInput.form)
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
