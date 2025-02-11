package repository

import (
	"context"
	"hangry/domain/models"

	"gorm.io/gorm"
)

type CheckItemInput struct {
	CartId    *uint
	ProductId uint
	UserId    *uint
}

type GetUserCartInput struct {
	UserId    uint
	Relations []string
}

//go:generate mockgen -source=./cart_repository.go -destination=./mocks/mock_cart_repository.go -package=mocks
type CartRepository interface {
	GetUserCart(ctx context.Context, tx *gorm.DB, input GetUserCartInput) (models.Cart, error)
	CreateCart(ctx context.Context, tx *gorm.DB, userId uint) (models.Cart, error)
	CheckItem(ctx context.Context, tx *gorm.DB, input CheckItemInput) (models.CartItem, error)
	AddToCart(ctx context.Context, tx *gorm.DB, input *models.CartItem) error
	RemoveCartItem(ctx context.Context, tx *gorm.DB, cartItemId []uint) error
}
