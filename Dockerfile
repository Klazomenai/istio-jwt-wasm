# Multi-stage Dockerfile for building Autonity JWT Revocation WasmPlugin
# Stage 1: Build the Wasm module with Go 1.24+

FROM golang:1.24-alpine AS builder

# Install git (required for go mod download)
RUN apk add --no-cache git

WORKDIR /build

# Copy all source code
COPY . .

# Download dependencies and generate go.sum
RUN go mod download && go mod tidy

# Build the Wasm module
# GOOS=wasip1 targets WebAssembly System Interface (WASI) preview 1
# GOARCH=wasm targets WebAssembly architecture
# -buildmode=c-shared creates a shared library (required for proxy-wasm)
RUN GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o autonity-jwt-revocation.wasm main.go

# Verify the build (file command not available in alpine, just check size)
RUN ls -lh autonity-jwt-revocation.wasm

# Stage 2: Package the Wasm module for OCI distribution
# Istio can fetch WasmPlugins from OCI registries

FROM scratch

# Copy the Wasm module from builder
COPY --from=builder /build/autonity-jwt-revocation.wasm /plugin.wasm

# Metadata labels
LABEL description="JWT Revocation WasmPlugin for Istio Gateway"
LABEL vendor="Klazomenai"
LABEL version="0.0.1-alpha"
