package main

import (
	"encoding/json"
	"moneyman-writer-go/internal/model"
	x "moneyman-writer-go/internal/utils/logger"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port int `default:"8080"`
}

func parseEnvConfig() *Config {
	c := Config{}
	err := envconfig.Process("moneyman-writer", &c)
	if err != nil {
		x.Logger().Fatalw("failed to process env vars", "error", err)
	}

	return &c
}

func main() {
	// c := parseEnvConfig()
	// r := rest.MakeRouter()
	// r.Serve(c.Port)

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		x.Logger().Fatalw("failed to read file", "error", err)
	}

	x.Logger().Infow("file read", "len", len(data))

	txns := []model.Transaction{}
	err = json.Unmarshal(data, &txns)
	if err != nil {
		x.Logger().Fatalw("failed to unmarshal json", "error", err)
	}

	x.Logger().Infow("transactions read", "len", len(txns))
}
