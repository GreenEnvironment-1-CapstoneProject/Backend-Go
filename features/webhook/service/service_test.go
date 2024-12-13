package service

import (
	"testing"

	"greenenvironment/features/transactions/repository"
	"greenenvironment/features/webhook"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMidtransNotificationRepository struct {
	mock.Mock
}

func (m *MockMidtransNotificationRepository) HandleNotification(notification webhook.PaymentNotification, transactionData repository.Transaction) error {
	args := m.Called(notification, transactionData)
	return args.Error(0)
}

func (m *MockMidtransNotificationRepository) UpdateStockFailedTransaction(orderID string) error {
	args := m.Called(orderID)
	return args.Error(0)
}

func (m *MockMidtransNotificationRepository) InsertUserCoin(transactionID string) error {
	args := m.Called(transactionID)
	return args.Error(0)
}
func TestHandleNotification(t *testing.T) {
	mockRepo := new(MockMidtransNotificationRepository)
	service := NewWebhookService(mockRepo)

	tests := []struct {
		name           string
		notification   webhook.PaymentNotification
		expectedStatus string
		mockError      error
	}{
		{
			name: "capture and accept",
			notification: webhook.PaymentNotification{
				TransactionStatus: "capture",
				FraudStatus:       "accept",
				OrderID:           "order123",
				PaymentType:       "credit_card",
			},
			expectedStatus: "capture",
			mockError:      nil,
		},
		{
			name: "settlement",
			notification: webhook.PaymentNotification{
				TransactionStatus: "settlement",
				OrderID:           "order124",
				PaymentType:       "credit_card",
			},
			expectedStatus: "settlement",
			mockError:      nil,
		},
		{
			name: "cancel",
			notification: webhook.PaymentNotification{
				TransactionStatus: "cancel",
				OrderID:           "order125",
				PaymentType:       "credit_card",
			},
			expectedStatus: "cancel",
			mockError:      nil,
		},
		{
			name: "pending",
			notification: webhook.PaymentNotification{
				TransactionStatus: "pending",
				OrderID:           "order126",
				PaymentType:       "credit_card",
			},
			expectedStatus: "pending",
			mockError:      nil,
		},
		// {
		// 	name: "update stock failed",
		// 	notification: webhook.PaymentNotification{
		// 		TransactionStatus: "cancel",
		// 		OrderID:           "order127",
		// 		PaymentType:       "credit_card",
		// 	},
		// 	expectedStatus: "cancel",
		// 	mockError:      errors.New("update stock failed"),
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactionData := repository.Transaction{
				ID:            tt.notification.OrderID,
				PaymentMethod: tt.notification.PaymentType,
				Status:        tt.expectedStatus,
			}

			// Set expectations based on the transaction status
			if tt.notification.TransactionStatus == "cancel" || tt.notification.TransactionStatus == "deny" || tt.notification.TransactionStatus == "expire" {
				mockRepo.On("UpdateStockFailedTransaction", tt.notification.OrderID).Return(tt.mockError).Once()
			}

			// Always expect HandleNotification to be called
			mockRepo.On("HandleNotification", tt.notification, transactionData).Return(tt.mockError).Once()

			// Call the function under test
			err := service.HandleNotification(tt.notification)

			// Assert the results
			if tt.mockError != nil {
				assert.EqualError(t, err, tt.mockError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Assert that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
