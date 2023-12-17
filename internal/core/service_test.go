package core_test

import (
	"context"
	"moneyman-writer-go/internal/adapter/db"
	cloud_storage "moneyman-writer-go/internal/adapter/google/cloud-storage"
	"moneyman-writer-go/internal/core"
	"moneyman-writer-go/internal/model"
	"moneyman-writer-go/internal/utils/test/containers"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
)

var url string

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer := containers.NewPostgresTestContainer()
	defer pgContainer.Terminate(ctx)

	err := db.RunMigrations(ctx, pgContainer.URL, "../adapter/db/migrations")
	if err != nil {
		panic(err)
	}

	url = pgContainer.URL

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestService_saveNewTransactionsFromObjectFile(t *testing.T) {
	ctx := context.Background()

	bucket := os.Getenv("GCS_TEST_BUCKET")
	path := os.Getenv("GCS_TEST_FILE")
	creds := os.Getenv("GCS_TEST_CREDS")
	require.NotEmptyf(t, bucket, "GCS_TEST_BUCKET is empty")
	require.NotEmptyf(t, path, "GCS_TEST_FILE is empty")
	require.NotEmptyf(t, creds, "GCS_TEST_CREDS is empty")

	client, err := db.NewClient(ctx, url)
	require.NoError(t, err)
	repo := db.NewPostgresTransactionRepo(client)

	downloader, err := cloud_storage.NewGcsDownloader(ctx, option.WithCredentialsJSON([]byte(creds)))
	require.NoError(t, err)

	svc := core.NewService(repo)
	err = svc.SaveNewTransactionsFromObjectFile(ctx, downloader, &model.TransactionsFileUploadedEvent{
		TimeCreated: time.Now(),
		Bucket:      bucket,
		Name:        path,
	})
	require.NoError(t, err)
}
