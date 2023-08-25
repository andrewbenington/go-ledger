package import_transactions

import (
	"strconv"

	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/excel"
	"github.com/andrewbenington/go-ledger/source"
)

var (
	ImportCmd = &command.Command{
		Name:  "import",
		Short: "Import transactions from sources",
		Long: `go-ledger will search for files for each source, and import
any transactions in those files from the specified year`,
		ExpectedArgs: []command.Argument{
			{Name: "Year"},
		},
		Run:      Import,
		ShowLogs: true,
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

func Import(args []string) ([]command.Output, error) {
	allSources, err := source.Get()
	if err != nil {
		return []command.Output{}, err
	}
	year, err := strconv.Atoi(args[0])
	if err != nil {
		return []command.Output{argErrorOutput}, err
	}
	l, err := excel.LedgerFromFile(year)
	if err != nil {
		return []command.Output{}, err
	}
	err = l.UpdateFromSources(allSources.List())
	if err != nil {
		return []command.Output{}, err
	}
	return []command.Output{successOutput}, excel.WriteLedger(*l)
}
