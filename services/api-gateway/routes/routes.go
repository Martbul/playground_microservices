package routes

import (
	"github.com/gorilla/mux"
	"github.com/martbul/playground_microservices/services/api-gateway/clients"
	"github.com/martbul/playground_microservices/services/api-gateway/handlers"
	"github.com/martbul/playground_microservices/services/api-gateway/middleware"
)

func SetupAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler) {
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	
	// Public routes (no authentication required)
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRouter.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
}

func SetupProtectedAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler, authClient *clients.AuthClient) {
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.Use(middleware.AuthMiddleware(authClient))
	
	// Protected routes (authentication required)
	authRouter.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
	authRouter.HandleFunc("/profile", authHandler.UpdateProfile).Methods("PUT")
	authRouter.HandleFunc("/change-password", authHandler.ChangePassword).Methods("POST")
}

func SetupProductRoutes(router *mux.Router, productHandler *handlers.ProductHandler, authClient *clients.AuthClient) {
	productRouter := router.PathPrefix("/api/products").Subrouter()
	
	// Public routes (no authentication required)
	publicRouter := productRouter.PathPrefix("").Subrouter()
	publicRouter.Use(middleware.OptionalAuthMiddleware(authClient))
	publicRouter.HandleFunc("", productHandler.ListProducts).Methods("GET")
	publicRouter.HandleFunc("/search", productHandler.SearchProducts).Methods("GET")
	publicRouter.HandleFunc("/categories", productHandler.GetCategories).Methods("GET")
	publicRouter.HandleFunc("/{id}", productHandler.GetProduct).Methods("GET")
	
	// Protected routes (authentication required)
	protectedRouter := productRouter.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware(authClient))
	protectedRouter.HandleFunc("", productHandler.CreateProduct).Methods("POST")
	protectedRouter.HandleFunc("/{id}", productHandler.UpdateProduct).Methods("PUT")
	protectedRouter.HandleFunc("/{id}", productHandler.DeleteProduct).Methods("DELETE")
	
	// Setup protected auth routes here too
	SetupProtectedAuthRoutes(router, &handlers.AuthHandler{}, authClient)
}