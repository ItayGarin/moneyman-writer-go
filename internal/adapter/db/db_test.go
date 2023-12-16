package db

import (
	"context"
	"moneyman-writer-go/internal/adapter/db/sqlc/sql"
	"moneyman-writer-go/internal/utils/test/containers"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

var url string

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer := containers.NewPostgresTestContainer()
	defer pgContainer.Terminate(ctx)

	err := RunMigrations(ctx, pgContainer.URL, "./migrations")
	if err != nil {
		panic(err)
	}

	url = pgContainer.URL

	exitCode := m.Run()
	os.Exit(exitCode)
}

func makeMockTransaction() *sql.InsertTransactionParams {
	mockDate := pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}
	return &sql.InsertTransactionParams{
		Identifier:       "TX123456",
		Type:             "Credit",
		Status:           "Completed",
		Date:             mockDate,
		ProcessedDate:    mockDate,
		OriginalAmount:   100.50,
		OriginalCurrency: "USD",
		ChargedAmount:    100.50,
		ChargedCurrency:  "USD",
		Description:      "Mock transaction description",
		Memo:             "Mock transaction memo",
		Category:         "Groceries",
		Account:          "1234-5678-9012",
		CompanyID:        "COMP12345",
		Hash:             "abcde12345hash",
	}
}

func TestDB_SingleWrite_NoConflict_shouldSucceed(t *testing.T) {
	ctx := context.Background()
	c, err := NewClient(ctx, url)
	require.NoError(t, err)

	mockTransaction := makeMockTransaction()
	err = c.InsertTransaction(ctx, *mockTransaction)
	require.NoError(t, err)

	txn, err := c.GetTransactionByHash(ctx, mockTransaction.Hash)
	require.NoError(t, err)
	require.Equal(t, mockTransaction.Identifier, txn.Identifier)
	require.Equal(t, mockTransaction.Type, txn.Type)
	require.Equal(t, mockTransaction.Status, txn.Status)
	require.Equal(t, mockTransaction.Hash, txn.Hash)
}

func TestDB_SingleWrite_HasConflict_shouldSucceed(t *testing.T) {
	ctx := context.Background()
	c, err := NewClient(ctx, url)
	require.NoError(t, err)

	for i := 0; i < 2; i++ {
		mockTransaction := makeMockTransaction()
		err = c.InsertTransaction(ctx, *mockTransaction)
		require.NoError(t, err)

		txn, err := c.GetTransactionByHash(ctx, mockTransaction.Hash)
		require.NoError(t, err)
		require.Equal(t, mockTransaction.Identifier, txn.Identifier)
	}
}
