// internal/repository/mock/mock_product_repository.go
package mock

import (
	"product-management/internal/domain"

	"github.com/stretchr/testify/mock"
)

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
