package routes

import (
	"greenenvironment/configs"
	"greenenvironment/constant/route"
	"greenenvironment/features/admin"
	"greenenvironment/features/impacts"
	"greenenvironment/features/products"
	"greenenvironment/features/users"
	"greenenvironment/helper"
	"greenenvironment/utils/storages"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Echo, uh users.UserControllerInterface, cfg configs.GEConfig) {
	e.POST(route.UserRegister, uh.Register)
	e.POST(route.UserLogin, uh.Login)

	e.GET(route.UserLoginGoogle, uh.GoogleLogin)
	e.GET(route.UserGoogleCallback, uh.GoogleCallback)

	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.GET(route.UserPath, uh.GetUserData, echojwt.WithConfig(jwtConfig))
	e.PUT(route.UserPath, uh.Update, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.UserPath, uh.Delete, echojwt.WithConfig(jwtConfig))

	// Admin
	e.GET(route.AdminManageUserPath, uh.GetAllUsersForAdmin, echojwt.WithConfig(jwtConfig))
	e.GET(route.AdminManageUserByID, uh.GetUserByIDForAdmin, echojwt.WithConfig(jwtConfig))
	e.PUT(route.AdminManageUserByID, uh.UpdateUserForAdmin, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.AdminManageUserByID, uh.DeleteUserForAdmin, echojwt.WithConfig(jwtConfig))
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

func RoutesProducts(e *echo.Echo, ph products.ProductControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}
	e.POST(route.ProductPath, ph.Create, echojwt.WithConfig(jwtConfig))
	e.GET(route.ProductPath, ph.GetAll, echojwt.WithConfig(jwtConfig))
	e.GET(route.ProductByID, ph.GetById, echojwt.WithConfig(jwtConfig))
	e.GET(route.CategoryProduct, ph.GetByCategory, echojwt.WithConfig(jwtConfig))
	e.PUT(route.ProductByID, ph.Update, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.ProductByID, ph.Delete, echojwt.WithConfig(jwtConfig))
}

func RouteImpacts(e *echo.Echo, ic impacts.ImpactControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST(route.ImpactCategoryPath, ic.Create, echojwt.WithConfig(jwtConfig))
	e.GET(route.ImpactCategoryPath, ic.GetAll, echojwt.WithConfig(jwtConfig))
	e.GET(route.ImpactCategoryByID, ic.GetByID, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.ImpactCategoryByID, ic.Delete, echojwt.WithConfig(jwtConfig))
}

func RouteStorage(e *echo.Echo, sc storages.StorageInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST("/api/v1/media/upload", sc.UploadFileHandler, echojwt.WithConfig(jwtConfig))
}
