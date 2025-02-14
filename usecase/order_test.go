package usecase

import (
	"context"
	"errors"
	"hangry/constants"
	"hangry/domain/dto"
	"hangry/domain/models"
	"hangry/repository"
	repo_mock "hangry/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
)

func Test_orderUsecase_CreateOrder(t *testing.T) {

	type args struct {
		ctx context.Context
		dto dto.OrderInput
	}
	cartData := models.Cart{
		ID: 1,
		CartItems: []models.CartItem{
			{
				ID:        1,
				CartID:    1,
				ProductID: 1,
				Quantity:  10,
				Cart:      &models.Cart{},
				Product: &models.Product{
					ID:    1,
					Price: 10000,
				},
			},
		},
	}

	buyProductId := uint(1)
	freeProductId := uint(2)
	promoBuyXGetY := models.Promo{
		ID:             1,
		Segmentation:   constants.PROMOSEGMENTATIONCITY,
		Type:           constants.PROMOTYPEBUYXGETY,
		BuyProductID:   &buyProductId,
		FreeProductID:  &freeProductId,
		BuyProductQty:  10,
		FreeProductQty: 1,
	}
	promoDiscount := models.Promo{
		ID:                2,
		Segmentation:      constants.PROMOSEGMENTATIONALL,
		Type:              constants.PROMOTYPEPERCENTAGE,
		MinOrderAmount:    10000,
		DiscountValue:     10,
		MaxDiscountAmount: 1000,
		PromoCities:       []models.PromoCity{},
		OrderPromos:       []models.OrderPromo{},
	}
	orderData := models.Order{
		UserID:      0,
		TotalAmount: 99000,
		OrderItems: []models.OrderItem{
			{
				ProductID:   buyProductId,
				Price:       10000,
				Quantity:    10,
				TotalAmount: 100000,
			},
		},
		OrderPromos: []models.OrderPromo{
			{
				PromoID:        1,
				DiscountAmount: 0,
				FreeProductID:  &freeProductId,
				FreeProductQty: 1,
			},
			{
				PromoID:        2,
				DiscountAmount: 1000,
			},
		},
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		mockRepo func(
			transaction *repo_mock.MockTransactionRepository,
			cart *repo_mock.MockCartRepository,
			user *repo_mock.MockUserRepository,
			promo *repo_mock.MockPromoRepository,
			order *repo_mock.MockOrderRepository,
		)
	}{
		{
			name: "err get user cart",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(models.Cart{}, errors.New("error"))
			},
		},
		{
			name: "cart not found",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(models.Cart{}, nil)
			},
		},
		{
			name: "err get promo by user cart",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{1},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(cartData, nil)

				isAvailable := true
				promo.EXPECT().GetPromoByUserCart(gomock.Any(), nil, repository.GetPromoByUserCartInput{
					Cart:        cartData,
					IsAvailable: &isAvailable,
					PromoIds:    []uint{1},
				}).Return([]models.Promo{}, int64(0), errors.New("error"))
			},
		},
		{
			name: "promo not found",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{1},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(cartData, nil)

				isAvailable := true
				promo.EXPECT().GetPromoByUserCart(gomock.Any(), nil, repository.GetPromoByUserCartInput{
					Cart:        cartData,
					IsAvailable: &isAvailable,
					PromoIds:    []uint{1},
				}).Return([]models.Promo{}, int64(0), nil)
			},
		},
		{
			name: "err make order",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{1},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(cartData, nil)

				isAvailable := true
				promo.EXPECT().GetPromoByUserCart(gomock.Any(), nil, repository.GetPromoByUserCartInput{
					Cart:        cartData,
					IsAvailable: &isAvailable,
					PromoIds:    []uint{1},
				}).Return([]models.Promo{promoBuyXGetY, promoDiscount}, int64(2), nil)

				order.EXPECT().MakeOrder(gomock.Any(), nil, &orderData).Return(errors.New("error"))
			},
		},
		{
			name: "err get user order count",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{1},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(cartData, nil)

				isAvailable := true
				promo.EXPECT().GetPromoByUserCart(gomock.Any(), nil, repository.GetPromoByUserCartInput{
					Cart:        cartData,
					IsAvailable: &isAvailable,
					PromoIds:    []uint{1},
				}).Return([]models.Promo{promoBuyXGetY, promoDiscount}, int64(2), nil)

				order.EXPECT().MakeOrder(gomock.Any(), nil, &orderData).Return(nil)
				order.EXPECT().GetUserOrderCount(gomock.Any(), nil, orderData.UserID).Return((0), errors.New("error"))
			},
		},
		{
			name: "err get user",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{1},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(cartData, nil)

				isAvailable := true
				promo.EXPECT().GetPromoByUserCart(gomock.Any(), nil, repository.GetPromoByUserCartInput{
					Cart:        cartData,
					IsAvailable: &isAvailable,
					PromoIds:    []uint{1},
				}).Return([]models.Promo{promoBuyXGetY, promoDiscount}, int64(2), nil)

				order.EXPECT().MakeOrder(gomock.Any(), nil, &orderData).Return(nil)
				order.EXPECT().GetUserOrderCount(gomock.Any(), nil, orderData.UserID).Return(4, nil)
				user.EXPECT().Get(gomock.Any(), nil, orderData.UserID).Return(&models.User{}, errors.New("error"))
			},
		},
		{
			name: "user not found",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{1},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(cartData, nil)

				isAvailable := true
				promo.EXPECT().GetPromoByUserCart(gomock.Any(), nil, repository.GetPromoByUserCartInput{
					Cart:        cartData,
					IsAvailable: &isAvailable,
					PromoIds:    []uint{1},
				}).Return([]models.Promo{promoBuyXGetY, promoDiscount}, int64(2), nil)

				order.EXPECT().MakeOrder(gomock.Any(), nil, &orderData).Return(nil)
				order.EXPECT().GetUserOrderCount(gomock.Any(), nil, orderData.UserID).Return(4, nil)
				user.EXPECT().Get(gomock.Any(), nil, orderData.UserID).Return(&models.User{}, nil)
			},
		},
		{
			name: "err save user",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{1},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(cartData, nil)

				isAvailable := true
				promo.EXPECT().GetPromoByUserCart(gomock.Any(), nil, repository.GetPromoByUserCartInput{
					Cart:        cartData,
					IsAvailable: &isAvailable,
					PromoIds:    []uint{1},
				}).Return([]models.Promo{promoBuyXGetY, promoDiscount}, int64(2), nil)

				order.EXPECT().MakeOrder(gomock.Any(), nil, &orderData).Return(nil)
				order.EXPECT().GetUserOrderCount(gomock.Any(), nil, orderData.UserID).Return(4, nil)
				userData := models.User{
					ID: 1,
				}
				user.EXPECT().Get(gomock.Any(), nil, orderData.UserID).Return(&userData, nil)
				userData.IsLoyal = true
				user.EXPECT().Save(gomock.Any(), nil, &userData).Return(errors.New("error"))
			},
		},
		{
			name: "err save promo",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{1},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(cartData, nil)

				isAvailable := true
				promoDatas := []models.Promo{promoBuyXGetY, promoDiscount}
				promo.EXPECT().GetPromoByUserCart(gomock.Any(), nil, repository.GetPromoByUserCartInput{
					Cart:        cartData,
					IsAvailable: &isAvailable,
					PromoIds:    []uint{1},
				}).Return(promoDatas, int64(2), nil)

				order.EXPECT().MakeOrder(gomock.Any(), nil, &orderData).Return(nil)
				order.EXPECT().GetUserOrderCount(gomock.Any(), nil, orderData.UserID).Return(4, nil)
				userData := models.User{
					ID: 1,
				}
				user.EXPECT().Get(gomock.Any(), nil, orderData.UserID).Return(&userData, nil)
				userData.IsLoyal = true
				user.EXPECT().Save(gomock.Any(), nil, &userData).Return(nil)
				for _, promoData := range promoDatas {
					promoData.CurrentUsageCount++
					promo.EXPECT().Save(gomock.Any(), nil, &promoData).Return(errors.New("error"))
					break
				}
			},
		},
		{
			name: "err remove cart item",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{1},
				},
			},
			wantErr: true,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(cartData, nil)

				isAvailable := true
				promoDatas := []models.Promo{promoBuyXGetY, promoDiscount}
				promo.EXPECT().GetPromoByUserCart(gomock.Any(), nil, repository.GetPromoByUserCartInput{
					Cart:        cartData,
					IsAvailable: &isAvailable,
					PromoIds:    []uint{1},
				}).Return(promoDatas, int64(2), nil)

				order.EXPECT().MakeOrder(gomock.Any(), nil, &orderData).Return(nil)
				order.EXPECT().GetUserOrderCount(gomock.Any(), nil, orderData.UserID).Return(4, nil)
				userData := models.User{
					ID: 1,
				}
				user.EXPECT().Get(gomock.Any(), nil, orderData.UserID).Return(&userData, nil)
				userData.IsLoyal = true
				user.EXPECT().Save(gomock.Any(), nil, &userData).Return(nil)
				for _, promoData := range promoDatas {
					promoData.CurrentUsageCount++
					promo.EXPECT().Save(gomock.Any(), nil, &promoData).Return(nil)
				}
				cartItemIds := []uint{}
				for _, cartItem := range cartData.CartItems {
					cartItemIds = append(cartItemIds, cartItem.ID)
				}
				cart.EXPECT().RemoveCartItem(gomock.Any(), nil, cartItemIds).Return(errors.New("error"))
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: dto.OrderInput{
					UserId:   1,
					PromoIds: []uint{1},
				},
			},
			wantErr: false,
			mockRepo: func(
				transaction *repo_mock.MockTransactionRepository,
				cart *repo_mock.MockCartRepository,
				user *repo_mock.MockUserRepository,
				promo *repo_mock.MockPromoRepository,
				order *repo_mock.MockOrderRepository,
			) {
				transaction.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})

				cart.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems", "CartItems.Product"},
				}).Return(cartData, nil)

				isAvailable := true
				promoDatas := []models.Promo{promoBuyXGetY, promoDiscount}
				promo.EXPECT().GetPromoByUserCart(gomock.Any(), nil, repository.GetPromoByUserCartInput{
					Cart:        cartData,
					IsAvailable: &isAvailable,
					PromoIds:    []uint{1},
				}).Return(promoDatas, int64(2), nil)

				order.EXPECT().MakeOrder(gomock.Any(), nil, &orderData).Return(nil)
				order.EXPECT().GetUserOrderCount(gomock.Any(), nil, orderData.UserID).Return(4, nil)
				userData := models.User{
					ID: 1,
				}
				user.EXPECT().Get(gomock.Any(), nil, orderData.UserID).Return(&userData, nil)
				userData.IsLoyal = true
				user.EXPECT().Save(gomock.Any(), nil, &userData).Return(nil)
				for _, promoData := range promoDatas {
					promoData.CurrentUsageCount++
					promo.EXPECT().Save(gomock.Any(), nil, &promoData).Return(nil)
				}
				cartItemIds := []uint{}
				for _, cartItem := range cartData.CartItems {
					cartItemIds = append(cartItemIds, cartItem.ID)
				}
				cart.EXPECT().RemoveCartItem(gomock.Any(), nil, cartItemIds).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			transactionRepo := repo_mock.NewMockTransactionRepository(ctrl)
			cartRepo := repo_mock.NewMockCartRepository(ctrl)
			promoRepo := repo_mock.NewMockPromoRepository(ctrl)
			orderRepo := repo_mock.NewMockOrderRepository(ctrl)
			userRepo := repo_mock.NewMockUserRepository(ctrl)

			tt.mockRepo(transactionRepo, cartRepo, userRepo, promoRepo, orderRepo)

			usecase := NewOrderUsecase(transactionRepo, orderRepo, userRepo, cartRepo, promoRepo)

			if err := usecase.CreateOrder(tt.args.ctx, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("orderUsecase.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
