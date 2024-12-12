package controller

type ProductRequest struct {
	Name            string   `json:"name" validate:"required"`
	Description     string   `json:"description" validate:"required"`
	Price           int      `json:"price" validate:"required,min=0"`
	Coin            int      `json:"coin" validate:"required,min=0"`
	Stock           int      `json:"stock" validate:"required,min=0"`
	CategoryProduct string   `json:"category_product" validate:"required"`
	CategoryImpact  []string `json:"category_impact" validate:"required"`
	Images          []string `json:"images" validate:"required"`
}
