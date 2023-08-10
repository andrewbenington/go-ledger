package ledger

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

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

func AllLabels() []Label {
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

func AddLabelKeywords(labelName string, toAdd []string) (l Label, created bool, err error) {
	var label *Label
	for i := range allLabels {
		if strings.EqualFold(labelName, allLabels[i].Name) {
			label = &allLabels[i]
			label.Keywords = append(label.Keywords, toAdd...)
			break
		}
	}
	if label == nil {
		created = true
		allLabels = append(allLabels, Label{Name: labelName, Keywords: toAdd})
	}
	err = saveLabels()
	if err != nil {
		return l, false, fmt.Errorf("error saving labels: %w", err)
	}
	return *label, created, err
}

// RemoveLabelKeywords removes the given keywords from the label with name labelName
// and saves the change to file
func RemoveLabelKeywords(labelName string, toRemove []string) (l Label, err error) {
	var label *Label
	for i := range allLabels {
		if strings.EqualFold(labelName, allLabels[i].Name) {
			label = &allLabels[i]
			break
		}
	}
	if label == nil {
		return l, fmt.Errorf("no label with name %s", labelName)
	}
	hash := make(map[string]bool)
	for _, r := range toRemove {
		stripped := strings.TrimSpace(r)
		if stripped == "" {
			continue
		}
		hash[stripped] = true
	}

	toKeep := []string{}
	for _, k := range label.Keywords {
		if _, ok := hash[k]; !ok {
			toKeep = append(toKeep, k)
		}
	}
	label.Keywords = toKeep
	err = saveLabels()
	if err != nil {
		return l, fmt.Errorf("error saving labels: %w", err)
	}
	return *label, nil
}
