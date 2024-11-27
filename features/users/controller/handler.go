package controller

import (
	"greenenvironment/constant"
	"greenenvironment/features/users"
	"greenenvironment/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService users.UserServiceInterface
	jwt         helper.JWTInterface
}

func NewUserController(u users.UserServiceInterface, j helper.JWTInterface) users.UserControllerInterface {
	return &UserHandler{
		userService: u,
		jwt:         j,
	}
}

func (h *UserHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var UserRegisterRequest UserRegisterRequest

		err := c.Bind(&UserRegisterRequest)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		user := users.User{
			Name:     UserRegisterRequest.Name,
			Email:    UserRegisterRequest.Email,
			Password: UserRegisterRequest.Password,
		}

		createdUser, err := h.userService.Register(user)
		if err != nil {
			return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, constant.UserSuccessRegister, createdUser))
	}
}

func (h *UserHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var UserLoginRequest UserLoginRequest

		err := c.Bind(&UserLoginRequest)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		user := users.User{
			Email:    UserLoginRequest.Email,
			Password: UserLoginRequest.Password,
		}

		userLogin, err := h.userService.Login(user)
		if err != nil {
			return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var response UserLoginResponse
		response.Token = userLogin.Token
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.UserSuccessLogin, response))
	}
}
