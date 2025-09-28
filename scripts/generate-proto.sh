# #!/bin/bash

# # Generate Protocol Buffer files for Go microservices

# set -e

# echo "🚀 Generating Protocol Buffer files..."

# # Check if protoc is installed
# if ! command -v protoc &> /dev/null; then
#     echo "❌ protoc is not installed. Please install Protocol Buffers compiler."
#     echo "   Visit: https://grpc.io/docs/protoc-installation/"
#     exit 1
# fi

# # Check if protoc-gen-go is installed
# if ! command -v protoc-gen-go &> /dev/null; then
#     echo "❌ protoc-gen-go is not installed. Installing..."
#     go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# fi

# # Check if protoc-gen-go-grpc is installed
# if ! command -v protoc-gen-go-grpc &> /dev/null; then
#     echo "❌ protoc-gen-go-grpc is not installed. Installing..."
#     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# fi

# # Create output directories
# mkdir -p proto/common
# mkdir -p proto/auth
# mkdir -p proto/product

# echo "📦 Generating common protobuf files..."
# protoc --go_out=. --go_opt=paths=source_relative \
#     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#     proto/common/common.proto

# echo "🔐 Generating auth protobuf files..."
# protoc --go_out=. --go_opt=paths=source_relative \
#     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#     --proto_path=. \
#     proto/auth/auth.proto

# echo "📦 Generating product protobuf files..."
# protoc --go_out=. --go_opt=paths=source_relative \
#     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#     --proto_path=. \
#     proto/product/product.proto

# echo "✅ Protocol Buffer files generated successfully!"
# echo ""
# echo "📁 Generated files:"
# echo "   - proto/common/common.pb.go"
# echo "   - proto/common/common_grpc.pb.go"
# echo "   - proto/auth/auth.pb.go"
# echo "   - proto/auth/auth_grpc.pb.go"
# echo "   - proto/product/product.pb.go"
# echo "   - proto/product/product_grpc.pb.go"
# echo ""
# echo "🎉 Ready to build your microservices!"

#!/bin/bash

# Simple fix - Use a single module approach instead of separate proto modules
# This is easier and avoids complex module dependencies

set -e

echo "🔧 Applying simple protobuf fix..."

MODULE_NAME="github.com/martbul/playground_microservices"

# Step 1: Create a single root go.mod for the entire project
cat > go.mod << EOF
module $MODULE_NAME

go 1.21

require (
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/gorilla/mux v1.8.1
    github.com/gorilla/sessions v1.2.2
    github.com/lib/pq v1.10.9
    golang.org/x/crypto v0.17.0
    google.golang.org/grpc v1.60.1
    google.golang.org/protobuf v1.31.0
)

require (
    github.com/golang/protobuf v1.5.3 // indirect
    github.com/gorilla/securecookie v1.1.2 // indirect
    golang.org/x/net v0.16.0 // indirect
    golang.org/x/sys v0.15.0 // indirect
    golang.org/x/text v0.14.0 // indirect
    google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
)
EOF

# Step 2: Remove all service-level go.mod files
rm -f services/*/go.mod
rm -f proto/*/go.mod

# Step 3: Update protobuf go_package options
sed -i.bak 's|option go_package = ".*";|option go_package = "'$MODULE_NAME'/proto/common";|' proto/common/common.proto
sed -i.bak 's|option go_package = ".*";|option go_package = "'$MODULE_NAME'/proto/auth";|' proto/auth/auth.proto  
sed -i.bak 's|option go_package = ".*";|option go_package = "'$MODULE_NAME'/proto/product";|' proto/product/product.proto

# Remove backup files
rm -f proto/*/*.bak

# Step 4: Regenerate protobuf files
echo "📦 Regenerating protobuf files..."

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/common/common.proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --proto_path=. \
    proto/auth/auth.proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --proto_path=. \
    proto/product/product.proto

# Step 5: Update all Go files with correct import paths
echo "🔄 Updating import statements..."

find . -name "*.go" -not -path "./proto/*" | while read file; do
    # Update proto imports
    sed -i.bak "s|github.com/microservices-tutorial/proto/|$MODULE_NAME/proto/|g" "$file"
    sed -i.bak "s|github.com/microservices-tutorial/services/|$MODULE_NAME/services/|g" "$file"
    
    # Remove backup
    rm -f "${file}.bak"
done

# Step 6: Download dependencies
echo "📥 Downloading dependencies..."
go mod tidy

echo ""
echo "✅ Simple protobuf fix applied successfully!"
echo ""
echo "Now you can:"
echo "  1. Build services: go build ./services/auth-service"
echo "  2. Or use Docker: docker-compose up --build"
echo ""
echo "All services now use a single module: $MODULE_NAME"