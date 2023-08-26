package sources

import (
	"github.com/andrewbenington/go-ledger/command"
)

var (
	EditCommand = &command.Command{
		Name:  "edit",
		Short: "edit source",
		SubCommands: []*command.Command{
			EditChaseCommand,
		},
	}
)
