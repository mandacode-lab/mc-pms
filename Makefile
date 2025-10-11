# =============================================================================
# MandaCode Service Hub - Makefile
# =============================================================================

-include .env.make

.PHONY: install-tools swagger-gen swagger-clean help \
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
# TOOLS INSTALLATION
# =============================================================================

# Install required tools
install-tools:
	@echo "Installing development tools..."
	@echo "  - swag"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed!"

# =============================================================================
# SWAGGER DOCUMENTATION
# =============================================================================

# Swagger documentation
swagger-gen:
	@echo "Generating Swagger documentation..."
	@swag init --parseDependency --parseInternal \
		--generalInfo ./cmd/client/main.go \
		--output ./docs/client \
		--exclude ./internal/adapter/handler/service_mgmt,./internal/adapter/handler/client_mgmt
	@swag init --parseDependency --parseInternal \
		--generalInfo ./cmd/management/main.go \
		--output ./docs/management \
		--exclude ./internal/adapter/handler/client_access

swagger-clean:
	@echo "Cleaning Swagger documentation..."
	@rm -rf docs/
	@echo "Swagger docs cleaned!"

.PHONY: deploy
deploy:
	helm upgrade --install $(RELEASE_NAME) $(CHART_PATH) \
		-n $(NAMESPACE) \
		--create-namespace \
		-f $(VALUES_FILE)

# =============================================================================
# HELP & DOCUMENTATION
# =============================================================================

# Help
help:
	@echo "MandaCode Service Hub - Available Commands"
	@echo "="`printf '%.0s' {1..50}`
	@echo ""
	@echo "Database Migration Commands:"
	@echo "  generate-migrate MIGRATE_NAME=name - Generate new migration"
	@echo "  apply-migrate MIGRATE_DEV_URL=url  - Apply migrations to database"
	@echo "  build-migrate MIGRATE_TAG=tag      - Build Docker migration image"
	@echo "  push-migrate MIGRATE_TAG=tag       - Push Docker migration image"
	@echo ""
	@echo "Documentation Commands:"
	@echo "  swagger-gen    - Generate Swagger documentation"
	@echo "  swagger-clean  - Clean generated Swagger files"
	@echo ""
	@echo "Setup Commands:"
	@echo "  install-tools  - Install required development tools"
	@echo ""
