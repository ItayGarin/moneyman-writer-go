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
	repo TransactionRepo
}

func NewService(repo TransactionRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) SaveNewTransactionsFromObjectFile(ctx context.Context, d ObjectDownloader, event *model.TransactionsFileUploadedEvent) error {
	x.Logger().Infow("downloading object", "bucket", event.Bucket, "name", event.Name)
	data, err := d.Download(ctx, event.Bucket, event.Name)
	if err != nil {
		return err
	}

	return s.SaveNewTransactionsFromBlob(ctx, data)
}

func (s *Service) SaveNewTransactionsFromBlob(ctx context.Context, blob []byte) error {
	x.Logger().Infow("parsing transactions", "file_size", len(blob))
	txns, err := parseTransactions(blob)
	if err != nil {
		return err
	}

	x.Logger().Infow("saving transactions", "count", len(txns))
	for _, txn := range txns {
		err = s.repo.Save(ctx, &txn)
		if err != nil {
			return fmt.Errorf("failed to save transaction: %s: %w", txn.Hash, err)
		}
	}

	x.Logger().Infow("saved all transactions", "count", len(txns))
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
