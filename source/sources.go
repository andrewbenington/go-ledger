package source

import (
	"fmt"
	"os"

	"github.com/andrewbenington/go-ledger/ledger"
	"gopkg.in/yaml.v2"
)

type Sources struct {
	Chase []ChaseSource `yaml:"chase"`
	Venmo []VenmoSource `yaml:"venmo"`
}

func Get() (*Sources, error) {
	s := &Sources{}
	yamlFile, err := os.ReadFile("config/sources.yaml")
	if err != nil {
		return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
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
