package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/martbul/playground_microservices/services/auth-service/config"
	"github.com/martbul/playground_microservices/services/auth-service/handlers"
	"github.com/martbul/playground_microservices/services/auth-service/repository"
	"github.com/martbul/playground_microservices/services/auth-service/service"
	pb "github.com/martbul/playground_microservices/proto/auth"

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
	userRepo := repository.NewUserRepository(db)

	// Initialize service
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	// Initialize handler
	authHandler := handlers.NewAuthHandler(authService)

	// Create gRPC server
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authHandler)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	log.Printf("Auth service starting on port %s", cfg.Port)
	if err := server.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}

func runMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		username VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		role VARCHAR(50) DEFAULT 'user',
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token_hash VARCHAR(255) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
	`

	_, err := db.Exec(query)
	return err
}