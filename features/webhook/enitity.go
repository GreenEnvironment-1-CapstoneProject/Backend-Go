package webhook

import (
	transactions "greenenvironment/features/transactions/repository"

	"github.com/labstack/echo/v4"
)

type PaymentNotification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
	SettlementTime    string `json:"settlement_time,omitempty"`
}

type MidtransNotificationController interface {
	HandleNotification(c echo.Context) error
}

type MidtransNotificationService interface {
	HandleNotification(notification PaymentNotification) error
}

type MidtransNotificationRepository interface {
	HandleNotification(notification PaymentNotification, transaction transactions.Transaction) error
	InsertUserCoin(transactionID string) error
	UpdateStockFailedTransaction(transactionId string) error
}
