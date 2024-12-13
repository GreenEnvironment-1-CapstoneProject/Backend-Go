package service

import (
	"errors"
	"testing"

	"greenenvironment/features/impacts"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockImpactRepository struct {
	mock.Mock
}

func (m *MockImpactRepository) GetAll() ([]impacts.ImpactCategory, error) {
	args := m.Called()
	return args.Get(0).([]impacts.ImpactCategory), args.Error(1)
}

func (m *MockImpactRepository) GetByID(ID string) (impacts.ImpactCategory, error) {
	args := m.Called(ID)
	return args.Get(0).(impacts.ImpactCategory), args.Error(1)
}

func (m *MockImpactRepository) Create(category impacts.ImpactCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockImpactRepository) Delete(category impacts.ImpactCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func TestGetAll(t *testing.T) {
	mockRepo := new(MockImpactRepository)
	mockRepo.On("GetAll").Return([]impacts.ImpactCategory{}, nil)

	service := NewNewImpactService(mockRepo)
	categories, err := service.GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	mockRepo.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	mockRepo := new(MockImpactRepository)
	mockCategory := impacts.ImpactCategory{ID: "1", Name: "Test"}
	mockRepo.On("GetByID", "1").Return(mockCategory, nil)

	service := NewNewImpactService(mockRepo)
	category, err := service.GetByID("1")

	assert.NoError(t, err)
	assert.Equal(t, "1", category.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	mockRepo := new(MockImpactRepository)
	mockRepo.On("GetByID", "2").Return(impacts.ImpactCategory{}, errors.New("not found"))

	service := NewNewImpactService(mockRepo)
	category, err := service.GetByID("2")

	assert.Error(t, err)
	assert.Equal(t, impacts.ImpactCategory{}, category)
	mockRepo.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	mockRepo := new(MockImpactRepository)
	mockCategory := impacts.ImpactCategory{ID: "1", Name: "Test"}
	mockRepo.On("Create", mockCategory).Return(nil)

	service := NewNewImpactService(mockRepo)
	err := service.Create(mockCategory)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	mockRepo := new(MockImpactRepository)
	mockCategory := impacts.ImpactCategory{ID: "1", Name: "Test"}
	mockRepo.On("Delete", mockCategory).Return(nil)

	service := NewNewImpactService(mockRepo)
	err := service.Delete(mockCategory)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
