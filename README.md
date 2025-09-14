# Serengeti

Multi-tenant OAuth authentication and service management platform built with hexagonal architecture.

## Architecture

This project follows **Hexagonal Architecture (Ports and Adapters)** pattern with clean separation of concerns:

- **Domain Layer** (`internal/domain/`): Pure business logic and entities
- **Usecase Layer** (`internal/usecase/`): Application business rules and orchestration
- **Port Layer** (`internal/port/`): Interface contracts between layers
- **Adapter Layer** (`internal/adapter/`): External system integrations

### Core Services

- **Auth Service** (port 8080): OAuth authentication and user identity management
- **Core Service** (port 8081): Service management, client apps, and WebOAuth configuration

### Key Features

- **Multi-tenant Architecture**: Service-based isolation with public ID system
- **OAuth 2.0 Support**: Google, Kakao, Naver providers with extensible design
- **Security**: KEK/DEK encryption pattern for sensitive data protection
- **Caching**: Redis-based caching for OAuth configurations
- **Database**: PostgreSQL with Ent ORM and Atlas migrations

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15
- Redis 7

### Development Setup

1. **Clone and setup**
   ```bash
   git clone <repository-url>
   cd serengeti
   go mod download
   ```

2. **Start infrastructure**
   ```bash
   docker compose -f docker-compose.dev.yml up -d
   ```

3. **Configure environment**
   ```bash
   # Generate KEK (256-bit hex key)
   openssl rand -hex 32

   # Update .env.dev.core and .env.dev.auth with your KEK
   ```

4. **Run services**
   ```bash
   # Core Service (port 8081)
   go run cmd/core/*.go

   # Auth Service (port 8080)
   go run cmd/auth/*.go
   ```

### Environment Configuration

Set the following environment variables or use `.env.dev.*` files:

```bash
# Required
ENV=dev
KEK=<64-character-hex-string>  # 256-bit encryption key

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=serengeti
POSTGRES_PASSWORD=serengeti123
POSTGRES_DB=serengeti_dev
POSTGRES_SSLMODE=disable

# Cache
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0

# Server
SERVER_HOST=localhost
SERVER_PORT=8080  # or 8081 for core
```

## API Endpoints

### Auth Service (port 8080)
- `GET /health` - Health check
- `GET /v1/auth/auth-url` - Get OAuth authorization URL
- `GET /v1/auth/code` - OAuth code flow authentication
- `GET /v1/auth/token` - OAuth token flow authentication

### Core Service (port 8081)
- `GET /health` - Health check
- `POST /v1/services/` - Create service
- `PUT /v1/services/:id` - Update service
- `POST /v1/client-apps/` - Create client application
- `PUT /v1/client-apps/:id` - Update client application
- `POST /v1/client-apps/:id/refresh-secret` - Refresh client secret
- `GET /v1/client-apps/` - List client applications
- `POST /v1/weboauth/` - Register OAuth configuration
- `GET /v1/weboauth/` - Get OAuth configuration
- `GET /v1/users/` - Find user information

## Development

### Build Commands

```bash
# Build all packages
go build ./...

# Run tests
go test ./...

# Generate database migrations
make generate-migrate MIGRATE_NAME=your_migration_name

# Apply database migrations
make apply-migrate MIGRATE_DEV_URL=postgres://user:pass@localhost/db

# Generate protobuf files
make proto-gen

# Generate Swagger documentation
make swagger-gen
```

### Project Structure

```
internal/
├── domain/          # Domain entities and business logic
│   ├── clientapp/   # OAuth client applications
│   ├── service/     # Multi-tenant services
│   ├── useridentity/# OAuth user identities
│   ├── userinfo/    # User profile information
│   └── weboauth/    # OAuth provider configurations
├── usecase/         # Application use cases
│   ├── clientapp_mgmt/  # Client app CRUD operations
│   ├── identify_user/   # OAuth authentication flows
│   ├── service_mgmt/    # Service management
│   ├── user_mgmt/       # User management
│   └── weboauth_mgmt/   # OAuth configuration management
├── port/            # Interface contracts
│   ├── in/          # Inbound ports (use case interfaces)
│   └── out/         # Outbound ports (adapter interfaces)
└── adapter/         # External integrations
    ├── handler/     # HTTP handlers
    ├── repository/  # Database repositories
    ├── oauth/       # OAuth provider clients
    ├── kek/         # Encryption services
    └── cache/       # Caching services
```

### Key Patterns

- **Error Handling**: Structured errors with `merr` library
- **Transactions**: Consistent transaction management across use cases
- **Security**: KEK/DEK encryption for sensitive data
- **Validation**: Input validation at adapter layer
- **Caching**: Selective caching for OAuth configurations

## Security

- **Encryption**: AES-256-GCM with KEK/DEK pattern
- **Authentication**: OAuth 2.0 with multiple providers
- **Secrets**: Environment-based KEK injection
- **Transport**: HTTPS recommended for production

## Contributing

1. Follow hexagonal architecture principles
2. Maintain clean separation between layers
3. Use existing patterns for consistency
4. Add tests for new functionality
5. Update documentation as needed

## License

[License information]
