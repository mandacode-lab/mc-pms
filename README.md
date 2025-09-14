# Serengeti

Multi-tenant OAuth authentication and service management platform built with hexagonal architecture.

## Architecture

This project follows **Hexagonal Architecture (Ports and Adapters)** pattern with clean separation of concerns:

- **Domain Layer**: Pure business logic and entities
- **Usecase Layer**: Application business rules and orchestration
- **Port Layer**: Interface contracts between layers
- **Adapter Layer**: External system integrations

### Key Features

- **Multi-tenant Architecture**: Service-based isolation with public ID system
- **OAuth 2.0 Support**: Multiple OAuth providers with extensible design
- **Security**: KEK/DEK encryption pattern for sensitive data protection
- **Database**: PostgreSQL with Ent ORM and Atlas migrations
- **Caching**: Redis-based caching layer

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL
- Redis

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
   # Update .env.dev.* files with your configuration
   ```

4. **Run services**
   ```bash
   # Auth Service
   go run cmd/auth/*.go

   # Core Service
   go run cmd/core/*.go
   ```

### Environment Configuration

Configure environment variables using `.env.dev.*` files:
- KEK (256-bit encryption key)
- Database connection settings
- Redis connection settings
- Server configuration

## API Overview

The platform provides two main services:
- **Auth Service**: OAuth authentication and user identity management
- **Core Service**: Service management, client apps, and OAuth configuration

See Swagger documentation at `/swagger/index.html` for detailed API specifications.

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

# Install development tools
make install-tools

# Generate Swagger documentation
make swagger-gen
```

### Project Structure

This project follows hexagonal architecture with clean separation of concerns:

```
internal/
├── domain/          # Domain entities and business logic
├── usecase/         # Application use cases and orchestration
├── port/            # Interface contracts
│   ├── in/          # Inbound ports (use case interfaces)
│   └── out/         # Outbound ports (adapter interfaces)
└── adapter/         # External integrations
    ├── handler/     # HTTP handlers
    ├── repository/  # Database repositories
    └── ...          # Other external adapters
```

### Key Patterns

- **Error Handling**: Structured error handling with consistent patterns
- **Transactions**: Consistent transaction management across use cases
- **Security**: KEK/DEK encryption for sensitive data protection
- **Validation**: Input validation at adapter layer

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
