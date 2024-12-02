package controller

import (
	"context"
	"encoding/json"
	"greenenvironment/constant"
	"greenenvironment/features/users"
	"greenenvironment/helper"
	"greenenvironment/utils/google"
	"net/http"
	"strconv"
	"time"

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

// Google Login
// @Summary      Google Login
// @Description  Redirect to Google's OAuth 2.0 authentication page
// @Tags         Users
// @Produce      json
// @Success      302  {string}  string  "Redirect to Google OAuth"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/login-google [get]
func (h *UserHandler) GoogleLogin(c echo.Context) error {
	url := google.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// Google Callback
// @Summary      Google OAuth Callback
// @Description  Handle the OAuth 2.0 callback from Google and authenticate the user
// @Tags         Users
// @Produce      json
// @Param        code  query     string  true  "Authorization code from Google"
// @Success      200   {object}  helper.Response{data=UserLoginResponse} "Login successful with JWT token"
// @Failure      400   {object}  helper.Response{data=string} "Invalid request or missing code"
// @Failure      500   {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/google-callback [get]
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
		ID:       createdUser.ID,
		Name:     createdUser.Name,
		Email:    createdUser.Email,
		Username: createdUser.Username,
		Role:     constant.RoleUser,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to generate token", nil))
	}

	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.UserSuccessLogin, map[string]string{"token": tokenString}))
}

// Admin
func (h *UserHandler) GetAllUsersForAdmin(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}

		token, err := h.jwt.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.jwt.ExtractAdminToken(token)
		role := adminData[constant.JWT_ROLE]

		if role != constant.RoleAdmin {
			helper.UnauthorizedError(c)
		}

		pageStr := c.QueryParam("page")
		page := 1
		if pageStr != "" {
			page, err = strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrPageInvalid.Error(), nil))
			}
		}
		var totalPages int
		var user []users.User
		user, totalPages, err = h.userService.GetAllByPageForAdmin(page)

		metadata := MetadataResponse{
			CurrentPage: page,
			TotalPage:   totalPages,
		}

		if err != nil {
			code, message := helper.HandleEchoError(err)
			return c.JSON(code, helper.FormatResponse(false, message, nil))
		}

		var response []UserbyAdminandPageResponse
		for _, user := range user {
			response = append(response, UserbyAdminandPageResponse{
				ID:            user.ID,
				Name:          user.Name,
				Email:         user.Email,
				Username:      user.Username,
				Address:       user.Address,
				Gender:        user.Gender,
				Phone:         user.Phone,
				Is_Membership: user.Is_Membership,
				AvatarURL:     user.AvatarURL,

				CreatedAt: user.CreatedAt.Format("02/01/06"),
				UpdatedAt: user.UpdatedAt.Format("02/01/06"),
			})
		}
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.AdminSuccessGetAllUser, metadata, response))
}

func (h *UserHandler) GetUserByIDForAdmin(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}

		token, err := h.jwt.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.jwt.ExtractAdminToken(token)
		role := adminData[constant.JWT_ROLE]

		if role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}

		userId := c.Param("id")
		if err != nil {
			code, message := helper.HandleEchoError(err)
			return c.JSON(code, helper.FormatResponse(false, message, nil))
		}

		users, err := h.userService.GetUserByIDForAdmin(userId)
		if err != nil {
			return c.JSON(http.StatusNotFound, helper.ObjectFormatResponse(false, constant.ErrUserIDNotFound.Error(), nil))
		}

		response := UserbyAdminResponse{
			ID:            users.ID,
			Name:          users.Name,
			Email:         users.Email,
			Username:      users.Username,
			Address:       users.Address,
			Gender:        users.Gender,
			Phone:         users.Phone,
			AvatarURL:     users.AvatarURL,
			Is_Membership: users.Is_Membership,
			CreatedAt:     users.CreatedAt.Format("02/01/06"),
			UpdatedAt:     users.UpdatedAt.Format("02/01/06"),
		}

		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.AdminSuccessGetUser, response))
}

func (h *UserHandler) UpdateUserForAdmin(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}

		token, err := h.jwt.ValidateToken(tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}

		adminData := h.jwt.ExtractAdminToken(token)
		role := adminData[constant.JWT_ROLE]

		if role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		_, err = h.userService.GetUserByIDForAdmin(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, helper.FormatResponse(false, string(constant.ErrUserIDNotFound.Error()), nil))
		}

		var userEdit UserbyAdminRequest
		if err := c.Bind(&userEdit); err != nil {
			code, message := helper.HandleEchoError(err)
			return c.JSON(code, helper.FormatResponse(false, message, nil))
		}

		if err := c.Validate(&userEdit); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
		}

		response := users.UpdateUserByAdmin{
			ID:       id,
			Name:     userEdit.Name,
			Address:  userEdit.Address,
			Gender:   userEdit.Gender,
			Phone:    userEdit.Phone,
			UpdateAt: time.Now(),
		}

		if err := h.userService.UpdateUserForAdmin(response); err != nil {
			return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.AdminSuccessUpdateUser, nil))
}

func (h *UserHandler) DeleteUserForAdmin(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}

		token, err := h.jwt.ValidateToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, constant.Unauthorized, nil))
		}

		adminData := h.jwt.ExtractAdminToken(token)
		role := adminData[constant.JWT_ROLE]

		if role != constant.RoleAdmin {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, constant.Unauthorized, nil))
		}

		id := c.Param("id")
		if err := h.userService.DeleteUserForAdmin(id); err != nil {
			return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.AdminSuccessDeleteUser, nil))
}
