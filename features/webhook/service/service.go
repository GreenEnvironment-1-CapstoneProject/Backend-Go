package service

import (
	transactions "greenenvironment/features/transactions/repository"
	"greenenvironment/features/webhook"
)

type WebhookService struct {
	d webhook.MidtransNotificationRepository
}

func NewWebhookService(data webhook.MidtransNotificationRepository) webhook.MidtransNotificationService {
	return &WebhookService{
		d: data,
	}
}

func (s *WebhookService) HandleNotification(notification webhook.PaymentNotification) error {
	// orderId := notification.OrderID
	transactionStatus := notification.TransactionStatus
	fraudStatus := notification.FraudStatus
	transactionData := transactions.Transaction{
		ID:            notification.OrderID,
		PaymentMethod: notification.PaymentType,
	}
	if transactionStatus == "capture" {
		if fraudStatus == "accept" {
			transactionData.Status = transactionStatus
		}
	} else if transactionStatus == "settlement" {
		transactionData.Status = transactionStatus
	} else if transactionStatus == "cancel" || transactionStatus == "deny" || transactionStatus == "expire" {
		transactionData.Status = transactionStatus
	} else if transactionStatus == "pending" {
		transactionData.Status = transactionStatus
	}
	return s.d.HandleNotification(notification, transactionData)
}
