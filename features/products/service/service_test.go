package service

import (
	"greenenvironment/constant"
	"greenenvironment/features/impacts"
	"greenenvironment/features/products"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepo struct {
	mock.Mock
}

type MockImpactRepo struct {
	mock.Mock
}

func (m *MockProductRepo) Create(product products.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepo) GetAllByPage(page int, search string, sort string) ([]products.Product, int, error) {
	args := m.Called(page, search, sort)
	return args.Get(0).([]products.Product), args.Int(1), args.Error(2)
}

func (m *MockProductRepo) GetByCategory(category string, page int, search string, sort string) ([]products.Product, int, error) {
	args := m.Called(category, page, search, sort)
	return args.Get(0).([]products.Product), args.Int(1), args.Error(2)
}

func (m *MockProductRepo) GetById(id string) (products.Product, error) {
	args := m.Called(id)
	return args.Get(0).(products.Product), args.Error(1)
}

func (m *MockProductRepo) Update(product products.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepo) Delete(productId string) error {
	args := m.Called(productId)
	return args.Error(0)
}

func (m *MockProductRepo) GetTotalProduct() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockImpactRepo) GetAll() ([]impacts.ImpactCategory, error) {
	args := m.Called()
	return args.Get(0).([]impacts.ImpactCategory), args.Error(1)
}

func (m *MockImpactRepo) GetByID(ID string) (impacts.ImpactCategory, error) {
	args := m.Called(ID)
	return args.Get(0).(impacts.ImpactCategory), args.Error(1)
}

func (m *MockImpactRepo) Create(category impacts.ImpactCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockImpactRepo) Delete(category impacts.ImpactCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func TestCreateProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockImpactRepo := new(MockImpactRepo)
	productService := NewProductService(mockProductRepo, mockImpactRepo)

	product := products.Product{
		ImpactCategories: []products.ProductImpactCategory{
			{ImpactCategoryID: "1"},
		},
		Images: []products.ProductImage{
			{AlbumsURL: "image1.jpg"},
		},
	}

	mockImpactRepo.On("GetByID", "1").Return(impacts.ImpactCategory{ID: "1"}, nil)
	mockProductRepo.On("Create", mock.Anything).Return(nil)

	err := productService.Create(product)
	assert.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
	mockImpactRepo.AssertExpectations(t)
}

func TestCreateProductInvalidImpact(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockImpactRepo := new(MockImpactRepo)
	productService := NewProductService(mockProductRepo, mockImpactRepo)

	product := products.Product{
		ImpactCategories: []products.ProductImpactCategory{
			{ImpactCategoryID: "1"},
		},
		Images: []products.ProductImage{
			{AlbumsURL: "image1.jpg"},
		},
	}

	mockImpactRepo.On("GetByID", "1").Return(impacts.ImpactCategory{}, nil)

	err := productService.Create(product)
	assert.Error(t, err)
	assert.Equal(t, constant.ErrCreateProduct, err)
	mockImpactRepo.AssertExpectations(t)
}

func TestGetAllByPage(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockImpactRepo := new(MockImpactRepo)
	productService := NewProductService(mockProductRepo, mockImpactRepo)

	mockProducts := []products.Product{
		{ID: "1"},
		{ID: "2"},
	}

	mockProductRepo.On("GetAllByPage", 1, "", "").Return(mockProducts, 2, nil)
	mockProductRepo.On("GetTotalProduct").Return(2, nil)

	products, total, totalProduct, err := productService.GetAllByPage(1, "", "")
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Equal(t, 2, totalProduct)
	assert.Equal(t, mockProducts, products)
	mockProductRepo.AssertExpectations(t)
}

func TestGetByCategory(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockImpactRepo := new(MockImpactRepo)
	productService := NewProductService(mockProductRepo, mockImpactRepo)

	mockProducts := []products.Product{
		{ID: "1"},
		{ID: "2"},
	}

	mockProductRepo.On("GetByCategory", "category1", 1, "", "").Return(mockProducts, 2, nil)
	mockProductRepo.On("GetTotalProduct").Return(2, nil)

	products, total, totalProduct, err := productService.GetByCategory("category1", 1, "", "")
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Equal(t, 2, totalProduct)
	assert.Equal(t, mockProducts, products)
	mockProductRepo.AssertExpectations(t)
}

func TestGetById(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockImpactRepo := new(MockImpactRepo)
	productService := NewProductService(mockProductRepo, mockImpactRepo)

	mockProduct := products.Product{ID: "1"}

	mockProductRepo.On("GetById", "1").Return(mockProduct, nil)

	product, err := productService.GetById("1")
	assert.NoError(t, err)
	assert.Equal(t, mockProduct, product)
	mockProductRepo.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockImpactRepo := new(MockImpactRepo)
	productService := NewProductService(mockProductRepo, mockImpactRepo)

	product := products.Product{
		ID: "1",
		ImpactCategories: []products.ProductImpactCategory{
			{ImpactCategoryID: "1"},
		},
		Images: []products.ProductImage{
			{AlbumsURL: "image1.jpg"},
		},
	}

	mockImpactRepo.On("GetByID", "1").Return(impacts.ImpactCategory{ID: "1"}, nil)
	mockProductRepo.On("Update", mock.Anything).Return(nil)

	err := productService.Update(product)
	assert.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
	mockImpactRepo.AssertExpectations(t)
}

func TestDeleteProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockImpactRepo := new(MockImpactRepo)
	productService := NewProductService(mockProductRepo, mockImpactRepo)

	mockProductRepo.On("Delete", "1").Return(nil)

	err := productService.Delete("1")
	assert.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
}
