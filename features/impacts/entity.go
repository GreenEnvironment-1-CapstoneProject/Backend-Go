package impacts

import "github.com/labstack/echo/v4"

type ImpactCategory struct {
	ID          string
	Name        string
	ImpactPoint int
	Description string
}

type ImpactServiceInterface interface {
	GetAll() ([]ImpactCategory, error)
	GetByID(ID string) (ImpactCategory, error)
	Create(ImpactCategory) error
	Delete(ImpactCategory) error
}

type ImpactRepositoryInterface interface {
	GetAll() ([]ImpactCategory, error)
	GetByID(ID string) (ImpactCategory, error)
	Create(ImpactCategory) error
	Delete(ImpactCategory) error
}

type ImpactControllerInterface interface {
	GetAll(c echo.Context) error
	GetByID(c echo.Context) error
	Create(c echo.Context) error
	Delete(c echo.Context) error
}
