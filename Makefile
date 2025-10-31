.PHONY: help build clean test docker-build docker-push docker-run

# Docker image configuration
DOCKER_REGISTRY ?= ghcr.io/klazomenai
IMAGE_NAME ?= istio-jwt-wasm
IMAGE_TAG ?= latest
FULL_IMAGE = $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Wasm module
	@echo "Building autonity-jwt-revocation.wasm..."
	GOOS=wasip1 GOARCH=wasm go build -o autonity-jwt-revocation.wasm main.go
	@echo "Build complete: autonity-jwt-revocation.wasm"
	@ls -lh autonity-jwt-revocation.wasm

clean: ## Clean build artifacts
	rm -f autonity-jwt-revocation.wasm

test: ## Run Go tests
	go test -v ./...

deps: ## Download Go dependencies
	go mod download
	go mod tidy

verify: ## Verify the Wasm module
	file autonity-jwt-revocation.wasm
	wasm-validate autonity-jwt-revocation.wasm || echo "wasm-validate not installed (optional)"

docker-build: ## Build Docker image with Wasm module
	@echo "Building Docker image: $(FULL_IMAGE)"
	docker build -t $(FULL_IMAGE) .
	@echo "Build complete: $(FULL_IMAGE)"

docker-push: docker-build ## Push Docker image to registry
	@echo "Pushing Docker image: $(FULL_IMAGE)"
	docker push $(FULL_IMAGE)

docker-run: docker-build ## Extract Wasm module from Docker image
	@echo "Extracting plugin.wasm from Docker image..."
	docker create --name wasm-extract $(FULL_IMAGE) /bin/sh || true
	docker cp wasm-extract:/plugin.wasm ./autonity-jwt-revocation.wasm || true
	docker rm wasm-extract || true
	@echo "Extracted Wasm module:"
	@ls -lh autonity-jwt-revocation.wasm
