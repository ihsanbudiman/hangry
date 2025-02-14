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

func Test_productRepository_Get(t *testing.T) {
	timeNow := time.Now()
	type args struct {
		ctx context.Context
		tx  *gorm.DB
		id  uint
	}
	tests := []struct {
		name    string
		args    args
		want    models.Product
		wantErr bool
		sqlMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				tx:  nil,
				id:  1,
			},
			want: models.Product{
				ID:        1,
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 ORDER BY "products"."id" LIMIT $2`)
				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(1, timeNow, timeNow))
			},
		},
		{
			name: "failed",
			args: args{
				ctx: context.Background(),
				tx:  nil,
				id:  2,
			},
			want:    models.Product{},
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 ORDER BY "products"."id" LIMIT $2`)
				mock.ExpectQuery(query).
					WithArgs(2, 1).
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

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.sqlMock(mock)

			repo := NewProductRepository(gormDB)
			got, err := repo.Get(tt.args.ctx, gormDB, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("productRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("productRepository.Get() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
