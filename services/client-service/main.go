package main

import (
	"log"
	"net/http"

	"github.com/martbul/playground_microservices/services/client-service/clients"
	"github.com/martbul/playground_microservices/services/client-service/config"
	"github.com/martbul/playground_microservices/services/client-service/handlers"
	"github.com/martbul/playground_microservices/services/client-service/middleware"
	"github.com/martbul/playground_microservices/services/client-service/routes"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize session store
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24, // 24 hours
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	// Initialize API Gateway client
	apiClient := clients.NewAPIClient(cfg.APIGatewayURL)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(apiClient, store)
	productHandler := handlers.NewProductHandler(apiClient, store)
	pageHandler := handlers.NewPageHandler(apiClient, store)

	// Create router
	router := mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Apply middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.SessionMiddleware(store))

	// Setup routes
	routes.SetupRoutesWithStore(router, authHandler, productHandler, pageHandler, store)

	log.Printf("Client service starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
