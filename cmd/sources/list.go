package sources

import (
	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/source"
)

var (
	ListCommand = &command.Command{
		Name:  "list",
		Short: "list existing sources",
		Run:   ListSources,
	}
)

func ListSources(args []string) ([]command.Output, error) {
	outputs := []command.Output{}
	allSources, err := source.Get()
	if err != nil {
		return nil, err
	}
	sourceList := allSources.List()
	if len(sourceList) == 0 {
		return []command.Output{{
			String:    "No sources configured",
			IsMessage: true,
		}}, nil
	}
	for _, source := range sourceList {
		outputs = append(outputs, command.Output{
			String:  source.Name(),
			Options: []command.OutputOption{},
		})
	}
	return outputs, nil
}
