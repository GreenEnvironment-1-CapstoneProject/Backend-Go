package guest

import (
	"greenenvironment/features/products"

	"github.com/labstack/echo/v4"
)

type Guest struct {
	Products []products.Product
}

type GuestController interface {
	GetGuestProduct(c echo.Context) error
}
type GuestService interface {
	GetGuest() (Guest, error)
}
type GuestRepostory interface {
	GetGuest() (Guest, error)
}
