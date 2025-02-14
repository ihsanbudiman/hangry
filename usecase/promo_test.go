package usecase

import (
	"context"
	"errors"
	"fmt"
	"hangry/constants"
	"hangry/domain/dto"
	"hangry/domain/models"
	"hangry/repository"
	repo_mock "hangry/repository/mocks"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
)

func Test_promoUsecase_GetPromo(t *testing.T) {

	type args struct {
		ctx context.Context
		dto dto.GetPromoInput
	}
	tests := []struct {
		name                string
		args                args
		want                []models.Promo
		want1               int64
		wantErr             bool
		transactionRepoMock func(*repo_mock.MockTransactionRepository)
		promoRepoMock       func(*repo_mock.MockPromoRepository)
		cartRepoMock        func(*repo_mock.MockCartRepository)
		productRepoMock     func(*repo_mock.MockProductRepository)
	}{
		// TODO: Add test cases.
		{
			name: "failed to get user cart",
			args: args{
				ctx: context.Background(),
				dto: dto.GetPromoInput{
					UserId:  1,
					Page:    1,
					PerPage: 10,
				},
			},
			want:    nil,
			want1:   0,
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
				r.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    uint(1),
					Relations: []string{"CartItems"},
				}).Return(models.Cart{}, errors.New("error"))
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
			},
		},
		{
			name: "cart not found",
			args: args{
				ctx: context.Background(),
				dto: dto.GetPromoInput{
					UserId:  1,
					Page:    1,
					PerPage: 10,
				},
			},
			want:    []models.Promo{},
			want1:   0,
			wantErr: false,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
				r.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    uint(1),
					Relations: []string{"CartItems"},
				}).Return(models.Cart{}, nil)
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
			},
		},
		{
			name: "err get promo by user cart",
			args: args{
				ctx: context.Background(),
				dto: dto.GetPromoInput{
					UserId:  1,
					Page:    1,
					PerPage: 10,
				},
			},
			want:    nil,
			want1:   0,
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {

				cart := models.Cart{
					ID:     1,
					UserID: 1,
				}
				isAvailable := true
				page := 1
				perPage := 10

				r.EXPECT().GetPromoByUserCart(gomock.Any(), nil, gomock.Eq(repository.GetPromoByUserCartInput{
					Cart:        cart,
					IsAvailable: &isAvailable,
					Page:        &page,
					PerPage:     &perPage,
				})).Return([]models.Promo{}, int64(0), errors.New("error"))
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
				r.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    uint(1),
					Relations: []string{"CartItems"},
				}).Return(models.Cart{
					ID:     1,
					UserID: 1,
				}, nil)
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: dto.GetPromoInput{
					UserId:  1,
					Page:    1,
					PerPage: 10,
				},
			},
			want: []models.Promo{
				{
					ID: 1,
				},
			},
			want1:   1,
			wantErr: false,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {

				cart := models.Cart{
					ID:     1,
					UserID: 1,
				}
				isAvailable := true
				page := 1
				perPage := 10

				r.EXPECT().GetPromoByUserCart(gomock.Any(), nil, gomock.Eq(repository.GetPromoByUserCartInput{
					Cart:        cart,
					IsAvailable: &isAvailable,
					Page:        &page,
					PerPage:     &perPage,
				})).Return([]models.Promo{{ID: 1}}, int64(1), nil)
			},
			cartRepoMock: func(r *repo_mock.MockCartRepository) {
				r.EXPECT().GetUserCart(gomock.Any(), nil, repository.GetUserCartInput{
					UserId:    uint(1),
					Relations: []string{"CartItems"},
				}).Return(models.Cart{
					ID:     1,
					UserID: 1,
				}, nil)
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
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

			promoRepo := repo_mock.NewMockPromoRepository(ctrl)
			tt.promoRepoMock(promoRepo)

			usecase := NewPromoUsecase(promoRepo, transactionRepo, cartRepo, productRepo)

			res, count, err := usecase.GetPromo(context.Background(), tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("promoUsecase.GetPromo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Println(reflect.DeepEqual(res, tt.want))

			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("promoUsecase.GetPromo() got = %v, want %v", res, tt.want)
			}

			if count != tt.want1 {
				t.Errorf("promoUsecase.GetPromo() got1 = %v, want %v", count, tt.want1)
			}

		})
	}
}

