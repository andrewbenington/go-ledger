package label

import (
	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/ledger"
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
	labelList := ledger.AllLabels()
	if len(labelList) == 0 {
		return []command.Output{{
			String:    "No labels configured",
			IsMessage: true,
		}}, nil
	}
	for _, label := range ledger.AllLabels() {
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
