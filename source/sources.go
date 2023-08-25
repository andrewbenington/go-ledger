package source

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/andrewbenington/go-ledger/ledger"
	"gopkg.in/yaml.v2"
)

type Sources struct {
	Chase []ChaseSource `yaml:"chase"`
	Venmo []VenmoSource `yaml:"venmo"`
}

const (
	configFolder = "config"
	configFile   = "sources.yaml"
)

func (s *Sources) Validate() error {
	for i := range s.Chase {
		c := s.Chase[i]
		err := c.Validate()
		if err != nil {
			return fmt.Errorf("chase[%d]: %w", i, err)
		}
	}
	for i := range s.Venmo {
		v := s.Venmo[i]
		err := v.Validate()
		if err != nil {
			return fmt.Errorf("venmo[%d]: %w", i, err)
		}
	}
	return nil
}

func (s *Sources) List() []ledger.Source {
	allSources := []ledger.Source{}
	for i := range s.Chase {
		allSources = append(allSources, &s.Chase[i])
	}
	for i := range s.Venmo {
		allSources = append(allSources, &s.Venmo[i])
	}
	return allSources
}

func ensureSourcesFile() error {
	_, err := os.Stat(filepath.Join(configFolder, configFile))
	if err == nil {
		return nil
	}
	err = os.MkdirAll(configFolder, 0755)
	if err != nil {
		return fmt.Errorf("mkdir '%s': %w", configFolder, err)
	}
	return os.WriteFile(filepath.Join(configFolder, configFile), []byte{}, 0644)
}

func Get() (*Sources, error) {
	err := ensureSourcesFile()
	if err != nil {
		return nil, fmt.Errorf("sources file: %w", err)
	}
	s := &Sources{}
	yamlFile, err := os.ReadFile(filepath.Join(configFolder, configFile))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w ", filepath.Join(configFolder, configFile), err)
	}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		return nil, fmt.Errorf("error parsing sources: %w", err)
	}
	err = s.Validate()
	if err != nil {
		return nil, fmt.Errorf("validate sources: %w", err)
	}
	return s, nil
}

func saveSources(s *Sources) error {
	err := ensureSourcesFile()
	if err != nil {
		return fmt.Errorf("sources file: %w", err)
	}
	rawBytes, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal sources: %w", err)
	}
	return os.WriteFile(filepath.Join(configFolder, configFile), rawBytes, 0755)
}
