// internal/repository/mysql/mysql_product_repository
package mysql

import (
	"database/sql"
	"product-management/internal/domain"
)

type MySQLProductRepository struct {
	db *sql.DB
}

func NewMySQLProductRepository(db *sql.DB) *MySQLProductRepository {
	return &MySQLProductRepository{db: db}
}

// Create method
func (r *MySQLProductRepository) Create(product *domain.Product) error {
	query := "INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)"
	_, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Stock)
	return err
}

// GetAllProducts method
// GetAllProducts method
func (r *MySQLProductRepository) GetAllProducts() ([]domain.Product, error) {
	rows, err := r.db.Query("SELECT id, name, description, price, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product) // Menggunakan product bukan &product
	}

	return products, nil
}

// GetProductById method
func (r *MySQLProductRepository) GetProductById(id string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateProduct method
func (r *MySQLProductRepository) UpdateProduct(id string, product *domain.Product) error {
	_, err := r.db.Exec(`
		UPDATE products
		SET name = ?, description = ?, price = ?, stock = ?
		WHERE id = ?`,
		product.Name, product.Description, product.Price, product.Stock, id)
	return err
}

// DeleteProduct method
func (r *MySQLProductRepository) DeleteProduct(id string) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id=?", id)
	return err
}
