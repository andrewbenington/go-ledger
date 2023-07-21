package label

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/andrewbenington/go-ledger/cmd/command"
	"github.com/rivo/tview"
)

var (
	RemoveKeywordCommand = &command.Command{
		Name:  "remove keywords",
		Short: "remove comma-separated keywords from label",
		ExpectedArgs: []command.ArgOptions{
			{Name: "Label", AutoComplete: autoCompleteLabel},
			{Name: "Keywords (Comma-separated)", AutoComplete: autoCompleteKeyword, OnAutoCompleted: onAutoCompletedKeyword},
		},
		Run: RemoveLabelKeyword,
	}
)

func RemoveLabelKeyword(args []string) ([]command.Output, error) {
	if len(args) < 2 {
		return []command.Output{}, errors.New("must specify label and at least one keyword")
	}
	for i, label := range allLabels {
		if strings.EqualFold(label.Name, args[0]) {
			removeKeywords(&allLabels[i], strings.Split(args[1], ","))
			err := saveLabels()
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

func removeKeywords(l *Label, toRemove []string) {
	fmt.Fprintln(os.Stderr, toRemove)
	hash := make(map[string]bool)
	for _, r := range toRemove {
		stripped := strings.TrimSpace(r)
		if stripped == "" {
			continue
		}
		hash[stripped] = true
	}

	toKeep := []string{}
	for _, k := range l.Keywords {
		if _, ok := hash[k]; !ok {
			toKeep = append(toKeep, k)
		}
	}
	l.Keywords = toKeep
}

func onAutoCompletedKeyword(text string, index int, field *tview.InputField) bool {
	fmt.Fprintln(os.Stderr, text, index, field.GetText())
	field.SetText(fmt.Sprintf("%s%s, ", field.GetText(), text))
	return true
}
