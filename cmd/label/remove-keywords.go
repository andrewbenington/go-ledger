package label

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/ledger"
)

var (
	RemoveKeywordCommand = &command.Command{
		Name:  "remove keywords",
		Short: "remove comma-separated keywords from label",
		ExpectedArgs: []command.ArgOptions{
			{Name: "Label", AutoComplete: autoCompleteLabel},
			{Name: "Keywords (Comma-separated)"},
			// auto complete logic is broken
			// {Name: "Keywords (Comma-separated)", AutoComplete: autoCompleteKeyword, OnAutoCompleted: onAutoCompletedKeyword},
		},
		Run: RemoveLabelKeyword,
	}
)

func RemoveLabelKeyword(args []string) ([]command.Output, error) {
	if len(args) < 2 {
		return []command.Output{}, errors.New("must specify label and at least one keyword")
	}
	allLabels := ledger.AllLabels()
	for _, label := range allLabels {
		if strings.EqualFold(label.Name, args[0]) {
			_, err := ledger.RemoveLabelKeywords(label.Name, strings.Split(args[1], ","))
			if err != nil {
				return nil, fmt.Errorf("error saving labels: %w", err)
			}
			return []command.Output{{
				String:    "successfully removed keywords",
				IsMessage: true,
			}}, nil
		}
	}
	return nil, fmt.Errorf("that label doesn't exist")
}

// func onAutoCompletedKeyword(text string, index int, field *tview.InputField) bool {
// 	fmt.Fprintln(os.Stderr, text, index, field.GetText())
// 	field.SetText(fmt.Sprintf("%s%s, ", field.GetText(), text))
// 	return true
// }
