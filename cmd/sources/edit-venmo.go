package sources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/source"
	"github.com/andrewbenington/go-ledger/util"
)

var (
	EditVenmoCommand = &command.Command{
		Name:  "venmo",
		Short: "update Venmo source",
		ExpectedArgs: []command.Argument{
			{Name: "Old Name", IsConstant: true},
			{Name: "New Name"},
			{Name: "Account Holder Name"},
			{
				Name: "Ignore Transfers",
				Type: command.BoolArg,
			},
			{Name: "Directories"},
		},
		PrefillForm: true,
		Run:         EditVenmoSource,
	}
)

func EditVenmoSource(args []string) ([]command.Output, error) {
	if len(args) < 5 {
		return []command.Output{}, errors.New("must specify old source name, new source name, account holder name, ignore transfers, and directories")
	}
	directories, err := util.SplitWordsIgnoreQuotes(args[4], ',')
	if err != nil {
		return []command.Output{}, fmt.Errorf("parse directories: %w", err)
	}
	updatedSource := source.VenmoSource{
		SourceName:        strings.TrimSpace(args[1]),
		AccountHolderName: strings.TrimSpace(args[2]),
		HideTransfers:     strings.EqualFold(args[3], "true"),
		Directories:       directories,
	}
	err = source.EditVenmoSource(args[0], updatedSource)
	if err != nil {
		return []command.Output{}, err
	}
	return []command.Output{
		{
			String:    fmt.Sprintf("Updated Venmo source '%s' (%s)\n", updatedSource.SourceName, updatedSource.AccountHolderName),
			IsMessage: true,
		},
	}, nil
}
