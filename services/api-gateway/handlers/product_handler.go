package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	commonPb "github.com/playground_microservices/proto/common"
	pb "github.com/playground_microservices/proto/product"
	"github.com/playground_microservices/services/api-gateway/clients"
)

type ProductHandler struct {
	productClient *clients.ProductClient
}

func NewProductHandler(productClient *clients.ProductClient) *ProductHandler {
	return &ProductHandler{
		productClient: productClient,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(string)
	if !ok {
		http.Error(w, "Token not found in context", http.StatusInternalServerError)
		return
	}

	var createReq struct {
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		Price         float64 `json:"price"`
		StockQuantity int32   `json:"stock_quantity"`
		Category      string  `json:"category"`
		ImageURL      string  `json:"image_url"`
		SKU           string  `json:"sku"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.CreateProductRequest{
		Token:         token,
		Name:          createReq.Name,
		Description:   createReq.Description,
		Price:         createReq.Price,
		StockQuantity: createReq.StockQuantity,
		Category:      createReq.Category,
		ImageUrl:      createReq.ImageURL,
		Sku:           createReq.SKU,
	}

	resp, err := h.productClient.CreateProduct(r.Context(), req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Response.Success {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	req := &pb.GetProductRequest{
		Id: productID,
	}

	resp, err := h.productClient.GetProduct(r.Context(), req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	token, ok := r.Context().Value("token").(string)
	if !ok {
		http.Error(w, "Token not found in context", http.StatusInternalServerError)
		return
	}

	var updateReq struct {
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		Price         float64 `json:"price"`
		StockQuantity int32   `json:"stock_quantity"`
		Category      string  `json:"category"`
		ImageURL      string  `json:"image_url"`
		SKU           string  `json:"sku"`
		IsActive      bool    `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.UpdateProductRequest{
		Token:         token,
		Id:            productID,
		Name:          updateReq.Name,
		Description:   updateReq.Description,
		Price:         updateReq.Price,
		StockQuantity: updateReq.StockQuantity,
		Category:      updateReq.Category,
		ImageUrl:      updateReq.ImageURL,
		Sku:           updateReq.SKU,
		IsActive:      updateReq.IsActive,
	}

	resp, err := h.productClient.UpdateProduct(r.Context(), req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	token, ok := r.Context().Value("token").(string)
	if !ok {
		http.Error(w, "Token not found in context", http.StatusInternalServerError)
		return
	}

	req := &pb.DeleteProductRequest{
		Token: token,
		Id:    productID,
	}

	_, err := h.productClient.DeleteProduct(r.Context(), req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Product deleted successfully",
	})
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	category := query.Get("category")
	sortBy := query.Get("sort_by")
	sortOrder := query.Get("sort_order")
	activeOnly := query.Get("active_only") == "true"

	req := &pb.ListProductsRequest{
		Pagination: &commonPb.PaginationRequest{
			Page:      int32(page),
			Limit:     int32(limit),
			SortBy:    sortBy,
			SortOrder: sortOrder,
		},
		Category:   category,
		ActiveOnly: activeOnly,
	}

	resp, err := h.productClient.ListProducts(r.Context(), req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	searchQuery := query.Get("q")
	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	category := query.Get("category")
	minPrice, _ := strconv.ParseFloat(query.Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(query.Get("max_price"), 64)

	req := &pb.SearchProductsRequest{
		Query: searchQuery,
		Pagination: &commonPb.PaginationRequest{
			Page:  int32(page),
			Limit: int32(limit),
		},
		Category: category,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}

	resp, err := h.productClient.SearchProducts(r.Context(), req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *ProductHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	req := &pb.GetCategoriesRequest{}

	resp, err := h.productClient.GetCategories(r.Context(), req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
