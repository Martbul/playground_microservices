#!/bin/bash

# Exit on first error
set -e

# Define proto paths
PROTO_DIR="./proto"
OUT_DIR="./proto"

# Check if protoc is installed
if ! command -v protoc &> /dev/null
then
    echo "protoc not found. Please install protoc first."
    exit 1
fi

# Generate Go code for common.proto
echo "Generating Go code for common.proto..."
protoc -I=$PROTO_DIR \
  --go_out=$OUT_DIR --go_opt=paths=source_relative \
  --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative \
  $PROTO_DIR/common/common.proto

# Generate Go code for auth.proto
echo "Generating Go code for auth.proto..."
protoc -I=$PROTO_DIR \
  --go_out=$OUT_DIR --go_opt=paths=source_relative \
  --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative \
  $PROTO_DIR/auth/auth.proto

echo "✅ Protobuf generation completed successfully."
