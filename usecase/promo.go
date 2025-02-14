package usecase

import (
	"context"
	"hangry/constants"
	"hangry/domain/dto"
	"hangry/domain/models"
	"hangry/generated"
	"hangry/repository"
	"hangry/utils"
	"net/http"
	"time"

	"gorm.io/gorm"
)

//go:generate mockgen -source=./promo.go -destination=./mocks/mock_promo.go -package=mocks
type PromoUsecase interface {
	CreatePromo(ctx context.Context, dto dto.CreatePromoInput) (uint, error)
	ExtendPromo(ctx context.Context, dto dto.ExtendPromoInput) error
	GetPromo(ctx context.Context, dto dto.GetPromoInput) ([]models.Promo, int64, error)
}

type promoUsecase struct {
	promoRepository       repository.PromoRepository
	transactionRepository repository.TransactionRepository
	cartRepository        repository.CartRepository
	productRepository     repository.ProductRepository
}

// GetPromo implements PromoUsecase.
func (p *promoUsecase) GetPromo(ctx context.Context, dto dto.GetPromoInput) ([]models.Promo, int64, error) {
	// get user cart
	cart, err := p.cartRepository.GetUserCart(ctx, nil, repository.GetUserCartInput{
		UserId:    uint(dto.UserId),
		Relations: []string{"CartItems"},
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	if cart.ID == 0 {
		return []models.Promo{}, 0, nil
	}

	isAvailable := true
	promos, total, err := p.promoRepository.GetPromoByUserCart(ctx, nil, repository.GetPromoByUserCartInput{
		Cart:        cart,
		IsAvailable: &isAvailable,
		Page:        &dto.Page,
		PerPage:     &dto.PerPage,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	return promos, total, nil
}

// ExtendPromo implements PromoUsecase.
func (p *promoUsecase) ExtendPromo(ctx context.Context, dto dto.ExtendPromoInput) error {
	return p.transactionRepository.Execute(ctx, func(tx *gorm.DB) error {
		promo, err := p.promoRepository.GetPromoByPromoID(ctx, tx, dto.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if promo.ID == 0 {
			return utils.NewCustomError("promo not found", nil, http.StatusNotFound)
		}

		// check if dto.EndDate is less than promo.EndDate
		if dto.EndDate.Before(promo.EndDate) || dto.EndDate.Before(time.Now()) {
			return utils.NewCustomError("new end date must be greater than current end date", nil, http.StatusBadRequest)
		}

		promo.EndDate = dto.EndDate
		if !dto.StartDate.IsZero() {
			promo.StartDate = dto.StartDate
		}
		return p.promoRepository.Save(ctx, tx, &promo)
	})

}

// CreatePromo implements PromoUsecase.
func (p *promoUsecase) CreatePromo(ctx context.Context, dto dto.CreatePromoInput) (uint, error) {

	var promoId uint

	err := p.transactionRepository.Execute(ctx, func(tx *gorm.DB) error {
		if dto.Type == string(constants.PROMOTYPEBUYXGETY) {
			// check buy product id
			product, err := p.productRepository.Get(ctx, tx, uint(*dto.BuyProductId))
			if err != nil {
				return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
			}
			if product.ID == 0 {
				return utils.NewCustomError("buy product not found", nil, http.StatusNotFound)
			}

			// check free product id
			product, err = p.productRepository.Get(ctx, tx, uint(*dto.FreeProductId))
			if err != nil {
				return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
			}
			if product.ID == 0 {
				return utils.NewCustomError("free product not found", nil, http.StatusNotFound)
			}
		}

		promo := dto.CreatePromoModel()

		if err := p.promoRepository.Save(ctx, tx, &promo); err != nil {
			return err
		}

		if promo.Segmentation == string(generated.PromoSegmentationCITY) {
			if err := p.promoRepository.SaveCities(ctx, tx, promo.ID, dto.Cities); err != nil {
				return err
			}
		}

		promoId = promo.ID

		return nil
	})

	if err != nil {
		return 0, err
	}

	return promoId, nil
}

func NewPromoUsecase(
	promoRepository repository.PromoRepository,
	transactionRepository repository.TransactionRepository,
	cartRepository repository.CartRepository,
	productRepository repository.ProductRepository,
) PromoUsecase {
	return &promoUsecase{
		promoRepository:       promoRepository,
		transactionRepository: transactionRepository,
		cartRepository:        cartRepository,
		productRepository:     productRepository,
	}
}
