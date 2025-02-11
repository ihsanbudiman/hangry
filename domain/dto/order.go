package dto

type OrderInput struct {
	UserId   uint   `json:"user_id"`
	PromoIds []uint `json:"promo_id"`
}
