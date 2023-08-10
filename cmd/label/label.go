/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package label

import (
	"strings"

	"github.com/andrewbenington/go-ledger/command"
	"github.com/andrewbenington/go-ledger/ledger"
)

var (
	LabelCmd = &command.Command{
		Name:  "label",
		Short: "Manage labels",
		Long: `go-ledger can automatically classify transactions using labels you define.
		You can set a list of keywords for a given label, and transactions with any of those
		keywords will be assigned that label.`,
		SubCommands: []*command.Command{
			ListCommand, CheckCommand, AddKeywordCommand, RemoveKeywordCommand,
		},
	}
)

func autoCompleteLabel(current string, _ *[]string) []string {
	labelNames := []string{}
	for _, l := range ledger.AllLabels() {
		if strings.HasPrefix(l.Name, current) {
			labelNames = append(labelNames, l.Name)
		}
	}
	return labelNames
}

func autoCompleteKeyword(current string, currentArgs *[]string) []string {
	currentLabel := (*currentArgs)[0]
	for _, l := range ledger.AllLabels() {
		if l.Name == currentLabel {
			return l.Keywords
		}
	}
	return []string{}
}
