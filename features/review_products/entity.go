package reviewproducts

import (
	"time"

	"github.com/labstack/echo/v4"
)

type ReviewProduct struct {
	Name      string
	Email     string
	Review    string
	Rate      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateReviewProduct struct {
	UserID    string
	ProductID string
	Review    string
	Rate      int
}

type ReviewProductRepositoryInterface interface {
	Create(createDto CreateReviewProduct) error
	GetProductReview(productID string) ([]ReviewProduct, error)
}
type ReviewProductServiceInterface interface {
	Create(createDto CreateReviewProduct) error
	GetProductReview(productID string) ([]ReviewProduct, error)
}
type ReviewProductControllerInterface interface {
	Create(c echo.Context) error
	GetProductReview(c echo.Context) error
}
