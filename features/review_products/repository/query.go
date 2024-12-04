package repository

import (
	reviewproducts "greenenvironment/features/review_products"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewProductRepository struct {
	DB *gorm.DB
}

func NewReviewProductRepository(db *gorm.DB) reviewproducts.ReviewProductRepositoryInterface {
	return &ReviewProductRepository{
		DB: db,
	}
}

func (rpr *ReviewProductRepository) Create(review reviewproducts.CreateReviewProduct) error {
	newReview := &ReviewProduct{
		ID:        uuid.New().String(),
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Review:    review.Review,
		Rate:      review.Rate,
	}

	err := rpr.DB.Create(newReview).Error
	if err != nil {
		return err
	}

	return nil
}

func (rpr *ReviewProductRepository) GetProductReview(productID string) ([]reviewproducts.ReviewProduct, error) {
	var reviews []ReviewProduct
	var result []reviewproducts.ReviewProduct

	err := rpr.DB.Model(&ReviewProduct{}).
		Preload("User").
		Where("product_id = ?", productID).Find(&reviews).Error

	if err != nil {
		return result, err
	}

	if len(reviews) == 0 {
		return []reviewproducts.ReviewProduct{}, nil
	}

	for _, review := range reviews {
		cartEntity := review.ToEntity()
		result = append(result, cartEntity)
	}

	return result, nil
}
