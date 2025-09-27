package routes

import (
	"github.com/microservices-tutorial/services/client-service/handlers"
	"github.com/microservices-tutorial/services/client-service/middleware"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func SetupRoutes(
	router *mux.Router,
	authHandler *handlers.AuthHandler,
	productHandler *handlers.ProductHandler,
	pageHandler *handlers.PageHandler,
) {
	store := authHandler.(*handlers.AuthHandler).(*handlers.AuthHandler) // This needs to be fixed - we need access to store

	// Public routes
	router.HandleFunc("/", pageHandler.Home).Methods("GET")
	
	// Auth routes for guests only
	guestRouter := router.NewRoute().Subrouter()
	// guestRouter.Use(middleware.GuestOnlyMiddleware(store)) // Commented out for now
	guestRouter.HandleFunc("/login", authHandler.ShowLogin).Methods("GET")
	guestRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	guestRouter.HandleFunc("/register", authHandler.ShowRegister).Methods("GET")
	guestRouter.HandleFunc("/register", authHandler.Register).Methods("POST")

	// Protected routes (require authentication)
	protectedRouter := router.NewRoute().Subrouter()
	// protectedRouter.Use(middleware.AuthRequiredMiddleware(store)) // Commented out for now
	protectedRouter.HandleFunc("/dashboard", pageHandler.Dashboard).Methods("GET")
	protectedRouter.HandleFunc("/logout", authHandler.Logout).Methods("GET", "POST")
	protectedRouter.HandleFunc("/profile", authHandler.ShowProfile).Methods("GET")
	protectedRouter.HandleFunc("/profile", authHandler.UpdateProfile).Methods("POST")
	protectedRouter.HandleFunc("/profile/change-password", authHandler.ChangePassword).Methods("POST")

	// Product routes (public + protected)
	router.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	router.HandleFunc("/products/{id}", productHandler.ShowProduct).Methods("GET")
	
	// Protected product routes
	protectedRouter.HandleFunc("/products/create", productHandler.ShowCreateProduct).Methods("GET")
	protectedRouter.HandleFunc("/products/create", productHandler.CreateProduct).Methods("POST")
	protectedRouter.HandleFunc("/products/{id}/edit", productHandler.ShowEditProduct).Methods("GET")
	protectedRouter.HandleFunc("/products/{id}/edit", productHandler.UpdateProduct).Methods("POST")
	protectedRouter.HandleFunc("/products/{id}/delete", productHandler.DeleteProduct).Methods("POST")

	// 404 handler
	router.NotFoundHandler = router.NewRoute().HandlerFunc(pageHandler.NotFound).GetHandler()
}

// SetupRoutesWithStore - proper setup with store access
func SetupRoutesWithStore(
	router *mux.Router,
	authHandler *handlers.AuthHandler,
	productHandler *handlers.ProductHandler,
	pageHandler *handlers.PageHandler,
	store *sessions.CookieStore,
) {
	// Public routes
	router.HandleFunc("/", pageHandler.Home).Methods("GET")
	
	// Auth routes for guests only
	guestRouter := router.NewRoute().Subrouter()
	guestRouter.Use(middleware.GuestOnlyMiddleware(store))
	guestRouter.HandleFunc("/login", authHandler.ShowLogin).Methods("GET")
	guestRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	guestRouter.HandleFunc("/register", authHandler.ShowRegister).Methods("GET")
	guestRouter.HandleFunc("/register", authHandler.Register).Methods("POST")

	// Protected routes (require authentication)
	protectedRouter := router.NewRoute().Subrouter()
	protectedRouter.Use(middleware.AuthRequiredMiddleware(store))
	protectedRouter.HandleFunc("/dashboard", pageHandler.Dashboard).Methods("GET")
	protectedRouter.HandleFunc("/logout", authHandler.Logout).Methods("GET", "POST")
	protectedRouter.HandleFunc("/profile", authHandler.ShowProfile).Methods("GET")
	protectedRouter.HandleFunc("/profile", authHandler.UpdateProfile).Methods("POST")
	protectedRouter.HandleFunc("/profile/change-password", authHandler.ChangePassword).Methods("POST")

	// Product routes (public + protected)
	router.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	router.HandleFunc("/products/{id}", productHandler.ShowProduct).Methods("GET")
	
	// Protected product routes
	protectedRouter.HandleFunc("/products/create", productHandler.ShowCreateProduct).Methods("GET")
	protectedRouter.HandleFunc("/products/create", productHandler.CreateProduct).Methods("POST")
	protectedRouter.HandleFunc("/products/{id}/edit", productHandler.ShowEditProduct).Methods("GET")
	protectedRouter.HandleFunc("/products/{id}/edit", productHandler.UpdateProduct).Methods("POST")
	protectedRouter.HandleFunc("/products/{id}/delete", productHandler.DeleteProduct).Methods("POST")

	// 404 handler
	router.NotFoundHandler = router.NewRoute().HandlerFunc(pageHandler.NotFound).GetHandler()
}