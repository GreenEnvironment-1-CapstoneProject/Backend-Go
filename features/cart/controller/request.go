package controller

type CreateCartRequest struct {
	ProductID string `json:"product_id" validate:"required"`
}

type UpdateCartRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=increment decrement"`
}
