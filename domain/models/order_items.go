package models

import "time"

// OrderItem represents the order_items table
type OrderItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderID     uint      `gorm:"not null" json:"order_id"`
	ProductID   uint      `gorm:"not null" json:"product_id"`
	Price       float64   `gorm:"not null" json:"price"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	TotalAmount float64   `gorm:"not null" json:"total_amount"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	Order   *Order   `gorm:"foreignKey:OrderID" json:"order"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product"`
}
