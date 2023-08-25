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
	AddChaseCommand = &command.Command{
		Name:  "chase",
		Short: "add Chase source",
		ExpectedArgs: []command.ArgOptions{
			{Name: "Name"},
			{Name: "Last Four Digits"},
			{Name: "Account Type"},
			{Name: "Directories"},
		},
		Run: AddChaseSource,
	}
)

func AddChaseSource(args []string) ([]command.Output, error) {
	if len(args) < 4 {
		return []command.Output{}, errors.New("must specify source name, last four account digits, account type, and directories")
	}

	directories, err := util.SplitWordsIgnoreQuotes(args[3])
	if err != nil {
		return []command.Output{}, fmt.Errorf("parse directories: %w", err)
	}
	newSource := source.ChaseSource{
		SourceName:  strings.TrimSpace(args[0]),
		LastDigits:  strings.TrimSpace(args[1]),
		AccountType: strings.ToLower(strings.TrimSpace(args[2])),
		Directories: directories,
	}
	err = source.AddChaseSource(newSource)
	if err != nil {
		return []command.Output{}, err
	}
	return []command.Output{
		{
			String:    fmt.Sprintf("Created Chase source '%s' (%s)\n", newSource.SourceName, newSource.LastDigits),
			IsMessage: true,
		},
	}, nil
}
