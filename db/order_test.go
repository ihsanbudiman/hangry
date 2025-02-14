package db

import (
	"context"
	"errors"
	"hangry/domain/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_orderRepository_GetUserOrderCount(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx    context.Context
		tx     *gorm.DB
		userID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			fields: fields{
				db: nil,
			},
			args: args{
				ctx:    context.Background(),
				tx:     nil,
				userID: 1,
			},
			want:    5,
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`SELECT count(*) FROM "orders" WHERE user_id = $1`)
				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
			},
		},
		{
			name: "error",
			fields: fields{
				db: nil,
			},
			args: args{
				ctx:    context.Background(),
				tx:     nil,
				userID: 1,
			},
			want:    0,
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`SELECT count(*) FROM "orders" WHERE user_id = $1`)
				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnError(errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock database: %v", err)
			}
			defer sqlDB.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm database: %v", err)
			}

			tt.sqlMock(mock)

			o := NewOrderRepository(gormDB)
			got, err := o.GetUserOrderCount(tt.args.ctx, tt.args.tx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderRepository.GetUserOrderCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("orderRepository.GetUserOrderCount() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_orderRepository_MakeOrder(t *testing.T) {
	timeNow := time.Now()
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx   context.Context
		tx    *gorm.DB
		order *models.Order
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
				order: &models.Order{
					UserID:     1,
					OrderItems: []models.OrderItem{{ProductID: 1, Quantity: 2}},
				},
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders" ("user_id","total_amount") VALUES ($1,$2) RETURNING "created_at","updated_at","id"`)).
					WithArgs(int64(1), float64(0)).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "id"}).AddRow(timeNow, timeNow, 1)).
					WillReturnError(nil)
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "order_items" ("order_id","product_id","price","quantity","total_amount") VALUES ($1,$2,$3,$4,$5) ON CONFLICT ("id") DO UPDATE SET "order_id"="excluded"."order_id","product_id"="excluded"."product_id","price"="excluded"."price","quantity"="excluded"."quantity","total_amount"="excluded"."total_amount" RETURNING "created_at","updated_at","id"`)).
					WithArgs(uint(1), uint(1), float64(0), 2, float64(0)).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "id"}).AddRow(timeNow, timeNow, 1)).
					WillReturnError(nil)
				mock.ExpectCommit()
			},
		},
		{
			name: "failed to insert order",
			fields: fields{
				db: nil,
			},
			args: args{
				ctx: context.Background(),
				tx:  nil,
				order: &models.Order{
					UserID: 1,
				},
			},
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders" ("user_id","total_amount") VALUES ($1,$2) RETURNING "created_at","updated_at","id"`)).
					WithArgs(int64(1), float64(0)).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "id"})).
					WillReturnError(errors.New("error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock database: %v", err)
			}
			defer sqlDB.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm database: %v", err)
			}

			tt.sqlMock(mock)

			o := NewOrderRepository(gormDB)

			if err := o.MakeOrder(tt.args.ctx, tt.args.tx, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("orderRepository.MakeOrder() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}
