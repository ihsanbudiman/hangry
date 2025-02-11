package repository

import (
	"context"
	"hangry/domain/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=./user_repository.go -destination=./mocks/mock_user_repository.go -package=mocks
type UserRepository interface {
	Save(ctx context.Context, tx *gorm.DB, user *models.User) error
	Get(ctx context.Context, tx *gorm.DB, userID uint) (*models.User, error)
}
