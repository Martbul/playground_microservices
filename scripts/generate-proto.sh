#!/bin/bash

set -e

echo "=== Fixing Microservices Proto Setup ==="
echo "Current directory: $(pwd)"
echo ""

# Step 1: Clean up any nested proto directories and old generated files
echo "🧹 Cleaning up old generated files..."
if [ -d "proto/proto" ]; then
    rm -rf proto/proto
    echo "Removed nested proto/proto directory"
fi

find proto -name "*.pb.go" -delete 2>/dev/null || true
find proto -name "*_grpc.pb.go" -delete 2>/dev/null || true
echo "Cleaned up old generated files"

# Step 2: Fix proto files to use relative imports
echo ""
echo "🔧 Fixing proto import paths..."

# Fix auth.proto
cat > proto/auth/auth.proto << 'EOF'
syntax = "proto3";

package auth;

import "common/common.proto";  // Relative import

option go_package = "github.com/martbul/playground_microservices/proto/auth";

service AuthService {
    rpc Register(RegisterRequest) returns (RegisterResponse);
    rpc Login(LoginRequest) returns (LoginResponse);
    rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc UpdateProfile(UpdateProfileRequest) returns (UpdateProfileResponse);
    rpc ChangePassword(ChangePasswordRequest) returns (ChangePasswordResponse);
    rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
    rpc HealthCheck(common.HealthCheckRequest) returns (common.HealthCheckResponse);
}

//User model
message User {
    string id = 1;
    string email = 2;
    string username = 3;
    string first_name = 4;
    string last_name = 5;
    string role = 6;
    bool is_active = 7;  // Fixed field name
    string created_at = 8;
    string updated_at = 9;
}

message RegisterRequest {
    string email = 1;
    string username = 2;
    string password = 3;
    string first_name = 4;  // Fixed typo
    string last_name = 5;
}

message RegisterResponse {
    common.Response response = 1;
    User user = 2;
    string token = 3;
}

message LoginRequest {
    string email = 1;
    string password = 2;
}

message LoginResponse {
    common.Response response = 1;
    User user = 2;
    string token = 3;
    string refresh_token = 4;
    int64 expires_at = 5;
}

message ValidateTokenRequest {
    string token = 1;
}

message ValidateTokenResponse {
    common.Response response = 1;
    bool valid = 2;
    User user = 3;
}

message GetUserRequest {
    string user_id = 1;
    string token = 2;
}

message GetUserResponse {
    common.Response response = 1;
    User user = 2;
}

message UpdateProfileRequest {
    string user_id = 1;
    string token = 2;
    string first_name = 3;
    string last_name = 4;
    string username = 5;
}

message UpdateProfileResponse {
    common.Response response = 1;
    User user = 2;
}

message ChangePasswordRequest {
    string user_id = 1;
    string token = 2;
    string current_password = 3;
    string new_password = 4;
}

message ChangePasswordResponse {
    common.Response response = 1;
}

message RefreshTokenRequest {
    string refresh_token = 1;
}

message RefreshTokenResponse {
    common.Response response = 1;
    string token = 2;
    string refresh_token = 3;
    int64 expires_at = 4;
}
EOF

# Fix product.proto
cat > proto/product/product.proto << 'EOF'
syntax = "proto3";

package product;

import "common/common.proto";  // Relative import

option go_package = "github.com/martbul/playground_microservices/proto/product";

service ProductService {
  rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  rpc UpdateProduct(UpdateProductRequest) returns (UpdateProductResponse);
  rpc DeleteProduct(DeleteProductRequest) returns (DeleteProductResponse);
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc SearchProducts(SearchProductsRequest) returns (SearchProductsResponse);
  rpc GetCategories(GetCategoriesRequest) returns (GetCategoriesResponse);
  rpc HealthCheck(common.HealthCheckRequest) returns (common.HealthCheckResponse);
}

// Product model
message Product {
  string id = 1;
  string name = 2;
  string description = 3;
  double price = 4;
  int32 stock_quantity = 5;
  string category = 6;
  string image_url = 7;
  string sku = 8;
  bool is_active = 9;
  string created_at = 10;
  string updated_at = 11;
  string created_by = 12;
}

// Category model
message Category {
  string id = 1;
  string name = 2;
  string description = 3;
  string parent_id = 4;
  bool is_active = 5;
}

// Create product
message CreateProductRequest {
  string token = 1;
  string name = 2;
  string description = 3;
  double price = 4;
  int32 stock_quantity = 5;
  string category = 6;
  string image_url = 7;
  string sku = 8;
}

