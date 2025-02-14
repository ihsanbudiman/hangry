package db

import (
	"context"
	"errors"
	"hangry/constants"
	"hangry/domain/models"
	"hangry/repository"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_promoRepostory_GetPromoByUserCart(t *testing.T) {
	isAvailable := true
	page := 1
	perPage := 10

	type args struct {
		ctx   context.Context
		tx    *gorm.DB
		input repository.GetPromoByUserCartInput
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Promo
		wantCnt int64
		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				tx:  nil,
				input: repository.GetPromoByUserCartInput{
					Cart: models.Cart{
						UserID: 1,
					},
					PromoIds:    []uint{1, 2},
					IsAvailable: &isAvailable,
					Page:        &page,
					PerPage:     &perPage,
				},
			},
			want: []models.Promo{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			wantCnt: 2,
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
					with user_data as ( select * from users u where u.id = $1 ), items as ( select p.price, ci.* from cart_items ci join carts c on c.id = ci.cart_id join products p on p.id = ci.product_id where c.user_id = $2 ), summary as ( select sum(i.price * i.quantity) total from items i ) select p.* from promos p where ( ( p."type" = 'PERCENTAGE_DISCOUNT' and (select total from summary) >= p.min_order_amount ) or ( p."type" = 'BUY_X_GET_Y_FREE' and (select coalesce(i.quantity,0) from items i where i.product_id = p.buy_product_id ) >= p.buy_product_qty ) ) and ( p.segmentation = 'ALL' or ( p.segmentation = 'CITY' and lower((select city from user_data)) in (select lower(city) from promo_cities pc where pc.promo_id = p.id) ) or ( p.segmentation = 'LOYAL_USER' and (select is_loyal from user_data) = true ) or ( p.segmentation = 'NEW_USER' and (select created_at from user_data) > (current_timestamp - interval '1 month') ) ) and p.id in ($3,$4) and p.start_date <= current_timestamp and p.end_date >= current_timestamp and p.current_usage_count < p.max_usage_limit group by p.id limit 10 offset 0;
				`)
				mock.ExpectQuery(query).
					WithArgs(1, 1, 1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1).
						AddRow(2)).
					WillReturnError(nil)

				countQuery := regexp.QuoteMeta(`
				with 
            		user_data as (
            				select 
            						*
            				from users u 
            				where u.id = $1
            		),
            		items as (
            				select
            						p.price,
            						ci.*
            				from cart_items ci 
            				join carts c on c.id = ci.cart_id
            				join products p on p.id = ci.product_id 
            				where c.user_id  = $2
            		),
            		summary as (
            				select
            						sum(i.price * i.quantity) total
            				from items i
            		)
            		select
            				count(*)
            		from promos p
            		where (
            						(
            								p."type" = 'PERCENTAGE_DISCOUNT' 
            								and (select total from summary) >= p.min_order_amount 
            						)
            						or (
            								p."type" = 'BUY_X_GET_Y_FREE' 
            								and (select 
            										coalesce(i.quantity,0)
            										from items i
            										where i.product_id = p.buy_product_id 
            								) >= p.buy_product_qty 
            						)
            				)
            				and (
            						p.segmentation = 'ALL'
            						or (
            								p.segmentation = 'CITY' 
            								and lower((select city from user_data)) in (select lower(city) from promo_cities pc where pc.promo_id = p.id)
            						)
            						or (
            								p.segmentation = 'LOYAL_USER'
            								and (select is_loyal from user_data) = true
            						)
            						or (
            								p.segmentation = 'NEW_USER'
            								and (select created_at from user_data) > (current_timestamp - interval '1 month')
            						)
            				)
            		 and p.id in ($3,$4) and p.start_date <= current_timestamp and p.end_date >= current_timestamp and p.current_usage_count < p.max_usage_limit
            		;
				`)
				mock.ExpectQuery(countQuery).
					WithArgs(1, 1, 1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(2)).
					WillReturnError(nil)

			},
		},
		{
			name: "error get promo",
			args: args{
				ctx: context.Background(),
				tx:  nil,
				input: repository.GetPromoByUserCartInput{
					Cart: models.Cart{
						UserID: 1,
					},
					PromoIds:    []uint{1, 2},
					IsAvailable: &isAvailable,
					Page:        &page,
					PerPage:     &perPage,
				},
			},
			want:    nil,
			wantCnt: 0,
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
					with user_data as ( select * from users u where u.id = $1 ), items as ( select p.price, ci.* from cart_items ci join carts c on c.id = ci.cart_id join products p on p.id = ci.product_id where c.user_id = $2 ), summary as ( select sum(i.price * i.quantity) total from items i ) select p.* from promos p where ( ( p."type" = 'PERCENTAGE_DISCOUNT' and (select total from summary) >= p.min_order_amount ) or ( p."type" = 'BUY_X_GET_Y_FREE' and (select coalesce(i.quantity,0) from items i where i.product_id = p.buy_product_id ) >= p.buy_product_qty ) ) and ( p.segmentation = 'ALL' or ( p.segmentation = 'CITY' and lower((select city from user_data)) in (select lower(city) from promo_cities pc where pc.promo_id = p.id) ) or ( p.segmentation = 'LOYAL_USER' and (select is_loyal from user_data) = true ) or ( p.segmentation = 'NEW_USER' and (select created_at from user_data) > (current_timestamp - interval '1 month') ) ) and p.id in ($3,$4) and p.start_date <= current_timestamp and p.end_date >= current_timestamp and p.current_usage_count < p.max_usage_limit group by p.id limit 10 offset 0;
				`)
				mock.ExpectQuery(query).
					WithArgs(1, 1, 1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id"})).
					WillReturnError(errors.New("error"))

			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				tx:  nil,
				input: repository.GetPromoByUserCartInput{
					Cart: models.Cart{
						UserID: 1,
					},
					PromoIds:    []uint{1, 2},
					IsAvailable: &isAvailable,
					Page:        &page,
					PerPage:     &perPage,
				},
			},
			want:    nil,
			wantCnt: 0,
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
					with user_data as ( select * from users u where u.id = $1 ), items as ( select p.price, ci.* from cart_items ci join carts c on c.id = ci.cart_id join products p on p.id = ci.product_id where c.user_id = $2 ), summary as ( select sum(i.price * i.quantity) total from items i ) select p.* from promos p where ( ( p."type" = 'PERCENTAGE_DISCOUNT' and (select total from summary) >= p.min_order_amount ) or ( p."type" = 'BUY_X_GET_Y_FREE' and (select coalesce(i.quantity,0) from items i where i.product_id = p.buy_product_id ) >= p.buy_product_qty ) ) and ( p.segmentation = 'ALL' or ( p.segmentation = 'CITY' and lower((select city from user_data)) in (select lower(city) from promo_cities pc where pc.promo_id = p.id) ) or ( p.segmentation = 'LOYAL_USER' and (select is_loyal from user_data) = true ) or ( p.segmentation = 'NEW_USER' and (select created_at from user_data) > (current_timestamp - interval '1 month') ) ) and p.id in ($3,$4) and p.start_date <= current_timestamp and p.end_date >= current_timestamp and p.current_usage_count < p.max_usage_limit group by p.id limit 10 offset 0;
				`)
				mock.ExpectQuery(query).
					WithArgs(1, 1, 1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1).
						AddRow(2)).
					WillReturnError(nil)

				countQuery := regexp.QuoteMeta(`
				with 
            		user_data as (
            				select 
            						*
            				from users u 
            				where u.id = $1
            		),
            		items as (
            				select
            						p.price,
            						ci.*
            				from cart_items ci 
            				join carts c on c.id = ci.cart_id
            				join products p on p.id = ci.product_id 
            				where c.user_id  = $2
            		),
            		summary as (
            				select
            						sum(i.price * i.quantity) total
            				from items i
            		)
            		select
            				count(*)
            		from promos p
            		where (
            						(
            								p."type" = 'PERCENTAGE_DISCOUNT' 
            								and (select total from summary) >= p.min_order_amount 
            						)
            						or (
            								p."type" = 'BUY_X_GET_Y_FREE' 
            								and (select 
            										coalesce(i.quantity,0)
            										from items i
            										where i.product_id = p.buy_product_id 
            								) >= p.buy_product_qty 
            						)
            				)
            				and (
            						p.segmentation = 'ALL'
            						or (
            								p.segmentation = 'CITY' 
            								and lower((select city from user_data)) in (select lower(city) from promo_cities pc where pc.promo_id = p.id)
            						)
            						or (
            								p.segmentation = 'LOYAL_USER'
            								and (select is_loyal from user_data) = true
            						)
            						or (
            								p.segmentation = 'NEW_USER'
            								and (select created_at from user_data) > (current_timestamp - interval '1 month')
            						)
            				)
            		 and p.id in ($3,$4) and p.start_date <= current_timestamp and p.end_date >= current_timestamp and p.current_usage_count < p.max_usage_limit
            		;
				`)
				mock.ExpectQuery(countQuery).
					WithArgs(1, 1, 1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"})).
					WillReturnError(errors.New("error"))

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer sqlDB.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.sqlMock(mock)

			r := NewPromoRepository(gormDB)

			got, gotCnt, err := r.GetPromoByUserCart(tt.args.ctx, tt.args.tx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("promoRepostory.GetPromoByUserCart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("promoRepostory.GetPromoByUserCart() = %v, want %v", got, tt.want)
			}
			if gotCnt != tt.wantCnt {
				t.Errorf("promoRepostory.GetPromoByUserCart() count = %v, want %v", gotCnt, tt.wantCnt)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_promoRepostory_GetPromoByPromoID(t *testing.T) {
	timeNow := time.Now()
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx     context.Context
		tx      *gorm.DB
		promoID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Promo
		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			fields: fields{
				db: nil,
			},
			args: args{
				ctx:     context.Background(),
				tx:      nil,
				promoID: 1,
			},
			want: models.Promo{
				ID:        1,
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`SELECT * FROM "promos" WHERE id = $1 ORDER BY "promos"."id" LIMIT $2`)
				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(1, timeNow, timeNow))
			},
		},
		{
			name: "err get promo",
			fields: fields{
				db: nil,
			},
			args: args{
				ctx:     context.Background(),
				tx:      nil,
				promoID: 1,
			},
			want:    models.Promo{},
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`SELECT * FROM "promos" WHERE id = $1 ORDER BY "promos"."id" LIMIT $2`)
				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnError(errors.New("error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer sqlDB.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.sqlMock(mock)

			r := NewPromoRepository(gormDB)
			got, err := r.GetPromoByPromoID(tt.args.ctx, tt.args.tx, tt.args.promoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("promoRepostory.GetPromoByPromoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("promoRepostory.GetPromoByPromoID() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_promoRepostory_SaveCities(t *testing.T) {
	timeNow := time.Now()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx     context.Context
		tx      *gorm.DB
		promoID uint
		cities  []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			fields: fields{
				db: nil,
			},
			args: args{
				ctx:     context.Background(),
				tx:      nil,
				promoID: 1,
				cities:  []string{"Jakarta", "Bandung"},
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				query := regexp.QuoteMeta(`INSERT INTO "promo_cities" ("promo_id","city") VALUES ($1,$2),($3,$4) RETURNING "created_at","updated_at","id"`)
				mock.ExpectQuery(query).
					WithArgs(1, "Jakarta", 1, "Bandung").
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "id"}).
						AddRow(timeNow, timeNow, 1).
						AddRow(timeNow, timeNow, 2))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer sqlDB.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.sqlMock(mock)

			r := NewPromoRepository(gormDB)
			if err := r.SaveCities(tt.args.ctx, tt.args.tx, tt.args.promoID, tt.args.cities); (err != nil) != tt.wantErr {
				t.Errorf("promoRepostory.SaveCities() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_promoRepostory_Save(t *testing.T) {
	timeNow := time.Now()
	buyProductID := uint(1)
	freeProductID := uint(1)
	maxUsageLimit := 1

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx   context.Context
		tx    *gorm.DB
		promo *models.Promo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			fields: fields{
				db: nil,
			},
			args: args{
				ctx: context.Background(),
				tx:  nil,
				promo: &models.Promo{
					ID:                1,
					Name:              "name",
					Description:       "description",
					Segmentation:      constants.PROMOSEGMENTATIONALL,
					Type:              constants.PROMOTYPEBUYXGETY,
					MinOrderAmount:    1,
					DiscountValue:     1,
					MaxDiscountAmount: 1,
					BuyProductID:      &buyProductID,
					FreeProductID:     &freeProductID,
					BuyProductQty:     1,
					FreeProductQty:    1,
					StartDate:         timeNow,
					EndDate:           timeNow,
					MaxUsageLimit:     &maxUsageLimit,
					CurrentUsageCount: 1,
					CreatedAt:         timeNow,
					UpdatedAt:         timeNow,
				},
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				query := regexp.QuoteMeta(`UPDATE "promos" SET "name"=$1,"description"=$2,"segmentation"=$3,"type"=$4,"min_order_amount"=$5,"discount_value"=$6,"max_discount_amount"=$7,"buy_product_id"=$8,"free_product_id"=$9,"buy_product_qty"=$10,"free_product_qty"=$11,"start_date"=$12,"end_date"=$13,"max_usage_limit"=$14,"current_usage_count"=$15,"created_at"=$16,"updated_at"=$17 WHERE "id" = $18`)
				mock.ExpectExec(query).
					WithArgs("name", "description", constants.PROMOSEGMENTATIONALL, constants.PROMOTYPEBUYXGETY,
						float64(1), float64(1), float64(1),
						buyProductID, freeProductID, 1, 1,
						sqlmock.AnyArg(), sqlmock.AnyArg(), // start_date, end_date
						1, 1,
						sqlmock.AnyArg(), sqlmock.AnyArg(), // created_at, updated_at
						uint(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			fields: fields{
				db: nil,
			},
			args: args{
				ctx: context.Background(),
				tx:  nil,
				promo: &models.Promo{
					ID:                1,
					Name:              "name",
					Description:       "description",
					Segmentation:      constants.PROMOSEGMENTATIONALL,
					Type:              constants.PROMOTYPEBUYXGETY,
					MinOrderAmount:    1,
					DiscountValue:     1,
					MaxDiscountAmount: 1,
					BuyProductID:      &buyProductID,
					FreeProductID:     &freeProductID,
					BuyProductQty:     1,
					FreeProductQty:    1,
					StartDate:         timeNow,
					EndDate:           timeNow,
					MaxUsageLimit:     &maxUsageLimit,
					CurrentUsageCount: 1,
					CreatedAt:         timeNow,
					UpdatedAt:         timeNow,
				},
			},
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				query := regexp.QuoteMeta(`UPDATE "promos" SET "name"=$1,"description"=$2,"segmentation"=$3,"type"=$4,"min_order_amount"=$5,"discount_value"=$6,"max_discount_amount"=$7,"buy_product_id"=$8,"free_product_id"=$9,"buy_product_qty"=$10,"free_product_qty"=$11,"start_date"=$12,"end_date"=$13,"max_usage_limit"=$14,"current_usage_count"=$15,"created_at"=$16,"updated_at"=$17 WHERE "id" = $18`)
				mock.ExpectExec(query).
					WithArgs("name", "description", constants.PROMOSEGMENTATIONALL, constants.PROMOTYPEBUYXGETY,
						float64(1), float64(1), float64(1),
						buyProductID, freeProductID, 1, 1,
						sqlmock.AnyArg(), sqlmock.AnyArg(), // start_date, end_date
						1, 1,
						sqlmock.AnyArg(), sqlmock.AnyArg(), // created_at, updated_at
						uint(1)).
					WillReturnError(errors.New("error"))

				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer sqlDB.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.sqlMock(mock)

			r := NewPromoRepository(gormDB)
			if err := r.Save(tt.args.ctx, gormDB, tt.args.promo); (err != nil) != tt.wantErr {
				t.Errorf("promoRepostory.Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
