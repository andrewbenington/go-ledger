package label

import (
	"errors"
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
			{Name: "Delete Label", Type: command.BoolArg},
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
	if len(args) == 3 && strings.EqualFold(args[2], "true") {
		err := ledger.RemoveLabel(args[0])
		if err != nil {
			return nil, err
		}
	}
	_, err := ledger.RemoveLabelKeywords(args[0], strings.Split(args[1], ","))
	if err != nil {
		return nil, err
	}
	return []command.Output{{
		String:    "successfully removed keywords",
		IsMessage: true,
	}}, nil
}

// func onAutoCompletedKeyword(text string, index int, field *tview.InputField) bool {
// 	fmt.Fprintln(os.Stderr, text, index, field.GetText())
// 	field.SetText(fmt.Sprintf("%s%s, ", field.GetText(), text))
// 	return true
// }
