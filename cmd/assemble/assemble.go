package assemble

import (
	"fmt"
	"sort"

	"github.com/andrewbenington/go-ledger/cmd/source"
	"github.com/andrewbenington/go-ledger/cmd/source/csv"
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
	fmt.Println("ledger:")
	l := LedgerFromSources()
	fmt.Printf("%+v\n", len(l.Entries))
	err := csv.WriteLedger(l, "ledger.csv")
	if err != nil {
		fmt.Println(err)
	}
	err = excel.WriteLedger(l)
	if err != nil {
		fmt.Println(err)
	}
}

func LedgerFromSources() ledger.Ledger {
	allSources := source.All()
	l := ledger.Ledger{}
	for _, source := range allSources {
		fmt.Printf("getting entries from %s...\n", source.Name())
		entries, err := source.GetLedgerEntries()
		if err != nil {
			fmt.Printf("Error getting entries from %s: %e", source.Name(), err)
			continue
		}
		l.Entries = append(l.Entries, entries...)
	}
	sort.Sort(l)
	return l
}
