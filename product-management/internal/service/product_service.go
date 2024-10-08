// internal/service/product_service.go
package service

import (
	"errors"
	"log"
	"product-management/internal/domain"
)

type ProductService struct {
	mysqlRepo domain.ProductRepository
	mongoRepo domain.ProductRepository
}

// func NewProductService(mysqlRepo, mongoRepo domain.ProductRepository) *ProductService {
// 	return &ProductService{
// 		mysqlRepo: mysqlRepo,
// 		mongoRepo: mongoRepo,
// 	}
// }

func NewProductService(mysqlRepo, mongoRepo domain.ProductRepository) *ProductService {
	if mysqlRepo == nil {
		log.Fatal("MySQL repository cannot be nil")
	}
	if mongoRepo == nil {
		log.Fatal("MongoDB repository cannot be nil")
	}
	return &ProductService{
		mysqlRepo: mysqlRepo,
		mongoRepo: mongoRepo,
	}
}

func (s *ProductService) CreateProduct(product *domain.Product) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}
	err := s.mysqlRepo.Create(product)
	if err != nil {
		return err
	}
	return s.mongoRepo.Create(product)
}

func (s *ProductService) GetAllProducts() ([]domain.Product, error) {
	mysqlProducts, err := s.mysqlRepo.GetAllProducts()
	if err != nil {
		log.Printf("Error getting MySQL products: %v", err)
		return nil, err
	}

	mongoProducts, err := s.mongoRepo.GetAllProducts()
	if err != nil {
		log.Printf("Error getting MongoDB products: %v", err)
		return nil, err
	}

	combinedProducts := append(mysqlProducts, mongoProducts...)
	return combinedProducts, nil
}
func (s *ProductService) GetProductById(id string) (*domain.Product, error) {
	// Coba ambil dari MySQL
	product, err := s.mysqlRepo.GetProductById(id)
	if err == nil { // Jika berhasil, kembalikan produk
		return product, nil
	}

	// Jika tidak ditemukan di MySQL, coba ambil dari MongoDB
	product, err = s.mongoRepo.GetProductById(id)
	if err != nil {
		return nil, err // Jika masih tidak ditemukan, kembalikan error
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(id string, product *domain.Product) error {
	err := s.mysqlRepo.UpdateProduct(id, product) // Update in MySQL
	if err != nil {
		return err
	}
	return s.mongoRepo.UpdateProduct(id, product) // Update in MongoDB
}

func (s *ProductService) DeleteProduct(id string) error {
	// Hapus dari MySQL
	err := s.mysqlRepo.DeleteProduct(id)
	if err != nil {
		return err
	}
	// Hapus dari MongoDB
	return s.mongoRepo.DeleteProduct(id)
}
func (s *ProductService) GetMySQLProducts() ([]domain.Product, error) {
	return s.mysqlRepo.GetAllProducts()
}

func (s *ProductService) GetMongoDBProducts() ([]domain.Product, error) {
	return s.mongoRepo.GetAllProducts()
}
