package handlers

import (
	"context"
	"log"
	"time"

	pb "github.com/martbul/playground_microservices/services/auth-service/genproto/auth"
	commonPb "github.com/martbul/playground_microservices/services/auth-service/genproto/common"
	"github.com/martbul/playground_microservices/services/auth-service/models"
	"github.com/martbul/playground_microservices/services/auth-service/service"
)

//here is where it creates a handler for the already implemented .proto services
//and i beleve(not sure) these are the actual methods that are exposed to grpc clients(but not sure)
type AuthGrpcHandler struct {
	//this has receiver methods
	pb.UnimplementedAuthServiceServer //embedig a auto generated struct into the server
	authService service.AuthService
}

func NewAuthGrpcHandler(authService service.AuthService) *AuthGrpcHandler {
	return &AuthGrpcHandler{
		authService: authService,
	}
}


//regiser is specifed in the auth.proto it is implemented here
func (h *AuthGrpcHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Register request for email: %s", req.Email)

	registerReq := &models.RegisterRequest{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	

	user, token, err := h.authService.Register(registerReq)
	if err != nil {
		log.Printf("Register error: %v", err)
		return &pb.RegisterResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}


	return &pb.RegisterResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "User registered successfully",
		},
		User:  h.userToProto(user),
		Token: token,
	}, nil
}

func (h *AuthGrpcHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Login request for email: %s", req.Email)

	loginReq := &models.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, refreshToken, err := h.authService.Login(loginReq)
	if err != nil {
		log.Printf("Login error: %v", err)
		return &pb.LoginResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.LoginResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Login successful",
		},
		User:         h.userToProto(user),
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func (h *AuthGrpcHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	log.Printf("Token validation request")

	user, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Token validation error: %v", err)
		return &pb.ValidateTokenResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
			Valid: false,
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Token is valid",
		},
		Valid: true,
		User:  h.userToProto(user),
	}, nil
}

func (h *AuthGrpcHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Get user request for ID: %s", req.UserId)

	// Validate token first
	_, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Token validation error: %v", err)
		return &pb.GetUserResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: "Invalid token",
			},
		}, nil
	}

	user, err := h.authService.GetUser(req.UserId)
	if err != nil {
		log.Printf("Get user error: %v", err)
		return &pb.GetUserResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.GetUserResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "User retrieved successfully",
		},
		User: h.userToProto(user),
	}, nil
}

func (h *AuthGrpcHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	log.Printf("Update profile request for user ID: %s", req.UserId)

	// Validate token first
	_, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Token validation error: %v", err)
		return &pb.UpdateProfileResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: "Invalid token",
			},
		}, nil
	}

	updateReq := &models.UpdateProfileRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
	}

	user, err := h.authService.UpdateProfile(req.UserId, updateReq)
	if err != nil {
		log.Printf("Update profile error: %v", err)
		return &pb.UpdateProfileResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.UpdateProfileResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Profile updated successfully",
		},
		User: h.userToProto(user),
	}, nil
}

func (h *AuthGrpcHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	log.Printf("Change password request for user ID: %s", req.UserId)

	// Validate token first
	_, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Token validation error: %v", err)
		return &pb.ChangePasswordResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: "Invalid token",
			},
		}, nil
	}

	changeReq := &models.ChangePasswordRequest{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	err = h.authService.ChangePassword(req.UserId, changeReq)
	if err != nil {
		log.Printf("Change password error: %v", err)
		return &pb.ChangePasswordResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.ChangePasswordResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Password changed successfully",
		},
	}, nil
}

func (h *AuthGrpcHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	log.Printf("Refresh token request")

	newToken, newRefreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		log.Printf("Refresh token error: %v", err)
		return &pb.RefreshTokenResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.RefreshTokenResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Token refreshed successfully",
		},
		Token:        newToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func (h *AuthGrpcHandler) HealthCheck(ctx context.Context, req *commonPb.HealthCheckRequest) (*commonPb.HealthCheckResponse, error) {
	return &commonPb.HealthCheckResponse{
		Status:    "healthy",
		Service:   "auth-service",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

func (h *AuthGrpcHandler) userToProto(user *models.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
