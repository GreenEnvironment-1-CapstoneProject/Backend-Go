package controller

import (
	"greenenvironment/constant"
	"greenenvironment/features/transactions"
	"greenenvironment/helper"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	transactionService transactions.TransactionServiceInterface
	jwtService         helper.JWTInterface
}

func NewTransactionController(s transactions.TransactionServiceInterface, j helper.JWTInterface) transactions.TransactionControllerInterface {
	return &TransactionController{
		transactionService: s,
		jwtService:         j,
	}
}

// Get User Transactions
// @Summary      Get user transactions
// @Description  Retrieve all transactions made by the logged-in user.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Success      200  {object}  helper.Response{data=[]TransactionUserResponse} "Transactions retrieved successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /transactions [get]
func (tc *TransactionController) GetUserTransaction(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	token, err := tc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	userData := tc.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID].(string)
	transactions, err := tc.transactionService.GetUserTransaction(userId)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	response := []TransactionUserResponse{}
	for _, transaction := range transactions {
		response = append(response, new(TransactionUserResponse).FromEntity(transaction))
	}
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success get user transaction", response))

}

// Create Transaction
// @Summary      Create a new transaction
// @Description  Create a new transaction using the specified cart items.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        request        body      TransactionRequest  true  "Transaction Request"
// @Success      200  {object}  helper.Response{data=TransactionResponse} "Transaction created successfully"
// @Failure      400  {object}  helper.Response{data=string} "Bad request"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /transactions [post]
func (tc *TransactionController) CreateTransaction(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	token, err := tc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	userData := tc.jwtService.ExtractUserToken(token)
	userId := userData[constant.JWT_ID].(string)

	var request TransactionRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	transactionData := transactions.CreateTransaction{
		UserID:    userId,
		CartID:    request.CartIds,
		UsingCoin: request.UsingCoin,
	}
	transaction, err := tc.transactionService.CreateTransaction(transactionData)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}
	transactionResponse := TransactionResponse{
		ID:      transaction.ID,
		Amount:  int(transaction.Total),
		SnapURL: transaction.SnapURL,
	}
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success create transaction", transactionResponse))
}

// Delete Transaction
// @Summary      Delete a transaction
// @Description  Delete a transaction by ID. Only accessible by admin users.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Transaction ID"
// @Success      200  {object}  helper.Response{data=string} "Transaction deleted successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /transactions/{id} [delete]
func (tc *TransactionController) DeleteTransaction(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	token, err := tc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	adminData := tc.jwtService.ExtractAdminToken(token)
	role := adminData[constant.JWT_ROLE].(string)
	if role != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	paramId := c.Param("id")
	transactionId, err := uuid.Parse(paramId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	err = tc.transactionService.DeleteTransaction(transactionId.String())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	return c.JSON(http.StatusCreated, helper.FormatResponse(true, "delete transaction successfully", nil))
}

// Get All Transactions
// @Summary      Get all transactions
// @Description  Retrieve all transactions in the system. Only accessible by admin users.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Success      200  {object}  helper.Response{data=[]TransactionAllUserResponses} "Transactions retrieved successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/transactions [get]
func (tc *TransactionController) GetAllTransaction(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	token, err := tc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	adminData := tc.jwtService.ExtractAdminToken(token)
	role := adminData[constant.JWT_ROLE].(string)
	if role != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}

	transactions, err := tc.transactionService.GetAllTransaction()
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}
	response := []TransactionAllUserResponses{}
	for _, transaction := range transactions {
		response = append(response, new(TransactionAllUserResponses).FromEntity(transaction))
	}
	return c.JSON(http.StatusOK, helper.FormatResponse(true, "Get all Transactions", response))
}

// Get Transaction By ID
// @Summary      Get transaction by ID
// @Description  Retrieve a specific transaction by its ID. Only accessible by admin users.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer Token"
// @Param        id             path      string  true   "Transaction ID"
// @Success      200  {object}  helper.Response{data=TransactionAllUserResponses} "Transaction retrieved successfully"
// @Failure      401  {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500  {object}  helper.Response{data=string} "Internal server error"
// @Router       /admin/transactions/{id} [get]
func (tc *TransactionController) GetTransactionByID(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	token, err := tc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	adminData := tc.jwtService.ExtractAdminToken(token)
	role := adminData[constant.JWT_ROLE].(string)
	if role != constant.RoleAdmin {
		return helper.UnauthorizedError(c)
	}
	paramId := c.Param("id")
	transactionId, err := uuid.Parse(paramId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
	}

	transaction, err := tc.transactionService.GetTransactionByID(transactionId.String())
	if err != nil {
		return c.JSON(helper.ConvertResponseCode(err), helper.FormatResponse(false, err.Error(), nil))
	}

	response := new(TransactionAllUserResponses).FromEntity(transaction)

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "get transactions by id successfully", response))
}
