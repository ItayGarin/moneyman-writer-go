package main

import (
	"context"
	"moneyman-writer-go/internal/adapter/db"
	"moneyman-writer-go/internal/core"
	x "moneyman-writer-go/internal/utils/logger"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port int `default:"8080"`

	PostgresUrl     string `required:"true" split_words:"true"`
	EnableMigration bool   `default:"false" split_words:"true"`
	MigrationsDir   string `required:"false" split_words:"true"`
}

func parseEnvConfig() *Config {
	c := Config{}
	err := envconfig.Process("", &c)
	if err != nil {
		x.Logger().Fatalw("failed to process env vars", "error", err)
	}

	return &c
}

func main() {
	x.InitDev()
	ctx := context.Background()
	c := parseEnvConfig()

	client, err := db.NewClient(ctx, c.PostgresUrl)
	if err != nil {
		x.Logger().Fatalw("failed to initialize db client", "error", err)
	}
	if c.EnableMigration {
		if c.MigrationsDir == "" {
			x.Logger().Fatalw("migration dir is required")
		} else {
			x.Logger().Infow("running migrations", "dir", c.MigrationsDir)
			err = db.RunMigrations(ctx, c.PostgresUrl, c.MigrationsDir)
			if err != nil {
				x.Logger().Fatalw("failed to run migrations", "error", err)
			}
			x.Logger().Infow("finished migrating", "dir", c.MigrationsDir)
		}
	}

	repo := db.NewPostgresTransactionRepo(client)
	svc := core.NewService(repo)

	if len(os.Args) < 2 {
		x.Logger().Fatalw("no files provided")
		return
	}

	files := os.Args[1:]
	for _, file := range files {
		x.Logger().Infow("reading file", "file", file)
		data, err := os.ReadFile(file)
		if err != nil {
			x.Logger().Fatalw("failed to read file", "error", err)
		}

		err = svc.SaveNewTransactionsFromBlob(ctx, data)
		if err != nil {
			x.Logger().Fatalw("failed to save transactions", "error", err)
		}
	}

	x.Logger().Infow("saved all transactions", "count", len(files))
}
