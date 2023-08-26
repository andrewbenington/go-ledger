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
	EditChaseCommand = &command.Command{
		Name:  "chase",
		Short: "update Chase source",
		ExpectedArgs: []command.Argument{
			{Name: "Old Name", IsConstant: true},
			{Name: "New Name"},
			{Name: "Last Four Digits"},
			{
				Name: "Account Type",
				Options: []command.ArgOption{
					{Label: "Credit", Value: "credit"},
					{Label: "Non-Credit", Value: "non-credit"},
				},
				Type: command.SelectArg,
			},
			{Name: "Directories"},
		},
		PrefillForm: true,
		Run:         EditChaseSource,
	}
)

func EditChaseSource(args []string) ([]command.Output, error) {
	if len(args) < 5 {
		return []command.Output{}, errors.New("must specify old source name, new source name, last four account digits, account type, and directories")
	}
	directories, err := util.SplitWordsIgnoreQuotes(args[4], ',')
	if err != nil {
		return []command.Output{}, fmt.Errorf("parse directories: %w", err)
	}
	updatedSource := source.ChaseSource{
		SourceName:  strings.TrimSpace(args[1]),
		LastDigits:  strings.TrimSpace(args[2]),
		AccountType: strings.ToLower(strings.TrimSpace(args[3])),
		Directories: directories,
	}
	err = source.EditChaseSource(args[0], updatedSource)
	if err != nil {
		return []command.Output{}, err
	}
	return []command.Output{
		{
			String:    fmt.Sprintf("Updated Chase source '%s' (%s)\n", updatedSource.SourceName, updatedSource.LastDigits),
			IsMessage: true,
		},
	}, nil
}
