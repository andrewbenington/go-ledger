package open

import (
	"os/exec"

	"github.com/andrewbenington/go-ledger/command"
)

var (
	OpenCmd = &command.Command{
		Name:  "open",
		Short: "Open Excel document",
		Run:   Open,
	}
)

func Open(args []string) ([]command.Output, error) {
	openCmd := exec.Command("open", "2023.xlsx")
	err := openCmd.Run()
	return []command.Output{}, err
}
