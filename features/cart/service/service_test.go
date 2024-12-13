package service

import (
	"greenenvironment/features/cart"
	"greenenvironment/features/products"
	"greenenvironment/features/users"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCartRepo struct {
	mock.Mock
}

func (m *MockCartRepo) IsCartExist(userID, productID string) (bool, error) {
	args := m.Called(userID, productID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCartRepo) GetCartQty(userID, productID string) (int, error) {
	args := m.Called(userID, productID)
	return args.Int(0), args.Error(1)
}

func (m *MockCartRepo) GetStock(productID string) (int, error) {
	args := m.Called(productID)
	return args.Int(0), args.Error(1)
}

func (m *MockCartRepo) InsertIncrement(userID, productID string, quantity int) error {
	args := m.Called(userID, productID, quantity)
	return args.Error(0)
}

func (m *MockCartRepo) Create(cart cart.NewCart) error {
	args := m.Called(cart)
	return args.Error(0)
}

func (m *MockCartRepo) Delete(userID, productID string) error {
	args := m.Called(userID, productID)
	return args.Error(0)
}

func (m *MockCartRepo) InsertDecrement(userID, productID string) error {
	args := m.Called(userID, productID)
	return args.Error(0)
}

func (m *MockCartRepo) InsertByQuantity(userID, productID string, quantity int) error {
	args := m.Called(userID, productID, quantity)
	return args.Error(0)
}

func (m *MockCartRepo) Get(userID string) (cart.Cart, error) {
	args := m.Called(userID)
	return args.Get(0).(cart.Cart), args.Error(1)
}
func (m *MockCartRepo) Update(cartUpdate cart.UpdateCart) error {
	args := m.Called(cartUpdate)
	return args.Error(0)
}

// TestCreateCart_Success tests the Create method of CartService
func TestCreateCart_Success(t *testing.T) {
	mockRepo := new(MockCartRepo)
	service := NewCartService(mockRepo)

	newCart := cart.NewCart{
		UserID:    "user1",
		ProductID: "product1",
		Quantity:  2,
	}

	mockRepo.On("IsCartExist", newCart.UserID, newCart.ProductID).Return(false, nil)
	mockRepo.On("Create", newCart).Return(nil)

	err := service.Create(newCart)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreateCart_ExistsAndStockSufficient tests the Create method when the cart exists and stock is sufficient
func TestCreateCart_ExistsAndStockSufficient(t *testing.T) {
	mockRepo := new(MockCartRepo)
	service := NewCartService(mockRepo)

	newCart := cart.NewCart{
		UserID:    "user1",
		ProductID: "product1",
		Quantity:  2,
	}

	mockRepo.On("IsCartExist", newCart.UserID, newCart.ProductID).Return(true, nil)
	mockRepo.On("GetCartQty", newCart.UserID, newCart.ProductID).Return(1, nil)
	mockRepo.On("GetStock", newCart.ProductID).Return(5, nil)
	mockRepo.On("InsertIncrement", newCart.UserID, newCart.ProductID, newCart.Quantity).Return(nil)

	err := service.Create(newCart)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreateCart_ExistsAndStockInsufficient tests the Create method when the cart exists but stock is insufficient
func TestCreateCart_ExistsAndStockInsufficient(t *testing.T) {
	mockRepo := new(MockCartRepo)
	service := NewCartService(mockRepo)

	newCart := cart.NewCart{
		UserID:    "user1",
		ProductID: "product1",
		Quantity:  2,
	}

	mockRepo.On("IsCartExist", newCart.UserID, newCart.ProductID).Return(true, nil)
	mockRepo.On("GetCartQty", newCart.UserID, newCart.ProductID).Return(1, nil)
	mockRepo.On("GetStock", newCart.ProductID).Return(2, nil)

	err := service.Create(newCart)

	assert.EqualError(t, err, "error quantity exceeds stock")
	mockRepo.AssertExpectations(t)
}

// TestUpdateCart_Increment tests the Update method for incrementing quantity
func TestUpdateCart_Increment(t *testing.T) {
	mockRepo := new(MockCartRepo)
	service := NewCartService(mockRepo)

	updateCart := cart.UpdateCart{
		UserID:    "user1",
		ProductID: "product1",
		Quantity:  1,
		Type:      "increment",
	}

	mockRepo.On("GetCartQty", updateCart.UserID, updateCart.ProductID).Return(1, nil)
	mockRepo.On("InsertIncrement", updateCart.UserID, updateCart.ProductID, updateCart.Quantity).Return(nil)

	err := service.Update(updateCart)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUpdateCart_Decrement tests the Update method for decrementing quantity
func TestUpdateCart_Decrement(t *testing.T) {
	mockRepo := new(MockCartRepo)
	service := NewCartService(mockRepo)

	updateCart := cart.UpdateCart{
		UserID:    "user1",
		ProductID: "product1",
		Quantity:  1,
		Type:      "decrement",
	}

	mockRepo.On("GetCartQty", updateCart.UserID, updateCart.ProductID).Return(2, nil)
	mockRepo.On("InsertDecrement", updateCart.UserID, updateCart.ProductID).Return(nil)

	err := service.Update(updateCart)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestDeleteCart tests the Delete method of CartService
func TestDeleteCart(t *testing.T) {
	mockRepo := new(MockCartRepo)
	service := NewCartService(mockRepo)

	userID := "user1"
	productID := "product1"

	mockRepo.On("Delete", userID, productID).Return(nil)

	err := service.Delete(userID, productID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestGetCart tests the Get method of CartService
func TestGetCart(t *testing.T) {
	mockRepo := new(MockCartRepo)
	service := NewCartService(mockRepo)

	userID := "user1"
	expectedCart := cart.Cart{
		User: users.User{
			ID:       "user1",
			Username: "user1",
		},
		Items: []cart.CartItem{{Product: products.Product{
			ID:   "product1",
			Name: "product1",
		}}},
	}

	mockRepo.On("Get", userID).Return(expectedCart, nil)

	cart, err := service.Get(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCart, cart)
	mockRepo.AssertExpectations(t)
}
