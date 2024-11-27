package admin

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Admin struct {
	ID        string
	Name      string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AdminLogin struct {
	Email    string
	Password string
	Token    string
}

type AdminUpdate struct {
	ID        string
	Name      string
	Username  string
	Email     string
	Password  string
	UpdatedAt time.Time
	Token     string
}

type AdminRepositoryInterface interface {
	Login(Admin) (Admin, error)
	Update(AdminUpdate) (Admin, error)
	Delete(Admin) error

	GetAdminByID(id string) (Admin, error)
	IsEmailExist(email string) bool
}

type AdminControllerInterface interface {
	Login(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	GetAdminData(c echo.Context) error
}

type AdminServiceInterface interface {
	Login(Admin) (AdminLogin, error)
	Update(AdminUpdate) (AdminUpdate, error)
	Delete(Admin) error

	GetAdminData(Admin) (Admin, error)
}
