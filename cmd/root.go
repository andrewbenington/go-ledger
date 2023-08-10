package cmd

import (
	"os"

	"github.com/andrewbenington/go-ledger/cmd/import_transactions"
	"github.com/andrewbenington/go-ledger/cmd/label"
	"github.com/andrewbenington/go-ledger/cmd/sources"
	"github.com/andrewbenington/go-ledger/command"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &command.Command{
	Name:  "go-ledger",
	Short: "application for managing transactions and budgets",
	Long: `Go-Ledger manages a history of transactions. It can automatically import, label, and de-duplicate transactions
	from source CSV files.`,
	SubCommands: []*command.Command{
		label.LabelCmd, import_transactions.ImportCmd, sources.SourceCmd,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.ToCobra().Execute()
	if err != nil {
		os.Exit(1)
	}
}

func GetCommand() *command.Command {
	return rootCmd
}
