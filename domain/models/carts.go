package models

import "time"

// Cart represents the carts table
type Cart struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	User      *User      `gorm:"foreignKey:UserID" json:"user"`
	CartItems []CartItem `gorm:"foreignKey:CartID" json:"cart_items"`
}
