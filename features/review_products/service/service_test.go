package service

import (
	"errors"
	"testing"
	"time"

	reviewproducts "greenenvironment/features/review_products"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReviewProductRepository struct {
	mock.Mock
}

func (m *MockReviewProductRepository) Create(createDto reviewproducts.CreateReviewProduct) error {
	args := m.Called(createDto)
	return args.Error(0)
}

func (m *MockReviewProductRepository) GetProductReview(productID string) ([]reviewproducts.ReviewProduct, error) {
	args := m.Called(productID)
	return args.Get(0).([]reviewproducts.ReviewProduct), args.Error(1)
}

func TestCreateReviewProduct_Success(t *testing.T) {
	mockRepo := new(MockReviewProductRepository)
	service := NewReviewProductService(mockRepo)

	createDto := reviewproducts.CreateReviewProduct{
		ProductID: "123",
		UserID:    "456",
		Rate:      5,
		Review:    "Excellent product!",
	}

	mockRepo.On("Create", createDto).Return(nil)

	err := service.Create(createDto)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateReviewProduct_Error(t *testing.T) {
	mockRepo := new(MockReviewProductRepository)
	service := NewReviewProductService(mockRepo)

	createDto := reviewproducts.CreateReviewProduct{
		ProductID: "123",
		UserID:    "456",
		Rate:      5,
		Review:    "Excellent product!",
	}

	mockRepo.On("Create", createDto).Return(errors.New("failed to create review"))

	err := service.Create(createDto)

	assert.Error(t, err)
	assert.Equal(t, "failed to create review", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetProductReview_Success(t *testing.T) {
	mockRepo := new(MockReviewProductRepository)
	service := NewReviewProductService(mockRepo)

	mockReviews := []reviewproducts.ReviewProduct{
		{
			Name:      "John Doe",
			Email:     "john@example.com",
			Review:    "Excellent!",
			Rate:      5,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Jane Smith",
			Email:     "jane@example.com",
			Review:    "Very good quality.",
			Rate:      4,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.On("GetProductReview", "123").Return(mockReviews, nil)

	reviews, err := service.GetProductReview("123")

	assert.NoError(t, err)
	assert.Len(t, reviews, 2)
	assert.Equal(t, "John Doe", reviews[0].Name)
	assert.Equal(t, "Excellent!", reviews[0].Review)
	mockRepo.AssertExpectations(t)
}

func TestGetProductReview_Error(t *testing.T) {
	mockRepo := new(MockReviewProductRepository)
	service := NewReviewProductService(mockRepo)

	mockRepo.On("GetProductReview", "123").Return([]reviewproducts.ReviewProduct{}, errors.New("reviews not found"))

	reviews, err := service.GetProductReview("123")

	assert.Error(t, err)
	assert.Equal(t, []reviewproducts.ReviewProduct{}, reviews)
	assert.Equal(t, "reviews not found", err.Error())
	mockRepo.AssertExpectations(t)
}
