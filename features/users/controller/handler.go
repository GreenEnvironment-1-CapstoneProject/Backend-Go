package controller

import (
	"context"
	"encoding/json"
	"greenenvironment/constant"
	"greenenvironment/features/users"
	"greenenvironment/helper"
	"greenenvironment/utils/google"
	"greenenvironment/utils/storages"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type UserHandler struct {
	userService users.UserServiceInterface
	jwt         helper.JWTInterface
	storage     storages.StorageInterface
}

func NewUserController(u users.UserServiceInterface, j helper.JWTInterface, s storages.StorageInterface) users.UserControllerInterface {
	return &UserHandler{
		userService: u,
		jwt:         j,
		storage:     s,
	}
}

// Request OTP for Registration
// @Summary      Request OTP for user registration
// @Description  Sends an OTP to the user's email for registration verification
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      controller.UserRegisterRequest  true  "User registration payload"
// @Success      201      {object}  helper.Response{data=string} "OTP sent successfully"
// @Failure      400      {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      500      {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/register/request-otp [post]
func (h *UserHandler) RequestRegisterOTP(c echo.Context) error {
	var req UserRegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid input", nil))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Validation error", nil))
	}

	err := h.userService.RequestRegisterOTP(req.Name, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "OTP sent successfully", nil))
}

// Verify OTP and Register User
// @Summary      Verify OTP and register user
// @Description  Verifies the OTP sent to the user's email and registers the user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      controller.UserVerifyRegisterRequest  true  "User verification payload"
// @Success      201      {object}  helper.Response{data=UserRegisterResponse} "User registered successfully"
// @Failure      400      {object}  helper.Response{data=string} "Invalid OTP or input"
// @Failure      500      {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/register/verify-otp [post]
func (h *UserHandler) VerifyRegisterOTP(c echo.Context) error {
	var req UserVerifyRegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid input", nil))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Validation error", nil))
	}

	user, err := h.userService.VerifyRegisterOTP(req.OTP)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "User registered successfully", user))
}

// Forgot Password Request
// @Summary      Request OTP for password reset
// @Description  Sends an OTP to the user's email for password reset verification
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      controller.ForgotPasswordRequest  true  "Forgot password request payload"
// @Success      200      {object}  helper.Response{data=string} "OTP sent successfully"
// @Failure      400      {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      404      {object}  helper.Response{data=string} "Email not found"
// @Failure      500      {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/forgot-password [post]
func (h *UserHandler) ForgotPasswordRequest(c echo.Context) error {
	var request ForgotPasswordRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrInvalidInput.Error(), nil))
	}

	if !h.userService.IsEmailExist(request.Email) {
		return c.JSON(http.StatusNotFound, helper.FormatResponse(false, constant.ErrEmailNotFound.Error(), nil))
	}

	err := h.userService.RequestPasswordResetOTP(request.Email)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.UserSuccessForgotPassword, nil))
}

// Verify Forgot Password OTP
// @Summary      Verify OTP for password reset
// @Description  Verifies the OTP sent to the user's email for password reset
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      controller.VerifyOTPRequest  true  "Verify OTP request payload"
// @Success      200      {object}  helper.Response{data=string} "OTP verified successfully"
// @Failure      400      {object}  helper.Response{data=string} "Invalid OTP"
// @Failure      500      {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/verify-otp [post]
func (h *UserHandler) VerifyForgotPasswordOTP(c echo.Context) error {
	var request VerifyOTPRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrInvalidInput.Error(), nil))
	}

	err := h.userService.VerifyPasswordResetOTP(request.OTP)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.UserSuccessOTPValidation, nil))
}

// Reset Password
// @Summary      Reset user password
// @Description  Resets the user's password after verifying OTP
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      controller.ResetPasswordRequest  true  "Reset password request payload"
// @Success      200      {object}  helper.Response{data=string} "Password reset successfully"
// @Failure      400      {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      500      {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/reset-password [put]
func (h *UserHandler) ResetPassword(c echo.Context) error {
	var request ResetPasswordRequest
	if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrInvalidInput.Error(), nil))
	}

	err := h.userService.ResetPassword(request.NewPassword)
	if err != nil {
			return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.UserSuccessResetPassword, nil))
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
// Update User Data
func (h *UserHandler) UpdateUserInfo(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
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

	user := users.UserUpdate{
		ID:      userId.(string),
		Name:    UserUpdateRequest.Name,
		Address: UserUpdateRequest.Address,
		Gender:  UserUpdateRequest.Gender,
		Phone:   UserUpdateRequest.Phone,
	}

	err = h.userService.UpdateUserInfo(user)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.UserSuccessUpdate, nil))
}

