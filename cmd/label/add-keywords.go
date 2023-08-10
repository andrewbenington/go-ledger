package label

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/ledger"
)

var (
	AddKeywordCommand = &command.Command{
		Name:  "add keyword",
		Short: "add comma-separated keywords to label",
		ExpectedArgs: []command.ArgOptions{
			{Name: "Label (new or existing)", AutoComplete: autoCompleteLabel},
			{Name: "Keywords (Comma-separated)"},
		},
		Run: AddLabelKeyword,
	}
)

func AddLabelKeyword(args []string) ([]command.Output, error) {
	if len(args) < 2 {
		return []command.Output{}, errors.New("must specify label and at least one keyword")
	}
	newKeywords := stripKeywords(strings.Split(args[1], ","))
	label, isNew, err := ledger.AddLabelKeywords(args[0], newKeywords)
	if err != nil {
		return []command.Output{}, err
	}
	outputString := fmt.Sprintf("Updated label '%s'\n", args[0])
	if isNew {
		outputString = fmt.Sprintf("Created label '%s'\n", args[0])
	}
	return []command.Output{
		{
			String:    fmt.Sprintf("%sKeywords:\n%s", outputString, strings.Join(label.Keywords, "\n")),
			IsMessage: true,
		},
	}, nil
}

func stripKeywords(inputs []string) []string {
	stripped := []string{}
	for _, input := range inputs {
		keyword := strings.TrimSpace(input)
		if input != "" {
			stripped = append(stripped, keyword)
		}
	}
	return stripped
}
