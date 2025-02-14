package db

import (
	"context"
	"errors"
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

func Test_cartRepositoryImpl_RemoveCartItem(t *testing.T) {

	type args struct {
		ctx         context.Context
		cartItemIds []uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			args: args{
				ctx:         context.Background(),
				cartItemIds: []uint{1},
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				// Expect transaction begin
				mock.ExpectBegin()

				// Prepare the DELETE query expectation
				query := regexp.QuoteMeta("DELETE FROM \"cart_items\" WHERE id IN ($1)")
				mock.ExpectExec(query).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1)).
					WillReturnError(nil)

				// Expect commit transaction
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

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})

			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.sqlMock(mock)

			repo := NewCartRepository(gormDB)

			err = repo.RemoveCartItem(tt.args.ctx, gormDB, tt.args.cartItemIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("cartRepositoryImpl.RemoveCartItem() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_cartRepositoryImpl_AddToCart(t *testing.T) {
	timeNow := time.Now()
	type args struct {
		ctx   context.Context
		tx    *gorm.DB
		input *models.CartItem
	}
	tests := []struct {
		name   string
		args   args
		withTx bool

		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock, isTx bool)
	}{
		{
			name:   "success",
			withTx: false,
			args: args{
				ctx:   context.Background(),
				tx:    nil,
				input: &models.CartItem{CartID: 1, ProductID: 1},
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock, isTx bool) {
				mock.ExpectBegin()
				qeury := regexp.QuoteMeta(`INSERT INTO "cart_items" ("cart_id","product_id","quantity") VALUES ($1,$2,$3) RETURNING "created_at","updated_at","id"`)
				mock.ExpectQuery(qeury).
					WithArgs(1, 1, 0).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "id"}).AddRow(timeNow, timeNow, 1)).WillReturnError(nil)
				mock.ExpectCommit()
			},
		},
		{
			name:   "error",
			withTx: false,
			args: args{
				ctx:   context.Background(),
				tx:    nil,
				input: &models.CartItem{CartID: 1, ProductID: 1},
			},
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock, isTx bool) {
				mock.ExpectBegin()
				qeury := regexp.QuoteMeta(`INSERT INTO "cart_items" ("cart_id","product_id","quantity") VALUES ($1,$2,$3) RETURNING "created_at","updated_at","id"`)
				mock.ExpectQuery(qeury).
					WithArgs(1, 1, 0).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "id"})).WillReturnError(errors.New("error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// setup sqlmock and gormDB
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.sqlMock(mock, tt.args.tx != nil)

			repo := NewCartRepository(gormDB)
			err = repo.AddToCart(tt.args.ctx, gormDB, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("cartRepositoryImpl.AddToCart() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_cartRepositoryImpl_CheckItem(t *testing.T) {

	cartId := uint(1)
	userId := uint(1)
	timeNow := time.Now()

	type args struct {
		ctx   context.Context
		input repository.CheckItemInput
	}
	tests := []struct {
		name    string
		args    args
		want    models.CartItem
		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				input: repository.CheckItemInput{
					CartId:    &cartId,
					ProductId: 1,
					UserId:    &userId,
				},
			},
			want: models.CartItem{
				ID:        1,
				CartID:    cartId,
				ProductID: 1,
				Quantity:  10,
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`SELECT "cart_items"."id","cart_items"."cart_id","cart_items"."product_id","cart_items"."quantity","cart_items"."created_at","cart_items"."updated_at" FROM "cart_items" inner join carts on cart_items.cart_id = carts.id WHERE carts.user_id = $1 AND cart_id = $2 AND product_id = $3 ORDER BY "cart_items"."id" LIMIT $4`)
				mock.ExpectQuery(query).
					WithArgs(userId, cartId, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product_id", "quantity", "created_at", "updated_at"}).
						AddRow(uint(1), cartId, 1, uint(10), timeNow, timeNow),
					).WillReturnError(nil)
			},
		},
		{
			name: "failed",
			args: args{
				ctx: context.Background(),
				input: repository.CheckItemInput{
					CartId:    &cartId,
					ProductId: 1,
					UserId:    &userId,
				},
			},
			want:    models.CartItem{},
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`SELECT "cart_items"."id","cart_items"."cart_id","cart_items"."product_id","cart_items"."quantity","cart_items"."created_at","cart_items"."updated_at" FROM "cart_items" inner join carts on cart_items.cart_id = carts.id WHERE carts.user_id = $1 AND cart_id = $2 AND product_id = $3 ORDER BY "cart_items"."id" LIMIT $4`)
				mock.ExpectQuery(query).
					WithArgs(userId, cartId, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product_id", "quantity", "created_at", "updated_at"})).WillReturnError(errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})

			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.sqlMock(mock)

			repo := NewCartRepository(gormDB)

			got, err := repo.CheckItem(tt.args.ctx, gormDB, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("cartRepositoryImpl.CheckItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cartRepositoryImpl.CheckItem() = %v, want %v", got, tt.want)
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_cartRepositoryImpl_CreateCart(t *testing.T) {
	timeNow := time.Now()
	type args struct {
		ctx    context.Context
		tx     *gorm.DB
		userId uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock, userId uint)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				tx:     nil,
				userId: 1,
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock, userId uint) {
				mock.ExpectBegin()
				query := regexp.QuoteMeta(`INSERT INTO "carts" ("user_id") VALUES ($1) RETURNING "created_at","updated_at","id"`)
				mock.ExpectQuery(query).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "id"}).AddRow(timeNow, timeNow, 1)).WillReturnError(nil)
				mock.ExpectCommit()
			},
		},
		{
			name: "failed",
			args: args{
				ctx:    context.Background(),
				tx:     nil,
				userId: 1,
			},
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock, userId uint) {
				mock.ExpectBegin()
				query := regexp.QuoteMeta(`INSERT INTO "carts" ("user_id") VALUES ($1) RETURNING "created_at","updated_at","id"`)
				mock.ExpectQuery(query).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "id"})).WillReturnError(errors.New("error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.sqlMock(mock, tt.args.userId)

			repo := NewCartRepository(gormDB)
			got, err := repo.CreateCart(tt.args.ctx, gormDB, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("cartRepositoryImpl.CreateCart() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && got.UserID != tt.args.userId {
				t.Errorf("cartRepositoryImpl.CreateCart() got = %v, want userId %v", got.UserID, tt.args.userId)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_cartRepositoryImpl_GetUserCart(t *testing.T) {
	timeNow := time.Now()
	type args struct {
		ctx   context.Context
		tx    *gorm.DB
		input repository.GetUserCartInput
	}
	tests := []struct {
		name    string
		args    args
		want    models.Cart
		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success with relation preload",
			args: args{
				ctx: context.Background(),
				tx:  nil,
				input: repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems"},
				},
			},
			want: models.Cart{
				ID:        1,
				UserID:    1,
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
				CartItems: []models.CartItem{
					{
						ID:        1,
						CartID:    1,
						ProductID: 1,
						Quantity:  10,
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					},
				},
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				// Expect query with preload does not change the SQL generated by gorm for the main query.
				query := regexp.QuoteMeta(`SELECT * FROM "carts" WHERE user_id = $1 ORDER BY "carts"."id" LIMIT $2`)
				// Return a successful row
				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at"}).
						AddRow(uint(1), uint(1), timeNow, timeNow),
					).WillReturnError(nil)

				cartItemQuery := regexp.QuoteMeta(`SELECT * FROM "cart_items" WHERE "cart_items"."cart_id" = $1`)
				mock.ExpectQuery(cartItemQuery).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product_id", "quantity", "created_at", "updated_at"}).
						AddRow(uint(1), uint(1), uint(1), uint(10), timeNow, timeNow),
					).WillReturnError(nil)

			},
		},
		{
			name: "failed",
			args: args{
				ctx: context.Background(),
				tx:  nil,
				input: repository.GetUserCartInput{
					UserId:    1,
					Relations: []string{"CartItems"},
				},
			},
			want:    models.Cart{},
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				// Expect query with preload does not change the SQL generated by gorm for the main query.
				query := regexp.QuoteMeta(`SELECT * FROM "carts" WHERE user_id = $1 ORDER BY "carts"."id" LIMIT $2`)
				// Return a successful row
				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at"})).WillReturnError(errors.New("error"))

			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.sqlMock(mock)

			repo := NewCartRepository(gormDB)
			got, err := repo.GetUserCart(tt.args.ctx, gormDB, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("cartRepositoryImpl.GetUserCart() error = %v, wantErr %v", err, tt.wantErr)
			}
			// In the "not found" case, got is zero value.
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cartRepositoryImpl.GetUserCart() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
