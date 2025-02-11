package dto

type AddToCartInput struct {
	ProductId uint `json:"product_id"`
	UserId    uint `json:"user_id"`
	Quantity  int  `json:"quantity"`
}

type RemoveFromCartInput struct {
	ProductId uint `json:"product_id"`
	UserId    uint `json:"user_id"`
}
