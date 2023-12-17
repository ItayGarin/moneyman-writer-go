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
	Port int `default:"8080"`

	PostgresUrl     string `required:"true" split_words:"true"`
	EnableMigration bool   `default:"false" split_words:"true"`
	MigrationsDir   string `required:"false" split_words:"true"`

	GcsCreds string `required:"false" split_words:"true"`
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
	opts := []option.ClientOption{}
	if c.GcsCreds != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(c.GcsCreds)))
	}
	downloader, err := cloud_storage.NewGcsDownloader(ctx, opts...)
	if err != nil {
		x.Logger().Fatalw("failed to initialize gcs downloader", "error", err)
	}

	svc := core.NewService(repo)
	r := rest.MakeRouter(svc, downloader)
	r.Serve(c.Port)
}
