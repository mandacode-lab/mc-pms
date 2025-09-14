# =============================================================================
# Serengeti Integrated - Makefile
# =============================================================================

-include .make.env

.PHONY: proto-gen proto-clean proto-deps install-tools swagger-gen swagger-clean help \
        generate-migrate apply-migrate build-migrate push-migrate check-migrate-name check-migrate-tag

# =============================================================================
# CONFIGURATION
# =============================================================================

# Docker settings
DOCKER_CONTEXT ?= 
MIGRATE_IMAGE ?= 
MIGRATE_TAG ?=

# Database migration settings
MIGRATE_NAME ?=
MIGRATE_DEV_URL ?=

# Directories
PROTO_DIR = pkg
PROTO_OUT_DIR = pkg
THIRD_PARTY_DIR = third_party

# Proto dependency versions
PROTOC_GEN_VALIDATE_VERSION = main
PROTOBUF_VERSION = main

# Proto dependency URLs
VALIDATE_PROTO_URL = https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/$(PROTOC_GEN_VALIDATE_VERSION)/validate/validate.proto
TIMESTAMP_PROTO_URL = https://raw.githubusercontent.com/protocolbuffers/protobuf/$(PROTOBUF_VERSION)/src/google/protobuf/timestamp.proto
EMPTY_PROTO_URL = https://raw.githubusercontent.com/protocolbuffers/protobuf/$(PROTOBUF_VERSION)/src/google/protobuf/empty.proto

# Go tool versions
PROTOC_GEN_GO_VERSION = latest
PROTOC_GEN_GO_GRPC_VERSION = latest
PROTOC_GEN_VALIDATE_VERSION_GO = latest

# =============================================================================
# DATABASE MIGRATIONS
# =============================================================================
# Generate new database migration
# Usage: make generate-migrate MIGRATE_NAME=your_migration_name
generate-migrate: check-migrate-name
	@echo "Generating migration: $(MIGRATE_NAME)"
	atlas migrate diff $(MIGRATE_NAME) \
		--dev-url "docker://postgres/17" \
		--dir "file://ent/migrate/migrations" \
		--to "ent://ent/schema"
	@echo "Migration generated successfully"

# Apply migrations to database
# Usage: make apply-migrate MIGRATE_DEV_URL=your_db_url
apply-migrate:
	@echo "Applying migrations to database"
	atlas migrate apply \
		--url $(MIGRATE_DEV_URL) \
		--dir "file://ent/migrate/migrations"
	@echo "Migrations applied successfully"

# Build Docker migration image
# Usage: make build-migrate MIGRATE_TAG=your_tag
build-migrate: check-migrate-tag
	@echo "Building migration image: $(MIGRATE_IMAGE):$(MIGRATE_TAG)"
	docker build -t $(MIGRATE_IMAGE):$(MIGRATE_TAG) -f $(DOCKER_CONTEXT)/Dockerfile.migration .
	@echo "Migration image built successfully"

# Push Docker migration image
# Usage: make push-migrate MIGRATE_TAG=your_tag
push-migrate: build-migrate
	@echo "Pushing migration image: $(MIGRATE_IMAGE):$(MIGRATE_TAG)"
	docker push $(MIGRATE_IMAGE):$(MIGRATE_TAG)
	@echo "Migration image pushed successfully"

# Internal: Check if MIGRATE_NAME is provided
check-migrate-name:
ifndef MIGRATE_NAME
	$(error MIGRATE_NAME is required. Usage: make generate-migrate MIGRATE_NAME=your_migration_name)
endif

# Internal: Check if MIGRATE_TAG is provided
check-migrate-tag:
ifndef MIGRATE_TAG
	$(error MIGRATE_TAG is required. Usage: make build-migrate MIGRATE_TAG=your_tag)
endif

