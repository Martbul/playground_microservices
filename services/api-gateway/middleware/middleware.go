package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/martbul/playground_microservices/services/api-gateway/clients"
	pb "github.com/martbul/playground_microservices/services/api-gateway/genproto/auth"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		// Log the request
		duration := time.Since(start)
		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			duration,
		)
	})
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(authClient *clients.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validate token with auth service
			validateReq := &pb.ValidateTokenRequest{
				Token: token,
			}

			resp, err := authClient.ValidateToken(r.Context(), validateReq)
			if err != nil {
				log.Printf("Token validation error: %v", err)
				http.Error(w, "Token validation failed", http.StatusUnauthorized)
				return
			}

			if !resp.Response.Success || !resp.Valid {
				http.Error(w, resp.Response.Message, http.StatusUnauthorized)
				return
			}

			// Add user info to request context
			ctx := context.WithValue(r.Context(), "user", resp.User)
			ctx = context.WithValue(ctx, "token", token)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuthMiddleware validates JWT tokens but allows requests without them
func OptionalAuthMiddleware(authClient *clients.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// Check Bearer token format
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token := parts[1]

					// Validate token with auth service
					validateReq := &pb.ValidateTokenRequest{
						Token: token,
					}

					resp, err := authClient.ValidateToken(r.Context(), validateReq)
					if err == nil && resp.Response.Success && resp.Valid {
						// Add user info to request context
						ctx := context.WithValue(r.Context(), "user", resp.User)
						ctx = context.WithValue(ctx, "token", token)
						r = r.WithContext(ctx)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}