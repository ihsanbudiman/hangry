package db

import (
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type NewDBConnOptions struct {
	Dsn string
}

func NewDBConn(opts NewDBConnOptions) *gorm.DB {
	db, err := gorm.Open(postgres.Open(opts.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	return db
}
