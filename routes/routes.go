package routes

import (
	"greenenvironment/constant/route"
	"greenenvironment/features/users"

	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Echo, uh users.UserControllerInterface) {
	e.POST(route.UserRegister, uh.Register())
	e.POST(route.UserLogin, uh.Login())
}
