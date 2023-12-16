package db

import (
	"context"
	"fmt"
	"log"
	"moneyman-writer-go/internal/adapter/db/sqlc/sql"
	"os"
	"testing"
	"time"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/docker/docker/api/types/container"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var dsn string
var url string

func TestMain(m *testing.M) {
	// Context for the container
	ctx := context.Background()

	// Define the PostgreSQL container request
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "test_db",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
		HostConfigModifier: func(config *container.HostConfig) {
			config.AutoRemove = true
		},
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container: %s", err)
	}
	defer postgresContainer.Terminate(ctx)

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Failed to get mapped port: %s", err)
	}

	// Build the database connection string
	dbHost, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container IP: %s", err)
	}
	dbPort := mappedPort.Port()
	dsn = fmt.Sprintf("host=%s port=%s user=user password=password dbname=test_db sslmode=disable", dbHost, dbPort)
	url = fmt.Sprintf("postgres://user:password@%s:%s/%s?sslmode=%s", dbHost, dbPort, "test_db", "disable")

	runMigrations()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func runMigrations() {
	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS("./migrations"),
		),
	)
	if err != nil {
		log.Fatalf("failed to load working directory: %v", err)
	}
	defer workdir.Close()

	fmt.Printf("Applying migrations to %s\n", url)
	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}
	status, err := client.MigrateStatus(context.Background(), &atlasexec.MigrateStatusParams{
		URL: url,
	})
	if err != nil {
		log.Fatalf("failed to get migration status: %v", err)
	}
	fmt.Printf("Status: %v\n", status)

	res, err := client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: url,
	})
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	fmt.Printf("Applied %d migrations\n", len(res.Applied))
}

func TestDB_Write(t *testing.T) {
	ctx := context.Background()
	c, err := NewClient(ctx, dsn)
	require.NoError(t, err)

	mockDate := pgtype.Timestamptz{Time: time.Now()}
	mockTransaction := sql.InsertTransactionParams{
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

	err = c.InsertTransaction(ctx, mockTransaction)
	require.NoError(t, err)

	txn, err := c.GetTransactionByHash(ctx, mockTransaction.Hash)
	require.NoError(t, err)
	require.Equal(t, mockTransaction, txn)
}
