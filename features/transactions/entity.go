package transactions

import (
	cart "greenenvironment/features/cart/repository"
	"greenenvironment/features/products"
	users "greenenvironment/features/users/repository"
	"time"

	"github.com/labstack/echo/v4"
)

type Transaction struct {
	ID            string
	UserID        string
	Address       string
	Total         float64
	Status        string
	PaymentMethod string
	SnapURL       string
	Coin          int
}

type CreateTransaction struct {
	UserID    string
	CartID    []string
	UsingCoin bool
}

type TransactionItems struct {
	ID            string
	TransactionID string
	ProductID     string
	Qty           int
	Product       products.Product
}

type UpdateTransaction struct {
	ID            string
	Status        string
	PaymentMethod string
}

type TransactionData struct {
	ID               string
	VoucherID        string
	Status           string
	Total            float64
	Coin             int
	SnapURL          string
	User             users.User
	TransactionItems []TransactionItems
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type TransactionRepositoryInterface interface {
	GetUserTransaction(userId string) ([]TransactionData, error)
	GetTransactionByID(transactionId string) (TransactionData, error)
	CreateTransactions(transaction Transaction) error
	DeleteTransaction(transactionId string) error
	GetUserData(userId string) (users.User, error)
	GetUserCoin(userId string) (int, error)
	DecreaseUserCoin(userId string, coin int, total float64) (float64, int, error)
	CreateTransactionItems(tansactionItems []TransactionItems) error
	GetAllTransaction() ([]TransactionData, error)
	GetDataCartTransaction(cartIds []string, userId string) ([]cart.Cart, error)
	UpdateStockByProductID(productId string, stock int) error
}

type TransactionServiceInterface interface {
	GetUserTransaction(userId string) ([]TransactionData, error)
	GetTransactionByID(transactionId string) (TransactionData, error)
	CreateTransaction(transaction CreateTransaction) (Transaction, error)
	DeleteTransaction(transactionId string) error
	GetAllTransaction() ([]TransactionData, error)
}

type TransactionControllerInterface interface {
	GetUserTransaction(c echo.Context) error
	CreateTransaction(c echo.Context) error
	DeleteTransaction(c echo.Context) error
	GetAllTransaction(c echo.Context) error
}
