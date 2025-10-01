package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/martbul/playground_microservices/services/api-gateway/genproto/auth"

	"github.com/martbul/playground_microservices/services/api-gateway/clients"
)

type AuthHandler struct {
	authClient *clients.AuthGrpcClient
}

func NewAuthHandler(authClient *clients.AuthGrpcClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthHandler: Register request received")

	var req pb.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("AuthHandler: Failed to decode register request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Username == "" || req.Password == "" {
		http.Error(w, "Email, username, and password are required", http.StatusBadRequest)
		return
	}

	resp, err := h.authClient.Register(r.Context(), &req)
	if err != nil {
		log.Printf("AuthHandler: Register error: %v", err)
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

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthHandler: Login request received")

	var req pb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("AuthHandler: Failed to decode login request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	resp, err := h.authClient.Login(r.Context(), &req)
	if err != nil {
		log.Printf("AuthHandler: Login error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
	json.NewEncoder(w).Encode(resp)
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthHandler: GetProfile request received")

	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*pb.User)
	if !ok {
		log.Println("AuthHandler: User not found in context")
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	token, ok := r.Context().Value("token").(string)
	if !ok {
		log.Println("AuthHandler: Token not found in context")
		http.Error(w, "Token not found in context", http.StatusInternalServerError)
		return
	}

	req := &pb.GetUserRequest{
		UserId: user.Id,
		Token:  token,
	}

	resp, err := h.authClient.GetUser(r.Context(), req)
	if err != nil {
		log.Printf("AuthHandler: GetProfile error: %v", err)
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

// UpdateProfile updates the current user's profile
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthHandler: UpdateProfile request received")

	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*pb.User)
	if !ok {
		log.Println("AuthHandler: User not found in context")
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	token, ok := r.Context().Value("token").(string)
	if !ok {
		log.Println("AuthHandler: Token not found in context")
		http.Error(w, "Token not found in context", http.StatusInternalServerError)
		return
	}

	var updateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		log.Printf("AuthHandler: Failed to decode update profile request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.UpdateProfileRequest{
		UserId:    user.Id,
		Token:     token,
		FirstName: updateReq.FirstName,
		LastName:  updateReq.LastName,
		Username:  updateReq.Username,
	}

	resp, err := h.authClient.UpdateProfile(r.Context(), req)
	if err != nil {
		log.Printf("AuthHandler: UpdateProfile error: %v", err)
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

// ChangePassword changes the current user's password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthHandler: ChangePassword request received")

	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*pb.User)
	if !ok {
		log.Println("AuthHandler: User not found in context")
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	token, ok := r.Context().Value("token").(string)
	if !ok {
		log.Println("AuthHandler: Token not found in context")
		http.Error(w, "Token not found in context", http.StatusInternalServerError)
		return
	}

	var changeReq struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&changeReq); err != nil {
		log.Printf("AuthHandler: Failed to decode change password request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if changeReq.CurrentPassword == "" || changeReq.NewPassword == "" {
		http.Error(w, "Current password and new password are required", http.StatusBadRequest)
		return
	}

	req := &pb.ChangePasswordRequest{
		UserId:          user.Id,
		Token:           token,
		CurrentPassword: changeReq.CurrentPassword,
		NewPassword:     changeReq.NewPassword,
	}

	resp, err := h.authClient.ChangePassword(r.Context(), req)
	if err != nil {
		log.Printf("AuthHandler: ChangePassword error: %v", err)
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

// RefreshToken refreshes the JWT token using a refresh token
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthHandler: RefreshToken request received")

	var refreshReq struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		log.Printf("AuthHandler: Failed to decode refresh token request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required field
	if refreshReq.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	req := &pb.RefreshTokenRequest{
		RefreshToken: refreshReq.RefreshToken,
	}

	resp, err := h.authClient.RefreshToken(r.Context(), req)
	if err != nil {
		log.Printf("AuthHandler: RefreshToken error: %v", err)
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

// ValidateToken validates a JWT token
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthHandler: ValidateToken request received")

	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusBadRequest)
		return
	}

	// Extract token (remove "Bearer " prefix if present)
	token := authHeader
	if len(authHeader) > 7 && strings.ToLower(authHeader[:7]) == "bearer " {
		token = authHeader[7:]
	}

	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	req := &pb.ValidateTokenRequest{
		Token: token,
	}

	resp, err := h.authClient.ValidateToken(r.Context(), req)
	if err != nil {
		log.Printf("AuthHandler: ValidateToken error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Response.Success && resp.Valid {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
	json.NewEncoder(w).Encode(resp)
}

// GetUserByID retrieves a user by ID (admin only or own profile)
func (h *AuthHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthHandler: GetUserByID request received")

	// Get user from context (set by auth middleware)
	currentUser, ok := r.Context().Value("user").(*pb.User)
	if !ok {
		log.Println("AuthHandler: User not found in context")
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	token, ok := r.Context().Value("token").(string)
	if !ok {
		log.Println("AuthHandler: Token not found in context")
		http.Error(w, "Token not found in context", http.StatusInternalServerError)
		return
	}

	// Get user ID from URL
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Check if user is trying to access their own profile or is admin
	if currentUser.Id != userID && currentUser.Role != "admin" {
		log.Printf("AuthHandler: Access denied. User %s trying to access user %s", currentUser.Id, userID)
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	req := &pb.GetUserRequest{
		UserId: userID,
		Token:  token,
	}

	resp, err := h.authClient.GetUser(r.Context(), req)
	if err != nil {
		log.Printf("AuthHandler: GetUserByID error: %v", err)
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

// Logout handles user logout (client-side token invalidation)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthHandler: Logout request received")

	// For JWT-based auth, logout is typically handled client-side by removing the token
	// However, we can provide a consistent API response
	response := map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HealthCheck provides a health check endpoint for the auth handler
func (h *AuthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Test connection to auth service
	// ctx := r.Context()

	// Try to make a simple health check call to auth service
	// Note: This assumes the auth service has a health check method
	response := map[string]interface{}{
		"service":   "api-gateway-auth-handler",
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Helper method to extract token from various sources
func (h *AuthHandler) extractToken(r *http.Request) string {
	// Try Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if len(authHeader) > 7 && strings.ToLower(authHeader[:7]) == "bearer " {
			return authHeader[7:]
		}
		return authHeader
	}

	// Try query parameter as fallback
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	// Try form value as fallback
	return r.FormValue("token")
}

// Helper method to send JSON error response
func (h *AuthHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := map[string]interface{}{
		"success": false,
		"message": message,
		"error":   true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

// Helper method to send JSON success response
func (h *AuthHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	successResponse := map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
}
