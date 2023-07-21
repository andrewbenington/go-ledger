package assemble

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andrewbenington/go-ledger/cmd/command"
	"github.com/andrewbenington/go-ledger/config"
	"github.com/andrewbenington/go-ledger/excel"
	"github.com/andrewbenington/go-ledger/ledger"
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
)

func Assemble(args []string) ([]command.Output, error) {
	l := &ledger.Ledger{}
	cfg := config.GetConfig()
	if len(cfg.Sources.Venmo) > 0 {

	}
	if len(args) > 0 {
		PopulateLedger(l, args)
	}
	err := l.UpdateFromSources(config.Sources())
	if err != nil {
		return []command.Output{}, err
	}
	successOutput := command.Output{
		String:    "Successfully Assembled",
		IsMessage: true,
	}
	return []command.Output{successOutput}, excel.WriteLedger(*l)
}

func PopulateLedger(l *ledger.Ledger, args []string) *ledger.Ledger {
	years := []int{}
	for _, arg := range args {
		year, err := strconv.Atoi(arg)
		if err == nil && year > 2000 {
			years = append(years, year)
		}
	}
	err := excel.PopulateLedger(l, years)
	if err != nil {
		fmt.Printf("error loading sources: %s", err)
		os.Exit(1)
	}
	return l
}
