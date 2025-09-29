package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/martbul/playground_microservices/services/client-service/clients"
)

// GetUserFromSession retrieves user info from session
func GetUserFromSession(r *http.Request, store *sessions.CookieStore) (*clients.User, string, error) {
	session, err := store.Get(r, "session")
	if err != nil {
		return nil, "", err
	}

	user, ok := session.Values["user"].(clients.User)
	if !ok {
		return nil, "", nil
	}

	token, ok := session.Values["token"].(string)
	if !ok {
		return nil, "", nil
	}

	return &user, token, nil
}

// SaveUserToSession saves user info to session
func SaveUserToSession(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore, user clients.User, token string) error {
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["user"] = user
	session.Values["token"] = token
	session.Values["authenticated"] = true

	return session.Save(r, w)
}

// ClearSession clears the session
func ClearSession(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) error {
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}

	// Clear all session values
	for key := range session.Values {
		delete(session.Values, key)
	}

	session.Options.MaxAge = -1 // Delete the session
	return session.Save(r, w)
}

// IsAuthenticated checks if user is authenticated
func IsAuthenticated(r *http.Request, store *sessions.CookieStore) bool {
	session, err := store.Get(r, "session")
	if err != nil {
		return false
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	return ok && authenticated
}

// GetFlashMessage gets and clears a flash message
func GetFlashMessage(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore, key string) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	flashes := session.Flashes(key)
	if len(flashes) > 0 {
		session.Save(r, w)
		return flashes[0].(string)
	}

	return ""
}

// SetFlashMessage sets a flash message
func SetFlashMessage(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore, key, message string) {
	session, err := store.Get(r, "session")
	if err != nil {
		return
	}

	session.AddFlash(message, key)
	session.Save(r, w)
}
