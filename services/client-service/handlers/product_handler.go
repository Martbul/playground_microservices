package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/martbul/playground_microservices/services/client-service/clients"
	"github.com/martbul/playground_microservices/services/client-service/utils"
)

type ProductHandler struct {
	apiClient *clients.APIClient
	store     *sessions.CookieStore
}

func NewProductHandler(apiClient *clients.APIClient, store *sessions.CookieStore) *ProductHandler {
	return &ProductHandler{
		apiClient: apiClient,
		store:     store,
	}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	category := r.URL.Query().Get("category")
	searchQuery := r.URL.Query().Get("q")

	page := 1
	limit := 12

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	var resp *clients.ProductListResponse
	var err error

	// Check if user is authenticated
	user, _, _ := utils.GetUserFromSession(r, h.store)

	// Search or list products
	if searchQuery != "" {
		params := clients.ProductSearchParams{
			Query: searchQuery,
			Page:  page,
			Limit: limit,
		}
		resp, err = h.apiClient.SearchProducts(r.Context(), params)
	} else {
		params := clients.ProductListParams{
			Page:     page,
			Limit:    limit,
			Category: category,
		}
		resp, err = h.apiClient.ListProducts(r.Context(), params)
	}

	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	// Get categories for filter
	categoriesResp, _ := h.apiClient.GetCategories(r.Context())
	var categories []clients.Category
	if categoriesResp != nil {
		categories = categoriesResp.Categories
	}

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/products/list.html",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// FIXED: Use correct field names from Pagination struct
	data := map[string]interface{}{
		"Title":            "Products",
		"Products":         resp.Products,
		"Page":             resp.Pagination.Page,
		"Limit":            resp.Pagination.Limit,
		"TotalPages":       resp.Pagination.TotalPages,
		"TotalCount":       resp.Pagination.TotalCount,
		"HasNext":          resp.Pagination.HasNext,
		"HasPrev":          resp.Pagination.HasPrev,
		"Categories":       categories,
		"SelectedCategory": category,
		"SearchQuery":      searchQuery,
		"User":             user,
	}

	tmpl.Execute(w, data)
}

func (h *ProductHandler) ShowProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Redirect(w, r, "/products", http.StatusFound)
		return
	}

	resp, err := h.apiClient.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if !resp.Response.Success {
		http.Error(w, resp.Response.Message, http.StatusNotFound)
		return
	}

	// Check if user is authenticated
	user, _, _ := utils.GetUserFromSession(r, h.store)

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/products/detail.html",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":   resp.Product.Name,
		"Product": resp.Product,
		"User":    user,
	}

	tmpl.Execute(w, data)
}

func (h *ProductHandler) ShowCreateProduct(w http.ResponseWriter, r *http.Request) {
	user, _, err := utils.GetUserFromSession(r, h.store)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get categories
	categoriesResp, _ := h.apiClient.GetCategories(r.Context())
	var categories []clients.Category
	if categoriesResp != nil {
		categories = categoriesResp.Categories
	}

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/products/create.html",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":      "Create Product",
		"Categories": categories,
		"User":       user,
		"Error":      r.URL.Query().Get("error"),
	}

	tmpl.Execute(w, data)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/products/create", http.StatusFound)
		return
	}

	user, token, err := utils.GetUserFromSession(r, h.store)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Parse form data
	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	category := r.FormValue("category")
	imageURL := r.FormValue("image_url")
	stockStr := r.FormValue("stock_quantity")

	// Validate and convert
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Redirect(w, r, "/products/create?error=invalid_price", http.StatusFound)
		return
	}

	stock, err := strconv.ParseInt(stockStr, 10, 32)
	if err != nil {
		http.Redirect(w, r, "/products/create?error=invalid_stock", http.StatusFound)
		return
	}

	req := clients.CreateProductRequest{
		Name:          name,
		Description:   description,
		Price:         price,
		Category:      category,
		ImageURL:      imageURL,
		StockQuantity: int32(stock),
	}

	resp, err := h.apiClient.CreateProduct(r.Context(), token, req)
	if err != nil {
		http.Redirect(w, r, "/products/create?error=server_error", http.StatusFound)
		return
	}

	if !resp.Response.Success {
		http.Redirect(w, r, "/products/create?error="+resp.Response.Message, http.StatusFound)
		return
	}

	http.Redirect(w, r, "/products/"+resp.Product.ID, http.StatusFound)
}

func (h *ProductHandler) ShowEditProduct(w http.ResponseWriter, r *http.Request) {
	user, token, err := utils.GetUserFromSession(r, h.store)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Redirect(w, r, "/products", http.StatusFound)
		return
	}

	// Get product details
	resp, err := h.apiClient.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if !resp.Response.Success {
		http.Error(w, resp.Response.Message, http.StatusNotFound)
		return
	}

	// Get categories
	categoriesResp, _ := h.apiClient.GetCategories(r.Context())
	var categories []clients.Category
	if categoriesResp != nil {
		categories = categoriesResp.Categories
	}

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/products/edit.html",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":      "Edit Product",
		"Product":    resp.Product,
		"Categories": categories,
		"User":       user,
		"Token":      token,
		"Error":      r.URL.Query().Get("error"),
		"Success":    r.URL.Query().Get("success"),
	}

	tmpl.Execute(w, data)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/products", http.StatusFound)
		return
	}

	user, token, err := utils.GetUserFromSession(r, h.store)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Redirect(w, r, "/products", http.StatusFound)
		return
	}

	// Parse form data
	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	category := r.FormValue("category")
	imageURL := r.FormValue("image_url")
	stockStr := r.FormValue("stock_quantity")
	isActiveStr := r.FormValue("is_active")

	// Build update request
	req := clients.UpdateProductRequest{}

	if name != "" {
		req.Name = name
	}
	if description != "" {
		req.Description = description
	}
	if priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			http.Redirect(w, r, "/products/"+id+"/edit?error=invalid_price", http.StatusFound)
			return
		}
		req.Price = price
	}
	if category != "" {
		req.Category = category
	}
	if imageURL != "" {
		req.ImageURL = imageURL
	}
	if stockStr != "" {
		stock, err := strconv.ParseInt(stockStr, 10, 32)
		if err != nil {
			http.Redirect(w, r, "/products/"+id+"/edit?error=invalid_stock", http.StatusFound)
			return
		}
		req.StockQuantity = int32(stock)
	}
	if isActiveStr != "" {
		isActive := isActiveStr == "true" || isActiveStr == "on"
		req.IsActive = &isActive
	}

	resp, err := h.apiClient.UpdateProduct(r.Context(), token, id, req)
	if err != nil {
		http.Redirect(w, r, "/products/"+id+"/edit?error=server_error", http.StatusFound)
		return
	}

	if !resp.Response.Success {
		http.Redirect(w, r, "/products/"+id+"/edit?error="+resp.Response.Message, http.StatusFound)
		return
	}

	http.Redirect(w, r, "/products/"+id+"/edit?success=product_updated", http.StatusFound)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/products", http.StatusFound)
		return
	}

	user, token, err := utils.GetUserFromSession(r, h.store)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Redirect(w, r, "/products", http.StatusFound)
		return
	}

	resp, err := h.apiClient.DeleteProduct(r.Context(), token, id)
	if err != nil {
		http.Redirect(w, r, "/products/"+id+"/edit?error=server_error", http.StatusFound)
		return
	}

	if !resp.Response.Success {
		http.Redirect(w, r, "/products/"+id+"/edit?error="+resp.Response.Message, http.StatusFound)
		return
	}

	http.Redirect(w, r, "/products?success=product_deleted", http.StatusFound)
}