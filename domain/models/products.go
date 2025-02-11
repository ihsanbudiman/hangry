package models

import "time"

// Product represents the products table
type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;size:255" json:"name"`
	Price     float64   `gorm:"not null;type:numeric(10,2)" json:"price"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	OrderItems []OrderItem `gorm:"foreignKey:ProductID" json:"order_items"`
	CartItems  []CartItem  `gorm:"foreignKey:ProductID" json:"cart_items"`
}
