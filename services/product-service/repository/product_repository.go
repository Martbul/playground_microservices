package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/martbul/playground_microservices/services/product-service/models"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id string) (*models.Product, error)
	GetBySKU(sku string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id string) error
	List(filter *models.ProductFilter, pagination *models.PaginationRequest) ([]*models.Product, *models.PaginationResponse, error)
	Search(filter *models.SearchFilter, pagination *models.PaginationRequest) ([]*models.Product, *models.PaginationResponse, error)
	GetCategories() ([]*models.Category, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock_quantity, category, image_url, sku, is_active, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.Category,
		product.ImageURL,
		product.SKU,
		product.IsActive,
		product.CreatedBy,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *productRepository) GetByID(id string) (*models.Product, error) {
	product := &models.Product{}
	query := `
		SELECT id, name, description, price, stock_quantity, category, image_url, sku, is_active, created_at, updated_at, created_by
		FROM products WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.Category,
		&product.ImageURL,
		&product.SKU,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	return product, nil
}

func (r *productRepository) GetBySKU(sku string) (*models.Product, error) {
	product := &models.Product{}
	query := `
		SELECT id, name, description, price, stock_quantity, category, image_url, sku, is_active, created_at, updated_at, created_by
		FROM products WHERE sku = $1
	`

	err := r.db.QueryRow(query, sku).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.Category,
		&product.ImageURL,
		&product.SKU,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	return product, nil
}

func (r *productRepository) Update(product *models.Product) error {
	query := `
		UPDATE products 
		SET name = $1, description = $2, price = $3, stock_quantity = $4, category = $5, 
		    image_url = $6, sku = $7, is_active = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.Category,
		product.ImageURL,
		product.SKU,
		product.IsActive,
		product.ID,
	).Scan(&product.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *productRepository) Delete(id string) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *productRepository) List(filter *models.ProductFilter, pagination *models.PaginationRequest) ([]*models.Product, *models.PaginationResponse, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	baseQuery := `
		SELECT id, name, description, price, stock_quantity, category, image_url, sku, is_active, created_at, updated_at, created_by
		FROM products
	`

	if filter != nil {
		if filter.Category != "" {
			conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
			args = append(args, filter.Category)
			argIndex++
		}

		if filter.MinPrice > 0 {
			conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
			args = append(args, filter.MinPrice)
			argIndex++
		}

		if filter.MaxPrice > 0 {
			conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
			args = append(args, filter.MaxPrice)
			argIndex++
		}

		if filter.IsActive != nil {
			conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
			args = append(args, *filter.IsActive)
			argIndex++
		}

		if filter.CreatedBy != "" {
			conditions = append(conditions, fmt.Sprintf("created_by = $%d", argIndex))
			args = append(args, filter.CreatedBy)
			argIndex++
		}
	}

	// Build WHERE clause
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := "SELECT COUNT(*) FROM products"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int64
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Add sorting
	if pagination.SortBy != "" {
		sortOrder := "ASC"
		if pagination.SortOrder == "desc" {
			sortOrder = "DESC"
		}
		baseQuery += fmt.Sprintf(" ORDER BY %s %s", pagination.SortBy, sortOrder)
	} else {
		baseQuery += " ORDER BY created_at DESC"
	}

	// Add pagination
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pagination.Limit, (pagination.Page-1)*pagination.Limit)

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.Category,
			&product.ImageURL,
			&product.SKU,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.CreatedBy,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	// Calculate pagination info
	totalPages := int32((totalCount + int64(pagination.Limit) - 1) / int64(pagination.Limit))
	paginationResponse := &models.PaginationResponse{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: totalPages,
		TotalCount: totalCount,
		HasNext:    pagination.Page < totalPages,
		HasPrev:    pagination.Page > 1,
	}

	return products, paginationResponse, nil
}

func (r *productRepository) Search(filter *models.SearchFilter, pagination *models.PaginationRequest) ([]*models.Product, *models.PaginationResponse, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	baseQuery := `
		SELECT id, name, description, price, stock_quantity, category, image_url, sku, is_active, created_at, updated_at, created_by
		FROM products
	`

	// Add search query condition
	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d OR sku ILIKE $%d)", argIndex, argIndex, argIndex))
		searchTerm := "%" + filter.Query + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	if filter.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, filter.Category)
		argIndex++
	}

	if filter.MinPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, filter.MinPrice)
		argIndex++
	}

	if filter.MaxPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, filter.MaxPrice)
		argIndex++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	// Build WHERE clause
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := "SELECT COUNT(*) FROM products"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int64
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Add sorting (relevance for search)
	baseQuery += " ORDER BY created_at DESC"

	// Add pagination
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pagination.Limit, (pagination.Page-1)*pagination.Limit)

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.Category,
			&product.ImageURL,
			&product.SKU,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.CreatedBy,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	// Calculate pagination info
	totalPages := int32((totalCount + int64(pagination.Limit) - 1) / int64(pagination.Limit))
	paginationResponse := &models.PaginationResponse{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: totalPages,
		TotalCount: totalCount,
		HasNext:    pagination.Page < totalPages,
		HasPrev:    pagination.Page > 1,
	}

	return products, paginationResponse, nil
}

func (r *productRepository) GetCategories() ([]*models.Category, error) {
	query := `
		SELECT id, name, description, parent_id, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = true
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.ParentID,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}	