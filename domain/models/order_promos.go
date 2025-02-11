package models

import "time"

// OrderPromo represents the order_promos table
type OrderPromo struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OrderID        uint      `gorm:"not null" json:"order_id"`
	PromoID        uint      `gorm:"not null" json:"promo_id"`
	DiscountAmount float64   `gorm:"type:numeric(10,2)" json:"discount_amount"`
	FreeProductID  *uint     `gorm:"index" json:"free_product_id"`
	FreeProductQty int       `json:"free_product_qty"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	Order *Order `gorm:"foreignKey:OrderID" json:"order"`
	Promo *Promo `gorm:"foreignKey:PromoID" json:"promo"`
}
