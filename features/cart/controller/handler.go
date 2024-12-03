package controller

import (
	"greenenvironment/constant"
	"greenenvironment/features/cart"
	products "greenenvironment/features/products/controller"
	"greenenvironment/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CartController struct {
	cartService cart.CartServiceInterface
	jwtService  helper.JWTInterface
}

func NewCartController(cartService cart.CartServiceInterface, jwtService helper.JWTInterface) cart.CartControllerInterface {
	return &CartController{
		cartService: cartService,
		jwtService:  jwtService,
	}
}

// @Summary      Create Cart
// @Description  Add a product to the user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        body           body      CreateCartRequest  true  "Request body for creating a cart"
// @Success      201  {object}  helper.Response "create card successfully"
// @Failure      400  {object}  helper.Response "Bad Request"
// @Failure      401  {object}  helper.Response "Unauthorized"
// @Failure      500  {object}  helper.Response "Internal Server Error"
// @Router       /cart [post]
func (cc *CartController) Create(c echo.Context) error {

	var cartRequest CreateCartRequest
	if err := c.Bind(&cartRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}
	if err := c.Validate(cartRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}
	token, err := cc.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}
	userData := cc.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID]
	role := userData[constant.JWT_ROLE]
	if role != constant.RoleUser {
		return helper.UnauthorizedError(c)
	}

	newCart := cart.NewCart{
		ProductID: cartRequest.ProductID,
		UserID:    userId.(string),
	}

	err = cc.cartService.Create(newCart)

	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "create card successfully", nil))
}

// @Summary      Update Cart
// @Description  Update a product's quantity or type in the user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        body           body      UpdateCartRequest  true  "Request body for updating a cart"
// @Success      200  {object}  helper.Response "update card successfully"
// @Failure      400  {object}  helper.Response "Bad Request"
// @Failure      401  {object}  helper.Response "Unauthorized"
// @Failure      500  {object}  helper.Response "Internal Server Error"
// @Router       /cart [put]
func (cc *CartController) Update(c echo.Context) error {
	var cartRequest UpdateCartRequest
	if err := c.Bind(&cartRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}
	if err := c.Validate(cartRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}
	token, err := cc.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}
	userData := cc.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID]
	role := userData[constant.JWT_ROLE]
	if role != constant.RoleUser {
		return helper.UnauthorizedError(c)
	}

	newCart := cart.UpdateCart{
		ProductID: cartRequest.ProductID,
		Type:      cartRequest.Type,
		UserID:    userId.(string),
	}

	err = cc.cartService.Update(newCart)

	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "update card successfully", nil))
}

// @Summary      Delete Cart
// @Description  Remove a product from the user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Product ID"
// @Success      200  {object}  helper.Response "delete card successfully"
// @Failure      400  {object}  helper.Response "Bad Request"
// @Failure      401  {object}  helper.Response "Unauthorized"
// @Failure      500  {object}  helper.Response "Internal Server Error"
// @Router       /cart/{id} [delete]
func (cc *CartController) Delete(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}
	token, err := cc.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}
	userData := cc.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID]
	role := userData[constant.JWT_ROLE]
	if role != constant.RoleUser {
		return helper.UnauthorizedError(c)
	}

	productID := c.Param("id")

	err = cc.cartService.Delete(userId.(string), productID)

	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "delete card successfully", nil))
}

// @Summary      Get Cart
// @Description  Get all items in the user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Success      200  {object}  helper.Response{data=CartResponse} "get cards successfully"
// @Failure      401  {object}  helper.Response "Unauthorized"
// @Failure      500  {object}  helper.Response "Internal Server Error"
// @Router       /cart [get]
func (cc *CartController) Get(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}
	token, err := cc.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}
	userData := cc.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID]

	carts, err := cc.cartService.Get(userId.(string))

	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	var response CartResponse
	response.User = User{
		ID:       carts.User.ID,
		Username: carts.User.Username,
		Email:    carts.User.Email,
		Address:  carts.User.Address,
	}

	for _, cart := range carts.Items {
		var images []products.ProductImage
		var impactCategories []products.ProductImpactCategory

		for _, img := range cart.Product.Images {
			images = append(images, products.ProductImage{
				ImageURL: img.AlbumsURL,
			})
		}

		for _, impact := range cart.Product.ImpactCategories {
			impactCategories = append(impactCategories, products.ProductImpactCategory{
				ImpactCategory: products.ImpactCategory{
					Name:        impact.ImpactCategory.Name,
					ImpactPoint: impact.ImpactCategory.ImpactPoint,
					Description: impact.ImpactCategory.Description,
				},
			})
		}

		response.Items = append(response.Items, CartItems{
			ID:       cart.ID,
			Quantity: cart.Quantity,
			Product: products.ProductResponse{
				ID:              cart.Product.ID,
				Name:            cart.Product.Name,
				Description:     cart.Product.Description,
				Price:           cart.Product.Price,
				Coin:            cart.Product.Coin,
				Stock:           cart.Product.Stock,
				CreatedAt:       cart.Product.CreatedAt.Format("2006-01-02"),
				UpdatedAt:       cart.Product.UpdatedAt.Format("2006-01-02"),
				Images:          images,
				CategoryImpact:  impactCategories,
				CategoryProduct: cart.Product.Category,
			},
		})
	}

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "get cards successfully", response))
}
