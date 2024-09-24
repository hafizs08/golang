package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product represents the product entity with validation tags
type Product struct {
	ID          string             `json:"id"`
	MongoID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" validate:"required,max=100"`
	Description string             `json:"description" validate:"required"`
	Price       float64            `json:"price" validate:"required,gt=0"`
	Stock       int                `json:"stock" validate:"required"`
}

// ProductRepository defines the methods for interacting with products in the repository

type ProductRepository interface {
	Create(product *Product) error
	GetAllProducts() ([]Product, error)
	GetProductById(id string) (*Product, error)
	UpdateProduct(id string, product *Product) error
	DeleteProduct(id string) error
}
