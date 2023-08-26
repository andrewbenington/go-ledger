package sources

import (
	"strings"

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
	for _, s := range sourceList {
		var editCommand *command.Command
		args := []string{}
		if chaseSource, ok := s.(*source.ChaseSource); ok {
			editCommand = EditChaseCommand
			args = []string{
				chaseSource.SourceName,
				chaseSource.SourceName,
				chaseSource.LastDigits,
				chaseSource.AccountType,
				strings.Join(chaseSource.Directories, ","),
			}
		}
		outputs = append(outputs, command.Output{
			String: s.Name(),
			Options: []command.OutputOption{
				{
					Name:   "Edit",
					Select: editCommand,
					Args:   args,
				},
			},
		})
	}
	return outputs, nil
}
