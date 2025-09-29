package clients

import (
	"context"
	"fmt"
)

// ProductClient handles product-related API calls
type ProductClient struct {
	apiClient *APIClient
}

// NewProductClient creates a new product client
func NewProductClient(apiClient *APIClient) *ProductClient {
	return &ProductClient{
		apiClient: apiClient,
	}
}

// ListProducts lists products with optional filters
func (c *ProductClient) ListProducts(ctx context.Context, params ProductListParams) (*ProductListResponse, error) {
	return c.apiClient.ListProducts(ctx, params)
}

// SearchProducts searches for products
func (c *ProductClient) SearchProducts(ctx context.Context, params ProductSearchParams) (*ProductListResponse, error) {
	return c.apiClient.SearchProducts(ctx, params)
}

// GetProduct gets a single product by ID
func (c *ProductClient) GetProduct(ctx context.Context, id string) (*ProductResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}
	return c.apiClient.GetProduct(ctx, id)
}

// CreateProduct creates a new product (requires authentication)
func (c *ProductClient) CreateProduct(ctx context.Context, token string, req CreateProductRequest) (*ProductResponse, error) {
	if token == "" {
		return nil, fmt.Errorf("authentication token is required")
	}
	return c.apiClient.CreateProduct(ctx, token, req)
}

// UpdateProduct updates an existing product (requires authentication)
func (c *ProductClient) UpdateProduct(ctx context.Context, token string, id string, req UpdateProductRequest) (*ProductResponse, error) {
	if token == "" {
		return nil, fmt.Errorf("authentication token is required")
	}
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}
	return c.apiClient.UpdateProduct(ctx, token, id, req)
}

// DeleteProduct deletes a product (requires authentication)
func (c *ProductClient) DeleteProduct(ctx context.Context, token string, id string) (*APIResponse, error) {
	if token == "" {
		return nil, fmt.Errorf("authentication token is required")
	}
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}
	return c.apiClient.DeleteProduct(ctx, token, id)
}

// GetCategories gets all available categories
func (c *ProductClient) GetCategories(ctx context.Context) (*CategoriesResponse, error) {
	return c.apiClient.GetCategories(ctx)
}

// Helper methods for common operations

// GetProductsByCategory gets products filtered by category
func (c *ProductClient) GetProductsByCategory(ctx context.Context, category string, page, limit int) (*ProductListResponse, error) {
	params := ProductListParams{
		Category: category,
		Page:     page,
		Limit:    limit,
	}
	return c.ListProducts(ctx, params)
}

// GetFeaturedProducts gets featured products (first N products)
func (c *ProductClient) GetFeaturedProducts(ctx context.Context, limit int) (*ProductListResponse, error) {
	if limit <= 0 {
		limit = 6 // Default featured products count
	}
	
	params := ProductListParams{
		Page:  1,
		Limit: limit,
	}
	return c.ListProducts(ctx, params)
}

// SearchProductsByName searches products by name
func (c *ProductClient) SearchProductsByName(ctx context.Context, query string, page, limit int) (*ProductListResponse, error) {
	params := ProductSearchParams{
		Query: query,
		Page:  page,
		Limit: limit,
	}
	return c.SearchProducts(ctx, params)
}

// GetProductsInPriceRange gets products within a price range
func (c *ProductClient) GetProductsInPriceRange(ctx context.Context, minPrice, maxPrice float64, page, limit int) (*ProductListResponse, error) {
	params := ProductSearchParams{
		Page:     page,
		Limit:    limit,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}
	return c.SearchProducts(ctx, params)
}

// IsProductAvailable checks if a product is available (in stock and active)
func (c *ProductClient) IsProductAvailable(ctx context.Context, productID string) (bool, error) {
	resp, err := c.GetProduct(ctx, productID)
	if err != nil {
		return false, err
	}

	if !resp.Response.Success {
		return false, nil
	}

	product := resp.Product
	return product.IsActive && product.StockQuantity > 0, nil
}

// GetProductStock gets the stock quantity for a product
func (c *ProductClient) GetProductStock(ctx context.Context, productID string) (int32, error) {
	resp, err := c.GetProduct(ctx, productID)
	if err != nil {
		return 0, err
	}

	if !resp.Response.Success {
		return 0, fmt.Errorf("product not found")
	}

	return resp.Product.StockQuantity, nil
}