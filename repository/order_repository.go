package repository

import (
	"context"
	"hangry/domain/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=./order_repository.go -destination=./mocks/mock_order_repository.go -package=mocks
type OrderRepository interface {
	MakeOrder(ctx context.Context, tx *gorm.DB, order *models.Order) error
	GetUserOrderCount(ctx context.Context, tx *gorm.DB, userID uint) (int, error)
}
