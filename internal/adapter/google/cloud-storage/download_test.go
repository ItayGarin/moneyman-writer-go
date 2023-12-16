package cloud_storage_test

import (
	"context"
	"fmt"
	cloud_storage "moneyman-writer-go/internal/adapter/google/cloud-storage"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func TestGcsDownloader_Download(t *testing.T) {
	ctx := context.Background()
	bucket := os.Getenv("GCS_TEST_BUCKET")
	path := os.Getenv("GCS_TEST_FILE")
	creds := os.Getenv("GCS_TEST_CREDS")

	fmt.Println(bucket)
	fmt.Println(path)
	d, err := cloud_storage.NewGcsDownloader(ctx, option.WithCredentialsJSON([]byte(creds)))
	assert.NoError(t, err)

	data, err := d.Download(ctx, bucket, path)
	assert.NoError(t, err)
	assert.Equal(t, len(data), 5526)
}
