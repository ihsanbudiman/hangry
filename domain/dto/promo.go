package dto

import (
	"hangry/domain/models"
	"time"
)

type CreatePromoInput struct {
	Name              string    `json:"name"`
	Description       *string   `json:"description,omitempty"`
	Segmentation      string    `json:"segmentation"`
	Type              string    `json:"type"`
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	MinOrderAmount    *float64  `json:"minOrderAmount,omitempty"`
	DiscountValue     *float64  `json:"discountValue,omitempty"`
	MaxDiscountAmount *float64  `json:"maxDiscountAmount,omitempty"`
	BuyProductId      *int      `json:"buyProductId,omitempty"`
	BuyItemCount      *int      `json:"buyItemCount,omitempty"`
	FreeProductId     *int      `json:"freeProductId,omitempty"`
	FreeItemCount     *int      `json:"freeItemCount,omitempty"`
	MaxUsageLimit     *int      `json:"maxUsageLimit,omitempty"`
	Cities            []string  `json:"cities,omitempty"`
}

func (c *CreatePromoInput) CreatePromoModel() models.Promo {
	promo := models.Promo{
		Name:              c.Name,
		Segmentation:      c.Segmentation,
		Type:              c.Type,
		StartDate:         c.StartDate,
		EndDate:           c.EndDate,
		CurrentUsageCount: 0,
	}

	if c.Description != nil {
		promo.Description = *c.Description
	}

	if c.MinOrderAmount != nil {
		promo.MinOrderAmount = *c.MinOrderAmount
	}

	if c.DiscountValue != nil {
		promo.DiscountValue = *c.DiscountValue
	}

	if c.MaxDiscountAmount != nil {
		promo.MaxDiscountAmount = *c.MaxDiscountAmount
	}

	if c.BuyProductId != nil {
		buyProductID := uint(*c.BuyProductId)
		promo.BuyProductID = &buyProductID
	}

	if c.BuyItemCount != nil {
		promo.BuyProductQty = *c.BuyItemCount
	}

	if c.FreeProductId != nil {
		freeProductID := uint(*c.FreeProductId)
		promo.FreeProductID = &freeProductID
	}

	if c.FreeItemCount != nil {
		promo.FreeProductQty = *c.FreeItemCount
	}

	if c.MaxUsageLimit != nil {
		promo.MaxUsageLimit = c.MaxUsageLimit
	}

	return promo
}

type ExtendPromoInput struct {
	ID        uint      `json:"id"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type GetPromoInput struct {
	UserId  uint
	Page    int
	PerPage int
}
