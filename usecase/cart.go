package usecase

import (
	"context"
	"hangry/domain/dto"
	"hangry/domain/models"
	"hangry/repository"
	"hangry/utils"
	"net/http"

	"gorm.io/gorm"
)

//go:generate mockgen -source=./cart.go -destination=./mocks/mock_cart.go -package=mocks
type CartUsecase interface {
	AddToCart(ctx context.Context, dto dto.AddToCartInput) error
	RemoveFromCart(ctx context.Context, dto dto.RemoveFromCartInput) error
}

type cartUsecase struct {
	transaction       repository.TransactionRepository
	cartRepository    repository.CartRepository
	productRepository repository.ProductRepository
}

// RemoveFromCart implements CartUsecase.
func (c *cartUsecase) RemoveFromCart(ctx context.Context, dto dto.RemoveFromCartInput) error {
	// check item exist in cart
	item, err := c.cartRepository.CheckItem(ctx, nil, repository.CheckItemInput{
		UserId:    &dto.UserId,
		ProductId: dto.ProductId,
	})
	if err != nil {
		return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
	}

	if item.ID == 0 {
		return utils.NewCustomError("item not found", nil, http.StatusNotFound)
	}

	err = c.cartRepository.RemoveCartItem(ctx, nil, []uint{item.ID})
	if err != nil {
		return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
	}

	return nil
}

// AddToCart implements CartUsecase.
func (c *cartUsecase) AddToCart(ctx context.Context, dto dto.AddToCartInput) error {
	return c.transaction.Execute(ctx, func(tx *gorm.DB) error {
		// check if product exist
		product, err := c.productRepository.Get(ctx, tx, dto.ProductId)
		if err != nil {
			return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
		}

		if product.ID == 0 {
			return utils.NewCustomError("product not found", nil, http.StatusNotFound)
		}

		// getting cart id if exist
		cart, err := c.cartRepository.GetUserCart(ctx, tx, repository.GetUserCartInput{
			UserId: dto.UserId,
		})
		if err != nil {
			return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
		}

		// if cart not exist, create one
		if cart.ID == 0 {
			cart, err = c.cartRepository.CreateCart(ctx, tx, dto.UserId)
			if err != nil {
				return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
			}
		}

		// check if item exist in cart
		item, err := c.cartRepository.CheckItem(ctx, tx, repository.CheckItemInput{
			CartId:    &cart.ID,
			ProductId: dto.ProductId,
		})
		if err != nil {
			return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
		}

		if item.ID == 0 {
			item = models.CartItem{
				CartID:    cart.ID,
				ProductID: dto.ProductId,
				Quantity:  dto.Quantity,
			}
		} else {
			item.Quantity += dto.Quantity
		}

		err = c.cartRepository.AddToCart(ctx, tx, &item)
		if err != nil {
			return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
		}

		return nil
	})
}

func NewCartUsecase(
	transaction repository.TransactionRepository,
	cartRepository repository.CartRepository,
	productRepository repository.ProductRepository,
) CartUsecase {
	return &cartUsecase{
		transaction:       transaction,
		cartRepository:    cartRepository,
		productRepository: productRepository,
	}
}
