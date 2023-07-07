package assemble

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andrewbenington/go-ledger/config"
	"github.com/andrewbenington/go-ledger/excel"
	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/spf13/cobra"
)

var (
	AssembleCmd = &cobra.Command{
		Use:   "assemble",
		Short: "Assemble ledger from sources",
		Long:  `go-ledger will compile a ledger from the sources specified`,
		Run:   Assemble,
	}
)

func Assemble(cmd *cobra.Command, args []string) {
	l := &ledger.Ledger{}
	if len(args) > 0 {
		l = ReadLedger(args)
	}
	l.UpdateFromSources(config.Sources())
	// err := csv.WriteLedger(l, "ledger.csv")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	err := excel.WriteLedger(*l)
	if err != nil {
		fmt.Println(err)
	}
}

func ReadLedger(args []string) *ledger.Ledger {
	years := []int{}
	for _, arg := range args {
		year, err := strconv.Atoi(arg)
		if err == nil && year > 2000 {
			years = append(years, year)
		}
	}
	l, err := excel.ReadLedger(years)
	if err != nil {
		fmt.Printf("error loading sources: %s", err)
		os.Exit(1)
	}
	return l
}
