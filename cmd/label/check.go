package label

import (
	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/ledger"
)

var (
	CheckCommand = &command.Command{
		Name:  "check",
		Short: "check label keywords",
		ExpectedArgs: []command.Argument{
			{Name: "Label", AutoComplete: autoCompleteLabel},
		},
		Run: CheckLabel,
	}
)

func CheckLabel(args []string) ([]command.Output, error) {
	var label ledger.Label
	for _, label = range ledger.AllLabels() {
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
