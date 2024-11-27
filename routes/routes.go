package routes

import (
	"greenenvironment/configs"
	"greenenvironment/constant/route"
	"greenenvironment/features/admin"
	"greenenvironment/features/users"
	"greenenvironment/helper"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Echo, uh users.UserControllerInterface) {
	e.POST(route.UserRegister, uh.Register)
	e.POST(route.UserLogin, uh.Login)
}

func RouteAdmin(e *echo.Echo, ah admin.AdminControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST(route.AdminLogin, ah.Login)

	e.GET(route.AdminPath, ah.GetAdminData, echojwt.WithConfig(jwtConfig))
	e.PUT(route.AdminPath, ah.Update, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.AdminPath, ah.Delete, echojwt.WithConfig(jwtConfig))
}
