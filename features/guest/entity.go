package guest

import (
	"greenenvironment/features/products"

	"github.com/labstack/echo/v4"
)

type Guest struct {
	TotalProduct             int
	TotalNewProductThisMonth int
	NewProduct               []products.Product
}

type GuestControllerInterface interface {
	GetGuestProduct(c echo.Context) error
}
type GuestServiceInterface interface {
	GetGuestProduct() (Guest, error)
}
type GuestRepostoryInterface interface {
	GetGuests() (Guest, error)
	GetTotalProduct() (int, error)
	GetTotalNewProductThisMonth() (int, error)
	GetNewProduct() ([]products.Product, error)
}
