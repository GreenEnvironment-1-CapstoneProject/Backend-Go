package service

import reviewproducts "greenenvironment/features/review_products"

type ReviewProductService struct {
	reviewRepo reviewproducts.ReviewProductRepositoryInterface
}

func NewReviewProductService(rp reviewproducts.ReviewProductRepositoryInterface) reviewproducts.ReviewProductServiceInterface {
	return &ReviewProductService{
		reviewRepo: rp,
	}
}

func (rps *ReviewProductService) Create(createDto reviewproducts.CreateReviewProduct) error {
	return rps.reviewRepo.Create(createDto)
}
func (rps *ReviewProductService) GetProductReview(productID string) ([]reviewproducts.ReviewProduct, error) {
	return rps.reviewRepo.GetProductReview(productID)
}
