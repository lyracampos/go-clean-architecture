package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/lyracampos/go-clean-architecture/config"
	"github.com/lyracampos/go-clean-architecture/internal/app"
)

const (
	defaultConfigFilePath = "../config/config.yaml"
	apiEntrypoint         = "api"
	workerEntrypoint      = "worker"
)

var errInvalidAppEntrypoint = errors.New("invalid entrypoint, must be one of [api, worker]")

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var appEntrypoint string
	var configFilePath string

	flag.StringVar(&appEntrypoint, "e", apiEntrypoint, "Entrypoint to define which application will be started. [api, worker]")
	flag.StringVar(&configFilePath, "c", defaultConfigFilePath, "File path with app configs file.")

	flag.Parse()

	config, err := config.NewConfig(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	switch appEntrypoint {
	case apiEntrypoint:
		app.RunAPI(config)
	case workerEntrypoint:
		app.RunWorker(config)
	default:
		return errInvalidAppEntrypoint
	}

	return nil
}
