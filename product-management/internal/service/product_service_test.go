package service_test

import (
	"product-management/internal/domain"
	"product-management/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock untuk MySQL dan MongoDB Repository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetAllProducts() ([]domain.Product, error) {
	args := m.Called()
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *MockProductRepository) GetProductById(id string) (*domain.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateProduct(id string, product *domain.Product) error {
	args := m.Called(id, product)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteProduct(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateProduct(t *testing.T) {
	mockMySQLRepo := new(MockProductRepository)
	mockMongoRepo := new(MockProductRepository)
	productService := service.NewProductService(mockMySQLRepo, mockMongoRepo)

	product := &domain.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		Stock:       10,
	}

	// Setup mocks untuk mengembalikan tidak ada error saat Create dipanggil
	mockMySQLRepo.On("Create", product).Return(nil)
	mockMongoRepo.On("Create", product).Return(nil)

	err := productService.CreateProduct(product)

	// Verifikasi hasil
	assert.NoError(t, err)
	mockMySQLRepo.AssertExpectations(t)
	mockMongoRepo.AssertExpectations(t)
}

func TestGetAllProducts(t *testing.T) {
	mockMySQLRepo := new(MockProductRepository)
	mockMongoRepo := new(MockProductRepository)
	productService := service.NewProductService(mockMySQLRepo, mockMongoRepo)

	mysqlProducts := []domain.Product{
		{Name: "Product 1", Description: "Desc 1", Price: 10.0, Stock: 10},
	}
	mongoProducts := []domain.Product{
		{Name: "Product 2", Description: "Desc 2", Price: 20.0, Stock: 5},
	}

	// Setup mocks
	mockMySQLRepo.On("GetAllProducts").Return(mysqlProducts, nil)
	mockMongoRepo.On("GetAllProducts").Return(mongoProducts, nil)

	result, err := productService.GetAllProducts()

	// Verifikasi hasil
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockMySQLRepo.AssertExpectations(t)
	mockMongoRepo.AssertExpectations(t)
}

func TestGetProductById(t *testing.T) {
	mockMySQLRepo := new(MockProductRepository)
	mockMongoRepo := new(MockProductRepository)
	productService := service.NewProductService(mockMySQLRepo, mockMongoRepo)

	product := &domain.Product{Name: "Test Product", Description: "Test Description", Price: 100.0, Stock: 10}

	// Test case untuk produk yang ditemukan di MySQL
	mockMySQLRepo.On("GetProductById", "1").Return(product, nil)

	result, err := productService.GetProductById("1")

	assert.NoError(t, err)
	assert.Equal(t, product, result)
	mockMySQLRepo.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	mockMySQLRepo := new(MockProductRepository)
	mockMongoRepo := new(MockProductRepository)
	productService := service.NewProductService(mockMySQLRepo, mockMongoRepo)

	product := &domain.Product{Name: "Updated Product", Description: "Updated Description", Price: 150.0, Stock: 8}

	// Setup mocks
	mockMySQLRepo.On("UpdateProduct", "1", product).Return(nil)
	mockMongoRepo.On("UpdateProduct", "1", product).Return(nil)

	err := productService.UpdateProduct("1", product)

	// Verifikasi hasil
	assert.NoError(t, err)
	mockMySQLRepo.AssertExpectations(t)
	mockMongoRepo.AssertExpectations(t)
}

func TestDeleteProduct(t *testing.T) {
	mockMySQLRepo := new(MockProductRepository)
	mockMongoRepo := new(MockProductRepository)
	productService := service.NewProductService(mockMySQLRepo, mockMongoRepo)

	// Setup mocks
	mockMySQLRepo.On("DeleteProduct", "1").Return(nil)
	mockMongoRepo.On("DeleteProduct", "1").Return(nil)

	err := productService.DeleteProduct("1")

	// Verifikasi hasil
	assert.NoError(t, err)
	mockMySQLRepo.AssertExpectations(t)
	mockMongoRepo.AssertExpectations(t)
}