# =============================================================================
# PROTOBUF GENERATION
# =============================================================================
proto-gen: proto-deps
	@echo "Generating protobuf files..."
	@find $(PROTO_DIR) -name "*.proto" | xargs protoc \
		--proto_path=$(PROTO_DIR) \
		--proto_path=$(THIRD_PARTY_DIR) \
		--go_out=$(PROTO_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		--validate_out="lang=go:$(PROTO_OUT_DIR)" \
		--validate_opt=paths=source_relative
	@echo "Proto generation completed!"

# Download proto dependencies
proto-deps:
	@echo "Downloading proto dependencies..."
	@mkdir -p $(THIRD_PARTY_DIR)/validate
	@mkdir -p $(THIRD_PARTY_DIR)/google/protobuf
	@echo "  - validate.proto ($(PROTOC_GEN_VALIDATE_VERSION))"
	@curl -sSL $(VALIDATE_PROTO_URL) -o $(THIRD_PARTY_DIR)/validate/validate.proto
	@echo "  - timestamp.proto ($(PROTOBUF_VERSION))"
	@curl -sSL $(TIMESTAMP_PROTO_URL) -o $(THIRD_PARTY_DIR)/google/protobuf/timestamp.proto
	@echo "  - empty.proto ($(PROTOBUF_VERSION))"
	@curl -sSL $(EMPTY_PROTO_URL) -o $(THIRD_PARTY_DIR)/google/protobuf/empty.proto
	@echo "Proto dependencies downloaded!"

# Clean generated files
proto-clean:
	@echo "Cleaning generated proto files..."
	@find $(PROTO_DIR) -name "*.pb.go" -delete
	@find $(PROTO_DIR) -name "*_grpc.pb.go" -delete
	@find $(PROTO_DIR) -name "*.validate.go" -delete
	@rm -rf $(THIRD_PARTY_DIR)
	@echo "Cleanup completed!"

# =============================================================================
# TOOLS INSTALLATION
# =============================================================================

# Install required tools
install-tools:
	@echo "Installing protoc plugins..."
	@echo "  - protoc-gen-go@$(PROTOC_GEN_GO_VERSION)"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	@echo "  - protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)"
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	@echo "  - protoc-gen-validate@$(PROTOC_GEN_VALIDATE_VERSION_GO)"
	@go install github.com/envoyproxy/protoc-gen-validate@$(PROTOC_GEN_VALIDATE_VERSION_GO)
	@echo "  - swag"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed!"

# =============================================================================
# SWAGGER DOCUMENTATION
# =============================================================================

# Swagger documentation
swagger-gen:
	@echo "Generating Swagger documentation..."
	@echo "  - Core service documentation"
	@swag init --parseDependency --parseInternal --generalInfo ./cmd/core/main.go --output ./docs/core
	@echo "  - Auth service documentation"
	@swag init --parseDependency --parseInternal --generalInfo ./cmd/auth/main.go --output ./docs/auth
	@echo "Swagger docs generated:"
	@echo "  - Core: /swagger/index.html (port 8081)"
	@echo "  - Auth: /swagger/index.html (port 8080)"

swagger-clean:
	@echo "Cleaning Swagger documentation..."
	@rm -rf docs/
	@echo "Swagger docs cleaned!"

# =============================================================================
# HELP & DOCUMENTATION
# =============================================================================

# Help
help:
	@echo "Serengeti Integrated - Available Commands"
	@echo "="`printf '%.0s' {1..50}`
	@echo ""
	@echo "Database Migration Commands:"
	@echo "  generate-migrate MIGRATE_NAME=name - Generate new migration"
	@echo "  apply-migrate MIGRATE_DEV_URL=url  - Apply migrations to database"
	@echo "  build-migrate MIGRATE_TAG=tag      - Build Docker migration image"
	@echo "  push-migrate MIGRATE_TAG=tag       - Push Docker migration image"
	@echo ""
	@echo "Proto Commands:"
	@echo "  proto-gen      - Generate Go code from proto files"
	@echo "  proto-deps     - Download proto dependencies"
	@echo "  proto-clean    - Clean generated proto files"
	@echo ""
	@echo "Swagger Commands:"
	@echo "  swagger-gen    - Generate Swagger documentation"
	@echo "  swagger-clean  - Clean generated Swagger files"
	@echo ""
	@echo "Setup Commands:"
	@echo "  install-tools  - Install required protoc plugins and tools"
	@echo ""
	@echo "Configuration:"
	@echo "  PROTO_DIR       = $(PROTO_DIR)"
	@echo "  PROTO_OUT_DIR   = $(PROTO_OUT_DIR)"
	@echo "  THIRD_PARTY_DIR = $(THIRD_PARTY_DIR)"
	@echo ""
	@echo "For more details on each command, check the Makefile comments."