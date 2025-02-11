package db

import (
	"context"
	"hangry/domain/models"
	"hangry/repository"

	"gorm.io/gorm"
)

type cartRepositoryImpl struct {
	db *gorm.DB
}

// RemoveCartItem implements repository.CartRepository.
func (c *cartRepositoryImpl) RemoveCartItem(ctx context.Context, tx *gorm.DB, cartItemIds []uint) error {
	db := c.db
	if tx != nil {
		db = tx
	}

	return db.Where("id IN (?)", cartItemIds).Delete(&models.CartItem{}).Error
}

// AddToCart implements repository.CartRepository.
func (c *cartRepositoryImpl) AddToCart(ctx context.Context, tx *gorm.DB, input *models.CartItem) error {
	db := c.db
	if tx != nil {
		db = tx
	}

	err := db.Save(input).Error
	if err != nil {
		return err
	}

	return nil
}

// CheckItem implements repository.CartRepository.
func (c *cartRepositoryImpl) CheckItem(ctx context.Context, tx *gorm.DB, input repository.CheckItemInput) (models.CartItem, error) {
	db := c.db
	if tx != nil {
		db = tx
	}

	var cartItem models.CartItem

	if input.UserId != nil {
		db = db.Joins("inner join carts on cart_items.cart_id = carts.id").Where("carts.user_id = ?", input.UserId)
	}

	if input.CartId != nil {
		db = db.Where("cart_id = ?", input.CartId)
	}

	db = db.Where("product_id = ?", input.ProductId)

	if err := db.First(&cartItem).Error; err != nil && err != gorm.ErrRecordNotFound {
		return models.CartItem{}, err
	}

	return cartItem, nil
}

// CreateCart implements repository.CartRepository.
func (c *cartRepositoryImpl) CreateCart(ctx context.Context, tx *gorm.DB, userId uint) (models.Cart, error) {
	db := c.db
	if tx != nil {
		db = tx
	}

	cart := models.Cart{
		UserID: userId,
	}

	err := db.Create(&cart).Error
	if err != nil {
		return models.Cart{}, err
	}

	return cart, nil
}

// GetUserCart implements repository.CartRepository.
func (c *cartRepositoryImpl) GetUserCart(ctx context.Context, tx *gorm.DB, input repository.GetUserCartInput) (models.Cart, error) {
	db := c.db
	if tx != nil {
		db = tx
	}

	if len(input.Relations) > 0 {
		for _, relation := range input.Relations {
			db = db.Preload(relation)
		}
	}

	var cart models.Cart
	if err := db.Where("user_id = ?", input.UserId).First(&cart).Error; err != nil && err != gorm.ErrRecordNotFound {
		return models.Cart{}, err
	}

	return cart, nil
}

func NewCartRepository(db *gorm.DB) repository.CartRepository {
	return &cartRepositoryImpl{
		db: db,
	}
}
