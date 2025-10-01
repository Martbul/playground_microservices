package handlers

import (
	"context"
	"log"
	"time"

	commonPb "github.com/martbul/playground_microservices/services/product-service/genproto/common"
	pb "github.com/martbul/playground_microservices/services/product-service/genproto/product"
	"github.com/martbul/playground_microservices/services/product-service/models"
	"github.com/martbul/playground_microservices/services/product-service/service"
)

type ProductGrpcHandler struct {
	pb.UnimplementedProductServiceServer
	productService service.ProductService
}

func NewProductGrpcHandler(productService service.ProductService) *ProductGrpcHandler {
	return &ProductGrpcHandler{
		productService: productService,
	}
}

func (h *ProductGrpcHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	log.Printf("Create product request: %s", req.Name)

	// TODO: Validate token and get user info from auth service
	// For now, we'll use a placeholder
	createdBy := "system" // This should come from token validation

	createReq := &models.CreateProductRequest{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Category:      req.Category,
		ImageURL:      req.ImageUrl,
		SKU:           req.Sku,
	}

	product, err := h.productService.CreateProduct(createReq, createdBy)
	if err != nil {
		log.Printf("Create product error: %v", err)
		return &pb.CreateProductResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.CreateProductResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Product created successfully",
		},
		Product: h.productToProto(product),
	}, nil
}

func (h *ProductGrpcHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	log.Printf("Get product request: %s", req.Id)

	product, err := h.productService.GetProduct(req.Id)
	if err != nil {
		log.Printf("Get product error: %v", err)
		return &pb.GetProductResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.GetProductResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Product retrieved successfully",
		},
		Product: h.productToProto(product),
	}, nil
}

func (h *ProductGrpcHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	log.Printf("Update product request: %s", req.Id)

	// TODO: Validate token
	updateReq := &models.UpdateProductRequest{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Category:      req.Category,
		ImageURL:      req.ImageUrl,
		SKU:           req.Sku,
	}

	if req.IsActive {
		isActive := req.IsActive
		updateReq.IsActive = &isActive
	}

	product, err := h.productService.UpdateProduct(req.Id, updateReq)
	if err != nil {
		log.Printf("Update product error: %v", err)
		return &pb.UpdateProductResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.UpdateProductResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Product updated successfully",
		},
		Product: h.productToProto(product),
	}, nil
}

func (h *ProductGrpcHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	log.Printf("Delete product request: %s", req.Id)

	// TODO: Validate token

	err := h.productService.DeleteProduct(req.Id)
	if err != nil {
		log.Printf("Delete product error: %v", err)
		return &pb.DeleteProductResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.DeleteProductResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Product deleted successfully",
		},
	}, nil
}

func (h *ProductGrpcHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	log.Printf("List products request")

	var filter *models.ProductFilter
	if req.Category != "" || req.ActiveOnly {
		filter = &models.ProductFilter{
			Category: req.Category,
		}
		if req.ActiveOnly {
			active := true
			filter.IsActive = &active
		}
	}

	pagination := &models.PaginationRequest{
		Page:      req.Pagination.Page,
		Limit:     req.Pagination.Limit,
		SortBy:    req.Pagination.SortBy,
		SortOrder: req.Pagination.SortOrder,
	}

	products, paginationResp, err := h.productService.ListProducts(filter, pagination)
	if err != nil {
		log.Printf("List products error: %v", err)
		return &pb.ListProductsResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	var protoProducts []*pb.Product
	for _, product := range products {
		protoProducts = append(protoProducts, h.productToProto(product))
	}

	return &pb.ListProductsResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Products retrieved successfully",
		},
		Products:   protoProducts,
		Pagination: h.paginationToProto(paginationResp),
	}, nil
}

func (h *ProductGrpcHandler) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	log.Printf("Search products request: %s", req.Query)

	filter := &models.SearchFilter{
		Query:    req.Query,
		Category: req.Category,
		MinPrice: req.MinPrice,
		MaxPrice: req.MaxPrice,
	}

	pagination := &models.PaginationRequest{
		Page:      req.Pagination.Page,
		Limit:     req.Pagination.Limit,
		SortBy:    req.Pagination.SortBy,
		SortOrder: req.Pagination.SortOrder,
	}

	products, paginationResp, err := h.productService.SearchProducts(filter, pagination)
	if err != nil {
		log.Printf("Search products error: %v", err)
		return &pb.SearchProductsResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	var protoProducts []*pb.Product
	for _, product := range products {
		protoProducts = append(protoProducts, h.productToProto(product))
	}

	return &pb.SearchProductsResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Search completed successfully",
		},
		Products:   protoProducts,
		Pagination: h.paginationToProto(paginationResp),
	}, nil
}

func (h *ProductGrpcHandler) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesResponse, error) {
	log.Printf("Get categories request")

	categories, err := h.productService.GetCategories()
	if err != nil {
		log.Printf("Get categories error: %v", err)
		return &pb.GetCategoriesResponse{
			Response: &commonPb.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	var protoCategories []*pb.Category
	for _, category := range categories {
		protoCategories = append(protoCategories, h.categoryToProto(category))
	}

	return &pb.GetCategoriesResponse{
		Response: &commonPb.Response{
			Success: true,
			Message: "Categories retrieved successfully",
		},
		Categories: protoCategories,
	}, nil
}

func (h *ProductGrpcHandler) HealthCheck(ctx context.Context, req *commonPb.HealthCheckRequest) (*commonPb.HealthCheckResponse, error) {
	return &commonPb.HealthCheckResponse{
		Status:    "healthy",
		Service:   "product-service",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

func (h *ProductGrpcHandler) productToProto(product *models.Product) *pb.Product {
	return &pb.Product{
		Id:            product.ID,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		Category:      product.Category,
		ImageUrl:      product.ImageURL,
		Sku:           product.SKU,
		IsActive:      product.IsActive,
		CreatedAt:     product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     product.UpdatedAt.Format(time.RFC3339),
		CreatedBy:     product.CreatedBy,
	}
}

func (h *ProductGrpcHandler) categoryToProto(category *models.Category) *pb.Category {
	var parentID string
	if category.ParentID != nil {
		parentID = *category.ParentID
	}

	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ParentId:    parentID,
		IsActive:    category.IsActive,
	}
}

func (h *ProductGrpcHandler) paginationToProto(pagination *models.PaginationResponse) *commonPb.PaginationResponse {
	return &commonPb.PaginationResponse{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: pagination.TotalPages,
		TotalCount: pagination.TotalCount,
		HasNext:    pagination.HasNext,
		HasPrev:    pagination.HasPrev,
	}
}
