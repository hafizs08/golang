package main

import (
	"context"
	"database/sql"
	"log"
	"product-management/internal/handler"
	"product-management/internal/repository/mongodb"
	"product-management/internal/repository/mysql"
	"product-management/internal/service"

	_ "github.com/go-sql-driver/mysql" // pastikan driver mysql diimport
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MySQL setup
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/produk")
	if err != nil {
		log.Fatal(err)
	}
	mysqlRepo := mysql.NewMySQLProductRepository(db)

	// MongoDB setup
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	mongoCollection := mongoClient.Database("productDB").Collection("products")
	mongoRepo := mongodb.NewMongoDBProductRepository(mongoCollection)

	// Service and handler setup
	productService := service.NewProductService(mysqlRepo, mongoRepo)
	productHandler := handler.NewProductHandler(productService)

	// Fiber setup
	app := fiber.New()

	// CRUD Routes
	app.Post("/products", productHandler.CreateProduct)
	app.Get("/products", productHandler.GetAllProducts)
	app.Get("/products/:id", productHandler.GetProductByID)
	app.Put("/products/:id", productHandler.UpdateProduct)
	app.Delete("/products/:id", productHandler.DeleteProduct)
	app.Get("/Get", productHandler.GetMySQLProducts)

	log.Fatal(app.Listen(":3000"))
}
