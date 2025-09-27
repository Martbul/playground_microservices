package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
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

// SessionMiddleware makes session store available in context
func SessionMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Session is handled in individual handlers as needed
			next.ServeHTTP(w, r)
		})
	}
}

// AuthRequiredMiddleware redirects to login if not authenticated
func AuthRequiredMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "session")
			if err != nil {
				http.Redirect(w, r, "/login?error=session_error", http.StatusFound)
				return
			}

			token, ok := session.Values["token"].(string)
			if !ok || token == "" {
				http.Redirect(w, r, "/login?error=not_authenticated", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GuestOnlyMiddleware redirects to dashboard if already authenticated
func GuestOnlyMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "session")
			if err == nil {
				token, ok := session.Values["token"].(string)
				if ok && token != "" {
					http.Redirect(w, r, "/dashboard", http.StatusFound)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}