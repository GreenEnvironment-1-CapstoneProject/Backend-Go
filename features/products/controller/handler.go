package controller

import (
	"greenenvironment/constant"
	"greenenvironment/features/products"
	"greenenvironment/helper"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductController struct {
	productService products.ProductServiceInterface
	jwtService     helper.JWTInterface
}

func NewProductController(s products.ProductServiceInterface, j helper.JWTInterface) products.ProductControllerInterface {
	return &ProductController{
		productService: s,
		jwtService:     j,
	}
}

// Create Product
// @Summary      Create a product
// @Description  Create a new product with associated images and categories. Requires admin role.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        body           body      ProductRequest  true   "Product data"
// @Success      201  {object}  helper.Response{data=string} "Product created successfully"
// @Failure      400  {object}  helper.Response{data=string} "Bad request"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /products [post]
func (pc *ProductController) Create(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
	}
	token, err := pc.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
	}

	adminData := pc.jwtService.ExtractUserToken(token)
	role := adminData[constant.JWT_ROLE]
	if role != constant.RoleAdmin {
		helper.UnauthorizedError(c)
	}

	var productInput ProductRequest
	if err := c.Bind(&productInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	if err := c.Validate(productInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	productData := products.Product{
		Name:        productInput.Name,
		Description: productInput.Description,
		Price:       productInput.Price,
		Coin:        productInput.Coin,
		Stock:       productInput.Stock,
		Category:    productInput.CategoryProduct,
	}

	for _, categoryID := range productInput.CategoryImpact {
		productData.ImpactCategories = append(productData.ImpactCategories, products.ProductImpactCategory{
			ProductID:        productData.ID,
			ImpactCategoryID: categoryID,
		})
	}

	for _, imageURL := range productInput.Images {
		productData.Images = append(productData.Images, products.ProductImage{
			ProductID: productData.ID,
			AlbumsURL: imageURL,
		})
	}

	err = pc.productService.Create(productData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "create product successfully", nil))

}

// Get All Products
// @Summary      Get all products
// @Description  Retrieve all products with pagination.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        pages          query     int     false  "Page number"
// @Success      200  {object}  helper.MetadataResponse{data=[]ProductResponse} "Products retrieved successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /products [get]
func (pc *ProductController) GetAll(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("pages"))
	if err != nil {
		page = 1
	}

	products, totalPages, err := pc.productService.GetAllByPage(page)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	var response []interface{}
	for _, product := range products {
		response = append(response, new(ProductResponse).ToResponse(product))
	}

	metadata := map[string]interface{}{
		"TotalPage": totalPages,
		"Page":      page,
	}

	return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "get all products successfully", metadata, response))
}

// Get Product by ID
// @Summary      Get a product by ID
// @Description  Retrieve a specific product by its unique identifier.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id             path      string  true   "Product ID"
// @Success      200  {object}  helper.Response{data=ProductResponse} "Product retrieved successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      404  {object}  helper.Response{data=string} "Product not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /products/{id} [get]
func (pc *ProductController) GetById(c echo.Context) error {
	paramId := c.Param("id")
	productId, err := uuid.Parse(paramId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}
	product, err := pc.productService.GetById(productId.String())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	response := new(ProductResponse).ToResponse(product)

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "get product successfully", response))
}

// Get Products by Category
// @Summary      Get products by category
// @Description  Retrieve products by a specific category name with pagination.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        category_name  path      string  true   "Category name"
// @Param        pages          query     int     false  "Page number"
// @Success      200  {object}  helper.MetadataResponse{data=[]ProductResponse} "Products retrieved successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /products/categories/{category_name} [get]
func (pc *ProductController) GetByCategory(c echo.Context) error {
	productCategory := c.Param("category_name")
	page, err := strconv.Atoi(c.QueryParam("pages"))
	if err != nil {
		page = 1
	}

	products, totalPages, err := pc.productService.GetByCategory(productCategory, page)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	var response []interface{}
	for _, product := range products {
		response = append(response, new(ProductResponse).ToResponse(product))
	}
	metadata := map[string]interface{}{
		"TotalPage": totalPages,
		"Page":      page,
	}

	return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "get all products by category name successfully", metadata, response))

}

// Update Product
// @Summary      Update a product
// @Description  Update product details including images and categories. Requires admin role.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Product ID"
// @Param        body           body      ProductRequest  true   "Updated product data"
// @Success      201  {object}  helper.Response{data=string} "Product updated successfully"
// @Failure      400  {object}  helper.Response{data=string} "Bad request"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /products/{id} [put]
func (pc *ProductController) Update(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
		return nil
	}

	token, err := pc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return helper.UnauthorizedError(c)
	}

	adminData := pc.jwtService.ExtractUserToken(token)
	role := adminData[constant.JWT_ROLE]
	if role != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	productID := c.Param("id")

	var productInput ProductRequest

	if err := c.Bind(&productInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "error bad request", nil))
	}

	if err := c.Validate(productInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	productData := products.Product{
		ID:          productID,
		Name:        productInput.Name,
		Description: productInput.Description,
		Price:       productInput.Price,
		Coin:        productInput.Coin,
		Stock:       productInput.Stock,
		Category:    productInput.CategoryProduct,
	}

	for _, categoryID := range productInput.CategoryImpact {
		productData.ImpactCategories = append(productData.ImpactCategories, products.ProductImpactCategory{
			ProductID:        productData.ID,
			ImpactCategoryID: categoryID,
		})
	}

	for _, imageURL := range productInput.Images {
		productData.Images = append(productData.Images, products.ProductImage{
			ProductID: productData.ID,
			AlbumsURL: imageURL,
		})
	}

	err = pc.productService.Update(productData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "update product successfully", nil))
}

// Delete Product
// @Summary      Delete a product
// @Description  Remove a product by its ID. Requires admin role.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Product ID"
// @Success      200  {object}  helper.Response{data=string} "Product deleted successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      404  {object}  helper.Response{data=string} "Product not found"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /products/{id} [delete]
func (pc *ProductController) Delete(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		helper.UnauthorizedError(c)
		return nil
	}

	token, err := pc.jwtService.ValidateToken(tokenString)
	if err != nil {
		helper.UnauthorizedError(c)
		return nil
	}

	adminData := pc.jwtService.ExtractUserToken(token)
	role := adminData[constant.JWT_ROLE]
	if role != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	productID := c.Param("id")

	err = pc.productService.Delete(productID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "delete product successfully", nil))
}
