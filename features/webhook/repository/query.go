package repository

import (
	productData "greenenvironment/features/products/repository"
	transactionsEntity "greenenvironment/features/transactions"
	transactionsData "greenenvironment/features/transactions/repository"
	userData "greenenvironment/features/users/repository"
	"greenenvironment/features/webhook"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebhookRepository struct {
	DB *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) webhook.MidtransNotificationRepository {
	return &WebhookRepository{
		DB: db,
	}
}

func (w *WebhookRepository) HandleNotification(notification webhook.PaymentNotification, transaction transactionsData.Transaction) error {
	transactionUpdate := transactionsEntity.UpdateTransaction{
		ID:            transaction.ID,
		Status:        transaction.Status,
		PaymentMethod: transaction.PaymentMethod,
	}

	var paymentNotif = PaymentNotification{
		ID:                uuid.New().String(),
		OrderID:           notification.OrderID,
		TransactionTime:   notification.TransactionTime,
		TransactionStatus: notification.TransactionStatus,
		TransactionID:     notification.TransactionID,
		SignatureKey:      notification.SignatureKey,
		PaymentType:       notification.PaymentType,
		MerchantID:        notification.MerchantID,
		GrossAmount:       notification.GrossAmount,
		FraudStatus:       notification.FraudStatus,
		Currency:          notification.Currency,
		SettlementTime:    notification.SettlementTime,
	}
	tx := w.DB.Begin()

	err := w.DB.Model(&transactionsData.Transaction{}).Where("id = ?", transaction.ID).Updates(&transactionUpdate).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = w.DB.Model(&PaymentNotification{}).Create(&paymentNotif).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if notification.TransactionStatus == "settlement" {
		err = w.InsertUserCoin(transaction.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (w *WebhookRepository) InsertUserCoin(transactionId string) error {
	var transaction transactionsData.Transaction
	err := w.DB.Where("id = ?", transactionId).First(&transaction).Error
	if err != nil {
		return err
	}
	var user userData.User
	err = w.DB.Where("id = ?", transaction.UserID).First(&user).Error
	if err != nil {
		return err
	}
	var transactionItem []transactionsData.TransactionItem
	err = w.DB.Where("transaction_id = ?", transaction.ID).Find(&transactionItem).Error
	if err != nil {
		return err
	}
	var totalCoinxQty int
	for _, item := range transactionItem {
		var product productData.Product
		err = w.DB.Where("id = ?", item.ProductID).First(&product).Error
		totalCoinxQty += product.Coin * item.Quantity
		if err != nil {
			return err
		}
	}
	userUpdate := userData.User{
		Coin: user.Coin + totalCoinxQty,
	}

	err = w.DB.Model(&userData.User{}).Where("id = ?", user.ID).Updates(&userUpdate).Error
	if err != nil {
		return err
	}
	return nil
}

func (w *WebhookRepository) UpdateStockFailedTransaction(transactionId string) error {
	var transactionsItems []transactionsData.TransactionItem
	err := w.DB.Where("transaction_id = ?", transactionId).Find(&transactionsItems).Error
	if err != nil {
		return err
	}

	for _, item := range transactionsItems {
		err := w.DB.Model(productData.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error
		if err != nil {
			return err
		}
	}

	return nil
}
