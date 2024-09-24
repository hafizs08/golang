package service

import "product-management/internal/domain"

type ProductService struct {
	mysqlRepo domain.ProductRepository
	mongoRepo domain.ProductRepository
}

func NewProductService(mysqlRepo, mongoRepo domain.ProductRepository) *ProductService {
	return &ProductService{
		mysqlRepo: mysqlRepo,
		mongoRepo: mongoRepo,
	}
}

func (s *ProductService) CreateProduct(product *domain.Product) error {
	err := s.mysqlRepo.Create(product)
	if err != nil {
		return err
	}
	return s.mongoRepo.Create(product)
}

func (s *ProductService) GetAllProducts() ([]domain.Product, error) {
	mysqlProducts, err := s.mysqlRepo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	mongoProducts, err := s.mongoRepo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	combinedProducts := append(mysqlProducts, mongoProducts...)
	return combinedProducts, nil
}

func (s *ProductService) GetProductById(id string) (*domain.Product, error) {
	// Ambil dari MySQL
	product, err := s.mysqlRepo.GetProductById(id)
	if err != nil {
		return nil, err
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
