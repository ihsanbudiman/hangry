package models

import "time"

// PromoCity represents the promo_cities table
type PromoCity struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PromoID   uint      `gorm:"not null" json:"promo_id"`
	City      string    `gorm:"not null;size:255" json:"city"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	Promo *Promo `gorm:"foreignKey:PromoID"`
}
