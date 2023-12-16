package db

import (
	"context"
	"fmt"
	"moneyman-writer-go/internal/adapter/db/sqlc/sql"
	"os"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/jackc/pgx/v5"
)

type DebugTracer struct{}

func (dt *DebugTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	return ctx
}
func (dt *DebugTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func NewClient(ctx context.Context, conn string) (*sql.Queries, error) {
	cfg, err := pgx.ParseConfig(conn)
	if err != nil {
		return nil, err
	}

	// cfg.Tracer = &DebugTracer{}
	c, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	err = c.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return sql.New(c), nil
}

func RunMigrations(ctx context.Context, conn string, migrationsDir string) error {
	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS(migrationsDir),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to load working directory: %w", err)
	}
	defer workdir.Close()

	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	_, err = client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: conn,
	})
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
