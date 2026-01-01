# =============================================================================
# MandaCode PMS - Makefile
# =============================================================================

-include .make.env

# Migration directory
MIGRATIONS_DIR := ent/migrate/migrations

# Ent schema path
SCHEMA_PATH := ent/schema

# Dev database URL for Atlas
DEV_URL := docker://postgres/17

# Ent directory (where generate.go is located)
ENT_DIR := ent

# =============================================================================
# Ent Code Generation Commands
# =============================================================================

.PHONY: generate-ent
generate-ent:
	@echo "🔧 Generating ent code..."
	@cd $(ENT_DIR) && go generate .
	@echo "✅ Ent code generated successfully!"

# =============================================================================
# Migration Commands
# =============================================================================

.PHONY: migrate
migrate:
ifndef MIGRATE_NAME
	$(error MIGRATE_NAME is required. Usage: make migrate MIGRATE_NAME=your_migration_name)
endif
	@echo "📦 Generating migration: $(MIGRATE_NAME)"
	@atlas migrate diff $(MIGRATE_NAME) \
		--dev-url "$(DEV_URL)" \
		--dir "file://$(MIGRATIONS_DIR)" \
		--to "ent://$(SCHEMA_PATH)"
	@echo "✅ Migration generated: $(MIGRATIONS_DIR)/$(MIGRATE_NAME).sql"

# =============================================================================
# Swagger Documentation Commands
# =============================================================================

.PHONY: swagger-gen
swagger-gen:
	@echo "📚 Generating Swagger documentation..."
	@echo "  - Admin API..."
	@swag init \
		--parseDependency \
		--parseInternal \
		--generalInfo cmd/admin/main.go \
		--tags "admin" \
		--output docs/admin
	@echo "  - Client API..."
	@swag init \
		--parseDependency \
		--parseInternal \
		--generalInfo cmd/client/main.go \
		--tags "client" \
		--output docs/client
	@echo "✅ Swagger documentation generated successfully!"

.PHONY: swagger-clean
swagger-clean:
	@echo "🧹 Cleaning Swagger documentation..."
	@rm -rf docs/
	@echo "✅ Swagger docs cleaned"

# =============================================================================
# Help
# =============================================================================

.PHONY: help
help:
	@echo "MandaCode PMS - Makefile Commands"
	@echo ""
	@echo "Ent Code Generation Commands:"
	@echo "  make generate-ent                       - Generate ent code"
	@echo ""
	@echo "Migration Commands:"
	@echo "  make migrate MIGRATE_NAME=name          - Generate migration"
	@echo ""
	@echo "Swagger Documentation Commands:"
	@echo "  make swagger-gen                       - Generate Swagger documentation for all services"
	@echo "  make swagger-clean                     - Clean generated Swagger files"
	@echo ""
	@echo "Examples:"
	@echo "  make generate-ent"
	@echo "  make migrate MIGRATE_NAME=add_user_column"
	@echo "  make swagger-gen"
	@echo ""
	@echo "Using .make.env file:"
	@echo "  Create .make.env with: MIGRATE_NAME=your_migration_name"
	@echo "  Then run: make migrate"
	@echo ""
