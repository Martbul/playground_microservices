package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	pb "github.com/playground_microservices/proto/product"
	"github.com/playground_microservices/services/product-service/config"
	"github.com/playground_microservices/services/product-service/handlers"
	"github.com/playground_microservices/services/product-service/repository"
	"github.com/playground_microservices/services/product-service/service"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repository
	productRepo := repository.NewProductRepository(db)

	// Initialize service
	productService := service.NewProductService(productRepo)

	// Initialize handler
	productHandler := handlers.NewProductHandler(productService)

	// Create gRPC server
	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, productHandler)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	log.Printf("Product service starting on port %s", cfg.Port)
	if err := server.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}

func runMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		parent_id UUID REFERENCES categories(id),
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10,2) NOT NULL,
		stock_quantity INTEGER DEFAULT 0,
		category VARCHAR(255),
		image_url VARCHAR(500),
		sku VARCHAR(100) UNIQUE,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(255)
	);

	INSERT INTO categories (name, description) VALUES 
		('Electronics', 'Electronic devices and accessories'),
		('Clothing', 'Apparel and fashion items'),
		('Books', 'Books and educational materials'),
		('Sports', 'Sports and fitness equipment'),
		('Home', 'Home and garden products')
	ON CONFLICT DO NOTHING;

	CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
	CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
	CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
	CREATE INDEX IF NOT EXISTS idx_products_active ON products(is_active);
	CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
	CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories(parent_id);
	`

	_, err := db.Exec(query)
	return err
}
