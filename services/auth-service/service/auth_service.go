package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/martbul/playground_microservices/services/auth-service/models"
	"github.com/martbul/playground_microservices/services/auth-service/repository"
	"github.com/martbul/playground_microservices/services/auth-service/utils"
)


//believe that the service is for interaction with the db and the handler is the actual proto service
type AuthService interface {
	Register(req *models.RegisterRequest) (*models.User, string, error)
	Login(req *models.LoginRequest) (*models.User, string, string, error)
	ValidateToken(token string) (*models.User, error)
	GetUser(userID string) (*models.User, error)
	UpdateProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error)
	ChangePassword(userID string, req *models.ChangePasswordRequest) error
	RefreshToken(refreshToken string) (string, string, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *utils.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: utils.NewJWTManager(jwtSecret),
	}
}

func (s *authService) Register(req *models.RegisterRequest) (*models.User, string, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, "", fmt.Errorf("user with email already exists")
	}

	// Check if username already exists
	existingUsername, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check existing username: %w", err)
	}
	if existingUsername != nil {
		return nil, "", fmt.Errorf("username already taken")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        strings.ToLower(req.Email),
		Username:     req.Username,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "user",
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username, user.Role, 24*time.Hour)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *authService) Login(req *models.LoginRequest) (*models.User, string, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(strings.ToLower(req.Email))
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, "", "", fmt.Errorf("account is disabled")
	}

	// Check password
	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username, user.Role, 24*time.Hour)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate refresh token
	refreshTokenStr, err := utils.GenerateRandomToken(32)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: utils.HashString(refreshTokenStr),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := s.userRepo.SaveRefreshToken(refreshToken); err != nil {
		return nil, "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return user, token, refreshTokenStr, nil
}

func (s *authService) ValidateToken(token string) (*models.User, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is disabled")
	}

	return user, nil
}

func (s *authService) GetUser(userID string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *authService) UpdateProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Username != "" && req.Username != user.Username {
		// Check if username is already taken
		existingUser, err := s.userRepo.GetByUsername(req.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing username: %w", err)
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, fmt.Errorf("username already taken")
		}
		user.Username = req.Username
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *authService) ChangePassword(userID string, req *models.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	if err := utils.CheckPassword(req.CurrentPassword, user.PasswordHash); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(userID, newPasswordHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	tokenHash := utils.HashString(refreshToken)

	// Get refresh token from database
	storedToken, err := s.userRepo.GetRefreshToken(tokenHash)
	if err != nil {
		return "", "", fmt.Errorf("failed to get refresh token: %w", err)
	}
	if storedToken == nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(storedToken.UserID)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return "", "", fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return "", "", fmt.Errorf("account is disabled")
	}

	// Generate new JWT token
	newToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username, user.Role, 24*time.Hour)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate new refresh token
	newRefreshTokenStr, err := utils.GenerateRandomToken(32)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Delete old refresh token
	if err := s.userRepo.DeleteRefreshToken(tokenHash); err != nil {
		return "", "", fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	// Save new refresh token
	newRefreshToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: utils.HashString(newRefreshTokenStr),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := s.userRepo.SaveRefreshToken(newRefreshToken); err != nil {
		return "", "", fmt.Errorf("failed to save new refresh token: %w", err)
	}

	return newToken, newRefreshTokenStr, nil
}