message CreateProductResponse {
  common.Response response = 1;
  Product product = 2;
}

// Get product
message GetProductRequest {
  string id = 1;
}

message GetProductResponse {
  common.Response response = 1;
  Product product = 2;
}

// Update product
message UpdateProductRequest {
  string token = 1;
  string id = 2;
  string name = 3;
  string description = 4;
  double price = 5;
  int32 stock_quantity = 6;
  string category = 7;
  string image_url = 8;
  string sku = 9;
  bool is_active = 10;
}

message UpdateProductResponse {
  common.Response response = 1;
  Product product = 2;
}

// Delete product
message DeleteProductRequest {
  string token = 1;
  string id = 2;
}

message DeleteProductResponse {
  common.Response response = 1;
}

// List products
message ListProductsRequest {
  common.PaginationRequest pagination = 1;
  string category = 2;
  bool active_only = 3;
}

message ListProductsResponse {
  common.Response response = 1;
  repeated Product products = 2;
  common.PaginationResponse pagination = 3;
}

// Search products
message SearchProductsRequest {
  string query = 1;
  common.PaginationRequest pagination = 2;
  string category = 3;
  double min_price = 4;
  double max_price = 5;
}

message SearchProductsResponse {
  common.Response response = 1;
  repeated Product products = 2;
  common.PaginationResponse pagination = 3;
}

// Get categories
message GetCategoriesRequest {}

message GetCategoriesResponse {
  common.Response response = 1;
  repeated Category categories = 2;
}
EOF

echo "✅ Fixed proto import paths"

# Step 3: Generate proto files
echo ""
echo "🔄 Generating proto files..."

PROTO_DIR="./proto"
OUT_DIR="./proto"

# Check protoc installation
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc not found. Please install protoc first."
    exit 1
fi

if ! command -v protoc-gen-go &> /dev/null || ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "❌ Go protobuf plugins not found. Installing..."
    echo "Please run:"
    echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# Generate common.proto first
echo "Generating common.proto..."
protoc -I=$PROTO_DIR \
  --go_out=$OUT_DIR --go_opt=paths=source_relative \
  --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative \
  $PROTO_DIR/common/common.proto

# Generate auth.proto
echo "Generating auth.proto..."
protoc -I=$PROTO_DIR \
  --go_out=$OUT_DIR --go_opt=paths=source_relative \
  --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative \
  $PROTO_DIR/auth/auth.proto

# Generate product.proto
echo "Generating product.proto..."
protoc -I=$PROTO_DIR \
  --go_out=$OUT_DIR --go_opt=paths=source_relative \
  --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative \
  $PROTO_DIR/product/product.proto

echo ""
echo "✅ Proto generation completed!"
echo ""
echo "📋 Generated files:"
find proto -name "*.pb.go" -o -name "*_grpc.pb.go" | sort

# Step 4: Setup service dependencies
echo ""
echo "🔧 Setting up service dependencies..."

setup_service_deps() {
    local service_dir="$1"
    local service_name=$(basename "$service_dir")
    
    echo ""
    echo "Setting up $service_name..."
    
    if [ ! -f "$service_dir/go.mod" ]; then
        echo "❌ No go.mod found in $service_name"
        return 1
    fi
    
    cd "$service_dir"
    
    # Add protobuf dependencies
    go get google.golang.org/protobuf@latest
    go get google.golang.org/grpc@latest
    go get google.golang.org/protobuf/runtime/protoimpl@latest
    go get google.golang.org/protobuf/reflect/protoreflect@latest
    go get google.golang.org/protobuf/types/known/timestamppb@latest
    go get google.golang.org/grpc/codes@latest
    go get google.golang.org/grpc/status@latest
    
    # Clean up
    go mod tidy
    
    echo "✅ $service_name dependencies updated"
    
    cd - > /dev/null
}

# Setup dependencies for each service
for service_dir in services/*/; do
    if [ -d "$service_dir" ]; then
        setup_service_deps "$service_dir"
    fi
done

echo ""
echo "🎉 Microservices proto setup completed!"
echo ""
echo "📋 Your services can now import proto packages like:"
echo "  import pb \"github.com/martbul/playground_microservices/proto/auth\""
echo "  import pb \"github.com/martbul/playground_microservices/proto/product\""
echo "  import pb \"github.com/martbul/playground_microservices/proto/common\"" 