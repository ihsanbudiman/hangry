package models

import "time"

// Order represents the orders table
type Order struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	TotalAmount float64   `gorm:"not null;type:numeric(10,2)" json:"total_amount"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	User        *User        `gorm:"foreignKey:UserID" json:"user"`
	OrderItems  []OrderItem  `gorm:"foreignKey:OrderID" json:"order_items"`
	OrderPromos []OrderPromo `gorm:"foreignKey:OrderID" json:"order_promos"`
}
