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

// Register User
// @Summary      Register a new user
// @Description  Create a new user account in the system
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      controller.UserRegisterRequest  true  "User registration payload"
// @Success      201      {object}  helper.Response{data=UserRegisterResponse}
// @Failure      400      {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      500      {object}  helper.Response{data=string} "Internal server error"
// @Router       /register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var UserRegisterRequest UserRegisterRequest

	err := c.Bind(&UserRegisterRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	if err := c.Validate(UserRegisterRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	user := users.User{
		Name:     UserRegisterRequest.Name,
		Email:    UserRegisterRequest.Email,
		Password: UserRegisterRequest.Password,
	}

	createdUser, err := h.userService.Register(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	userResponse := UserRegisterResponse{
		ID:           createdUser.ID,
		Username:     createdUser.Username,
		Name:         createdUser.Name,
		Email:        createdUser.Email,
		Address:      createdUser.Address,
		Gender:       createdUser.Gender,
		Phone:        createdUser.Phone,
		Exp:          createdUser.Exp,
		Coin:         createdUser.Coin,
		AvatarURL:    createdUser.AvatarURL,
		IsMembership: createdUser.IsMembership,
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, constant.UserSuccessRegister, userResponse))
}

// Login User
// @Summary      User login
// @Description  Authenticate user and generate JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      controller.UserLoginRequest  true  "User login payload"
// @Success      200      {object}  helper.Response{data=UserLoginResponse}
// @Failure      400      {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      500      {object}  helper.Response{data=string} "Internal server error"
// @Router       /login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var UserLoginRequest UserLoginRequest

	err := c.Bind(&UserLoginRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	if err := c.Validate(UserLoginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	user := users.User{
		Email:    UserLoginRequest.Email,
		Password: UserLoginRequest.Password,
	}

	userLogin, err := h.userService.Login(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	var response UserLoginResponse
	response.Token = userLogin.Token
	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.UserSuccessLogin, response))
}
