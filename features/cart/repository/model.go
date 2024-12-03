package repository

import (
	"greenenvironment/features/cart"
	"greenenvironment/features/impacts"
	"greenenvironment/features/products"
	productModel "greenenvironment/features/products/repository"
	"greenenvironment/features/users"
	userModel "greenenvironment/features/users/repository"

	"gorm.io/gorm"
)

type Cart struct {
	*gorm.Model
	ID                    string                               `gorm:"primary_key;type:varchar(50);column:id"`
	UserID                string                               `gorm:"not null;type:varchar(50);column:user_id"`
	ProductID             string                               `gorm:"not null;type:varchar(50);column:product_id"`
	Quantity              int                                  `gorm:"not null;column:quantity"`
	Product               productModel.Product                 `gorm:"foreignKey:ProductID;references:ID"`
	ProductImage          []productModel.ProductImage          `gorm:"foreignKey:ProductID;references:ID"`
	ProductImpactCategory []productModel.ProductImpactCategory `gorm:"foreignKey:ProductID;references:ID"`
	User                  userModel.User                       `gorm:"foreignKey:UserID;references:ID"`
}

func (c *Cart) TableName() string {
	return "carts"
}

func (c *Cart) ToEntity() cart.Cart {
	var images []products.ProductImage
	var impactCategories []products.ProductImpactCategory

	for _, img := range c.Product.Images {
		images = append(images, products.ProductImage{
			ID:        img.ID,
			ProductID: img.ProductID,
			AlbumsURL: img.AlbumsURL,
		})
	}

	for _, impact := range c.Product.ImpactCategories {
		impactCategories = append(impactCategories, products.ProductImpactCategory{
			ID:               impact.ID,
			ProductID:        impact.ProductID,
			ImpactCategoryID: impact.ImpactCategoryID,
			ImpactCategory: impacts.ImpactCategory{
				ID:          impact.ImpactCategory.ID,
				Name:        impact.ImpactCategory.Name,
				ImpactPoint: impact.ImpactCategory.ImpactPoint,
				Description: impact.ImpactCategory.Description,
			},
		})
	}

	return cart.Cart{
		User: users.User{
			ID:       c.UserID,
			Username: c.User.Username,
			Email:    c.User.Email,
			Address:  c.User.Address,
			Phone:    c.User.Phone,
		},

		Items: []cart.CartItem{
			{
				Quantity: c.Quantity,
				Product: products.Product{
					ID:               c.Product.ID,
					Name:             c.Product.Name,
					Description:      c.Product.Description,
					Price:            c.Product.Price,
					Coin:             c.Product.Coin,
					Stock:            c.Product.Stock,
					CreatedAt:        c.Product.CreatedAt,
					UpdatedAt:        c.Product.UpdatedAt,
					Images:           images,
					ImpactCategories: impactCategories,
				},
			},
		},
	}
}
