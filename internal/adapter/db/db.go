package db

import (
	"context"
	"moneyman-writer-go/internal/adapter/db/sqlc/sql"

	"github.com/jackc/pgx/v5"
)

func NewClient(ctx context.Context, conn string) (*sql.Queries, error) {
	c, err := pgx.Connect(ctx, conn)
	if err != nil {
		return nil, err
	}

	return sql.New(c), nil
}
