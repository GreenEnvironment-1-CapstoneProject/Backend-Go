package repository

import (
	"errors"
	"greenenvironment/features/dashboard"
	"greenenvironment/features/transactions"
	"time"

	"gorm.io/gorm"
)

type DashboardData struct {
	DB *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) dashboard.DashboardRepositoryInterface {
	return &DashboardData{DB: db}
}

func (dd *DashboardData) GetDashboardData(filter string) (dashboard.DashboardData, error) {
	var current dashboard.DashboardData
	var previous dashboard.DashboardData

	query := dd.DB.Model(&transactions.Transaction{}).
		Where("status = ?", "settlement").
		Where("deleted_at IS NULL")
	prevQuery := dd.DB.Model(&transactions.Transaction{}).
		Where("status = ?", "settlement").
		Where("deleted_at IS NULL")

	switch filter {
	case "weekly":
		query = query.Where("YEARWEEK(created_at, 1) = YEARWEEK(NOW(), 1)")
		prevQuery = prevQuery.Where("YEARWEEK(created_at, 1) = YEARWEEK(DATE_SUB(NOW(), INTERVAL 1 WEEK), 1)")
	case "monthly":
		query = query.Where("YEAR(created_at) = YEAR(NOW()) AND MONTH(created_at) = MONTH(NOW())")
		prevQuery = prevQuery.Where("YEAR(created_at) = YEAR(NOW()) AND MONTH(created_at) = MONTH(DATE_SUB(NOW(), INTERVAL 1 MONTH))")
	case "yearly":
		query = query.Where("YEAR(created_at) = YEAR(NOW())")
		prevQuery = prevQuery.Where("YEAR(created_at) = YEAR(DATE_SUB(NOW(), INTERVAL 1 YEAR))")
	default:
		return current, errors.New("invalid filter")
	}

	if err := query.Select("IFNULL(SUM(total), 0)").Scan(&current.TotalTransactions).Error; err != nil {
		return current, err
	}
	if err := query.Count(&current.TotalOrders).Error; err != nil {
		return current, err
	}
	if err := query.Distinct("user_id").Count(&current.TotalCustomers).Error; err != nil {
		return current, err
	}

	if err := prevQuery.Select("IFNULL(SUM(total), 0)").Scan(&previous.TotalTransactions).Error; err != nil {
		return current, err
	}
	if err := prevQuery.Count(&previous.TotalOrders).Error; err != nil {
		return current, err
	}
	if err := prevQuery.Distinct("user_id").Count(&previous.TotalCustomers).Error; err != nil {
		return current, err
	}

	current.TransactionChange = calculateChange(current.TotalTransactions, previous.TotalTransactions)
	current.OrderChange = calculateChange(float64(current.TotalOrders), float64(previous.TotalOrders))
	current.CustomerChange = calculateChange(float64(current.TotalCustomers), float64(previous.TotalCustomers))

	return current, nil
}

func (dd *DashboardData) GetTopCategories(filter string) ([]dashboard.TopCategory, error) {
	var categories []dashboard.TopCategory

	query := dd.DB.Model(&transactions.TransactionItems{}).
		Select("p.category AS category_name, COUNT(transaction_items.id) AS items_sold, SUM(transaction_items.quantity * p.price) AS total_sales").
		Joins("JOIN products p ON p.id = transaction_items.product_id").
		Joins("JOIN transactions t ON t.id = transaction_items.transaction_id").
		Where("t.status = ?", "settlement").
		Where("t.deleted_at IS NULL").
		Group("p.category").
		Order("total_sales DESC").
		Limit(6)

	switch filter {
	case "weekly":
		query = query.Where("YEARWEEK(t.created_at, 1) = YEARWEEK(NOW(), 1)")
	case "monthly":
		query = query.Where("YEAR(t.created_at) = YEAR(NOW()) AND MONTH(t.created_at) = MONTH(NOW())")
	case "yearly":
		query = query.Where("YEAR(t.created_at) = YEAR(NOW())")
	default:
		return nil, errors.New("invalid filter")
	}

	if err := query.Scan(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func calculateChange(current, previous float64) dashboard.ChangeData {
	var change dashboard.ChangeData
	if previous == 0 {
		change.Percentage = 0
		change.Absolute = current
	} else {
		change.Percentage = ((current - previous) / previous) * 100
		change.Absolute = current - previous
	}
	return change
}

func (dd *DashboardData) GetLastTransactions() ([]dashboard.LastTransaction, error) {
	var lastTransactions []dashboard.LastTransaction

	rows, err := dd.DB.Table("transactions").
		Select(`transactions.id, 
				transactions.created_at AS transaction_date,
				u.name AS customer_name, 
				transactions.total, 
				transactions.payment_method, 
				transactions.status, 
				p.name AS product_name`).
		Joins("JOIN transaction_items ti ON ti.transaction_id = transactions.id").
		Joins("JOIN products p ON p.id = ti.product_id").
		Joins("JOIN users u ON u.id = transactions.user_id").
		Where("transactions.deleted_at IS NULL").
		Order("transactions.created_at DESC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactionMap := make(map[string]*dashboard.LastTransaction)

	for rows.Next() {
		var (
			id, customerName, paymentMethod, status, productName string
			transactionDate                                      time.Time
			total                                                float64
		)
		if err := rows.Scan(&id, &transactionDate, &customerName, &total, &paymentMethod, &status, &productName); err != nil {
			return nil, err
		}

		if _, exists := transactionMap[id]; !exists {
			transactionMap[id] = &dashboard.LastTransaction{
				ID:              id,
				TransactionDate: transactionDate,
				CustomerName:    customerName,
				Total:           total,
				PaymentMethod:   paymentMethod,
				Status:          status,
				Products:        []string{},
			}
		}
		transactionMap[id].Products = append(transactionMap[id].Products, productName)
	}

	for _, transaction := range transactionMap {
		lastTransactions = append(lastTransactions, *transaction)
	}

	return lastTransactions, nil
}
