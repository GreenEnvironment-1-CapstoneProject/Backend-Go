package controller

import (
	"context"
	"encoding/json"
	"greenenvironment/constant"
	"greenenvironment/features/users"
	"greenenvironment/helper"
	"greenenvironment/utils/google"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
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
// @Router       /users/register [post]
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
		ID:            createdUser.ID,
		Username:      createdUser.Username,
		Name:          createdUser.Name,
		Email:         createdUser.Email,
		Address:       createdUser.Address,
		Gender:        createdUser.Gender,
		Phone:         createdUser.Phone,
		Exp:           createdUser.Exp,
		Coin:          createdUser.Coin,
		AvatarURL:     createdUser.AvatarURL,
		Is_Membership: createdUser.Is_Membership,
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
// @Router       /users/login [post]
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

// Update User
// @Summary      Update user data
// @Description  Update the authenticated user's information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                   true  "Bearer token"
// @Param        request        body      controller.UserUpdateRequest  true  "User update payload"
// @Success      200            {object}  helper.Response{data=UserLoginResponse}
// @Failure      400            {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/update [put]
func (h *UserHandler) Update(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	userId := userData[constant.JWT_ID]

	var UserUpdateRequest UserUpdateRequest
	err = c.Bind(&UserUpdateRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrUpdateUser.Error(), nil))
	}

	if err := c.Validate(&UserUpdateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	currentUser, err := h.userService.GetUserByIDForAdmin(userId.(string))
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	if UserUpdateRequest.Email == "" {
		UserUpdateRequest.Email = currentUser.Email
	}
	if UserUpdateRequest.Username == "" {
		UserUpdateRequest.Username = currentUser.Username
	}
	if UserUpdateRequest.AvatarURL == "" {
		UserUpdateRequest.AvatarURL = currentUser.AvatarURL
	}

	user := users.UserUpdate{
		ID:        userId.(string),
		Username:  UserUpdateRequest.Username,
		Password:  UserUpdateRequest.Password,
		Name:      UserUpdateRequest.Name,
		Email:     UserUpdateRequest.Email,
		Address:   UserUpdateRequest.Address,
		Gender:    UserUpdateRequest.Gender,
		Phone:     UserUpdateRequest.Phone,
		AvatarURL: UserUpdateRequest.AvatarURL,
	}

	FromUserService, err := h.userService.Update(user)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	var UserToken UserLoginResponse
	UserToken.Token = FromUserService.Token
	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.UserSuccessUpdate, UserToken))
}

// Get User Data
// @Summary      Get user data
// @Description  Retrieve the authenticated user's profile information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                   true  "Bearer token"
// @Success      200            {object}  helper.Response{data=UserInfoResponse}
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/profile [get]
func (h *UserHandler) GetUserData(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	userId := userData[constant.JWT_ID]
	var user users.User
	user.ID = userId.(string)

	user, err = h.userService.GetUserData(user)

	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	var response UserInfoResponse
	response.ID = user.ID
	response.Name = user.Name
	response.Email = user.Email
	response.Username = user.Username
	response.Address = user.Address
	response.Gender = user.Gender
	response.Phone = user.Phone
	response.Coin = user.Coin
	response.Exp = user.Exp
	response.Is_Membership = user.Is_Membership
	response.AvatarURL = user.AvatarURL
	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.UserSuccessGetUser, response))
}

// Delete User
// @Summary      Delete user account
// @Description  Delete the authenticated user's account
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer token"
// @Success      200            {object}  helper.Response{data=string}
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/delete [delete]
func (h *UserHandler) Delete(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	userId := userData[constant.JWT_ID]
	var user users.User
	user.ID = userId.(string)

	err = h.userService.Delete(user)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.UserSuccessDelete, nil))
}

func (h *UserHandler) GoogleLogin(c echo.Context) error {
	url := google.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *UserHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "No code provided", nil))
	}

	token, err := google.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to exchange token", nil))
	}

	client := google.GoogleOauthConfig.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to get user info", nil))
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to parse user info", nil))
	}

	user := users.User{
		Name:  userInfo.Name,
		Email: userInfo.Email,
	}
	createdUser, err := h.userService.RegisterOrLoginGoogle(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to process user", nil))
	}

	tokenString, err := h.jwt.GenerateUserJWT(helper.UserJWT{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to generate token", nil))
	}

	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "Login successful", map[string]string{"token": tokenString}))
}
