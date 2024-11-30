package controller

import (
	"greenenvironment/constant"
	"greenenvironment/features/impacts"
	"greenenvironment/helper"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ImpactController struct {
	impactService impacts.ImpactServiceInterface
	jwtService    helper.JWTInterface
}

func NewImpactController(i impacts.ImpactServiceInterface, j helper.JWTInterface) impacts.ImpactControllerInterface {
	return &ImpactController{
		impactService: i,
		jwtService:    j,
	}
}

// Get All Impacts
// @Summary      Retrieve all impacts
// @Description  Get a list of all impact categories
// @Tags         Impact
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Success      200  {object}  helper.Response{data=[]ImpactCategoryResponse} "Retrieve impacts successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /impacts [get]
func (ic *ImpactController) GetAll(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}
	_, err := ic.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	impactCategories, err := ic.impactService.GetAll()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	var impactCategoryResponses []ImpactCategoryResponse
	for _, impactCategory := range impactCategories {
		impactCategoryResponses = append(impactCategoryResponses, ImpactCategoryResponse{
			ID:          impactCategory.ID,
			Name:        impactCategory.Name,
			ImpactPoint: impactCategory.ImpactPoint,
			Description: impactCategory.Description,
		})
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "get all impacts successfully", []interface{}{impactCategoryResponses}))
}

// Get Impact By ID
// @Summary      Retrieve impact by ID
// @Description  Get a specific impact category by its ID
// @Tags         Impact
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Impact ID"
// @Success      200  {object}  helper.Response{data=ImpactCategoryResponse} "Retrieve impact successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      404  {object}  helper.Response{data=string} "Impact not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /impacts/{id} [get]
func (ic *ImpactController) GetByID(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}
	_, err := ic.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	impactId := c.Param("id")
	impactCategory, err := ic.impactService.GetByID(impactId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ObjectFormatResponse(false, err.Error(), nil))
	}

	impactCategoryResponse := ImpactCategoryResponse{
		ID:          impactCategory.ID,
		Name:        impactCategory.Name,
		ImpactPoint: impactCategory.ImpactPoint,
		Description: impactCategory.Description,
	}

	return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "get impact successfully", impactCategoryResponse))
}

// Create Impact
// @Summary      Create a new impact
// @Description  Add a new impact category to the system
// @Tags         Impact
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        request        body      CreateImpactRequest  true  "Create impact payload"
// @Success      201  {object}  helper.Response{data=string} "Impact created successfully"
// @Failure      400  {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /impacts [post]
func (ic *ImpactController) Create(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}
	_, err := ic.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	var impactRequest CreateImpactRequest
	err = c.Bind(&impactRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	if err := c.Validate(impactRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	impact := impacts.ImpactCategory{
		ID:          uuid.New().String(),
		Name:        impactRequest.Name,
		ImpactPoint: impactRequest.ImpactPoint,
		Description: impactRequest.Description,
	}

	err = ic.impactService.Create(impact)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.ObjectFormatResponse(true, "create impact successfully", nil))
}

// Delete Impact
// @Summary      Delete an impact
// @Description  Remove an impact category by its ID
// @Tags         Impact
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Impact ID"
// @Success      200  {object}  helper.Response{data=string} "Impact deleted successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      403  {object}  helper.Response{data=string} "Forbidden"
// @Failure      404  {object}  helper.Response{data=string} "Impact not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /impacts/{id} [delete]
func (ic *ImpactController) Delete(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}
	token, err := ic.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}
	adminData := ic.jwtService.ExtractUserToken(token)
	role := adminData[constant.JWT_ROLE]
	if role != constant.RoleAdmin {
		helper.UnauthorizedError(c)
	}

	impactId := c.Param("id")
	var impactCategory impacts.ImpactCategory
	impactCategory.ID = impactId
	err = ic.impactService.Delete(impactCategory)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "delete impact successfully", nil))
}
