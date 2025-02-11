package usecase

import (
	"context"
	"hangry/constants"
	"hangry/domain/dto"
	"hangry/domain/models"
	"hangry/repository"
	"hangry/utils"
	"net/http"

	"gorm.io/gorm"
)

//go:generate mockgen -source=./order.go -destination=./mocks/mock_order.go -package=mocks
type OrderUsecase interface {
	CreateOrder(ctx context.Context, dto dto.OrderInput) error
}

type orderUsecase struct {
	transactionRepository repository.TransactionRepository
	orderRepository       repository.OrderRepository
	userRepository        repository.UserRepository
	cartRepository        repository.CartRepository
	promoRepository       repository.PromoRepository
}

// CreateOrder implements OrderUsecase.
func (o *orderUsecase) CreateOrder(ctx context.Context, dto dto.OrderInput) error {
	return o.transactionRepository.Execute(ctx, func(tx *gorm.DB) error {
		// get user cart
		cart, err := o.cartRepository.GetUserCart(ctx, tx, repository.GetUserCartInput{
			UserId:    dto.UserId,
			Relations: []string{"CartItems", "CartItems.Product"},
		})
		if err != nil && err != gorm.ErrRecordNotFound {
			return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
		}

		if cart.ID == 0 || len(cart.CartItems) == 0 {
			return utils.NewCustomError("cart not found", nil, http.StatusNotFound)
		}

		// get promo
		var promos []models.Promo
		if dto.PromoIds != nil && len(dto.PromoIds) > 0 {
			// check if promo still valid
			isAvailable := true
			promos, _, err = o.promoRepository.GetPromoByUserCart(ctx, tx, repository.GetPromoByUserCartInput{
				Cart:        cart,
				IsAvailable: &isAvailable,
				PromoIds:    dto.PromoIds,
			})
			if err != nil && err != gorm.ErrRecordNotFound {
				return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
			}
			if len(promos) == 0 {
				return utils.NewCustomError("promo not found", nil, http.StatusNotFound)
			}
		}

		// create order
		order := models.Order{
			UserID:      cart.UserID,
			TotalAmount: 0,
		}

		cartItemIds := make([]uint, 0)

		for _, cartItem := range cart.CartItems {
			order.TotalAmount += cartItem.Product.Price * float64(cartItem.Quantity)
			orderItems := models.OrderItem{
				ProductID:   cartItem.ProductID,
				Quantity:    cartItem.Quantity,
				Price:       cartItem.Product.Price,
				TotalAmount: cartItem.Product.Price * float64(cartItem.Quantity),
			}
			order.OrderItems = append(order.OrderItems, orderItems)
			cartItemIds = append(cartItemIds, cartItem.ID)
		}

		for i, promo := range promos {
			orderPromo := models.OrderPromo{
				PromoID:        promo.ID,
				DiscountAmount: 0,
			}
			if promo.Type == constants.PROMOTYPEBUYXGETY {
				productId := promo.FreeProductID
				orderPromo.FreeProductID = productId
				orderPromo.FreeProductQty = promo.FreeProductQty
			} else if promo.Type == constants.PROMOTYPEPERCENTAGE {
				total := order.TotalAmount
				discount := promo.DiscountValue
				orderPromo.DiscountAmount = total * discount / 100
				if orderPromo.DiscountAmount > promo.MaxDiscountAmount {
					orderPromo.DiscountAmount = promo.MaxDiscountAmount
				}
				order.TotalAmount -= orderPromo.DiscountAmount
			}

			order.OrderPromos = append(order.OrderPromos, orderPromo)
			promos[i].CurrentUsageCount++
		}

		if err := o.orderRepository.MakeOrder(ctx, tx, &order); err != nil {
			return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
		}

		// check if user has ordered more than 3 times
		totalOrder, err := o.orderRepository.GetUserOrderCount(ctx, tx, order.UserID)
		if err != nil {
			return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
		}

		if totalOrder > 3 {
			user, err := o.userRepository.Get(ctx, tx, order.UserID)
			if err != nil && err != gorm.ErrRecordNotFound {
				return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
			}
			if user.ID == 0 {
				return utils.NewCustomError("user not found", nil, http.StatusNotFound)
			}

			user.IsLoyal = true

			if err := o.userRepository.Save(ctx, tx, user); err != nil {
				return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
			}
		}

		// update promo
		for _, promo := range promos {
			if err := o.promoRepository.Save(ctx, tx, &promo); err != nil {
				return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
			}
		}

		// delete cart items
		if err := o.cartRepository.RemoveCartItem(ctx, tx, cartItemIds); err != nil {
			return utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
		}

		return nil
	})

}

func NewOrderUsecase(
	transactionRepository repository.TransactionRepository,
	orderRepository repository.OrderRepository,
	userRepository repository.UserRepository,
	cartRepository repository.CartRepository,
	promoRepository repository.PromoRepository,
) OrderUsecase {
	return &orderUsecase{
		transactionRepository: transactionRepository,
		orderRepository:       orderRepository,
		userRepository:        userRepository,
		cartRepository:        cartRepository,
		promoRepository:       promoRepository,
	}
}
