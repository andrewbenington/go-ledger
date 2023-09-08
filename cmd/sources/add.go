package sources

import (
	"github.com/andrewbenington/go-ledger/command"
)

var (
	AddCommand = &command.Command{
		Name:  "add",
		Short: "add source",
		SubCommands: []*command.Command{
			AddChaseCommand, AddVenmoCommand,
		},
	}
)
