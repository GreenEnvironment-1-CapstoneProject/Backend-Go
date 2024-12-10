package controller

import "greenenvironment/features/dashboard"

type DashboardResponse struct {
	TotalTransactions float64                   `json:"total_transactions"`
	TransactionChange ChangeResponse            `json:"transaction_change"`
	TotalOrders       int64                     `json:"total_orders"`
	OrderChange       ChangeResponse            `json:"order_change"`
	TotalCustomers    int64                     `json:"total_customers"`
	CustomerChange    ChangeResponse            `json:"customer_change"`
	TopCategories     []TopCategoryResponse     `json:"top_categories"`
	LastTransactions  []LastTransactionResponse `json:"last_transactions"`
}

type ChangeResponse struct {
	Percentage float64 `json:"percentage"`
	Absolute   float64 `json:"absolute"`
}

type TopCategoryResponse struct {
	CategoryName string  `json:"category_name"`
	ItemsSold    int64   `json:"items_sold"`
	TotalSales   float64 `json:"total_sales"`
}

type LastTransactionResponse struct {
	ID              string   `json:"id"`
	Products        []string `json:"products"`
	TransactionDate string   `json:"transaction_date"`
	CustomerName    string   `json:"customer_name"`
	Total           float64  `json:"total"`
	PaymentMethod   string   `json:"payment_method"`
	Status          string   `json:"status"`
}

func (d DashboardResponse) FromEntity(entity dashboard.DashboardData) DashboardResponse {
	lastTransactions := make([]LastTransactionResponse, len(entity.LastTransactions))
	for i, lt := range entity.LastTransactions {
		lastTransactions[i] = LastTransactionResponse{
			ID:              lt.ID,
			Products:        lt.Products,
			TransactionDate: lt.TransactionDate.Format("2006-01-02 15:04:05"),
			CustomerName:    lt.CustomerName,
			Total:           lt.Total,
			PaymentMethod:   lt.PaymentMethod,
			Status:          lt.Status,
		}
	}

	topCategories := make([]TopCategoryResponse, len(entity.TopCategories))
	for i, tc := range entity.TopCategories {
		topCategories[i] = TopCategoryResponse{
			CategoryName: tc.CategoryName,
			ItemsSold:    tc.ItemsSold,
			TotalSales:   tc.TotalSales,
		}
	}

	return DashboardResponse{
		TotalTransactions: entity.TotalTransactions,
		TransactionChange: ChangeResponse{
			Percentage: entity.TransactionChange.Percentage,
			Absolute:   entity.TransactionChange.Absolute,
		},
		TotalOrders: entity.TotalOrders,
		OrderChange: ChangeResponse{
			Percentage: entity.OrderChange.Percentage,
			Absolute:   entity.OrderChange.Absolute,
		},
		TotalCustomers: entity.TotalCustomers,
		CustomerChange: ChangeResponse{
			Percentage: entity.CustomerChange.Percentage,
			Absolute:   entity.CustomerChange.Absolute,
		},
		TopCategories:    topCategories,
		LastTransactions: lastTransactions,
	}
}
