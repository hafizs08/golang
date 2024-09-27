package handler

import (
	"log"
	"net/http" // Tambahkan ini untuk memperbaiki error 'undefined: http'
	"product-management/internal/domain"
	"product-management/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) GetMySQLProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetMySQLProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to fetch products",
		})
	}
	return c.JSON(products)
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Create product in database
	err := h.productService.CreateProduct(&product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	// Use messenger to send success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product successfully added",
		"product": product,
	})
}

// GetAllProducts retrieves all products from MySQL and MongoDB
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		log.Printf("Error retrieving products: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

// GetProductByID retrieves a product by its ID
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := h.productService.GetProductById(id)
	if err != nil {
		log.Printf("Error retrieving product by ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve product",
		})
	}

	// Tambahkan kondisi jika produk tidak ditemukan
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product with ID " + id + " not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// UpdateProduct updates a product by its ID
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product domain.Product

	if err := c.BodyParser(&product); err != nil {
		log.Printf("Error parsing product input: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Cek apakah produk dengan ID yang diberikan ada
	existingProduct, err := h.productService.GetProductById(id)
	if err != nil {
		log.Printf("Error retrieving product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve product",
		})
	}
	if existingProduct == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product with ID " + id + " not found",
		})
	}

	err = h.productService.UpdateProduct(id, &product)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product successfully updated",
		"product": product,
	})
}

// DeleteProduct deletes a product by its ID
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	// Cek apakah produk dengan ID yang diberikan ada
	existingProduct, err := h.productService.GetProductById(id)
	if err != nil {
		log.Printf("Error retrieving product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve product",
		})
	}
	if existingProduct == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product with ID " + id + " not found",
		})
	}

	err = h.productService.DeleteProduct(id)
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product successfully deleted",
	})
}

func (h *ProductHandler) GetMongoDBProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetMongoDBProducts()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(products)
}
