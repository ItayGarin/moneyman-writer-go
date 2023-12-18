package db

import (
	"context"
	"moneyman-writer-go/internal/adapter/db/sqlc/sql"
	"moneyman-writer-go/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type PostgresTransactionRepo struct {
	s *sql.Queries
}

func NewPostgresTransactionRepo(s *sql.Queries) *PostgresTransactionRepo {
	return &PostgresTransactionRepo{
		s: s,
	}
}

func toTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func txnToInsertParams(txn *model.Transaction) *sql.InsertTransactionParams {
	return &sql.InsertTransactionParams{
		Identifier:       txn.Identifier,
		Type:             txn.Type,
		Status:           txn.Status,
		Date:             toTimestamptz(txn.Date),
		ProcessedDate:    toTimestamptz(txn.ProcessedDate),
		OriginalAmount:   txn.OriginalAmount,
		OriginalCurrency: txn.OriginalCurrency,
		ChargedAmount:    txn.ChargedAmount,
		ChargedCurrency:  txn.ChargedCurrency,
		Description:      txn.Description,
		Memo:             txn.Memo,
		Category:         txn.Category,
		Account:          txn.Account,
		CompanyID:        txn.CompanyId,
		Hash:             txn.Hash,
	}
}

func (r *PostgresTransactionRepo) Save(ctx context.Context, txn *model.Transaction) error {
	params := txnToInsertParams(txn)
	return r.s.InsertTransaction(ctx, *params)
}

func (r *PostgresTransactionRepo) GetUncategorizedDescriptions(ctx context.Context) ([]string, error) {
	return r.s.GetUncategorizedDescriptions(ctx)
}
