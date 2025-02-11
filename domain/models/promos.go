package models

import "time"

// Promo represents the promos table
type Promo struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"not null;size:255" json:"name"`
	Description       string    `gorm:"type:text" json:"description"`
	Segmentation      string    `gorm:"not null;size:255;check:segmentation IN ('ALL', 'LOYAL_USER', 'NEW_USER', 'CITY')" json:"segmentation"`
	Type              string    `gorm:"not null;size:50;check:type IN ('PERCENTAGE_DISCOUNT', 'BUY_X_GET_Y_FREE')" json:"type"`
	MinOrderAmount    float64   `gorm:"type:numeric(10,2)" json:"min_order_amount"`
	DiscountValue     float64   `gorm:"type:numeric(10,2)" json:"discount_value"`
	MaxDiscountAmount float64   `gorm:"type:numeric(10,2)" json:"max_discount_amount"`
	BuyProductID      *uint     `gorm:"index" json:"buy_product_id"`
	FreeProductID     *uint     `gorm:"index" json:"free_product_id"`
	BuyProductQty     int       `json:"buy_product_qty"`
	FreeProductQty    int       `json:"free_product_qty"`
	StartDate         time.Time `gorm:"not null" json:"start_date"`
	EndDate           time.Time `gorm:"not null" json:"end_date"`
	MaxUsageLimit     *int      `json:"max_usage_limit"`
	CurrentUsageCount int       `gorm:"default:0" json:"current_usage_count"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	PromoCities []PromoCity  `gorm:"foreignKey:PromoID" json:"promo_cities"`
	OrderPromos []OrderPromo `gorm:"foreignKey:PromoID" json:"order_promos"`
}
