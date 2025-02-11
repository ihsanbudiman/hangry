package repository

import (
	"context"
	"hangry/domain/models"

	"gorm.io/gorm"
)

type GetPromoByUserCartInput struct {
	Cart        models.Cart
	PromoIds    []uint
	IsAvailable *bool
	Page        *int
	PerPage     *int
}

//go:generate mockgen -source=./promo_repository.go -destination=./mocks/mock_promo_repository.go -package=mocks
type PromoRepository interface {
	GetPromoByPromoID(ctx context.Context, tx *gorm.DB, promoID uint) (models.Promo, error)
	Save(ctx context.Context, tx *gorm.DB, promo *models.Promo) error
	SaveCities(ctx context.Context, tx *gorm.DB, promoID uint, cities []string) error
	GetPromoByUserCart(ctx context.Context, tx *gorm.DB, input GetPromoByUserCartInput) ([]models.Promo, int64, error)
}
