	package clients

	import (
		"context"
		"fmt"
		"time"

			pb "github.com/martbul/playground_microservices/services/api-gateway/genproto/product"

		"google.golang.org/grpc"
		"google.golang.org/grpc/credentials/insecure"
	)

type ProductClient struct {
	client pb.ProductServiceClient
	conn   *grpc.ClientConn
}

func NewProductClient(address string) (*ProductClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	client := pb.NewProductServiceClient(conn)

	return &ProductClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *ProductClient) Close() error {
	return c.conn.Close()
}

func (c *ProductClient) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.CreateProduct(ctx, req)
}

func (c *ProductClient) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.GetProduct(ctx, req)
}

func (c *ProductClient) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.UpdateProduct(ctx, req)
}

//!TODO: Check
// func (c *ProductClient) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductRequest, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	resp, err := c.client.DeleteProduct(ctx, req)
// 	return req, err
// }

func (c *ProductClient) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.client.DeleteProduct(ctx, req)
	return req, err
}

func (c *ProductClient) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.ListProducts(ctx, req)
}

func (c *ProductClient) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.SearchProducts(ctx, req)
}

func (c *ProductClient) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.GetCategories(ctx, req)
}