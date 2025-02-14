package db

import (
	"context"
	"errors"
	"hangry/domain/models"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_userRepository_Get(t *testing.T) {
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
		want    *models.User
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
			want: &models.User{
				ID: 1,
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)
				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
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
			want:    nil,
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)
				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"})).
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

			u := NewUserRepository(gormDB)

			got, err := u.Get(tt.args.ctx, tt.args.tx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepository.Get() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_userRepository_Save(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx  context.Context
		tx   *gorm.DB
		user *models.User
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
				user: &models.User{
					ID:   1,
					Name: "John Doe",
				},
			},
			wantErr: false,
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				query := regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"email"=$2,"city"=$3,"is_loyal"=$4,"created_at"=$5,"updated_at"=$6 WHERE "id" = $7`)
				mock.ExpectExec(query).
					WithArgs("John Doe", "", "", false, sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
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
				user: &models.User{
					ID:   1,
					Name: "John Doe",
				},
			},
			wantErr: true,
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				query := regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"email"=$2,"city"=$3,"is_loyal"=$4,"created_at"=$5,"updated_at"=$6 WHERE "id" = $7`)
				mock.ExpectExec(query).
					WithArgs("John Doe", "", "", false, sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
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

			u := NewUserRepository(gormDB)

			if err := u.Save(tt.args.ctx, tt.args.tx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}
