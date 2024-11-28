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
	Is_Membership bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserLogin struct {
	Email    string
	Password string
	Token    string
}

type UserUpdate struct {
	ID        string
	Username  string
	Password  string
	Name      string
	Email     string
	Address   string
	Gender    string
	Phone     string
	AvatarURL string
	Token     string
}

type UserRepoInterface interface {
	Register(User) (User, error)
	Login(User) (User, error)
	Update(UserUpdate) (User, error)
	Delete(User) error
	GetUserByID(id string) (User, error)

	// Admin
	GetUserByIDForAdmin(id string) (User, error)
}

type UserServiceInterface interface {
	Register(User) (User, error)
	Login(User) (UserLogin, error)
	Update(UserUpdate) (UserUpdate, error)
	GetUserData(User) (User, error)
	Delete(User) error

	// Admin
	GetUserByIDForAdmin(id string) (User, error)
}

type UserControllerInterface interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	Update(c echo.Context) error
	GetUserData(c echo.Context) error
	Delete(c echo.Context) error
}
