package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/martbul/playground_microservices/services/product-service/models"
	"github.com/martbul/playground_microservices/services/product-service/repository"
)

type ProductService interface {
	CreateProduct(req *models.CreateProductRequest, createdBy string) (*models.Product, error)
	GetProduct(id string) (*models.Product, error)
	UpdateProduct(id string, req *models.UpdateProductRequest) (*models.Product, error)
	DeleteProduct(id string) error
	ListProducts(filter *models.ProductFilter, pagination *models.PaginationRequest) ([]*models.Product, *models.PaginationResponse, error)
	SearchProducts(filter *models.SearchFilter, pagination *models.PaginationRequest) ([]*models.Product, *models.PaginationResponse, error)
	GetCategories() ([]*models.Category, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) CreateProduct(req *models.CreateProductRequest, createdBy string) (*models.Product, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	if req.StockQuantity < 0 {
		return nil, fmt.Errorf("stock quantity cannot be negative")
	}

	// Check if SKU already exists (if provided)
	if req.SKU != "" {
		existingProduct, err := s.productRepo.GetBySKU(req.SKU)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing SKU: %w", err)
		}
		if existingProduct != nil {
			return nil, fmt.Errorf("product with SKU already exists")
		}
	}

	// Generate SKU if not provided
	sku := req.SKU
	if sku == "" {
		sku = s.generateSKU(req.Name)
	}

	product := &models.Product{
		Name:          strings.TrimSpace(req.Name),
		Description:   strings.TrimSpace(req.Description),
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Category:      req.Category,
		ImageURL:      req.ImageURL,
		SKU:           sku,
		IsActive:      true,
		CreatedBy:     createdBy,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *productService) GetProduct(id string) (*models.Product, error) {
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	return product, nil
}

func (s *productService) UpdateProduct(id string, req *models.UpdateProductRequest) (*models.Product, error) {
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	// Get existing product
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Update fields if provided
	if req.Name != "" {
		product.Name = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		product.Description = strings.TrimSpace(req.Description)
	}
	if req.Price >= 0 {
		product.Price = req.Price
	}
	if req.StockQuantity >= 0 {
		product.StockQuantity = req.StockQuantity
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.SKU != "" && req.SKU != product.SKU {
		// Check if new SKU already exists
		existingProduct, err := s.productRepo.GetBySKU(req.SKU)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing SKU: %w", err)
		}
		if existingProduct != nil && existingProduct.ID != id {
			return nil, fmt.Errorf("SKU already exists")
		}
		product.SKU = req.SKU
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *productService) DeleteProduct(id string) error {
	if id == "" {
		return fmt.Errorf("product ID is required")
	}

	// Check if product exists
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	if product == nil {
		return fmt.Errorf("product not found")
	}

	if err := s.productRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (s *productService) ListProducts(filter *models.ProductFilter, pagination *models.PaginationRequest) ([]*models.Product, *models.PaginationResponse, error) {
	// Set default pagination values
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}
	if pagination.Limit > 100 {
		pagination.Limit = 100 // Max limit
	}

	// Default to show only active products if not specified
	if filter == nil {
		filter = &models.ProductFilter{}
	}
	if filter.IsActive == nil {
		active := true
		filter.IsActive = &active
	}

	products, paginationResp, err := s.productRepo.List(filter, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list products: %w", err)
	}

	return products, paginationResp, nil
}

func (s *productService) SearchProducts(filter *models.SearchFilter, pagination *models.PaginationRequest) ([]*models.Product, *models.PaginationResponse, error) {
	// Set default pagination values
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}
	if pagination.Limit > 100 {
		pagination.Limit = 100 // Max limit
	}

	// Default to show only active products if not specified
	if filter == nil {
		filter = &models.SearchFilter{}
	}
	if filter.IsActive == nil {
		active := true
		filter.IsActive = &active
	}

	products, paginationResp, err := s.productRepo.Search(filter, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search products: %w", err)
	}

	return products, paginationResp, nil
}

func (s *productService) GetCategories() ([]*models.Category, error) {
	categories, err := s.productRepo.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

func (s *productService) generateSKU(productName string) string {
	// Simple SKU generation: first 3 letters + timestamp
	name := strings.ReplaceAll(strings.ToUpper(productName), " ", "")
	if len(name) > 3 {
		name = name[:3]
	}
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s%d", name, timestamp)
}
