package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"product-management/internal/domain"
	"product-management/internal/handler"
	"product-management/internal/service"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository for testing
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetAllProducts() ([]domain.Product, error) {
	args := m.Called()
	products, ok := args.Get(0).([]domain.Product)
	if !ok {
		return nil, args.Error(1)
	}
	return products, args.Error(1)
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

func TestGetAllProducts(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo, nil)
	productHandler := handler.NewProductHandler(productService)

	app := fiber.New()
	app.Get("/products", productHandler.GetAllProducts)

	// Test: Mengembalikan daftar produk
	mockRepo.On("GetAllProducts").Return([]domain.Product{
		{ID: "1", Name: "Product 1", Description: "Desc 1", Price: 10.0, Stock: 5},
		{ID: "2", Name: "Product 2", Description: "Desc 2", Price: 20.0, Stock: 10},
	}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test: Mengembalikan daftar kosong jika tidak ada produk
	mockRepo.On("GetAllProducts").Return([]domain.Product{}, nil).Once()

	req = httptest.NewRequest(http.MethodGet, "/products", nil)
	resp, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
func TestCreateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo, nil)
	productHandler := handler.NewProductHandler(productService)

	app := fiber.New()
	app.Post("/products", productHandler.CreateProduct)

	// Test: Menghasilkan produk baru
	product := &domain.Product{
		ID: "1", Name: "Product 1", Description: "Desc 1", Price: 10.0, Stock: 5,
	}
	mockRepo.On("Create", product).Return(nil)

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Test: Mengembalikan error jika input tidak valid
	req = httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte(`{"name":""}`)))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req, -1)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
func TestUpdateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo, nil)
	productHandler := handler.NewProductHandler(productService)

	app := fiber.New()
	app.Put("/products/:id", productHandler.UpdateProduct)

	// Test: Mengupdate produk yang ada
	product := &domain.Product{
		ID: "1", Name: "Product Updated", Description: "Desc Updated", Price: 15.0, Stock: 8,
	}
	mockRepo.On("UpdateProduct", "1", product).Return(nil)

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test: Mengembalikan error jika produk tidak ditemukan
	mockRepo.On("UpdateProduct", "2", product).Return(fiber.ErrNotFound)
	req = httptest.NewRequest(http.MethodPut, "/products/2", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req, -1)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetProductById(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo, nil)
	productHandler := handler.NewProductHandler(productService)

	app := fiber.New()
	app.Get("/products/:id", productHandler.GetProductByID)

	// Test: Mengembalikan produk berdasarkan ID
	product := &domain.Product{
		ID: "1", Name: "Product 1", Description: "Desc 1", Price: 10.0, Stock: 5,
	}
	mockRepo.On("GetProductById", "1").Return(product, nil)

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test: Mengembalikan error jika produk tidak ditemukan
	mockRepo.On("GetProductById", "2").Return(nil, fiber.ErrNotFound)

	req = httptest.NewRequest(http.MethodGet, "/products/2", nil)
	resp, _ = app.Test(req, -1)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDeleteProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo, nil)
	productHandler := handler.NewProductHandler(productService)

	app := fiber.New()
	app.Delete("/products/:id", productHandler.DeleteProduct)

	// Test: Menghapus produk berdasarkan ID
	mockRepo.On("DeleteProduct", "1").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test: Mengembalikan error jika produk tidak ditemukan
	mockRepo.On("DeleteProduct", "2").Return(fiber.ErrNotFound)

	req = httptest.NewRequest(http.MethodDelete, "/products/2", nil)
	resp, _ = app.Test(req, -1)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
