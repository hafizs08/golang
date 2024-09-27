package product_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"product-management/internal/domain"
	"product-management/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock MySQL Repository
type mockMySQLRepo struct{}

func (m *mockMySQLRepo) GetAllProducts() ([]domain.Product, error) {
	return nil, nil
}

func (m *mockMySQLRepo) Create(product *domain.Product) error {
	return nil
}

func (m *mockMySQLRepo) GetProductById(id string) (*domain.Product, error) {
	return nil, nil
}

func (m *mockMySQLRepo) UpdateProduct(id string, product *domain.Product) error {
	return nil
}

func (m *mockMySQLRepo) DeleteProduct(id string) error {
	return nil
}

// Mock MongoDB Repository
type mockMongoRepo struct{}

func (m *mockMongoRepo) GetAllProducts() ([]domain.Product, error) {
	return []domain.Product{
		{
			ID:          "66f268c84b7253868ad8312e",
			Name:        "kecap",
			Description: "asus",
			Price:       10000,
			Stock:       20,
		},
		{
			ID:          "66f269094b7253868ad8312f",
			Name:        "ritonga",
			Description: "New product description",
			Price:       49.99,
			Stock:       20,
		},
		{
			ID:          "66f28741deae451485368e74",
			Name:        "kiki",
			Description: "New product description",
			Price:       49.99,
			Stock:       20,
		},
	}, nil
}

func (m *mockMongoRepo) Create(product *domain.Product) error {
	return nil
}

func (m *mockMongoRepo) GetProductById(id string) (*domain.Product, error) {
	if id == "66f268c84b7253868ad8312e" {
		return &domain.Product{
			ID:          "66f268c84b7253868ad8312e",
			Name:        "kecap",
			Description: "asus",
			Price:       10000,
			Stock:       20,
		}, nil
	}
	return nil, errors.New("product not found")
}

func (m *mockMongoRepo) UpdateProduct(id string, product *domain.Product) error {
	return nil
}

func (m *mockMongoRepo) DeleteProduct(id string) error {
	return nil
}

// Test untuk mendapatkan semua produk dari MongoDB
func TestGetMongoDBProducts(t *testing.T) {
	// Membuat mock repository MySQL dan MongoDB
	mockMySQLRepository := &mockMySQLRepo{}
	mockMongoRepository := &mockMongoRepo{}

	// Membuat service dengan repository mock
	productService := service.NewProductService(mockMySQLRepository, mockMongoRepository)

	// Membuat permintaan HTTP
	req, err := http.NewRequest("GET", "/mongodb-products", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Menggunakan httptest untuk membuat response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products, err := productService.GetMongoDBProducts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(products)
	})

	// Menjalankan handler dengan permintaan di atas
	handler.ServeHTTP(rr, req)

	// Memeriksa apakah status kode yang dikembalikan adalah 200
	assert.Equal(t, http.StatusOK, rr.Code)

	// Memeriksa apakah hasil JSON sesuai dengan yang diharapkan
	expected := `[{"_id":"66f268c84b7253868ad8312e","name":"kecap","description":"asus","price":10000,"stock":20},{"_id":"66f269094b7253868ad8312f","name":"ritonga","description":"New product description","price":49.99,"stock":20},{"_id":"66f28741deae451485368e74","name":"kiki","description":"New product description","price":49.99,"stock":20}]`
	assert.JSONEq(t, expected, rr.Body.String())
}
