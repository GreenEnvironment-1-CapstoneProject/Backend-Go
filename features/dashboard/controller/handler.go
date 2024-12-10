package controller

import (
	"greenenvironment/constant"
	"greenenvironment/features/dashboard"
	"greenenvironment/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	dashboardService dashboard.DashboardServiceInterface
	jwtService       helper.JWTInterface
}

func NewDashboardController(service dashboard.DashboardServiceInterface, jwt helper.JWTInterface) dashboard.DashboardControllerInterface {
	return &DashboardHandler{
		dashboardService: service,
		jwtService:       jwt,
	}
}

// GetDashboard retrieves dashboard data filtered by time period
// @Summary      Get dashboard data for admin
// @Description  Retrieve dashboard data filtered by weekly, monthly, or yearly
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        filter query string true "Filter value (weekly, monthly, yearly)"
// @Success      200 {object} helper.Response{data=DashboardResponse} "Dashboard data retrieved successfully"
// @Failure      400 {object} helper.Response{data=string} "Invalid filter value"
// @Failure      401 {object} helper.Response{data=string} "Unauthorized"
// @Failure      500 {object} helper.Response{data=string} "Internal server error"
// @Router       /admin/dashboard [get]
func (dc *DashboardHandler) GetDashboard(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	token, err := dc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	adminData := dc.jwtService.ExtractAdminToken(token)
	role := adminData[constant.JWT_ROLE].(string)
	if role != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	filter := c.QueryParam("filter")
	if filter != "weekly" && filter != "monthly" && filter != "yearly" {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid filter value. Use 'weekly', 'monthly', or 'yearly'.", nil))
	}

	data, err := dc.dashboardService.GetDashboardData(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	response := DashboardResponse{}.FromEntity(data)
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Dashboard data retrieved successfully", response))
}
