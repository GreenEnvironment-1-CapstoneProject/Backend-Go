package cart

import (
	"greenenvironment/features/products"
	"greenenvironment/features/users"

	"github.com/labstack/echo/v4"
)

type Cart struct {
	User  users.User
	Items []CartItem
}

type NewCart struct {
	UserID    string
	ProductID string
}

type CartItem struct {
	Quantity int
	Product  products.Product
}

type UpdateCart struct {
	UserID    string
	ProductID string
	Type      string
}

type CartRepositoryInterface interface {
	Create(cart NewCart) error
	Update(cart UpdateCart) error
	Delete(userId string, productId string) error
	Get(userId string) (Cart, error)
	IsCartExist(userId string, productId string) (bool, error)
	InsertIncrement(userId string, productId string) error
	InsertDecrement(userId string, productId string) error
	GetCartQty(userId string, productId string) (int, error)
}

type CartServiceInterface interface {
	Create(cart NewCart) error
	Update(cart UpdateCart) error
	Delete(userId string, productId string) error
	Get(userId string) (Cart, error)
}

type CartControllerInterface interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Get(c echo.Context) error
}
