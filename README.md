# Wallet Top-up System

A scalable wallet top-up system built with Go, using Clean Architecture principles, SQLC for database operations, and Redis for caching.

## Tech Stack

- **Go** - Programming language
- **Fiber** - Web framework
- **PostgreSQL** - Database
- **SQLC** - Type-safe SQL code generation
- **Redis** - Caching
- **Docker** - Containerization

## Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL 16+
- Redis 7+
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/hydr0g3nz/top_up_wallet.git
   cd top_up_wallet
   ```

2. **Install dependencies:**
   ```bash
   make install-deps
   ```

3. **Install tools:**
   ```bash
   make install-tools
   ```

4. **Create environment file:**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` with your configuration.

5. **Setup database:**
   ```bash
   make db-setup
   ```
   This will run migrations and generate SQLC code.

### Running the Application

#### With Docker Compose (Recommended)

```bash
make docker-run
```

#### Local Development

1. **Start PostgreSQL and Redis:**
   ```bash
   # Using Docker
   docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=pass postgres:16.2
   docker run -d --name redis -p 6379:6379 redis:7.2
   ```

2. **Run migrations:**
   ```bash
   ./scripts/migrate.sh
   ```

3. **Start the application:**
   ```bash
   make run
   ```

### API Endpoints

#### User Management

- `POST /api/v1/users/register` - Register a new user
- `POST /api/v1/users/login` - Login user
- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update current user profile
- `PUT /api/v1/users/me/password` - Change current user password

#### Admin Endpoints

- `GET /api/v1/users/` - List users by role
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `PUT /api/v1/users/:id/activate` - Activate user
- `PUT /api/v1/users/:id/deactivate` - Deactivate user

### Development

#### Generate SQLC Code

```bash
make sqlc-gen
```

#### Run Tests

```bash
make test
```

#### Hot Reload (with Air)

```bash
make dev
```

### Database Schema

The application uses PostgreSQL with the following main tables:

- `users` - User accounts with roles (candidate, company_hr, admin)

### Configuration

Environment variables:

```env
# Server
PORT=8080
SERVER_HOST=localhost
SERVER_READ_TIMEOUT=10
SERVER_WRITE_TIMEOUT=10

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=pass
DB_NAME=topup_wallet
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=pass
REDIS_DB=0

# Application
LOG_LEVEL=info
MAX_ACCEPTED_AMOUNT=100000.0
```

### Project Structure

```
.
├── cmd/
│   └── main.go                     # Application entry point
├── config/
│   └── config.go                   # Configuration management
├── internal/
│   ├── adapter/
│   │   ├── controller/             # HTTP handlers
│   │   ├── repository/             # Data access layer
│   │   │   └── sqlc/              # SQLC repositories
│   │   └── sqlc/                  # SQLC generated code
│   │       ├── generated/         # Generated SQLC code
│   │       ├── queries/           # SQL queries
│   │       └── schema/            # Database schema
│   ├── application/               # Use cases/business logic
│   ├── domain/                    # Domain entities and interfaces
│   │   ├── entity/               # Domain entities
│   │   ├── repository/           # Repository interfaces
│   │   ├── vo/                   # Value objects
│   │   └── infra/               # Infrastructure interfaces
│   └── infrastructure/           # External dependencies
├── scripts/
│   └── migrate.sh                # Database migration script
├── docker-compose.yaml
├── Dockerfile
├── sqlc.yaml                     # SQLC configuration
└── Makefile
```

### Architecture

This project follows Clean Architecture principles:

- **Domain Layer**: Contains business entities, value objects, and interfaces
- **Application Layer**: Contains use cases and business logic
- **Infrastructure Layer**: Contains external dependencies (database, cache, HTTP)
- **Adapter Layer**: Contains controllers, repositories, and data transformation

### Features

#### Completed

- User registration and authentication
- Role-based access control (candidate, company_hr, admin)
- User profile management
- Password management
- Email verification
- User activation/deactivation
- Clean Architecture implementation
- SQLC integration for type-safe database operations
- Redis caching
- Docker support
- Graceful shutdown

#### Planned

- JWT authentication
- Wallet top-up functionality
- Transaction management
- Rate limiting
- API documentation with Swagger
- Unit tests
- Integration tests

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

### License

This project is licensed under the MIT License.