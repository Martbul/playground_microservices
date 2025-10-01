package clients

import (
	"context"
	"fmt"
	"time"

	pb "github.com/martbul/playground_microservices/services/api-gateway/genproto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//the client implementation of the grpc .proto service
type AuthGrpcClient struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

//! use secure credentials
func NewAuthGrpcClient(address string) (*AuthGrpcClient, error) {
	//! always close a conn with defer(can cause memory leak)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}
	defer conn.Close()

	//creating new authServiceClient and passing the connection
	client := pb.NewAuthServiceClient(conn)

	return &AuthGrpcClient{
		client: client,
		conn:   conn,
	}, nil
}

// believe this methods here are for interaction with the server handler
func (c *AuthGrpcClient) Close() error {
	return c.conn.Close()
}

func (c *AuthGrpcClient) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel() //cancel if req takes to long

	return c.client.Register(ctx, req)
}

func (c *AuthGrpcClient) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.Login(ctx, req)
}

func (c *AuthGrpcClient) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.ValidateToken(ctx, req)
}

func (c *AuthGrpcClient) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.GetUser(ctx, req)
}

func (c *AuthGrpcClient) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.UpdateProfile(ctx, req)
}

func (c *AuthGrpcClient) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.ChangePassword(ctx, req)
}

func (c *AuthGrpcClient) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.RefreshToken(ctx, req)
}
