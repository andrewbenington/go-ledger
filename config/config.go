package config

import (
	"fmt"
	"log"

	"github.com/andrewbenington/go-ledger/ledger"
)

type Config struct {
	Sources SourcesConfig
}

var (
	config Config
)

func GetConfig() Config {
	err := ReadConfig()
	if err != nil {
		log.Fatalf("could not read config: %s", err)
	}
	return config
}

func ReadConfig() error {
	err := config.Sources.read()
	if err != nil {
		return fmt.Errorf("read sources: %s", err)
	}
	return nil
}

func Sources() []ledger.Source {
	err := ReadConfig()
	if err != nil {
		log.Fatalf("could not read config: %s", err)
	}
	return config.Sources.all()
}
