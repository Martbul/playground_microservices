#!/bin/bash

# Deploy script for Kubernetes

set -e

echo "🚀 Deploying Go Microservices to Kubernetes..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install kubectl."
    exit 1
fi

# Check if kubectl can connect to cluster
if ! kubectl cluster-info &> /dev/null; then
    print_error "Cannot connect to Kubernetes cluster. Please check your kubectl configuration."
    exit 1
fi

print_success "Connected to Kubernetes cluster"

# Create namespace
print_status "Creating namespace..."
kubectl apply -f k8s/namespace.yaml

# Apply ConfigMaps
print_status "Applying ConfigMaps..."
kubectl apply -f k8s/configmaps/

# Apply Services
print_status "Applying Services..."
kubectl apply -f k8s/services/

# Apply Deployments
print_status "Applying Deployments..."
kubectl apply -f k8s/deployments/

# Apply Ingress
print_status "Applying Ingress..."
kubectl apply -f k8s/ingress/

print_success "Deployment completed!"

# Wait for deployments to be ready
print_status "Waiting for deployments to be ready..."

DEPLOYMENTS=("postgres" "auth-service" "product-service" "api-gateway" "client-service")

for deployment in "${DEPLOYMENTS[@]}"; do
    print_status "Waiting for $deployment to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment/$deployment -n microservices-tutorial
    if [[ $? -eq 0 ]]; then
        print_success "$deployment is ready"
    else
        print_error "$deployment failed to become ready"
    fi
done

# Show status
print_status "Deployment Status:"
kubectl get all -n microservices-tutorial

echo ""
print_success "🎉 Deployment completed successfully!"
echo ""
echo "📋 Useful commands:"
echo "   View pods: kubectl get pods -n microservices-tutorial"
echo "   View services: kubectl get services -n microservices-tutorial"
echo "   View logs: kubectl logs -f deployment/<service-name> -n microservices-tutorial"
echo "   Port forward: kubectl port-forward service/client-service 8083:80 -n microservices-tutorial"
echo ""
echo "🌐 Access your application:"
echo "   If using port-forward: http://localhost:8083"
echo "   If using Ingress: Check your ingress controller configuration"