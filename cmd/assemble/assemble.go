package assemble

import (
	"strconv"

	"github.com/andrewbenington/go-ledger/cmd/command"
	"github.com/andrewbenington/go-ledger/config"
	"github.com/andrewbenington/go-ledger/excel"
)

var (
	AssembleCmd = &command.Command{
		Name:  "assemble",
		Short: "Assemble ledger from sources",
		Long:  `go-ledger will compile a ledger from the sources specified`,
		ExpectedArgs: []command.ArgOptions{
			{Name: "Years (Comma-separated)"},
		},
		Run:        Assemble,
		ShowOutput: true,
	}
	successOutput = command.Output{
		String:    "Successfully Assembled",
		IsMessage: true,
	}
	argErrorOutput = command.Output{
		String:    "Year must be provided",
		IsMessage: true,
	}
)

func Assemble(args []string) ([]command.Output, error) {
	cfg := config.GetConfig()
	if len(cfg.Sources.Venmo) > 0 {
	}
	year, err := strconv.Atoi(args[0])
	if err != nil {
		return []command.Output{argErrorOutput}, err
	}
	l, err := excel.LedgerFromFile(year)
	if err != nil {
		return []command.Output{}, err
	}
	err = l.UpdateFromSources(config.Sources())
	if err != nil {
		return []command.Output{}, err
	}
	return []command.Output{successOutput}, excel.WriteLedger(*l)
}
