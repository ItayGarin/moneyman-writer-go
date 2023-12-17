package main

import (
	"context"
	"moneyman-writer-go/internal/adapter/db"
	cloud_storage "moneyman-writer-go/internal/adapter/google/cloud-storage"
	"moneyman-writer-go/internal/core"
	"moneyman-writer-go/internal/driver/rest"
	x "moneyman-writer-go/internal/utils/logger"

	"github.com/kelseyhightower/envconfig"
	"google.golang.org/api/option"
)

type Config struct {
	Port            int    `default:"8080"`
	PostgresUrl     string `required:"true" split_words:"true"`
	MigrationsDir   string `required:"true" split_words:"true"`
	GcsCreds        string `required:"true" split_words:"true"`
	EnableMigration bool   `default:"false" split_words:"true"`
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
	ctx := context.Background()
	c := parseEnvConfig()

	client, err := db.NewClient(ctx, c.PostgresUrl)
	if err != nil {
		x.Logger().Fatalw("failed to initialize db client", "error", err)
	}
	if c.EnableMigration {
		x.Logger().Infow("running migrations", "dir", c.MigrationsDir)
		err = db.RunMigrations(ctx, c.PostgresUrl, c.MigrationsDir)
		if err != nil {
			x.Logger().Fatalw("failed to run migrations", "error", err)
		}
		x.Logger().Infow("finished migrating", "dir", c.MigrationsDir)
	}

	repo := db.NewPostgresTransactionRepo(client)
	downloader, err := cloud_storage.NewGcsDownloader(ctx, option.WithCredentialsJSON([]byte(c.GcsCreds)))
	if err != nil {
		x.Logger().Fatalw("failed to initialize gcs downloader", "error", err)
	}

	svc := core.NewService(repo, downloader)
	r := rest.MakeRouter(svc)
	r.Serve(c.Port)
}
