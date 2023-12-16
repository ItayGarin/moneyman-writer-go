package core

import (
	"context"
	"encoding/json"
	"fmt"
	"moneyman-writer-go/internal/model"
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
	data, err := s.downloader.Download(ctx, event.Bucket, event.Name)
	if err != nil {
		return err
	}

	txns, err := parseTransactions(data)
	if err != nil {
		return err
	}

	for _, txn := range txns {
		err = s.repo.Save(ctx, &txn)
		if err != nil {
			return fmt.Errorf("failed to save transaction: %s: %w", txn.Hash, err)
		}
	}
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
