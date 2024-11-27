package users

import (
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID           string
	Username     string
	Password     string
	Name         string
	Email        string
	Address      string
	Gender       string
	Phone        string
	Exp          int
	Coin         int
	AvatarURL    string
	IsMembership bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserLogin struct {
	Email    string
	Password string
	Token    string
}

type UserRepoInterface interface {
	Register(User) (User, error)
	Login(User) (User, error)
}

type UserServiceInterface interface {
	Register(User) (User, error)
	Login(User) (UserLogin, error)
}

type UserControllerInterface interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}
