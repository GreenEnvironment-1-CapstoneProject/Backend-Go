package controller

type CreateRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Review    string `json:"review" validate:"required"`
	Rate      int    `json:"rate" validate:"required,max=5"`
}
