package cloud_storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GcsDownloader struct {
	client *storage.Client
}

func NewGcsDownloader(ctx context.Context, opts ...option.ClientOption) (*GcsDownloader, error) {
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &GcsDownloader{
		client: client,
	}, nil
}

func (d *GcsDownloader) Download(ctx context.Context, bucket string, path string) ([]byte, error) {
	b := d.client.Bucket(bucket)
	obj := b.Object(path)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create object reader: %w", err)
	}
	defer reader.Close()

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object attributes: %w", err)
	}

	contents := make([]byte, attrs.Size)
	_, err = reader.Read(contents)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return contents, nil
}

func (d *GcsDownloader) Close() error {
	return d.client.Close()
}