// Request Password Update OTP
// @Summary      Request OTP for password update
// @Description  Sends an OTP to the authenticated user's email for password update verification
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer token"
// @Success      201            {object}  helper.Response{data=string} "OTP sent successfully"
// @Failure      400            {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/request-otp [post]
func (h *UserHandler) RequestPasswordUpdateOTP(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	email := userData[constant.JWT_EMAIL]

	err = h.userService.RequestPasswordUpdateOTP(email.(string))
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, constant.UserSuccessForgotPassword, nil))
}

// Update User Password
// @Summary      Update user password
// @Description  Updates the authenticated user's password using OTP and old password verification
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer token"
// @Param        request        body      controller.UserPasswordUpdateRequest  true  "Password update payload"
// @Success      200            {object}  helper.Response{data=string} "Password updated successfully"
// @Failure      400            {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/update-password [put]
func (h *UserHandler) UpdateUserPassword(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	userData := h.jwt.ExtractUserToken(token)
	email := userData[constant.JWT_EMAIL]

	var UserPasswordUpdateRequest UserPasswordUpdateRequest
	err = c.Bind(&UserPasswordUpdateRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrUpdateUser.Error(), nil))
	}

	updateData := users.PasswordUpdate{
		Email:       email.(string),
		OldPassword: UserPasswordUpdateRequest.OldPassword,
		NewPassword: UserPasswordUpdateRequest.NewPassword,
		OTP:         UserPasswordUpdateRequest.OTP,
	}

	err = h.userService.UpdatePassword(updateData)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.UserSuccessUpdatePassword, nil))
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

// Update Avatar
// @Summary      Update User Avatar
// @Description  Upload a new avatar for the authenticated user
// @Tags         Users
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        avatar         formData  file    true  "Avatar image"
// @Success      200            {object}  helper.Response{data=string} "Avatar updated successfully"
// @Failure      400            {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /users/avatar [put]
func (h *UserHandler) UpdateAvatar(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	token, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}
	userData := h.jwt.ExtractUserToken(token)
	userID := userData[constant.JWT_ID].(string)

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Avatar file is required", nil))
	}

	src, err := h.storage.ImageValidation(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	avatarURL, err := h.storage.UploadImageToCloudinary(src, "ecomate/avatars/")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to upload avatar", nil))
	}

	err = h.userService.UpdateAvatar(userID, avatarURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to update avatar", nil))
	}

	// Respon sukses
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Avatar updated successfully", avatarURL))
}

// Admin

// Get All Users
// @Summary      Get all users
// @Description  Retrieve a paginated list of all users (admin access only)
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        page           query     int     false "Page number (default is 1)"
// @Param        limit          query     int     false "Number of items per page (default is 20)"
// @Success      200            {object}  helper.MetadataResponse{data=[]UserbyAdminandPageResponse}
// @Failure      400            {object}  helper.Response{data=string} "Invalid page number"
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/users [get]
func (h *UserHandler) GetAllUsersForAdmin(c echo.Context) error {
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

	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrPageInvalid.Error(), nil))
		}
	}

	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid limit value", nil))
		}
	}

	var totalPages int
	var user []users.User
	user, totalPages, err = h.userService.GetAllByPageForAdmin(page, limit)

	metadata := MetadataResponse{
		CurrentPage: page,
		TotalPage:   totalPages,
	}

	if err != nil {
		code, message := helper.HandleEchoError(err)
		return c.JSON(code, helper.FormatResponse(false, message, nil))
	}

	if len(user) == 0 {
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.AdminSuccessGetAllUser, metadata, []UserbyAdminandPageResponse{}))
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
			CreatedAt:     user.CreatedAt.Format("02/01/06"),
			UpdatedAt:     user.UpdatedAt.Format("02/01/06"),

		})
	}
	return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.AdminSuccessGetAllUser, metadata, response))
}

// Get User by ID
// @Summary      Get user by ID
// @Description  Retrieve user information by ID (admin access only)
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        id             path      string  true  "User ID"
// @Success      200            {object}  helper.Response{data=UserbyAdminResponse}
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      404            {object}  helper.Response{data=string} "User not found"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/users/{id} [get]
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

// Update User
// @Summary      Update user data
// @Description  Update user information (admin access only)
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string             true  "Bearer token"
// @Param        id             path      string             true  "User ID"
// @Param        data           body      UserbyAdminRequest true  "User update data"
// @Success      200            {object}  helper.Response{data=string} "User updated successfully"
// @Failure      400            {object}  helper.Response{data=string} "Bad request"
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      404            {object}  helper.Response{data=string} "User not found"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/users/{id} [put]
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

// Delete User
// @Summary      Delete user
// @Description  Delete user account by ID (admin access only)
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        id             path      string  true  "User ID"
// @Success      200            {object}  helper.Response{data=string} "User deleted successfully"
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      404            {object}  helper.Response{data=string} "User not found"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/users/{id} [delete]
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
