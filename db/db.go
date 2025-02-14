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

func NewDBConn(options NewDBConnOptions) *gorm.DB {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(options.Dsn), config)
	if err != nil {
		panic(err)
	}

	return db
}
