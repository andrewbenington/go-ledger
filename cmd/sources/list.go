package sources

import (
	"github.com/andrewbenington/go-ledger/cmd/command"
	"github.com/andrewbenington/go-ledger/config"
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
	for _, source := range config.Sources() {
		outputs = append(outputs, command.Output{
			String:  source.Name(),
			Options: []command.OutputOption{},
		})
	}
	return outputs, nil
}
