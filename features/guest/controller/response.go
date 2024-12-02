package controller

import (
	"greenenvironment/features/guest"
	productHandler "greenenvironment/features/products/controller"
)

type GuestResponse struct {
	TotalProduct             int                              `json:"total_product"`
	TotalNewProductThisMonth int                              `json:"total_new_product_this_month"`
	NewProducts              []productHandler.ProductResponse `json:"new_products"`
}

func (d *GuestResponse) ToResponse(data guest.Guest) *GuestResponse {
	newProducts := make([]productHandler.ProductResponse, len(data.NewProduct))
	for i, product := range data.NewProduct {
		newProducts[i] = new(productHandler.ProductResponse).ToResponse(product)
	}

	return &GuestResponse{
		TotalProduct:             data.TotalProduct,
		TotalNewProductThisMonth: data.TotalNewProductThisMonth,
		NewProducts:              newProducts,
	}
}
