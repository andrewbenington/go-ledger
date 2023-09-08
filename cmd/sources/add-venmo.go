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
	AddVenmoCommand = &command.Command{
		Name:  "venmo",
		Short: "add Venmo source",
		ExpectedArgs: []command.Argument{
			{Name: "Name"},
			{Name: "Account Holder Name"},
			{
				Name: "Ignore Transfers",
				Type: command.BoolArg,
			},
			{Name: "Directories"},
		},
		Run: AddVenmoSource,
	}
)

func AddVenmoSource(args []string) ([]command.Output, error) {
	if len(args) < 4 {
		return []command.Output{}, errors.New("must specify source name, account holder name, ignore transfers, and directories")
	}

	directories, err := util.SplitWordsIgnoreQuotes(args[3], ',')
	if err != nil {
		return []command.Output{}, fmt.Errorf("parse directories: %w", err)
	}
	newSource := source.VenmoSource{
		SourceName:        strings.TrimSpace(args[0]),
		AccountHolderName: strings.TrimSpace(args[1]),
		HideTransfers:     strings.EqualFold(args[2], "true"),
		Directories:       directories,
	}
	err = source.AddVenmoSource(newSource)
	if err != nil {
		return []command.Output{}, err
	}
	return []command.Output{
		{
			String:    fmt.Sprintf("Created Venmo source '%s' (%s)\n", newSource.SourceName, newSource.AccountHolderName),
			IsMessage: true,
		},
	}, nil
}
