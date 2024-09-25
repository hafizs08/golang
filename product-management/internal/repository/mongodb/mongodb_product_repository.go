// internal/repository/mongodb/mongodb_product_repository
package mongodb

import (
	"context"
	"product-management/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDBProductRepository struct
type MongoDBProductRepository struct {
	db *mongo.Collection // Menggunakan field db
}

// NewMongoDBProductRepository creates a new instance of MongoDBProductRepository
func NewMongoDBProductRepository(db *mongo.Collection) *MongoDBProductRepository {
	return &MongoDBProductRepository{db: db}
}

// Create method
func (r *MongoDBProductRepository) Create(product *domain.Product) error {
	_, err := r.db.InsertOne(context.TODO(), product)
	return err
}

// GetAllProducts method
func (r *MongoDBProductRepository) GetAllProducts() ([]domain.Product, error) {
	var products []domain.Product
	cursor, err := r.db.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// GetProductById method
// Mengubah penggunaan r.collection menjadi r.db
func (r *MongoDBProductRepository) GetProductById(id string) (*domain.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product domain.Product
	// Mengganti r.collection dengan r.db
	err = r.db.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// UpdateProduct method
func (r *MongoDBProductRepository) UpdateProduct(id string, product *domain.Product) error {
	_, err := r.db.UpdateOne(context.TODO(), bson.M{"id": product.ID}, bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
		},
	})
	return err
}

// DeleteProduct method
func (r *MongoDBProductRepository) DeleteProduct(id string) error {
	_, err := r.db.DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}
