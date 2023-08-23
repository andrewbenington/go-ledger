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

func initialize() error {
	err := os.MkdirAll(configFolder, 0755)
	if err != nil {
		return fmt.Errorf("mkdir '%s': %w", configFolder, err)
	}
	return os.WriteFile(filepath.Join(configFolder, configFile), []byte{}, 0644)
}

func Get() (*Sources, error) {
	_, err := os.Stat(filepath.Join(configFolder, configFile))
	if err != nil {
		err = initialize()
		if err != nil {
			return nil, fmt.Errorf("initialize labels: %w", err)
		}
	}
	s := &Sources{}
	yamlFile, err := os.ReadFile(filepath.Join(configFolder, configFile))
	if err != nil {
		return nil, fmt.Errorf("read %s: %v ", filepath.Join(configFolder, configFile), err)
	}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		return nil, fmt.Errorf("error parsing sources: %s", err)
	}
	for i := range s.Chase {
		s := s.Chase[i]
		err := s.Validate()
		if err != nil {
			return nil, err
		}
	}
	for i := range s.Venmo {
		s := s.Venmo[i]
		err := s.Validate()
		if err != nil {
			return nil, err
		}
	}
	return s, nil
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
