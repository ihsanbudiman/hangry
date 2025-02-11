package repository

import (
	"context"
	"hangry/domain/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=./product_repository.go -destination=./mocks/mock_product_repository.go -package=mocks
type ProductRepository interface {
	Get(ctx context.Context, tx *gorm.DB, id uint) (models.Product, error)
}
