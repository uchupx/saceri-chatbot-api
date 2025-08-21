.PHONY: build run proto vendor 

build:
		@echo "Building the application..."
		@go build -o bin/api cmd/api/main.go

run:
		@echo "Running the application..."
		@go run cmd/api/main.go

proto:
		@echo "Generating protobuf code..."
		@protoc --go_out=pkg/grpc/proto/gen  \
        --go-grpc_out=pkg/grpc/proto/gen \
        pkg/grpc/proto/*.proto

vendor:
		@echo "Installing dependencies..."
		@go mod tidy
		@go mod vendor

.DEFAULT_GOAL := build
