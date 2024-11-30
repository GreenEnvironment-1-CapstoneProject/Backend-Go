package controller

type ProductRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Price       float64  `json:"price" validate:"required,min=0"`
	Coin        int      `json:"coin" validate:"required,min=0"`
	Stock       int      `json:"stock" validate:"required,min=0"`
	Category    []string `json:"category" validate:"required"`
	Images      []string `json:"images" validate:"required"`
}
