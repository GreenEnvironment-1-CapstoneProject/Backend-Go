package service

import (
	"testing"

	cart "greenenvironment/features/cart/repository"
	products "greenenvironment/features/products/repository"
	"greenenvironment/features/transactions"
	users "greenenvironment/features/users/repository"
	"greenenvironment/utils/midtrans"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepo struct {
	mock.Mock
}

type MockMidtransService struct {
	mock.Mock
}

func (m *MockTransactionRepo) GetUserTransaction(userId string, page int) ([]transactions.TransactionData, int, int, error) {
	args := m.Called(userId, page)
	return args.Get(0).([]transactions.TransactionData), args.Int(1), args.Int(2), args.Error(3)
}

func (m *MockTransactionRepo) GetTransactionByID(transactionId string) (transactions.TransactionData, error) {
	args := m.Called(transactionId)
	return args.Get(0).(transactions.TransactionData), args.Error(1)
}

func (m *MockTransactionRepo) GetUserData(userId string) (users.User, error) {
	args := m.Called(userId)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *MockTransactionRepo) GetDataCartTransaction(cartId []string, userId string) ([]cart.Cart, error) {
	args := m.Called(cartId, userId)
	return args.Get(0).([]cart.Cart), args.Error(1)
}

func (m *MockTransactionRepo) UpdateStockByProductID(productId string, quantity int) error {
	args := m.Called(productId, quantity)
	return args.Error(0)
}

func (m *MockTransactionRepo) GetUserCoin(userId string) (int, error) {
	args := m.Called(userId)
	return args.Int(0), args.Error(1)
}

func (m *MockTransactionRepo) DecreaseUserCoin(userId string, coin int, total float64) (float64, int, error) {
	args := m.Called(userId, coin, total)
	return args.Get(0).(float64), args.Int(1), args.Error(2)
}

func (m *MockTransactionRepo) CreateTransactions(transaction transactions.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepo) CreateTransactionItems(items []transactions.TransactionItems) error {
	args := m.Called(items)
	return args.Error(0)
}

func (m *MockTransactionRepo) DeleteTransaction(transactionId string) error {
	args := m.Called(transactionId)
	return args.Error(0)
}

func (m *MockTransactionRepo) GetAllTransaction(page int) ([]transactions.TransactionData, int, int, error) {
	args := m.Called(page)
	return args.Get(0).([]transactions.TransactionData), args.Int(1), args.Int(2), args.Error(3)
}

func (m *MockMidtransService) InitializeClientMidtrans() {
	m.Called()
}

func (m *MockMidtransService) CreateTransaction(req midtrans.CreatePaymentGateway) string {
	args := m.Called(req)
	return args.String(0)
}

func (m *MockMidtransService) CreateUrlTransactionWithGateway(snap midtrans.CreatePaymentGateway) string {
	args := m.Called(snap)
	return args.String(0)
}

func (m *MockMidtransService) CancelTransaction(orderId string) error {
	args := m.Called(orderId)
	return args.Error(0)
}
func TestGetUserTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepo)
	mockMidtrans := new(MockMidtransService)
	service := NewTransactionService(mockRepo, mockMidtrans)

	mockRepo.On("GetUserTransaction", "user1", 1).Return([]transactions.TransactionData{}, 1, 1, nil)

	transactions, totalPage, totalData, err := service.GetUserTransaction("user1", 1)

	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, 1, totalPage)
	assert.Equal(t, 1, totalData)
	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepo)
	mockMidtrans := new(MockMidtransService)
	service := NewTransactionService(mockRepo, mockMidtrans)

	mockRepo.On("GetUserData", "user1").Return(users.User{Address: "address"}, nil)
	mockRepo.On("GetDataCartTransaction", []string{"cart1"}, "user1").Return([]cart.Cart{
		{ID: "1", ProductID: "1", Product: products.Product{ID: "1", Price: 1000, Name: "product1"}, Quantity: 1},
	}, nil)
	mockRepo.On("UpdateStockByProductID", "1", 1).Return(nil)
	mockRepo.On("CreateTransactions", mock.Anything).Return(nil)
	mockRepo.On("CreateTransactionItems", mock.Anything).Return(nil)
	mockMidtrans.On("InitializeClientMidtrans").Return()
	mockMidtrans.On("CreateTransaction", mock.Anything).Return("snap_url")

	transaction := transactions.CreateTransaction{
		UserID: "user1",
		CartID: []string{"cart1"},
	}

	result, err := service.CreateTransaction(transaction)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "snap_url", result.SnapURL)
	mockRepo.AssertExpectations(t)
	mockMidtrans.AssertExpectations(t)
}

func TestDeleteTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepo)
	mockMidtrans := new(MockMidtransService)
	service := NewTransactionService(mockRepo, mockMidtrans)

	mockRepo.On("GetTransactionByID", "transaction1").Return(transactions.TransactionData{}, nil)
	mockRepo.On("DeleteTransaction", "transaction1").Return(nil)

	err := service.DeleteTransaction("transaction1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAllTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepo)
	mockMidtrans := new(MockMidtransService)
	service := NewTransactionService(mockRepo, mockMidtrans)

	mockRepo.On("GetAllTransaction", 1).Return([]transactions.TransactionData{}, 1, 1, nil)

	transactions, totalPage, totalData, err := service.GetAllTransaction(1)

	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, 1, totalPage)
	assert.Equal(t, 1, totalData)
	mockRepo.AssertExpectations(t)
}
