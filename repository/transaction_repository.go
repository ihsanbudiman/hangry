package repository

import (
	"context"

	"gorm.io/gorm"
)

//go:generate mockgen -source=./transaction_repository.go -destination=./mocks/mock_transaction_repository.go -package=mocks
type TransactionRepository interface {
	Execute(ctx context.Context, fn func(tx *gorm.DB) error) error
}
