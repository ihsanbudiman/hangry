package models

import (
	"time"
)

// User represents the users table
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;size:255" json:"name"`
	Email     string    `gorm:"unique;not null;size:255" json:"email"`
	City      string    `gorm:"size:255" json:"city"`
	IsLoyal   bool      `gorm:"default:false" json:"is_loyal"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	Orders []Order `gorm:"foreignKey:UserID" json:"orders"`
	Carts  []Cart  `gorm:"foreignKey:UserID" json:"carts"`
}
