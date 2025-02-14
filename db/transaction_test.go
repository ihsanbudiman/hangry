package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_transactionRepository_Execute(t *testing.T) {

	type args struct {
		ctx context.Context
		fn  func(tx *gorm.DB) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(mock sqlmock.Sqlmock)
	}{
		{
			name: "successful transaction",
			args: args{
				ctx: context.Background(),
				fn: func(tx *gorm.DB) error {
					return nil
				},
			},
			wantErr: false,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
		},
		{
			name: "failed begin transaction",
			args: args{
				ctx: context.Background(),
				fn: func(tx *gorm.DB) error {
					return nil
				},
			},
			wantErr: true,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
		},
		{
			name: "err on fn",
			args: args{
				ctx: context.Background(),
				fn: func(tx *gorm.DB) error {
					return gorm.ErrInvalidTransaction
				},
			},
			wantErr: true,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
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

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatalf("failed to open gorm connection: %v", err)
			}

			tt.mockFn(mock)

			r := NewTransactionRepository(gormDB)
			if err := r.Execute(tt.args.ctx, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("transactionRepository.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
