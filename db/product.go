package db

import (
	"context"
	"hangry/domain/models"
	"hangry/repository"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// Get implements repository.ProductRepository.
func (p *productRepository) Get(ctx context.Context, tx *gorm.DB, id uint) (models.Product, error) {
	db := p.db
	if tx != nil {
		db = tx
	}

	var product models.Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil && err != gorm.ErrRecordNotFound {
		return models.Product{}, err
	}

	return product, nil
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{
		db: db,
	}
}
