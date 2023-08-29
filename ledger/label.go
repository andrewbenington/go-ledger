package ledger

import (
	"fmt"
	"os"
	"path/filepath"
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
	configFolder = "config"
	configFile   = "labels.yaml"
)

var (
	allLabels []Label = []Label{}
)

func init() {
	labels, err := loadLabels()
	if err != nil {
		fmt.Printf("Error loading labels: %s\n", err)
		return
	}
	allLabels = labels
	for i := range allLabels {
		re, err := allLabels[i].RegExp()
		if err != nil {
			fmt.Printf("Error loading labels: %s\n", err)
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
		escaped := strings.ReplaceAll(keyword, "/", "//")
		if i > 0 {
			reString = fmt.Sprintf("%s|", reString)
		}
		reString = fmt.Sprintf("%s%s", reString, strings.ToLower(escaped))
	}
	return regexp.Compile(reString)
}

func initialize() error {
	err := os.MkdirAll(configFolder, 0755)
	if err != nil {
		return fmt.Errorf("mkdir '%s': %w", configFolder, err)
	}
	return os.WriteFile(filepath.Join(configFolder, configFile), []byte{}, 0644)
}

func loadLabels() ([]Label, error) {
	_, err := os.Stat(filepath.Join(configFolder, configFile))
	if err != nil {
		err = initialize()
		if err != nil {
			return nil, fmt.Errorf("initialize labels: %w", err)
		}
	}
	yamlFile, err := os.ReadFile(filepath.Join(configFolder, configFile))
	if err != nil {
		return nil, fmt.Errorf("read %s: %v ", filepath.Join(configFolder, configFile), err)
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
	return os.WriteFile(filepath.Join(configFolder, configFile), rawString, 0755)
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
		label = &Label{Name: labelName, Keywords: toAdd}
		allLabels = append(allLabels, *label)
	}
	err = saveLabels()
	if err != nil {
		return l, false, fmt.Errorf("save labels: %w", err)
	}
	return *label, created, err
}

func RemoveLabel(labelName string) error {
	for i := range allLabels {
		if strings.EqualFold(labelName, allLabels[i].Name) {
			allLabels = append(allLabels[:i], allLabels[i+1:]...)
			err := saveLabels()
			if err != nil {
				return fmt.Errorf("save labels: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("no label with name %s", labelName)
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
		return l, fmt.Errorf("save labels: %w", err)
	}
	return *label, nil
}
