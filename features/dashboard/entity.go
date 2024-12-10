package dashboard

import (
	"time"

	"github.com/labstack/echo/v4"
)

type DashboardData struct {
	TotalTransactions float64
	TransactionChange ChangeData
	TotalOrders       int64
	OrderChange       ChangeData
	TotalCustomers    int64
	CustomerChange    ChangeData
	TopCategories     []TopCategory
	LastTransactions  []LastTransaction
}

type TopCategory struct {
	CategoryName string
	ItemsSold    int64
	TotalSales   float64
}

type ChangeData struct {
	Percentage float64
	Absolute   float64
}

type LastTransaction struct {
	ID              string    `json:"id"`
	Products        []string  `json:"products"`
	TransactionDate time.Time `json:"transaction_date"`
	CustomerName    string    `json:"customer_name"`
	Total           float64   `json:"total"`
	PaymentMethod   string    `json:"payment_method"`
	Status          string    `json:"status"`
}

type DashboardRepositoryInterface interface {
	GetDashboardData(filter string) (DashboardData, error)
	GetTopCategories(filter string) ([]TopCategory, error)
	GetLastTransactions() ([]LastTransaction, error)
}

type DashboardServiceInterface interface {
	GetDashboardData(filter string) (DashboardData, error)
}

type DashboardControllerInterface interface {
	GetDashboard(c echo.Context) error
}
