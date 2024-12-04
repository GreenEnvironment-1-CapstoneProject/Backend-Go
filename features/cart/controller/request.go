package controller

type CreateCartRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity"`
}

type UpdateCartRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=increment decrement qty"`
	Quantity  int    `json:"quantity"`
}
