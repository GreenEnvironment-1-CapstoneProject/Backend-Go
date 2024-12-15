package repository

import (
	productModel "greenenvironment/features/products/repository"
	reviewproducts "greenenvironment/features/review_products"
	userModel "greenenvironment/features/users/repository"

	"gorm.io/gorm"
)

type ReviewProduct struct {
	*gorm.Model
	ID        string               `gorm:"primary_key;type:varchar(50);column:id"`
	UserID    string               `gorm:"not null;type:varchar(50);column:user_id"`
	ProductID string               `gorm:"not null;type:varchar(50);column:product_id"`
	Review    string               `gorm:"not null;type:text;column:review"`
	Rate      int                  `gorm:"not null;column:rate"`
	Product   productModel.Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User      userModel.User       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (c *ReviewProduct) TableName() string {
	return "review_products"
}

func (c *ReviewProduct) ToEntity() reviewproducts.ReviewProduct {
	return reviewproducts.ReviewProduct{
		Name:      c.User.Name,
		Email:     c.User.Email,
		Review:    c.Review,
		Rate:      c.Rate,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
