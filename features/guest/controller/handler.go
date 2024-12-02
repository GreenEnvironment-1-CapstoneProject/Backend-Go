package controller

import (
	"greenenvironment/features/guest"
	"greenenvironment/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GuestController struct {
	guestService guest.GuestServiceInterface
}

func NewGuestController(g guest.GuestServiceInterface) guest.GuestControllerInterface {
	return &GuestController{
		guestService: g,
	}
}

// Get Guest Product
// @Summary      Get Guest Product
// @Description  Retrieve a list of products available for guest users.
// @Tags         Guest
// @Accept       json
// @Produce      json
// @Success      200  {object}  helper.Response{data=GuestResponse} "Successful response with product data"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /guest/products [get]
func (gc *GuestController) GetGuestProduct(c echo.Context) error {
	data, err := gc.guestService.GetGuestProduct()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	response := new(GuestResponse).ToResponse(data)
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", response))
}
