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

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Label struct {
	Name     string
	Keywords []string
}

var (
	LabelCmd = &cobra.Command{
		Use:   "label",
		Short: "Manage labels",
		Long: `go-ledger can automatically classify transactions using labels you define.
		You can set a list of keywords for a given label, and transactions with any of those
		keywords will be assigned that label.`,
	}
	allLabels       map[string]Label          = make(map[string]Label)
	allLabelRegExps map[string]*regexp.Regexp = make(map[string]*regexp.Regexp)
)

func init() {
	labels, err := loadLabels()
	if err != nil {
		fmt.Printf("Error loading labels: %s", err)
		return
	}
	for _, label := range *labels {
		allLabels[label.Name] = label
		re, err := label.RegExp()
		if err != nil {
			fmt.Printf("Error loading labels: %s", err)
			continue
		}
		allLabelRegExps[label.Name] = re
	}
}

func All() map[string]Label {
	return allLabels
}

func (l *Label) RegExp() (*regexp.Regexp, error) {
	reString := ""
	for i, keyword := range l.Keywords {
		if i > 0 {
			reString = fmt.Sprintf("%s|", reString)
		}
		reString = fmt.Sprintf("%s%s", reString, keyword)
	}
	return regexp.Compile(reString)
}

func loadLabels() (*[]Label, error) {
	yamlFile, err := os.ReadFile("labels.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	labels := &[]Label{}
	err = yaml.Unmarshal(yamlFile, labels)
	if err != nil {
		return nil, fmt.Errorf("error parsing labels: %s", err)
	}

	return labels, nil
}

func FindLabel(memo string) string {
	for labelName, labelRegExp := range allLabelRegExps {
		if labelRegExp.MatchString(strings.ToLower(memo)) {
			return labelName
		}
	}
	return ""
}
