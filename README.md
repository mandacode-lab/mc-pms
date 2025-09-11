# Serengeti Integrated

A microservices architecture implementing authentication, client management, and service management using hexagonal architecture patterns.

## Services

### Auth Service (Port: 8080)
OAuth-based authentication service supporting Google, Kakao, and Naver providers.

**Endpoints:**
- `POST /v1/auth/login/code` - Login with OAuth code
- `POST /v1/auth/login/access` - Login with OAuth access token  
- `POST /v1/auth/verify` - Verify login code
- `GET /v1/auth/url/{provider}` - Get OAuth authorization URL

### Client Service (Port: 9090 - gRPC)
Client credential management service.

**gRPC Services:**
- `CreateClient` - Create new service client
- `GetClient` - Retrieve client information
- `VerifyClient` - Verify client credentials
- `ListClients` - List clients for service
- `UpdateClient` - Update client information
- `ActivateClient` - Activate client
- `DeactivateClient` - Deactivate client

### Core Service (Port: 8081)
Service management and administration.

**Endpoints:**
- `POST /v1/core/services` - Create new service
- `GET /v1/core/services` - List services
- `GET /v1/core/services/{id}` - Get service details
- `PUT /v1/core/services/{id}` - Update service
- `POST /v1/core/services/{id}/activate` - Activate service
- `POST /v1/core/services/{id}/deactivate` - Deactivate service

## Architecture

### Hexagonal Architecture
- **Domain Layer**: Core business logic and entities
- **Application Layer**: Use cases and business rules
- **Adapter Layer**: External interfaces (HTTP, gRPC, Database, Cache)
- **Port Layer**: Interfaces defining contracts between layers

### Project Structure
```
serengeti-integrated/
├── cmd/                    # Application entry points
│   ├── auth/              # Auth service
│   ├── client/            # Client service  
│   ├── core/              # Core service
│   └── shared/            # Shared server utilities
├── configs/               # Configuration management
├── internal/
│   ├── adapter/           # External adapters
│   │   ├── cipher/        # Encryption adapters
│   │   ├── code_store/    # Code storage adapters
│   │   ├── hasher/        # Password hashing adapters
│   │   ├── handler/       # HTTP/gRPC handlers
│   │   ├── oauth/         # OAuth provider adapters
│   │   ├── random/        # Random generation adapters
│   │   └── repository/    # Database adapters
│   ├── domain/            # Domain entities and business logic
│   │   ├── auth/          # Authentication domain
│   │   ├── client/        # Client management domain
│   │   ├── core/          # Service management domain
│   │   └── shared/        # Shared domain concepts
│   ├── port/              # Port interfaces
│   │   ├── in/            # Inbound ports (use cases)
│   │   └── out/           # Outbound ports (adapters)
│   └── usecase/           # Use case implementations
├── pkg/                   # Shared utilities
└── third_party/           # External proto dependencies
```

## Technologies

- **Language**: Go 1.25
- **Framework**: Gin (HTTP), gRPC
- **Database**: PostgreSQL with Ent ORM
- **Cache**: Redis
- **Configuration**: Environment variables with validation
- **Logging**: Zerolog
- **Authentication**: OAuth 2.0 (Google, Kakao, Naver)
- **Encryption**: AES-256-GCM
- **Containerization**: Docker
- **Orchestration**: Kubernetes with Helm

## Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL 15+
- Redis 7+
- Docker (optional)

### Environment Variables

**Common:**
```bash
ENV=dev
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=serengeti
REDIS_ADDR=localhost:6379
```

**Auth Service:**
```bash
HTTP_PORT=8080
KEK=your-32-byte-encryption-key-here
CLIENT_SERVICE_ENDPOINT=localhost:9090
```

**Client Service:**
```bash
GRPC_PORT=9090
BCRYPT_COST=12
```

**Core Service:**
```bash
HTTP_PORT=8081
CLIENT_SERVICE_ENDPOINT=localhost:9090
```

### Running Services

1. **Install dependencies:**
```bash
go mod download
```

2. **Run database migrations:**
```bash
make apply-migrate MIGRATE_DEV_URL=postgres://user:pass@localhost/db
```

3. **Start services:**

```bash
# Auth Service
go run cmd/auth/*.go

# Client Service  
go run cmd/client/*.go

# Core Service
go run cmd/core/*.go
```

### Development

**Generate migrations:**
```bash
make generate-migrate MIGRATE_NAME=your_migration_name
```

**Install development tools:**
```bash
make install-tools
```

**Generate protobuf files:**
```bash
make proto-gen
```

## License

This project is licensed under the MIT License.