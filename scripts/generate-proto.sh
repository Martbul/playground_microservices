#!/bin/bash

# Generate Protocol Buffer files for Go microservices

set -e

echo "🚀 Generating Protocol Buffer files..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc is not installed. Please install Protocol Buffers compiler."
    echo "   Visit: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "❌ protoc-gen-go is not installed. Installing..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "❌ protoc-gen-go-grpc is not installed. Installing..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Create output directories
mkdir -p proto/common
mkdir -p proto/auth
mkdir -p proto/product

echo "📦 Generating common protobuf files..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/common/common.proto

echo "🔐 Generating auth protobuf files..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --proto_path=. \
    proto/auth/auth.proto

echo "📦 Generating product protobuf files..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --proto_path=. \
    proto/product/product.proto

echo "✅ Protocol Buffer files generated successfully!"
echo ""
echo "📁 Generated files:"
echo "   - proto/common/common.pb.go"
echo "   - proto/common/common_grpc.pb.go"
echo "   - proto/auth/auth.pb.go"
echo "   - proto/auth/auth_grpc.pb.go"
echo "   - proto/product/product.pb.go"
echo "   - proto/product/product_grpc.pb.go"
echo ""
echo "🎉 Ready to build your microservices!"