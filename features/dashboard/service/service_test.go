package service

import (
	"errors"
	"greenenvironment/features/dashboard"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) GetDashboardData(filter string) (dashboard.DashboardData, error) {
	args := m.Called(filter)

	return args.Get(0).(dashboard.DashboardData), args.Error(1)
}

func (m *MockDashboardRepository) GetTopCategories(filter string) ([]dashboard.TopCategory, error) {
	args := m.Called(filter)

	return args.Get(0).([]dashboard.TopCategory), args.Error(1)
}

func (m *MockDashboardRepository) GetLastTransactions() ([]dashboard.LastTransaction, error) {
	args := m.Called()

	return args.Get(0).([]dashboard.LastTransaction), args.Error(1)
}

func TestGetDashboardData(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo)

	filter := "monthly"
	fixedTime := time.Date(2024, time.December, 13, 11, 15, 0, 0, time.Local) // Waktu tetap untuk pengujian

	expectedData := dashboard.DashboardData{
		TotalTransactions: 100,
		TransactionChange: dashboard.ChangeData{Percentage: 10, Absolute: 10},
		TotalOrders:       50,
		OrderChange:       dashboard.ChangeData{Percentage: 5, Absolute: 2},
		TotalCustomers:    30,
		CustomerChange:    dashboard.ChangeData{Percentage: 8, Absolute: 3},
		TopCategories: []dashboard.TopCategory{
			{CategoryName: "Electronics", ItemsSold: 30, TotalSales: 50000},
			{CategoryName: "Clothing", ItemsSold: 20, TotalSales: 30000},
		},
		LastTransactions: []dashboard.LastTransaction{
			{
				ID:              "1",
				Products:        []string{"Laptop", "Mouse"},
				TransactionDate: fixedTime,
				CustomerName:    "John Doe",
				Total:           50000,
				PaymentMethod:   "Credit Card",
				Status:          "Completed",
			},
			{
				ID:              "2",
				Products:        []string{"Shirt", "Pants"},
				TransactionDate: fixedTime,
				CustomerName:    "Jane Smith",
				Total:           30000,
				PaymentMethod:   "Debit Card",
				Status:          "Pending",
			},
		},
	}

	mockRepo.On("GetDashboardData", filter).Return(expectedData, nil)
	mockRepo.On("GetTopCategories", filter).Return(expectedData.TopCategories, nil)
	mockRepo.On("GetLastTransactions").Return(expectedData.LastTransactions, nil)

	result, err := service.GetDashboardData(filter)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
	mockRepo.AssertExpectations(t)
}

func TestGetDashboardData_Error(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo)

	filter := "monthly"

	mockRepo.On("GetDashboardData", filter).Return(dashboard.DashboardData{}, errors.New("data error"))

	_, err := service.GetDashboardData(filter)

	assert.Error(t, err)
	assert.Equal(t, "data error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetTopCategories_Error(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo)

	filter := "monthly"

	mockRepo.On("GetDashboardData", filter).Return(dashboard.DashboardData{
		TotalTransactions: 100,
		TransactionChange: dashboard.ChangeData{Percentage: 10, Absolute: 10},
		TotalOrders:       50,
		OrderChange:       dashboard.ChangeData{Percentage: 5, Absolute: 2},
		TotalCustomers:    30,
		CustomerChange:    dashboard.ChangeData{Percentage: 8, Absolute: 3},
	}, nil)
	mockRepo.On("GetTopCategories", filter).Return([]dashboard.TopCategory{}, errors.New("categories error"))

	_, err := service.GetDashboardData(filter)

	assert.Error(t, err)
	assert.Equal(t, "categories error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetLastTransactions_Error(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo)

	filter := "monthly"

	mockRepo.On("GetDashboardData", filter).Return(dashboard.DashboardData{
		TotalTransactions: 100,
		TransactionChange: dashboard.ChangeData{Percentage: 10, Absolute: 10},
		TotalOrders:       50,
		OrderChange:       dashboard.ChangeData{Percentage: 5, Absolute: 2},
		TotalCustomers:    30,
		CustomerChange:    dashboard.ChangeData{Percentage: 8, Absolute: 3},
	}, nil)
	mockRepo.On("GetTopCategories", filter).Return([]dashboard.TopCategory{
		{CategoryName: "Electronics", ItemsSold: 30, TotalSales: 50000},
		{CategoryName: "Clothing", ItemsSold: 20, TotalSales: 30000},
	}, nil)
	mockRepo.On("GetLastTransactions").Return([]dashboard.LastTransaction{}, errors.New("transactions error"))

	_, err := service.GetDashboardData(filter)

	assert.Error(t, err)
	assert.Equal(t, "transactions error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetDashboardData_Empty(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo)

	filter := "monthly"

	mockRepo.On("GetDashboardData", filter).Return(dashboard.DashboardData{
		TotalTransactions: 0,
		TransactionChange: dashboard.ChangeData{Percentage: 0, Absolute: 0},
		TotalOrders:       0,
		OrderChange:       dashboard.ChangeData{Percentage: 0, Absolute: 0},
		TotalCustomers:    0,
		CustomerChange:    dashboard.ChangeData{Percentage: 0, Absolute: 0},
	}, nil)
	mockRepo.On("GetTopCategories", filter).Return([]dashboard.TopCategory{}, nil)
	mockRepo.On("GetLastTransactions").Return([]dashboard.LastTransaction{}, nil)

	result, err := service.GetDashboardData(filter)

	assert.NoError(t, err)
	assert.Equal(t, dashboard.DashboardData{
		TotalTransactions: 0,
		TransactionChange: dashboard.ChangeData{Percentage: 0, Absolute: 0},
		TotalOrders:       0,
		OrderChange:       dashboard.ChangeData{Percentage: 0, Absolute: 0},
		TotalCustomers:    0,
		CustomerChange:    dashboard.ChangeData{Percentage: 0, Absolute: 0},
		TopCategories:     []dashboard.TopCategory{},
		LastTransactions:  []dashboard.LastTransaction{},
	}, result)
	mockRepo.AssertExpectations(t)
}
