package controller

import (
	"greenenvironment/constant"
	reviewproducts "greenenvironment/features/review_products"
	"greenenvironment/helper"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReviewProductController struct {
	reviewService reviewproducts.ReviewProductServiceInterface
	jwtService    helper.JWTInterface
}

func NewReviewProductController(rs reviewproducts.ReviewProductServiceInterface, js helper.JWTInterface) reviewproducts.ReviewProductControllerInterface {
	return &ReviewProductController{
		reviewService: rs,
		jwtService:    js,
	}
}

// Create Review
// @Summary      Create a new product review
// @Description  Add a review for a specific product.
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Request payload for creating a review"
// @Success      201   {object}  helper.Response{data=string} "Review created successfully"
// @Failure      400   {object}  helper.Response{data=string} "Bad request or validation error"
// @Failure      401   {object}  helper.Response{data=string} "Unauthorized access"
// @Failure      500   {object}  helper.Response{data=string} "Internal server error"
// @Router       /reviews [post]
func (rpc *ReviewProductController) Create(c echo.Context) error {
	var reviewRequest CreateRequest
	if err := c.Bind(&reviewRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}
	if err := c.Validate(reviewRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}
	token, err := rpc.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}
	userData := rpc.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID]
	role := userData[constant.JWT_ROLE]
	if role != constant.RoleUser {
		return helper.UnauthorizedError(c)
	}

	newReview := reviewproducts.CreateReviewProduct{
		UserID:    userId.(string),
		ProductID: reviewRequest.ProductID,
		Review:    reviewRequest.Review,
		Rate:      reviewRequest.Rate,
	}

	err = rpc.reviewService.Create(newReview)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "create review successfully", nil))
}

// Get Product Reviews
// @Summary      Retrieve reviews for a specific product
// @Description  Get a list of reviews for a product by its ID.
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Param        id   path      string  true   "Product ID"
// @Success      200  {object}  helper.Response{data=[]ResponseReviewProduct} "Reviews retrieved successfully"
// @Failure      400  {object}  helper.Response{data=string} "Bad request or invalid ID format"
// @Failure      404  {object}  helper.Response{data=string} "Product not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /reviews/products/{id} [get]
func (rpc *ReviewProductController) GetProductReview(c echo.Context) error {
	paramId := c.Param("id")
	productId, err := uuid.Parse(paramId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	reviews, err := rpc.reviewService.GetProductReview(productId.String())

	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}
	var response []interface{}
	for _, review := range reviews {
		response = append(response, new(ResponseReviewProduct).ToResponse(review))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "get reviews successfully", response))
}
