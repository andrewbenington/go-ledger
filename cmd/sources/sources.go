package sources

import "github.com/andrewbenington/go-ledger/command"

var (
	SourceCmd = &command.Command{
		Name:  "source",
		Short: "Manage sources",
		Long:  `go-ledger can pull transaction data from multiple sources`,
		SubCommands: []*command.Command{
			ListCommand, AddCommand,
		},
	}
)
