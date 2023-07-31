package config

import (
	"fmt"
	"log"
	"os"

	"github.com/andrewbenington/go-ledger/ledger"
)

type Config struct {
	Sources SourcesConfig
	General GeneralConfig
}

var (
	config Config
)

func GetConfig() Config {
	err := ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read config: %s", err)
		os.Exit(1)
	}
	return config
}

func ReadConfig() error {
	err := config.General.read()
	if err != nil {
		return fmt.Errorf("read general config: %s", err)
	}
	err = config.Sources.read()
	if err != nil {
		return fmt.Errorf("read sources config: %s", err)
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
