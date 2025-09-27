package handlers

import (
	"html/template"
	"net/http"

	"github.com/microservices-tutorial/services/client-service/clients"
	"github.com/microservices-tutorial/services/client-service/utils"
	"github.com/gorilla/sessions"
)

type PageHandler struct {
	apiClient *clients.APIClient
	store     *sessions.CookieStore
}

func NewPageHandler(apiClient *clients.APIClient, store *sessions.CookieStore) *PageHandler {
	return &PageHandler{
		apiClient: apiClient,
		store:     store,
	}
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	// Get featured products (first 6 products)
	params := clients.ProductListParams{
		Page:  1,
		Limit: 6,
	}

	resp, err := h.apiClient.ListProducts(r.Context(), params)
	if err != nil {
		// Continue without products if there's an error
		resp = &clients.ProductListResponse{
			Products: []clients.Product{},
		}
	}

	listResp := resp.(*clients.ProductListResponse)

	// Check if user is authenticated
	user, _, _ := utils.GetUserFromSession(r, h.store)

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/pages/home.html",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":            "Welcome to MicroStore",
		"FeaturedProducts": listResp.Products,
		"User":             user,
	}

	tmpl.Execute(w, data)
}

func (h *PageHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user, token, err := utils.GetUserFromSession(r, h.store)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get user's recent products (if they're an admin or have created products)
	// For now, just show recent products
	params := clients.ProductListParams{
		Page:  1,
		Limit: 10,
	}

	resp, err := h.apiClient.ListProducts(r.Context(), params)
	if err != nil {
		resp = &clients.ProductListResponse{
			Products: []clients.Product{},
		}
	}

	listResp := resp.(*clients.ProductListResponse)

	// Get fresh user profile
	profileResp, err := h.apiClient.GetProfile(r.Context(), token)
	if err != nil {
		http.Redirect(w, r, "/login?error=token_expired", http.StatusFound)
		return
	}

	profile := profileResp.(*clients.ProfileResponse)
	if !profile.Response.Success {
		http.Redirect(w, r, "/login?error="+profile.Response.Message, http.StatusFound)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/pages/dashboard.html",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":          "Dashboard",
		"User":           profile.User,
		"RecentProducts": listResp.Products,
		"Success":        r.URL.Query().Get("success"),
		"Error":          r.URL.Query().Get("error"),
	}

	tmpl.Execute(w, data)
}

func (h *PageHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	user, _, _ := utils.GetUserFromSession(r, h.store)

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/pages/error.html",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":       "Page Not Found",
		"ErrorCode":   "404",
		"ErrorTitle":  "Page Not Found",
		"ErrorMessage": "The page you're looking for doesn't exist.",
		"User":        user,
	}

	w.WriteHeader(http.StatusNotFound)
	tmpl.Execute(w, data)
}