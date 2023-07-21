/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package label

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/andrewbenington/go-ledger/cmd/command"
	"gopkg.in/yaml.v2"
)

type Label struct {
	Name     string
	Keywords []string
	re       *regexp.Regexp
}

const (
	filename = "config/labels.yaml"
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
	allLabels []Label = []Label{}
)

func init() {
	labels, err := loadLabels()
	if err != nil {
		fmt.Printf("Error loading labels: %s", err)
		return
	}
	allLabels = labels
	for i := range allLabels {
		re, err := allLabels[i].RegExp()
		if err != nil {
			fmt.Printf("Error loading labels: %s", err)
			continue
		}
		allLabels[i].re = re
	}
}

func All() []Label {
	return allLabels
}

func (l *Label) RegExp() (*regexp.Regexp, error) {
	if len(l.Keywords) == 0 {
		return nil, nil
	}
	reString := ""
	for i, keyword := range l.Keywords {
		if i > 0 {
			reString = fmt.Sprintf("%s|", reString)
		}
		reString = fmt.Sprintf("%s%s", reString, keyword)
	}
	return regexp.Compile(reString)
}

func loadLabels() ([]Label, error) {
	yamlFile, err := os.ReadFile("config/labels.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	labels := []Label{}
	err = yaml.Unmarshal(yamlFile, &labels)
	if err != nil {
		return nil, fmt.Errorf("error parsing labels: %s", err)
	}

	return labels, nil
}

func saveLabels() error {
	rawString, err := yaml.Marshal(allLabels)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, rawString, 0755)
}

func FindLabel(memo string) string {
	for _, label := range allLabels {
		if label.re != nil && label.re.MatchString(strings.ToLower(memo)) {
			return label.Name
		}
	}
	return ""
}

func autoCompleteLabel(current string, _ *[]string) []string {
	labelNames := []string{}
	for _, l := range allLabels {
		if strings.HasPrefix(l.Name, current) {
			labelNames = append(labelNames, l.Name)
		}
	}
	return labelNames
}

func autoCompleteKeyword(current string, currentArgs *[]string) []string {
	currentLabel := (*currentArgs)[0]
	for _, l := range allLabels {
		if l.Name == currentLabel {
			return l.Keywords
		}
	}
	return []string{}
}
