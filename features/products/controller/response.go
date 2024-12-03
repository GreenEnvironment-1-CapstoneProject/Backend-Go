package controller

import (
	"greenenvironment/features/products"
)

type ProductResponse struct {
	ID              string                  `json:"product_id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	Price           float64                 `json:"price"`
	Coin            int                     `json:"coin"`
	Stock           int                     `json:"stock"`
	CreatedAt       string                  `json:"created_at"`
	UpdatedAt       string                  `json:"updated_at"`
	CategoryProduct string                  `json:"category_product"`
	CategoryImpact  []ProductImpactCategory `json:"category_impact"`
	Images          []ProductImage          `json:"images"`
}

type ProductImpactCategory struct {
	ImpactCategory ImpactCategory `json:"impact_category"`
}

type ProductImage struct {
	ImageURL string `json:"image_url"`
}

type ImpactCategory struct {
	Name        string `json:"name"`
	ImpactPoint int    `json:"impact_point"`
	Description string `json:"description"`
}

func (p ProductResponse) ToResponse(product products.Product) ProductResponse {
	images := make([]ProductImage, len(product.Images))
	for i, image := range product.Images {
		images[i] = ProductImage{
			ImageURL: image.AlbumsURL,
		}
	}
	impactCategories := make([]ProductImpactCategory, len(product.ImpactCategories))
	for i, impactCategory := range product.ImpactCategories {
		impactCategories[i] = ProductImpactCategory{
			ImpactCategory: ImpactCategory{
				Name:        impactCategory.ImpactCategory.Name,
				ImpactPoint: impactCategory.ImpactCategory.ImpactPoint,
			},
		}
	}
	return ProductResponse{
		ID:              product.ID,
		Name:            product.Name,
		Description:     product.Description,
		Price:           product.Price,
		Coin:            product.Coin,
		Stock:           product.Stock,
		CategoryProduct: product.Category,
		CreatedAt:       product.CreatedAt.Format("02/01/2006"),
		UpdatedAt:       product.UpdatedAt.Format("02/01/2006"),
		Images:          images,
		CategoryImpact:  impactCategories,
	}
}
