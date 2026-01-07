#!/bin/bash
set -euo pipefail

# Color output for better visibility
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() { echo -e "${GREEN}ℹ️  $1${NC}"; }
log_warn() { echo -e "${YELLOW}⚠️  $1${NC}"; }
log_error() { echo -e "${RED}❌ $1${NC}"; }
log_debug() { echo -e "${BLUE}🔍 $1${NC}"; }

# Validate required environment variables
required_vars="USERNAME PASSWORD HOST PORT DB_NAME SSL_MODE"

log_info "Validating environment variables..."
for var in $required_vars; do
  if [ -z "${!var:-}" ]; then
    log_error "Required environment variable $var is not set"
    exit 1
  fi
done
log_info "✓ All required variables present"

# Construct DATABASE_URL
DATABASE_URL="postgresql://${USERNAME}:${PASSWORD}@${HOST}:${PORT}/${DB_NAME}?sslmode=${SSL_MODE}"
export DATABASE_URL

# Print configuration (mask password)
MASKED_URL="${DATABASE_URL/${PASSWORD}/****}"
log_info "Atlas Migration Tool v$(atlas version | head -1 | awk '{print $3}')"
log_info "Target: ${HOST}:${PORT}/${DB_NAME}"
log_info "SSL Mode: ${SSL_MODE}"
log_info "Connection: ${MASKED_URL}"

# Check Atlas CLI availability
if ! command -v atlas &> /dev/null; then
  log_error "Atlas CLI not found in PATH"
  exit 1
fi

# Test database connectivity
log_info "Testing database connectivity..."
if ! atlas schema inspect --url "$DATABASE_URL" > /dev/null 2>&1; then
  log_error "Failed to connect to database"
  log_error "Please verify database credentials and network connectivity"
  exit 1
fi
log_info "✓ Database connection successful"

# Check current migration status
log_info "Checking current migration status..."
if atlas migrate status --dir file://migrations --url "$DATABASE_URL"; then
  log_info "✓ Migration status check complete"
else
  log_warn "Unable to get migration status (may be first run)"
fi

# Dry-run mode (for testing)
if [ "${DRY_RUN:-false}" = "true" ]; then
  log_warn "=== DRY-RUN MODE ==="
  log_warn "No changes will be applied to the database"

  log_info "Simulating migration..."
  atlas migrate apply --dir file://migrations --url "$DATABASE_URL" --dry-run

  log_info "=== DRY-RUN COMPLETE ==="
  log_info "Review the output above to see what would be applied"
  exit 0
fi

# Apply migrations
log_info "Applying migrations..."
ARGS="--dir file://migrations --url $DATABASE_URL"

# Allow dirty state (use with caution)
if [ "${ALLOW_DIRTY:-false}" = "true" ]; then
  log_warn "ALLOW_DIRTY is enabled - applying migrations on dirty database state"
  ARGS="$ARGS --allow-dirty"
fi

# Execute migration with timeout
log_info "Executing: atlas migrate apply $ARGS"

if timeout 300s atlas migrate apply $ARGS; then
  log_info "✓ Migration apply successful"
else
  EXIT_CODE=$?
  log_error "Migration failed with exit code: $EXIT_CODE"

  # Print final status for debugging
  log_info "Final database status:"
  atlas migrate status --dir file://migrations --url "$DATABASE_URL" || true

  exit $EXIT_CODE
fi

# Verify final state
log_info "Verifying final migration status..."
atlas migrate status --dir file://migrations --url "$DATABASE_URL"

log_info "======================================="
log_info "✅ Migration completed successfully"
log_info "======================================="

exit 0
