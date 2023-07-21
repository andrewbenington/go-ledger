package label

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andrewbenington/go-ledger/cmd/command"
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
	existingLabel := false
	var label *Label
	for i := range allLabels {
		label = &allLabels[i]
		if strings.EqualFold(label.Name, args[0]) {
			existingLabel = true
			label.Keywords = append(label.Keywords, newKeywords...)
			break
		}
	}
	outputString := ""
	if !existingLabel {
		allLabels = append(allLabels, Label{Name: args[0], Keywords: newKeywords})
		outputString = fmt.Sprintf("Created label '%s'\n", args[0])
	}
	err := saveLabels()
	if err != nil {
		return []command.Output{}, fmt.Errorf("error saving labels: %w", err)
	}
	if err != nil {
		return []command.Output{}, err
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