func Test_promoUsecase_ExtendPromo(t *testing.T) {
	type args struct {
		ctx context.Context
		dto dto.ExtendPromoInput
	}
	startDate := time.Now().Add(time.Hour * 24 * 1)
	endDate := time.Now().Add(time.Hour * 24 * 2)

	tests := []struct {
		name                string
		args                args
		wantErr             bool
		transactionRepoMock func(*repo_mock.MockTransactionRepository)
		promoRepoMock       func(*repo_mock.MockPromoRepository)
	}{
		{
			name: "err get promo by id",
			args: args{
				ctx: context.Background(),
				dto: dto.ExtendPromoInput{
					ID:        1,
					StartDate: time.Now().Add(-time.Hour * 24),
					EndDate:   time.Now(),
				},
			},
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
				r.EXPECT().GetPromoByPromoID(gomock.Any(), nil, uint(1)).Return(models.Promo{}, errors.New("error"))
			},
		},
		{
			name: "promo not found",
			args: args{
				ctx: context.Background(),
				dto: dto.ExtendPromoInput{
					ID:        1,
					StartDate: time.Now().Add(-time.Hour * 24),
					EndDate:   time.Now(),
				},
			},
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
				r.EXPECT().GetPromoByPromoID(gomock.Any(), nil, uint(1)).Return(models.Promo{}, nil)
			},
		},
		{
			name: "new end date must be greater than current end date",
			args: args{
				ctx: context.Background(),
				dto: dto.ExtendPromoInput{
					ID:        1,
					StartDate: time.Now().Add(-time.Hour * 24 * 30),
					EndDate:   time.Now().Add(-time.Hour * 24 * 1),
				},
			},
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
				r.EXPECT().GetPromoByPromoID(gomock.Any(), nil, uint(1)).Return(models.Promo{ID: 1, EndDate: time.Now()}, nil)
			},
		},
		{
			name: "err save promo",
			args: args{
				ctx: context.Background(),
				dto: dto.ExtendPromoInput{
					ID:        1,
					StartDate: startDate,
					EndDate:   endDate,
				},
			},
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
				r.EXPECT().GetPromoByPromoID(gomock.Any(), nil, uint(1)).Return(
					models.Promo{
						ID:      1,
						EndDate: time.Now().Add(-time.Hour * 24 * 1)},
					nil)
				promo := gomock.Eq(&models.Promo{
					ID:        1,
					StartDate: startDate,
					EndDate:   endDate,
				})
				r.EXPECT().Save(gomock.Any(), gomock.Any(), promo).Return(errors.New("error"))
			},
		},
		{
			name: "err save promo",
			args: args{
				ctx: context.Background(),
				dto: dto.ExtendPromoInput{
					ID:        1,
					StartDate: startDate,
					EndDate:   endDate,
				},
			},
			wantErr: false,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
						return fn(nil)
					})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
				r.EXPECT().GetPromoByPromoID(gomock.Any(), nil, uint(1)).Return(
					models.Promo{
						ID:      1,
						EndDate: time.Now().Add(-time.Hour * 24 * 1)},
					nil)
				promo := gomock.Eq(&models.Promo{
					ID:        1,
					StartDate: startDate,
					EndDate:   endDate,
				})
				r.EXPECT().Save(gomock.Any(), gomock.Any(), promo).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			transactionRepo := repo_mock.NewMockTransactionRepository(ctrl)
			tt.transactionRepoMock(transactionRepo)

			promoRepo := repo_mock.NewMockPromoRepository(ctrl)
			tt.promoRepoMock(promoRepo)

			usecase := NewPromoUsecase(promoRepo, transactionRepo, nil, nil)

			err := usecase.ExtendPromo(context.Background(), tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("promoUsecase.GetPromo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_promoUsecase_CreatePromo(t *testing.T) {
	type args struct {
		ctx context.Context
		dto dto.CreatePromoInput
	}
	buyProductId := int(1)
	freeProductId := int(2)
	cities := []string{"Jakarta", "Bogor"}
	tests := []struct {
		name                string
		args                args
		want                uint
		wantErr             bool
		transactionRepoMock func(*repo_mock.MockTransactionRepository)
		promoRepoMock       func(*repo_mock.MockPromoRepository)
		productRepoMock     func(*repo_mock.MockProductRepository)
	}{
		{
			name: "err get product",
			args: args{
				ctx: context.Background(),
				dto: dto.CreatePromoInput{
					Type:         constants.PROMOTYPEBUYXGETY,
					BuyProductId: &buyProductId,
				},
			},
			want:    0,
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
					return fn(nil)
				})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), nil, uint(buyProductId)).Return(models.Product{}, errors.New("error"))
			},
		},
		{
			name: "err buy product not found",
			args: args{
				ctx: context.Background(),
				dto: dto.CreatePromoInput{
					Type:         constants.PROMOTYPEBUYXGETY,
					BuyProductId: &buyProductId,
				},
			},
			want:    0,
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
					return fn(nil)
				})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), nil, uint(buyProductId)).Return(models.Product{}, nil)
			},
		},
		{
			name: "err get free product",
			args: args{
				ctx: context.Background(),
				dto: dto.CreatePromoInput{
					Type:          constants.PROMOTYPEBUYXGETY,
					BuyProductId:  &buyProductId,
					FreeProductId: &freeProductId,
				},
			},
			want:    0,
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
					return fn(nil)
				})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), nil, uint(buyProductId)).Return(models.Product{
					ID: uint(buyProductId),
				}, nil)

				r.EXPECT().Get(gomock.Any(), nil, uint(freeProductId)).Return(models.Product{}, errors.New("error"))
			},
		},
		{
			name: "err free product not found",
			args: args{
				ctx: context.Background(),
				dto: dto.CreatePromoInput{
					Type:          constants.PROMOTYPEBUYXGETY,
					BuyProductId:  &buyProductId,
					FreeProductId: &freeProductId,
				},
			},
			want:    0,
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
					return fn(nil)
				})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), nil, uint(buyProductId)).Return(models.Product{
					ID: uint(buyProductId),
				}, nil)

				r.EXPECT().Get(gomock.Any(), nil, uint(freeProductId)).Return(models.Product{}, nil)
			},
		},
		{
			name: "err save promo",
			args: args{
				ctx: context.Background(),
				dto: dto.CreatePromoInput{
					Type:          constants.PROMOTYPEBUYXGETY,
					BuyProductId:  &buyProductId,
					FreeProductId: &freeProductId,
				},
			},
			want:    0,
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
					return fn(nil)
				})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
				buyProductId := uint(buyProductId)
				freeProductId := uint(freeProductId)
				promo := gomock.Eq(&models.Promo{
					Type:              constants.PROMOTYPEBUYXGETY,
					BuyProductID:      &buyProductId,
					FreeProductID:     &freeProductId,
					CurrentUsageCount: 0,
				})
				r.EXPECT().Save(gomock.Any(), gomock.Any(), promo).Return(errors.New("error"))
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), nil, uint(buyProductId)).Return(models.Product{
					ID: uint(buyProductId),
				}, nil)

				r.EXPECT().Get(gomock.Any(), nil, uint(freeProductId)).Return(models.Product{
					ID: uint(freeProductId),
				}, nil)
			},
		},
		{
			name: "err save promo cities",
			args: args{
				ctx: context.Background(),
				dto: dto.CreatePromoInput{
					Type:          constants.PROMOTYPEBUYXGETY,
					BuyProductId:  &buyProductId,
					FreeProductId: &freeProductId,
					Segmentation:  constants.PROMOSEGMENTATIONCITY,
					Cities:        cities,
				},
			},
			want:    0,
			wantErr: true,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
					return fn(nil)
				})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
				buyProductId := uint(buyProductId)
				freeProductId := uint(freeProductId)
				promo := models.Promo{
					Type:              constants.PROMOTYPEBUYXGETY,
					BuyProductID:      &buyProductId,
					FreeProductID:     &freeProductId,
					CurrentUsageCount: 0,
					Segmentation:      constants.PROMOSEGMENTATIONCITY,
				}
				promoMock := gomock.Eq(&promo)
				r.EXPECT().Save(gomock.Any(), gomock.Any(), promoMock).Return(nil)
				r.EXPECT().SaveCities(gomock.Any(), gomock.Any(), promo.ID, cities).Return(errors.New("error"))
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), nil, uint(buyProductId)).Return(models.Product{
					ID: uint(buyProductId),
				}, nil)

				r.EXPECT().Get(gomock.Any(), nil, uint(freeProductId)).Return(models.Product{
					ID: uint(freeProductId),
				}, nil)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: dto.CreatePromoInput{
					Type:          constants.PROMOTYPEBUYXGETY,
					BuyProductId:  &buyProductId,
					FreeProductId: &freeProductId,
					Segmentation:  constants.PROMOSEGMENTATIONCITY,
					Cities:        cities,
				},
			},
			want:    1,
			wantErr: false,
			transactionRepoMock: func(r *repo_mock.MockTransactionRepository) {
				r.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
					return fn(nil)
				})
			},
			promoRepoMock: func(r *repo_mock.MockPromoRepository) {
				buyProductId := uint(buyProductId)
				freeProductId := uint(freeProductId)
				promo := models.Promo{
					Type:              constants.PROMOTYPEBUYXGETY,
					BuyProductID:      &buyProductId,
					FreeProductID:     &freeProductId,
					CurrentUsageCount: 0,
					Segmentation:      constants.PROMOSEGMENTATIONCITY,
				}
				promoMock := gomock.AssignableToTypeOf(&promo)
				r.EXPECT().Save(gomock.Any(), gomock.Any(), promoMock).DoAndReturn(func(ctx context.Context, tx any, p *models.Promo) error {
					p.ID = 1
					return nil
				})
				promo.ID = 1
				r.EXPECT().SaveCities(gomock.Any(), gomock.Any(), promo.ID, cities).Return(nil)
			},
			productRepoMock: func(r *repo_mock.MockProductRepository) {
				r.EXPECT().Get(gomock.Any(), nil, uint(buyProductId)).Return(models.Product{
					ID: uint(buyProductId),
				}, nil)

				r.EXPECT().Get(gomock.Any(), nil, uint(freeProductId)).Return(models.Product{
					ID: uint(freeProductId),
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

			productRepo := repo_mock.NewMockProductRepository(ctrl)
			tt.productRepoMock(productRepo)

			promoRepo := repo_mock.NewMockPromoRepository(ctrl)
			tt.promoRepoMock(promoRepo)

			usecase := NewPromoUsecase(promoRepo, transactionRepo, nil, productRepo)

			res, err := usecase.CreatePromo(context.Background(), tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("promoUsecase.GetPromo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("promoUsecase.GetPromo() got = %v, want %v", res, tt.want)
			}

		})
	}
}
