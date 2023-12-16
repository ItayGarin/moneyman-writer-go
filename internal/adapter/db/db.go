package db

import (
	"context"
	"moneyman-writer-go/internal/adapter/db/sqlc/sql"

	"github.com/jackc/pgx/v5"
)

type DebugTracer struct{}

func (dt *DebugTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	// fmt.Println(data.SQL)
	// fmt.Println(data.Args)
	return ctx
}
func (dt *DebugTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func NewClient(ctx context.Context, conn string) (*sql.Queries, error) {
	cfg, err := pgx.ParseConfig(conn)
	if err != nil {
		return nil, err
	}

	cfg.Tracer = &DebugTracer{}

	c, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return sql.New(c), nil
}
