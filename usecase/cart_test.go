package usecase

import (
	"context"
	"errors"
	"hangry/domain/dto"
	"hangry/domain/models"
	"hangry/repository"
	repo_mock "hangry/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
)

// TestAddCart uses table-driven tests to verify CartUsecase.AddCart behavior.
func TestAddCart(t *testing.T) {
	type args struct {
		ctx context.Context
		dto dto.AddToCartInput
	}
	tests := []struct {
		name                string
		args                args
		wantErr             bool
		transactionRepoMock func(*repo_mock.MockTransactionRepository)
		cartRepoMock        func(*repo_mock.MockCartRepository)
		productRepoMock     func(*repo_mock.MockProductRepository)
	}{
		{
			name: "failed to product",
			args: args{
				ctx: context.Background(),
				dto: dto.AddToCartInput{
					ProductId: 1,
					UserId:    1,
					Quantity:  1,
				},
			},
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), gomock.Any(), uint(1)).Return(models.Product{}, errors.New("something went wrong"))
			},
		},
		{
			name: "product not found",
			args: args{
				ctx: context.Background(),
				dto: dto.AddToCartInput{
					ProductId: 1,
					UserId:    1,
					Quantity:  1,
				},
			},
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						// Simulate transaction execution by calling the callback
						return fn(nil)
					})
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), gomock.Any(), uint(1)).Return(models.Product{}, nil)
			},
		},
		{
			name: "failed to get user cart",
			args: args{
				ctx: context.Background(),
				dto: dto.AddToCartInput{
					ProductId: 1,
					UserId:    1,
					Quantity:  1,
				},
			},
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						// Simulate transaction execution by calling the callback
						return fn(nil)
					})
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
				r.EXPECT().GetUserCart(gomock.Any(), gomock.Any(), repository.GetUserCartInput{UserId: 1}).Return(models.Cart{}, errors.New("something went wrong"))
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), gomock.Any(), uint(1)).Return(models.Product{
					ID: 1,
				}, nil)
			},
		},
		{
			name: "cart not found",
			args: args{
				ctx: context.Background(),
				dto: dto.AddToCartInput{
					ProductId: 1,
					UserId:    1,
					Quantity:  1,
				},
			},
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						// Simulate transaction execution by calling the callback
						return fn(nil)
					})
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
				r.EXPECT().GetUserCart(gomock.Any(), gomock.Any(), repository.GetUserCartInput{UserId: 1}).Return(models.Cart{}, nil)
				r.EXPECT().CreateCart(gomock.Any(), gomock.Any(), uint(1)).Return(models.Cart{}, errors.New("something went wrong"))
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), gomock.Any(), uint(1)).Return(models.Product{
					ID: 1,
				}, nil)
			},
		},
		{
			name: "err check item",
			args: args{
				ctx: context.Background(),
				dto: dto.AddToCartInput{
					ProductId: 1,
					UserId:    1,
					Quantity:  1,
				},
			},
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						// Simulate transaction execution by calling the callback
						return fn(nil)
					})
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
				r.EXPECT().GetUserCart(gomock.Any(), gomock.Any(), repository.GetUserCartInput{UserId: 1}).Return(models.Cart{}, nil)
				cart := models.Cart{}
				r.EXPECT().CreateCart(gomock.Any(), gomock.Any(), uint(1)).Return(cart, nil)
				r.EXPECT().CheckItem(gomock.Any(), gomock.Any(), repository.CheckItemInput{
					CartId:    &cart.ID,
					ProductId: 1,
				}).Return(models.CartItem{}, errors.New("something went wrong"))
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), gomock.Any(), uint(1)).Return(models.Product{
					ID: 1,
				}, nil)
			},
		},
		{
			name: "err add to cart",
			args: args{
				ctx: context.Background(),
				dto: dto.AddToCartInput{
					ProductId: 1,
					UserId:    1,
					Quantity:  1,
				},
			},
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						// Simulate transaction execution by calling the callback
						return fn(nil)
					})
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
				r.EXPECT().GetUserCart(gomock.Any(), gomock.Any(), repository.GetUserCartInput{UserId: 1}).Return(models.Cart{}, nil)
				cart := models.Cart{
					ID: 1,
				}
				cartItem := models.CartItem{}
				r.EXPECT().CreateCart(gomock.Any(), gomock.Any(), uint(1)).Return(cart, nil)
				r.EXPECT().CheckItem(gomock.Any(), gomock.Any(), repository.CheckItemInput{
					CartId:    &cart.ID,
					ProductId: 1,
				}).Return(cartItem, nil)
				cartItem.CartID = cart.ID
				cartItem.ProductID = 1
				cartItem.Quantity = 1
				r.EXPECT().AddToCart(gomock.Any(), gomock.Any(), &cartItem).Return(errors.New("something went wrong"))
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), gomock.Any(), uint(1)).Return(models.Product{
					ID: 1,
				}, nil)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: dto.AddToCartInput{
					ProductId: 1,
					UserId:    1,
					Quantity:  1,
				},
			},
			wantErr: false,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						// Simulate transaction execution by calling the callback
						return fn(nil)
					})
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
				r.EXPECT().GetUserCart(gomock.Any(), gomock.Any(), repository.GetUserCartInput{UserId: 1}).Return(models.Cart{}, nil)
				cart := models.Cart{
					ID: 1,
				}
				cartItem := models.CartItem{}
				r.EXPECT().CreateCart(gomock.Any(), gomock.Any(), uint(1)).Return(cart, nil)
				r.EXPECT().CheckItem(gomock.Any(), gomock.Any(), repository.CheckItemInput{
					CartId:    &cart.ID,
					ProductId: 1,
				}).Return(cartItem, nil)
				cartItem.CartID = cart.ID
				cartItem.ProductID = 1
				cartItem.Quantity = 1
				r.EXPECT().AddToCart(gomock.Any(), gomock.Any(), &cartItem).Return(nil)
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), gomock.Any(), uint(1)).Return(models.Product{
					ID: 1,
				}, nil)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: dto.AddToCartInput{
					ProductId: 1,
					UserId:    1,
					Quantity:  1,
				},
			},
			wantErr: false,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						// Simulate transaction execution by calling the callback
						return fn(nil)
					})
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
				r.EXPECT().GetUserCart(gomock.Any(), gomock.Any(), repository.GetUserCartInput{UserId: 1}).Return(models.Cart{}, nil)
				cart := models.Cart{
					ID: 1,
				}
				cartItem := models.CartItem{
					ID:        1,
					Quantity:  1,
					ProductID: 1,
				}
				r.EXPECT().CreateCart(gomock.Any(), gomock.Any(), uint(1)).Return(cart, nil)
				r.EXPECT().CheckItem(gomock.Any(), gomock.Any(), repository.CheckItemInput{
					CartId:    &cart.ID,
					ProductId: 1,
				}).Return(cartItem, nil)
				cartItem.Quantity += 1
				r.EXPECT().AddToCart(gomock.Any(), gomock.Any(), &cartItem).Return(nil)
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), gomock.Any(), uint(1)).Return(models.Product{
					ID: 1,
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			transactionRepo := repo_mock.NewMockTransactionRepository(ctrl)
			tt.transactionRepoMock(transactionRepo)

			cartRepo := repo_mock.NewMockCartRepository(ctrl)
			tt.cartRepoMock(cartRepo)

			productRepo := repo_mock.NewMockProductRepository(ctrl)
			tt.productRepoMock(productRepo)

			usecase := NewCartUsecase(transactionRepo, cartRepo, productRepo)

			if err := usecase.AddToCart(tt.args.ctx, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("cartUsecase.AddToCart() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func Test_cartUsecase_RemoveFromCart(t *testing.T) {

	type args struct {
		ctx context.Context
		dto dto.RemoveFromCartInput
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		cartRepoMock func(*repo_mock.MockCartRepository)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: dto.RemoveFromCartInput{
					ProductId: 1,
					UserId:    1,
				},
			},
			wantErr: false,
			cartRepoMock: func(mockCartRepo *repo_mock.MockCartRepository) {
				mockCartRepo.EXPECT().CheckItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CartItem{
					ID: 1,
				}, nil)
				mockCartRepo.EXPECT().RemoveCartItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		// failed remove cart item
		{
			name: "failed remove cart item",
			args: args{
				ctx: context.Background(),
				dto: dto.RemoveFromCartInput{
					ProductId: 1,
					UserId:    1,
				},
			},
			wantErr: true,
			cartRepoMock: func(mockCartRepo *repo_mock.MockCartRepository) {
				mockCartRepo.EXPECT().CheckItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CartItem{
					ID: 1,
				}, nil)
				mockCartRepo.EXPECT().RemoveCartItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed remove cart item"))
			},
		},
		// item not found
		{
			name: "item not found",
			args: args{
				ctx: context.Background(),
				dto: dto.RemoveFromCartInput{
					ProductId: 1,
					UserId:    1,
				},
			},
			wantErr: true,
			cartRepoMock: func(mockCartRepo *repo_mock.MockCartRepository) {
				mockCartRepo.EXPECT().CheckItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CartItem{}, nil)
			},
		},
		// failed check item
		{
			name: "failed check item",
			args: args{
				ctx: context.Background(),
				dto: dto.RemoveFromCartInput{
					ProductId: 1,
					UserId:    1,
				},
			},
			wantErr: true,
			cartRepoMock: func(mockCartRepo *repo_mock.MockCartRepository) {
				mockCartRepo.EXPECT().CheckItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CartItem{}, errors.New("failed check item"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cartRepo := repo_mock.NewMockCartRepository(ctrl)
			tt.cartRepoMock(cartRepo)

			usecase := NewCartUsecase(nil, cartRepo, nil)

			if err := usecase.RemoveFromCart(tt.args.ctx, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("cartUsecase.RemoveFromCart() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
