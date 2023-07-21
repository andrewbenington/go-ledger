package label

import (
	"github.com/andrewbenington/go-ledger/cmd/command"
)

var (
	CheckCommand = &command.Command{
		Name:  "check",
		Short: "check label keywords",
		ExpectedArgs: []command.ArgOptions{
			{Name: "Label", AutoComplete: autoCompleteLabel},
		},
		Run: CheckLabel,
	}
)

func CheckLabel(args []string) ([]command.Output, error) {
	var label Label
	for _, label = range allLabels {
		if label.Name == args[0] {
			break
		}
	}
	outputs := []command.Output{}
	for _, keyword := range label.Keywords {
		outputs = append(outputs, command.Output{String: keyword})
	}
	return outputs, nil
}
