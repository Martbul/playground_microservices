package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// APIClient handles all API calls to the gateway
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Common types
type APIResponse struct {
	Response Response `json:"response"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}


type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	Category      string    `json:"category"`
	ImageURL      string    `json:"image_url"`
	StockQuantity int32     `json:"stock_quantity"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Auth request/response types

type RegisterRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RegisterResponse struct {
	Response  Response `json:"response"`
	User      User        `json:"user"`
	Token     string      `json:"token"`
	ExpiresAt int64       `json:"expires_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Response     Response `json:"response"`
	User         User     `json:"user"`
	Token        string   `json:"token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    int64    `json:"expires_at"`
}

type ProfileResponse struct {
	Response Response `json:"response"`
	User     User     `json:"user"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// Product request/response types

type ProductListParams struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Category string `json:"category,omitempty"`
}

type ProductSearchParams struct {
	Query    string  `json:"query"`
	Page     int     `json:"page"`
	Limit    int     `json:"limit"`
	MinPrice float64 `json:"min_price,omitempty"`
	MaxPrice float64 `json:"max_price,omitempty"`
}

type ProductListResponse struct {
	Response   Response   `json:"response"`
	Products   []Product  `json:"products"`
	Pagination Pagination `json:"pagination"`
}


type ProductResponse struct {
	Response Response `json:"response"`
	Product  Product  `json:"product"`
}


type CreateProductRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	Category      string  `json:"category"`
	ImageURL      string  `json:"image_url"`
	StockQuantity int32   `json:"stock_quantity"`
}

type UpdateProductRequest struct {
	Name          string  `json:"name,omitempty"`
	Description   string  `json:"description,omitempty"`
	Price         float64 `json:"price,omitempty"`
	Category      string  `json:"category,omitempty"`
	ImageURL      string  `json:"image_url,omitempty"`
	StockQuantity int32   `json:"stock_quantity,omitempty"`
	IsActive      *bool   `json:"is_active,omitempty"`
}

type CategoriesResponse struct {
	Response   Response   `json:"response"`
	Categories []Category `json:"categories"`
}

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
	IsActive    bool   `json:"is_active"`
}


type Pagination struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	TotalPages int32 `json:"total_pages"`
	TotalCount int64 `json:"total_count"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}


func (c *APIClient) post(ctx context.Context, path string, body interface{}, result interface{}) (interface{}, error) {
	return c.doRequest(ctx, "POST", path, body, result, "")
}

func (c *APIClient) postWithAuth(ctx context.Context, path string, body interface{}, result interface{}, token string) (interface{}, error) {
	return c.doRequest(ctx, "POST", path, body, result, token)
}

func (c *APIClient) get(ctx context.Context, path string, result interface{}) (interface{}, error) {
	return c.doRequest(ctx, "GET", path, nil, result, "")
}

func (c *APIClient) getWithAuth(ctx context.Context, path string, result interface{}, token string) (interface{}, error) {
	return c.doRequest(ctx, "GET", path, nil, result, token)
}

func (c *APIClient) putWithAuth(ctx context.Context, path string, body interface{}, result interface{}, token string) (interface{}, error) {
	return c.doRequest(ctx, "PUT", path, body, result, token)
}

func (c *APIClient) deleteWithAuth(ctx context.Context, path string, result interface{}, token string) (interface{}, error) {
	return c.doRequest(ctx, "DELETE", path, nil, result, token)
}

func (c *APIClient) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}, token string) (interface{}, error) {
	fullURL := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return result, nil
	}

	return nil, nil
}

// Auth API methods

func (c *APIClient) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	result := &RegisterResponse{}
	_, err := c.post(ctx, "/api/auth/register", req, result)
	return result, err
}

func (c *APIClient) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	result := &LoginResponse{}
	_, err := c.post(ctx, "/api/auth/login", req, result)
	return result, err
}

func (c *APIClient) GetProfile(ctx context.Context, token string) (*ProfileResponse, error) {
	result := &ProfileResponse{}
	_, err := c.getWithAuth(ctx, "/api/auth/profile", result, token)
	return result, err
}

func (c *APIClient) UpdateProfile(ctx context.Context, token string, req UpdateProfileRequest) (*ProfileResponse, error) {
	result := &ProfileResponse{}
	_, err := c.putWithAuth(ctx, "/api/auth/profile", req, result, token)
	return result, err
}

func (c *APIClient) ChangePassword(ctx context.Context, token string, req ChangePasswordRequest) (*APIResponse, error) {
	result := &APIResponse{}
	_, err := c.postWithAuth(ctx, "/api/auth/change-password", req, result, token)
	return result, err
}

// Product API methods

func (c *APIClient) ListProducts(ctx context.Context, params ProductListParams) (*ProductListResponse, error) {
	query := url.Values{}
	if params.Page > 0 {
		query.Set("page", strconv.Itoa(params.Page))
	}
	if params.Limit > 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Category != "" {
		query.Set("category", params.Category)
	}

	path := "/api/products"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	result := &ProductListResponse{}
	_, err := c.get(ctx, path, result)
	return result, err
}

func (c *APIClient) SearchProducts(ctx context.Context, params ProductSearchParams) (*ProductListResponse, error) {
	query := url.Values{}
	if params.Query != "" {
		query.Set("q", params.Query)
	}
	if params.Page > 0 {
		query.Set("page", strconv.Itoa(params.Page))
	}
	if params.Limit > 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.MinPrice > 0 {
		query.Set("min_price", fmt.Sprintf("%.2f", params.MinPrice))
	}
	if params.MaxPrice > 0 {
		query.Set("max_price", fmt.Sprintf("%.2f", params.MaxPrice))
	}

	path := "/api/products/search?" + query.Encode()

	result := &ProductListResponse{}
	_, err := c.get(ctx, path, result)
	return result, err
}

func (c *APIClient) GetProduct(ctx context.Context, id string) (*ProductResponse, error) {
	result := &ProductResponse{}
	_, err := c.get(ctx, "/api/products/"+id, result)
	return result, err
}

func (c *APIClient) CreateProduct(ctx context.Context, token string, req CreateProductRequest) (*ProductResponse, error) {
	result := &ProductResponse{}
	_, err := c.postWithAuth(ctx, "/api/products", req, result, token)
	return result, err
}

func (c *APIClient) UpdateProduct(ctx context.Context, token string, id string, req UpdateProductRequest) (*ProductResponse, error) {
	result := &ProductResponse{}
	_, err := c.putWithAuth(ctx, "/api/products/"+id, req, result, token)
	return result, err
}

func (c *APIClient) DeleteProduct(ctx context.Context, token string, id string) (*APIResponse, error) {
	result := &APIResponse{}
	_, err := c.deleteWithAuth(ctx, "/api/products/"+id, result, token)
	return result, err
}

func (c *APIClient) GetCategories(ctx context.Context) (*CategoriesResponse, error) {
	result := &CategoriesResponse{}
	_, err := c.get(ctx, "/api/products/categories", result)
	return result, err
}