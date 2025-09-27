package handlers

import (
	"html/template"
	"net/http"

	"github.com/microservices-tutorial/services/client-service/clients"
	"github.com/microservices-tutorial/services/client-service/utils"
	"github.com/gorilla/sessions"
)

type AuthHandler struct {
	apiClient *clients.APIClient
	store     *sessions.CookieStore
}

func NewAuthHandler(apiClient *clients.APIClient, store *sessions.CookieStore) *AuthHandler {
	return &AuthHandler{
		apiClient: apiClient,
		store:     store,
	}
}

func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout/base.html", "templates/auth/login.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Login",
		"Error": r.URL.Query().Get("error"),
	}

	tmpl.Execute(w, data)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	req := clients.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := h.apiClient.Login(r.Context(), req)
	if err != nil {
		http.Redirect(w, r, "/login?error=server_error", http.StatusFound)
		return
	}

	loginResp := resp.(*clients.LoginResponse)
	if !loginResp.Response.Success {
		http.Redirect(w, r, "/login?error="+loginResp.Response.Message, http.StatusFound)
		return
	}

	// Save user to session
	err = utils.SaveUserToSession(w, r, h.store, loginResp.User, loginResp.Token)
	if err != nil {
		http.Redirect(w, r, "/login?error=session_error", http.StatusFound)
		return
	}

	// Redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (h *AuthHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout/base.html", "templates/auth/register.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Register",
		"Error": r.URL.Query().Get("error"),
	}

	tmpl.Execute(w, data)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")

	// Basic validation
	if password != confirmPassword {
		http.Redirect(w, r, "/register?error=passwords_do_not_match", http.StatusFound)
		return
	}

	req := clients.RegisterRequest{
		Email:     email,
		Username:  username,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}

	resp, err := h.apiClient.Register(r.Context(), req)
	if err != nil {
		http.Redirect(w, r, "/register?error=server_error", http.StatusFound)
		return
	}

	registerResp := resp.(*clients.RegisterResponse)
	if !registerResp.Response.Success {
		http.Redirect(w, r, "/register?error="+registerResp.Response.Message, http.StatusFound)
		return
	}

	// Save user to session
	err = utils.SaveUserToSession(w, r, h.store, registerResp.User, registerResp.Token)
	if err != nil {
		http.Redirect(w, r, "/register?error=session_error", http.StatusFound)
		return
	}

	// Redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := utils.ClearSession(w, r, h.store)
	if err != nil {
		// Log error but still redirect
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AuthHandler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	user, token, err := utils.GetUserFromSession(r, h.store)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get fresh user data from API
	resp, err := h.apiClient.GetProfile(r.Context(), token)
	if err != nil {
		http.Redirect(w, r, "/login?error=token_expired", http.StatusFound)
		return
	}

	profileResp := resp.(*clients.ProfileResponse)
	if !profileResp.Response.Success {
		http.Redirect(w, r, "/login?error="+profileResp.Response.Message, http.StatusFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/layout/base.html", "templates/auth/profile.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":   "Profile",
		"User":    profileResp.User,
		"Success": r.URL.Query().Get("success"),
		"Error":   r.URL.Query().Get("error"),
	}

	tmpl.Execute(w, data)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	user, token, err := utils.GetUserFromSession(r, h.store)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	username := r.FormValue("username")

	req := clients.UpdateProfileRequest{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
	}

	resp, err := h.apiClient.UpdateProfile(r.Context(), token, req)
	if err != nil {
		http.Redirect(w, r, "/profile?error=server_error", http.StatusFound)
		return
	}

	updateResp := resp.(*clients.ProfileResponse)
	if !updateResp.Response.Success {
		http.Redirect(w, r, "/profile?error="+updateResp.Response.Message, http.StatusFound)
		return
	}

	// Update user in session
	err = utils.SaveUserToSession(w, r, h.store, updateResp.User, token)
	if err != nil {
		http.Redirect(w, r, "/profile?error=session_error", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/profile?success=profile_updated", http.StatusFound)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	user, token, err := utils.GetUserFromSession(r, h.store)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	if newPassword != confirmPassword {
		http.Redirect(w, r, "/profile?error=passwords_do_not_match", http.StatusFound)
		return
	}

	req := clients.ChangePasswordRequest{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}

	resp, err := h.apiClient.ChangePassword(r.Context(), token, req)
	if err != nil {
		http.Redirect(w, r, "/profile?error=server_error", http.StatusFound)
		return
	}

	changeResp := resp.(*clients.APIResponse)
	if !changeResp.Response.Success {
		http.Redirect(w, r, "/profile?error="+changeResp.Response.Message, http.StatusFound)
		return
	}

	http.Redirect(w, r, "/profile?success=password_changed", http.StatusFound)
}