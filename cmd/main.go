package main

import (
	"moneyman-writer-go/internal/driver/rest"
	x "moneyman-writer-go/internal/utils/logger"

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
	c := parseEnvConfig()
	r := rest.MakeRouter()
	r.Serve(c.Port)
}
