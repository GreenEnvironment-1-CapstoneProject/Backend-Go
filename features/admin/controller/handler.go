package controller

import (
	"greenenvironment/constant"
	"greenenvironment/features/admin"
	"greenenvironment/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AdminController struct {
	s admin.AdminServiceInterface
	j helper.JWTInterface
}

func NewAdminController(u admin.AdminServiceInterface, j helper.JWTInterface) admin.AdminControllerInterface {
	return &AdminController{
		s: u,
		j: j,
	}
}

// Login Admin
// @Summary      Admin login
// @Description  Authenticate admin and generate JWT token
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        request  body      controller.AdminLoginRequest  true  "Admin login payload"
// @Success      200      {object}  helper.Response{data=AdminLoginResponse} "Login successful"
// @Failure      400      {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      500      {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/login [post]
func (h *AdminController) Login(c echo.Context) error {

	var AdminLoginRequest AdminLoginRequest

	err := c.Bind(&AdminLoginRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	if err := c.Validate(AdminLoginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	admin := admin.Admin{
		Email:    AdminLoginRequest.Email,
		Password: AdminLoginRequest.Password,
	}

	adminLogin, err := h.s.Login(admin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	var response AdminLoginResponse
	response.Token = adminLogin.Token
	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "login successfully", response))

}

// Update Admin
// @Summary      Update admin profile
// @Description  Update admin details such as name, username, email, or password
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                         true  "Bearer token"
// @Param        request        body      controller.AdminUpdateRequest  true  "Admin update payload"
// @Success      200            {object}  helper.Response{data=string}   "Update successful"
// @Failure      400            {object}  helper.Response{data=string}   "Invalid input or validation error"
// @Failure      401            {object}  helper.Response{data=string}   "Unauthorized"
// @Failure      500            {object}  helper.Response{data=string}   "Internal server error"
// @Router       /admin [put]
func (h *AdminController) Update(c echo.Context) error {

	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "error unathorized", nil))
	}

	token, err := h.j.ValidateToken(tokenString)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "error unathorized", nil))
	}

	adminData := h.j.ExtractAdminToken(token)
	adminId := adminData[constant.JWT_ID]

	var AdminUpdateRequest AdminUpdateRequest
	err = c.Bind(&AdminUpdateRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	if err := c.Validate(AdminUpdateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	admin := admin.AdminUpdate{
		ID:       adminId.(string),
		Username: AdminUpdateRequest.Username,
		Name:     AdminUpdateRequest.Name,
		Email:    AdminUpdateRequest.Email,
		Password: AdminUpdateRequest.Password,
	}
	_, err = h.s.Update(admin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "update admin successfully", nil))

}

// Delete Admin
// @Summary      Delete admin account
// @Description  Remove an admin account from the system
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                        true  "Bearer token"
// @Success      200            {object}  helper.Response{data=string}  "Delete successful"
// @Failure      500            {object}  helper.Response{data=string}  "Internal server error"
// @Router       /admin [delete]
func (h *AdminController) Delete(c echo.Context) error {
	admin := c.Get("admin").(admin.Admin)
	err := h.s.Delete(admin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))

	}
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "delete admin successfully", nil))

}

// Get Admin Data
// @Summary      Retrieve admin details
// @Description  Get admin details based on the JWT token provided
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer token"
// @Success      200            {object}  helper.Response{data=AdminInfoResponse} "Admin data retrieved successfully"
// @Failure      401            {object}  helper.Response{data=string}    "Unauthorized"
// @Failure      500            {object}  helper.Response{data=string}    "Internal server error"
// @Router       /admin [get]
func (h *AdminController) GetAdminData(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "error unathorized", nil))
	}

	token, err := h.j.ValidateToken(tokenString)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "error unathorized", nil))
	}

	adminData := h.j.ExtractAdminToken(token)
	userId := adminData[constant.JWT_ID]
	var admin admin.Admin
	admin.ID = userId.(string)

	admin, err = h.s.GetAdminData(admin)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	var response AdminInfoResponse
	response.ID = admin.ID
	response.Name = admin.Name
	response.Email = admin.Email
	response.Username = admin.Username
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "get admin data successfully", response))
}
