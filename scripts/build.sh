#!/bin/bash

# Build script for Go microservices

set -e

echo "🚀 Building Go Microservices..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✅${NC} $1"
}

print_error() {
    echo -e "${RED}❌${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠️${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [[ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]]; then
    print_error "Go version $REQUIRED_VERSION or later is required. Current version: $GO_VERSION"
    exit 1
fi

print_success "Go version $GO_VERSION detected"

# Generate protobuf files
print_status "Generating Protocol Buffer files..."
./scripts/generate-proto.sh

# Build services
SERVICES=("auth-service" "product-service" "api-gateway" "client-service")

for service in "${SERVICES[@]}"; do
    print_status "Building $service..."
    
    cd "services/$service"
    
    # Initialize go module if go.mod doesn't exist
    if [[ ! -f "go.mod" ]]; then
        print_warning "Initializing Go module for $service..."
        go mod init "github.com/playground_microservices/services/$service"
    fi
    
    # Download dependencies
    print_status "Downloading dependencies for $service..."
    go mod download
    go mod tidy
    
    # Build the service
    print_status "Compiling $service..."
    go build -o "../../build/$service" .
    
    if [[ $? -eq 0 ]]; then
        print_success "$service built successfully"
    else
        print_error "Failed to build $service"
        exit 1
    fi
    
    cd "../.."
done

# Create build directory if it doesn't exist
mkdir -p build

print_success "All services built successfully!"
echo ""
echo "📁 Built binaries are available in the 'build' directory:"
for service in "${SERVICES[@]}"; do
    echo "   - build/$service"
done
echo ""
echo "🐳 To build Docker images, run: docker-compose build"
echo "🚀 To start all services, run: docker-compose up"