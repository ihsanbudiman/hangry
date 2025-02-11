package db

import (
	"context"
	"hangry/domain/models"
	"hangry/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// Get implements repository.UserRepository.
func (u *userRepository) Get(ctx context.Context, tx *gorm.DB, userID uint) (*models.User, error) {
	db := tx
	if db == nil {
		db = u.db.WithContext(ctx)
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Save implements repository.UserRepository.
func (u *userRepository) Save(ctx context.Context, tx *gorm.DB, user *models.User) error {
	db := tx
	if db == nil {
		db = u.db.WithContext(ctx)
	}

	if err := db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}
