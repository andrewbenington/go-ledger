package config

import (
	"fmt"
	"os"

	"github.com/andrewbenington/go-ledger/chase"
	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/andrewbenington/go-ledger/venmo"
	"gopkg.in/yaml.v2"
)

type SourcesConfig struct {
	Chase []chase.Source `yaml:"chase"`
	Venmo []venmo.Source `yaml:"venmo"`
}

func (s *SourcesConfig) read() error {
	yamlFile, err := os.ReadFile("config/sources.yaml")
	if err != nil {
		return fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		return fmt.Errorf("error parsing sources: %s", err)
	}
	for i := range s.Chase {
		s := s.Chase[i]
		err := s.Validate()
		if err != nil {
			return err
		}
	}
	for i := range s.Venmo {
		s := s.Venmo[i]
		err := s.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SourcesConfig) all() []ledger.Source {
	allSources := []ledger.Source{}
	for i := range s.Chase {
		allSources = append(allSources, &s.Chase[i])
	}
	for i := range s.Venmo {
		allSources = append(allSources, &s.Venmo[i])
	}
	return allSources
}
