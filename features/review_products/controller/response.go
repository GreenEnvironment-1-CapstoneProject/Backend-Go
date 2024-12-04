package controller

import reviewproducts "greenenvironment/features/review_products"

type ResponseReviewProduct struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Review    string `json:"review"`
	Rate      int    `json:"rate"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (r *ResponseReviewProduct) ToResponse(review reviewproducts.ReviewProduct) ResponseReviewProduct {
	return ResponseReviewProduct{
		Name:      review.Name,
		Email:     review.Email,
		Review:    review.Review,
		Rate:      review.Rate,
		CreatedAt: review.CreatedAt.Format("02/01/2006"),
		UpdatedAt: review.UpdatedAt.Format("02/01/2006"),
	}
}
