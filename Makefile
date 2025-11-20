APP_NAME := fujin-grpc-gateway

ifeq ($(OS),Windows_NT)
	VERSION ?= v0.1.7
else
VERSION ?= $(shell cat VERSION 2>/dev/null | head -n 1 | tr -d '\r\n' || echo "")
endif

# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY_EXT := .exe
    RM := del /Q /F
    RMDIR := rmdir /S /Q
    MKDIR := mkdir
    PATHSEP := \\
else
    DETECTED_OS := $(shell uname -s)
    BINARY_EXT :=
    RM := rm -f
    RMDIR := rm -rf
    MKDIR := mkdir -p
    PATHSEP := /
endif

BIN_DIR := bin
BINARY := $(BIN_DIR)$(PATHSEP)$(APP_NAME)$(BINARY_EXT)

.PHONY: generate
generate:
	@echo "==> Downloading and generating gRPC Gateway code..."
# ifeq ($(OS),Windows_NT)
# 	@powershell -Command "New-Item -ItemType Directory -Force -Path 'internal/v1' | Out-Null; Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/fujin-io/fujin/refs/heads/${VERSION}/public/proto/grpc/v1/fujin-grpc-gateway.proto' -OutFile 'internal/v1/fujin-grpc-gateway.proto'"
# else
# 	@mkdir -p internal/v1 
# 	@curl -sSL "https://raw.githubusercontent.com/fujin-io/fujin/refs/heads/${VERSION}/public/proto/grpc/v1/fujin-grpc-gateway.proto" -o internal/v1/fujin-grpc-gateway.proto
# endif
	@cd internal/v1 && buf dep update && buf generate --template ../../buf.gen.yaml . && cd ../..

.PHONY: build
build:
	@echo "==> Building ${APP_NAME} for ${DETECTED_OS} (Version: ${VERSION})"
	@go build -ldflags "-s -w -X main.Version=${VERSION}" -o ./$(BINARY) ./cmd/...
	@echo "==> Binary created: $(BINARY)"

.PHONY: clean
clean:
	@echo "==> Cleaning"
ifeq ($(OS),Windows_NT)
	@if exist $(BIN_DIR) $(RMDIR) $(BIN_DIR) 2>nul
else
	@$(RMDIR) $(BIN_DIR)
endif