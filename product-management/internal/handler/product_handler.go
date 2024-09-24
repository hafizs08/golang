package handler

import (
	"log"
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

	err := h.productService.UpdateProduct(id, &product)
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

	err := h.productService.DeleteProduct(id)
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
