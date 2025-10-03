package clients

import (
	"context"
	"fmt"
	"time"

)


// AuthClient handles authentication-related API calls
type AuthClient struct {
	apiClient *APIClient
}

// NewAuthClient creates a new auth client
func NewAuthClient(apiClient *APIClient) *AuthClient {
	return &AuthClient{
		apiClient: apiClient,
	}
}


// Register registers a new user
func (c *AuthClient) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	return c.apiClient.Register(ctx, req)
}

// Login authenticates a user
func (c *AuthClient) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	return c.apiClient.Login(ctx, req)
}

// GetProfile gets the current user's profile
func (c *AuthClient) GetProfile(ctx context.Context, token string) (*ProfileResponse, error) {
	return c.apiClient.GetProfile(ctx, token)
}

// UpdateProfile updates the user's profile
func (c *AuthClient) UpdateProfile(ctx context.Context, token string, req UpdateProfileRequest) (*ProfileResponse, error) {
	return c.apiClient.UpdateProfile(ctx, token, req)
}

// ChangePassword changes the user's password
func (c *AuthClient) ChangePassword(ctx context.Context, token string, req ChangePasswordRequest) (*APIResponse, error) {
	return c.apiClient.ChangePassword(ctx, token, req)
}

// ValidateToken validates a JWT token
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Make HTTP request to validate token
	req := ValidateTokenRequest{
		Token: token,
	}

	resp, err := c.apiClient.post(ctx, "/api/auth/validate", req, &ValidateTokenResponse{})
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	validateResp := resp.(*ValidateTokenResponse)
	return validateResp, nil
}

// RefreshToken refreshes an expired JWT token
func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	resp, err := c.apiClient.post(ctx, "/api/auth/refresh", req, &RefreshTokenResponse{})
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	refreshResp := resp.(*RefreshTokenResponse)
	return refreshResp, nil
}

// Logout logs out the current user (client-side)
func (c *AuthClient) Logout(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.apiClient.post(ctx, "/api/auth/logout", nil, &APIResponse{})
	return err
}

// Additional request/response types specific to auth

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Response Response `json:"response"`
	Valid    bool        `json:"valid"`
	User     User        `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Response     bool `json:"response"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    int64       `json:"expires_at"`
}

// Helper methods for token management

// IsTokenValid checks if a token is valid by making an API call
func (c *AuthClient) IsTokenValid(ctx context.Context, token string) bool {
	if token == "" {
		return false
	}

	resp, err := c.ValidateToken(ctx, token)
	if err != nil {
		return false
	}

	return resp.Response.Success && resp.Valid
}

// GetUserFromToken validates token and returns user info
func (c *AuthClient) GetUserFromToken(ctx context.Context, token string) (*User, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	resp, err := c.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	if !resp.Response.Success || !resp.Valid {
		return nil, fmt.Errorf("token is invalid: %s", resp.Response.Message)
	}

	return &resp.User, nil
}

// ExtractUserID extracts user ID from a valid token
func (c *AuthClient) ExtractUserID(ctx context.Context, token string) (string, error) {
	user, err := c.GetUserFromToken(ctx, token)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}