package core

import (
	"context"
	"encoding/json"
	"fmt"
	"moneyman-writer-go/internal/model"
	x "moneyman-writer-go/internal/utils/logger"
)

type TransactionRepo interface {
	Save(ctx context.Context, txn *model.Transaction) error
}

type ObjectDownloader interface {
	Download(ctx context.Context, bucket, path string) ([]byte, error)
}

type Service struct {
	repo       TransactionRepo
	downloader ObjectDownloader
}

func NewService(repo TransactionRepo, downloader ObjectDownloader) *Service {
	return &Service{
		repo:       repo,
		downloader: downloader,
	}
}

func (s *Service) SaveNewTransactionsFromObjectFile(ctx context.Context, event *model.TransactionsFileUploadedEvent) error {
	x.Logger().Infow("downloading object", "bucket", event.Bucket, "name", event.Name)
	data, err := s.downloader.Download(ctx, event.Bucket, event.Name)
	if err != nil {
		return err
	}

	x.Logger().Infow("parsing transactions", "bucket", event.Bucket, "name", event.Name, "file_size", len(data))
	txns, err := parseTransactions(data)
	if err != nil {
		return err
	}

	x.Logger().Infow("saving transactions", "bucket", event.Bucket, "name", event.Name, "count", len(txns))
	for _, txn := range txns {
		err = s.repo.Save(ctx, &txn)
		if err != nil {
			return fmt.Errorf("failed to save transaction: %s: %w", txn.Hash, err)
		}
	}

	x.Logger().Infow("saved all transactions", "bucket", event.Bucket, "name", event.Name, "count", len(txns))
	return nil
}

func parseTransactions(data []byte) ([]model.Transaction, error) {
	out := []model.Transaction{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
