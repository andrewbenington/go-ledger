package sources

import (
	"strings"

	"github.com/andrewbenington/go-ledger/command"
)

var (
	AddCommand = &command.Command{
		Name:  "add",
		Short: "add source",
		SubCommands: []*command.Command{
			AddChaseCommand,
		},
	}
)

func stripKeywords(inputs []string) []string {
	stripped := []string{}
	for _, input := range inputs {
		keyword := strings.TrimSpace(input)
		if input != "" {
			stripped = append(stripped, keyword)
		}
	}
	return stripped
}
