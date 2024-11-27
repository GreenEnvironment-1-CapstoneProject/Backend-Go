package controller

import (
	"greenenvironment/constant"
	"greenenvironment/features/admin"
	"greenenvironment/helper"
	"net/http"

	"github.com/labstack/echo"
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
func (h *AdminController) Delete(c echo.Context) error {
	admin := c.Get("admin").(admin.Admin)
	err := h.s.Delete(admin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))

	}
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "delete admin successfully", nil))

}

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
