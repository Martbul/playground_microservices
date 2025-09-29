package models

import (
	"time"
)

type Product struct {
	ID            string    `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	Price         float64   `json:"price" db:"price"`
	StockQuantity int32     `json:"stock_quantity" db:"stock_quantity"`
	Category      string    `json:"category" db:"category"`
	ImageURL      string    `json:"image_url" db:"image_url"`
	SKU           string    `json:"sku" db:"sku"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy     string    `json:"created_by" db:"created_by"`
}

type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	ParentID    *string   `json:"parent_id" db:"parent_id"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateProductRequest struct {
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description"`
	Price         float64 `json:"price" validate:"required,min=0"`
	StockQuantity int32   `json:"stock_quantity" validate:"min=0"`
	Category      string  `json:"category"`
	ImageURL      string  `json:"image_url"`
	SKU           string  `json:"sku"`
}

type UpdateProductRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price" validate:"min=0"`
	StockQuantity int32   `json:"stock_quantity" validate:"min=0"`
	Category      string  `json:"category"`
	ImageURL      string  `json:"image_url"`
	SKU           string  `json:"sku"`
	IsActive      *bool   `json:"is_active"`
}

type ProductFilter struct {
	Category  string
	MinPrice  float64
	MaxPrice  float64
	IsActive  *bool
	CreatedBy string
}

type SearchFilter struct {
	Query     string
	Category  string
	MinPrice  float64
	MaxPrice  float64
	IsActive  *bool
}

type PaginationRequest struct {
	Page      int32
	Limit     int32
	SortBy    string
	SortOrder string
}

type PaginationResponse struct {
	Page       int32
	Limit      int32
	TotalPages int32
	TotalCount int64
	HasNext    bool
	HasPrev    bool
}