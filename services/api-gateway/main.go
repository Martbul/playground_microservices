package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/martbul/playground_microservices/services/api-gateway/clients"
	"github.com/martbul/playground_microservices/services/api-gateway/config"
	"github.com/martbul/playground_microservices/services/api-gateway/handlers"
	"github.com/martbul/playground_microservices/services/api-gateway/middleware"
	"github.com/martbul/playground_microservices/services/api-gateway/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()  //! use viper for configuratin instead

	// Initialize gRPC clients
	authClient, err := clients.NewAuthGrpcClient(cfg.AuthService)
	if err != nil {
		log.Fatal("Failed to connect to auth service:", err)
	}
	defer authClient.Close()

	productClient, err := clients.NewProductClient(cfg.ProductService)
	if err != nil {
		log.Fatal("Failed to connect to product service:", err)
	}
	defer productClient.Close()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authClient)
	productHandler := handlers.NewProductHandler(productClient)

	// Create router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)

	// Setup routes
	routes.SetupAuthRoutes(router, authHandler)
	routes.SetupProductRoutes(router, productHandler, authClient)

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "api-gateway"}`))
	}).Methods("GET")

	log.Printf("API Gateway starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}