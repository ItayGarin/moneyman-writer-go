package main

import (
	"context"
	"fmt"
	"moneyman-writer-go/internal/adapter/db"
	x "moneyman-writer-go/internal/utils/logger"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresUrl     string `required:"true" split_words:"true"`
	EnableMigration bool   `default:"false" split_words:"true"`
	MigrationsDir   string `required:"false" split_words:"true"`
	Inserts         bool   `default:"false" split_words:"true"`
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

	if len(os.Args) < 2 {
		x.Logger().Fatalw("missing output file")
	}

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
	uncategorized, err := repo.GetUncategorizedDescriptions(ctx)
	if err != nil {
		x.Logger().Fatalw("failed to get uncategorized descriptions", "error", err)
	}

	var lines []string
	if c.Inserts {
		for _, desc := range uncategorized {
			lines = append(lines, "INSERT INTO exp_businesses (name, category, sub_category) VALUES ('', '', '');")
			lines = append(lines, fmt.Sprintf("INSERT INTO exp_desc_to_business (description, business_name) VALUES ('%s', '');", desc))
			lines = append(lines, "")
		}
	} else {
		lines = uncategorized
	}

	out := strings.Join(lines, "\n")
	outFile := os.Args[1]
	err = os.WriteFile(outFile, []byte(out), 0644)
	if err != nil {
		x.Logger().Fatalw("failed to write output file", "error", err)
	}

	x.Logger().Infow("finished writing output file", "file", outFile)
}
