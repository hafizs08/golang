package service_test

import (
	"errors"
	"product-management/internal/domain"
	"product-management/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock of ProductRepository for testing purposes
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockRepository) GetAllProducts() ([]domain.Product, error) {
	args := m.Called()
	if products, ok := args.Get(0).([]domain.Product); ok {
		return products, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetProductById(id string) (*domain.Product, error) {
	args := m.Called(id)
	if product, ok := args.Get(0).(*domain.Product); ok {
		return product, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) UpdateProduct(id string, product *domain.Product) error {
	args := m.Called(id, product)
	return args.Error(0)
}

func (m *MockRepository) DeleteProduct(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test CreateProduct function
func TestCreateProduct(t *testing.T) {
	mockMySQLRepo := new(MockRepository)
	mockMongoRepo := new(MockRepository)
	productService := service.NewProductService(mockMySQLRepo, mockMongoRepo)

	product := &domain.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
		Stock:       10,
	}

	t.Run("successfully creates product", func(t *testing.T) {
		mockMySQLRepo.On("Create", product).Return(nil)
		mockMongoRepo.On("Create", product).Return(nil)

		err := productService.CreateProduct(product)
		assert.NoError(t, err)
		mockMySQLRepo.AssertExpectations(t)
		mockMongoRepo.AssertExpectations(t)
	})

	t.Run("fails to create product in MySQL", func(t *testing.T) {
		mockMySQLRepo.On("Create", product).Return(errors.New("MySQL error"))

		err := productService.CreateProduct(product)
		assert.Error(t, err)
		mockMySQLRepo.AssertExpectations(t)
	})
}

// Test GetAllProducts function
func TestGetAllProducts(t *testing.T) {
	mockMySQLRepo := new(MockRepository)
	mockMongoRepo := new(MockRepository)
	productService := service.NewProductService(mockMySQLRepo, mockMongoRepo)

	t.Run("successfully gets all products", func(t *testing.T) {
		mysqlProducts := []domain.Product{
			{Name: "MySQL Product", Price: 100},
		}
		mongoProducts := []domain.Product{
			{Name: "Mongo Product", Price: 200},
		}

		mockMySQLRepo.On("GetAllProducts").Return(mysqlProducts, nil)
		mockMongoRepo.On("GetAllProducts").Return(mongoProducts, nil)

		products, err := productService.GetAllProducts()
		assert.NoError(t, err)
		assert.Len(t, products, 2)
		mockMySQLRepo.AssertExpectations(t)
		mockMongoRepo.AssertExpectations(t)
	})

	t.Run("fails to get products from MySQL", func(t *testing.T) {
		mockMySQLRepo.On("GetAllProducts").Return(nil, errors.New("MySQL error"))

		_, err := productService.GetAllProducts()
		assert.Error(t, err)
		mockMySQLRepo.AssertExpectations(t)
	})
}

// Test GetProductById function
func TestGetProductById(t *testing.T) {
	mockMySQLRepo := new(MockRepository)
	mockMongoRepo := new(MockRepository)
	productService := service.NewProductService(mockMySQLRepo, mockMongoRepo)

	productID := "12345"
	product := &domain.Product{
		ID:    productID,
		Name:  "Test Product",
		Price: 100,
	}

	t.Run("successfully gets product from MySQL", func(t *testing.T) {
		mockMySQLRepo.On("GetProductById", productID).Return(product, nil)

		result, err := productService.GetProductById(productID)
		assert.NoError(t, err)
		assert.Equal(t, product, result)
		mockMySQLRepo.AssertExpectations(t)
	})

	t.Run("gets product from MongoDB when not found in MySQL", func(t *testing.T) {
		mockMySQLRepo.On("GetProductById", productID).Return(nil, errors.New("not found"))
		mockMongoRepo.On("GetProductById", productID).Return(product, nil)

		result, err := productService.GetProductById(productID)
		assert.NoError(t, err)
		assert.Equal(t, product, result)
		mockMongoRepo.AssertExpectations(t)
	})
}

// Test DeleteProduct function
func TestDeleteProduct(t *testing.T) {
	mockMySQLRepo := new(MockRepository)
	mockMongoRepo := new(MockRepository)
	productService := service.NewProductService(mockMySQLRepo, mockMongoRepo)

	productID := "12345"

	t.Run("successfully deletes product", func(t *testing.T) {
		mockMySQLRepo.On("DeleteProduct", productID).Return(nil)
		mockMongoRepo.On("DeleteProduct", productID).Return(nil)

		err := productService.DeleteProduct(productID)
		assert.NoError(t, err)
		mockMySQLRepo.AssertExpectations(t)
		mockMongoRepo.AssertExpectations(t)
	})

	t.Run("fails to delete product in MySQL", func(t *testing.T) {
		mockMySQLRepo.On("DeleteProduct", productID).Return(errors.New("MySQL error"))

		err := productService.DeleteProduct(productID)
		assert.Error(t, err)
		mockMySQLRepo.AssertExpectations(t)
	})
}
