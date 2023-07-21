package label

import (
	"github.com/andrewbenington/go-ledger/cmd/command"
)

var (
	ListCommand = &command.Command{
		Name:  "list",
		Short: "list existing labels",
		Run:   ListLabels,
	}
)

func ListLabels(args []string) ([]command.Output, error) {
	outputs := []command.Output{}
	for _, label := range allLabels {
		outputs = append(outputs, command.Output{
			String: label.Name,
			Options: []command.OutputOption{
				{
					Name:   "View",
					Select: CheckCommand,
					Args:   []string{label.Name},
				},
			},
		})

	}
	return outputs, nil
}
