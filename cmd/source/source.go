package source

import (
	"fmt"
	"log"
	"os"

	"github.com/andrewbenington/go-ledger/cmd/source/chase"
	"github.com/andrewbenington/go-ledger/cmd/source/venmo"
	"github.com/andrewbenington/go-ledger/ledger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	sources, err := LoadSources()
	if err != nil {
		log.Fatalf("Error loading sources: %s\n", err)
	}
	for _, source := range sources {
		allSources[source.Name()] = source
	}
}

var (
	SourceCmd = &cobra.Command{
		Use:   "source",
		Short: "Manage sources",
		Long:  `go-ledger can pull transaction data from multiple sources`,
	}
	allSources map[string]Source = make(map[string]Source)
)

type Source interface {
	Name() string
	GetLedgerEntries() ([]ledger.LedgerEntry, error)
}

type FileSource interface {
}

type SourcesConfig struct {
	Chase []chase.Source `yaml:"chase"`
	Venmo []venmo.Source `yaml:"venmo"`
}

func All() []Source {
	sources := make([]Source, 0, len(allSources))

	for _, value := range allSources {
		sources = append(sources, value)
	}
	return sources
}

func LoadSources() ([]Source, error) {
	yamlFile, err := os.ReadFile("sources.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	sourcesConfig := SourcesConfig{}
	err = yaml.Unmarshal(yamlFile, &sourcesConfig)
	if err != nil {
		return nil, fmt.Errorf("error parsing sources: %s", err)
	}
	allSources := []Source{}
	for i := range sourcesConfig.Chase {
		s := sourcesConfig.Chase[i]
		err = s.Validate()
		if err != nil {
			return nil, err
		}
		allSources = append(allSources, &s)
	}
	for i := range sourcesConfig.Venmo {
		s := sourcesConfig.Venmo[i]
		err = s.Validate()
		if err != nil {
			return nil, err
		}
		allSources = append(allSources, &s)
	}

	return allSources, nil
}
