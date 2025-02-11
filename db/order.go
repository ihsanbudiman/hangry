package db

import (
	"context"
	"hangry/domain/models"
	"hangry/repository"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

// GetUserOrderCount implements repository.OrderRepository.
func (o *orderRepository) GetUserOrderCount(ctx context.Context, tx *gorm.DB, userID uint) (int, error) {
	db := tx
	if db == nil {
		db = o.db.WithContext(ctx)
	}

	var count int64
	if err := db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// MakeOrder implements repository.OrderRepository.
func (o *orderRepository) MakeOrder(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	db := tx
	if db == nil {
		db = o.db.WithContext(ctx)
	}

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(order).Error; err != nil {
		return err
	}
	return nil
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepository{db: db}
}
