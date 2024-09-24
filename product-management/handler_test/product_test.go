package handler_test

import (
	"bytes"
	"net/http/httptest"
	"product-management/internal/domain"
	"testing"

	"github.com/hafizs08/golang/product-management/internal/handler"
	mockRepo "github.com/hafizs08/golang/product-management/internal/repository/mock" // Ensure this path is correct
	"github.com/hafizs08/golang/product-management/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_CreateProduct(t *testing.T) {
	repoMock := new(mockRepo.MockProductRepository) // Updated variable name
	productService := service.NewProductService(repoMock)
	productHandler := handler.NewProductHandler(productService)

	app := fiber.New()
	app.Post("/products", productHandler.CreateProduct)

	repoMock.On("Create", mock.Anything).Return(nil)

	// Simulate a request to create a product
	req := httptest.NewRequest("POST", "/products", bytes.NewBuffer([]byte(`{"name":"Product 1","description":"Desc 1","price":10,"stock":5}`)))
	resp, _ := app.Test(req)

	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestProductHandler_GetAllProducts(t *testing.T) {
	repoMock := new(mockRepo.MockProductRepository) // Updated variable name
	productService := service.NewProductService(repoMock)
	productHandler := handler.NewProductHandler(productService)

	app := fiber.New()
	app.Get("/products", productHandler.GetAllProducts)

	// Simulate a request to get all products
	repoMock.On("GetAll").Return([]domain.Product{{ID: "1", Name: "Product 1", Price: 10, Stock: 5}}, nil)

	req := httptest.NewRequest("GET", "/products", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
