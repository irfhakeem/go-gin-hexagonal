# Go Gin Hexagonal Architecture

## Project Structure

```
go-gin-hexagonal/
├── cmd/
│   ├── api/
│   │   └── main.go              # API server entry point
│   └── migrate/
│       └── main.go              # Database migration entry point
├── internal/
│   ├── adapter/
│   │   ├── database/
│   │   │   ├── model/
│   │   │   │   ├── refresh_token.go
│   │   │   │   └── user.go
│   │   │   ├── refresh_token_repository.go
│   │   │   └── user_repository.go
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   │   ├── auth_handler.go
│   │   │   │   └── user_handler.go
│   │   │   ├── message/
│   │   │   │   ├── error.go
│   │   │   │   └── success.go
│   │   │   ├── middleware/
│   │   │   │   ├── cors.go
│   │   │   │   └── middleware.go
│   │   │   ├── routes/
│   │   │   │   └── router.go
│   │   │   └── response.go
│   │   ├── mailer/
│   │   │   ├── smtp.go          # SMTP mailer implementation
│   │   │   └── template/
│   │   │       └── new_user.html # Email template
│   │   └── security/
│   │       ├── bcrypt.go        # Password hashing adapter
│   │       └── jwt.go           # JWT token adapter
│   ├── application/
│   │   └── service/
│   │       ├── auth_service.go
│   │       ├── email_service.go
│   │       └── user_service.go
│   └── domain/
│       ├── dto/
│       │   ├── auth_dto.go
│       │   ├── email_dto.go
│       │   ├── token_dto.go
│       │   └── user_dto.go
│       ├── entity/
│       │   ├── refresh_token.go
│       │   └── user.go
│       └── ports/
│           ├── mailer.go
│           ├── repositories.go
│           ├── security.go
│           └── services.go
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── json/
│   │   │   └── user.json        # Seed data
│   │   ├── seeder/
│   │   │   └── user_seeder.go
│   │   └── postgres.go
│   └── utils/
│       ├── number_utils.go
│       └── string_utils.go
├── test/
│   ├── mailer_test.go
│   ├── security_test.go
│   └── user_test.go
├── tmp/
│   ├── build-errors.log
│   └── main.exe
├── .env.example
├── .air.toml                    # Air hot reload config
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Features

- **Clean Architecture**: Hexagonal Architecture pattern with clear separation of concerns
- **Authentication & Authorization**: JWT-based authentication with refresh tokens
- **Email System**: SMTP-based email service with HTML templates
- **Database**: PostgreSQL with GORM ORM and database migrations
- **Testing**: Comprehensive unit and integration tests
- **Configuration**: Environment-based configuration management
- **Hot Reload**: Development with Air for automatic restart
- **Docker Support**: Containerized deployment ready

## Architecture Overview

This project follows the **Hexagonal Architecture** (Ports and Adapters) pattern:

- **Domain Layer**: Contains business entities, DTOs, and ports (interfaces)
- **Application Layer**: Contains business logic and use cases (services)
- **Adapter Layer**: Contains external integrations (database, HTTP, mailer, security)
- **PKG Layer**: Contains configuration, utilities, and external dependencies

## Email System

The project includes a comprehensive email system with:

- **SMTP Mailer**: Direct SMTP email sending
- **Template Engine**: HTML email templates with data injection
- **Email Service**: High-level email operations (welcome, password reset)
- **Async Sending**: Non-blocking email delivery

### Email Templates

- `new_user.html`: Welcome email for new user registration

## How to Run

### Prerequisites

1. Go 1.20 or higher
2. PostgreSQL database
3. SMTP server configuration (Gmail, SendGrid, etc.)

### Environment Setup

1. Copy environment file:

```bash
cp .env.example .env
```

2. Update `.env` with your configuration:

```env
# Server
SERVER_HOST=localhost
SERVER_PORT=3000
APP_ENV=development
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name

# JWT
JWT_SECRET=your_jwt_secret
JWT_ACCESS_EXPIRY=1h
JWT_REFRESH_EXPIRY=168h

# Mailer
MAILER_HOST=smtp.gmail.com
MAILER_PORT=587
MAILER_SENDER=Your App <no-reply@yourapp.com>
MAILER_AUTH=your_email@gmail.com
MAILER_PASSWORD=your_app_password
```

### 1. Run Database Migration

Migrate only

```bash
go run cmd/migrate/main.go
```

Migrate and seeds

```bash
go run cmd/migrate/main.go --seed
```

Fresh Migrate

```bash
go run cmd/migrate/main.go --fresh
```

### 2. Run API Server

Using Air (for hot reload - recommended for development):

```bash
air
```

Using Go run:

```bash
go run cmd/api/main.go
```

## Dependencies

### Core Dependencies

- **Gin**: HTTP web framework
- **GORM**: ORM library for database operations
- **JWT-Go**: JWT token handling
- **Bcrypt**: Password hashing
- **UUID**: Unique identifier generation
- **Godotenv**: Environment variable loading
- **Gomail**: Email sending library

### Development Dependencies

- **Testify**: Testing framework with assertions and mocks
- **Air**: Hot reload for development

### Database

- **PostgreSQL**: Primary database
- **GORM PostgreSQL Driver**: Database driver

## Development

### Project Setup

1. Clone the repository:

```bash
git clone <repository-url>
cd go-gin-hexagonal
```

2. Install dependencies:

```bash
go mod download
```

3. Install Air for hot reload:

```bash
go install github.com/cosmtrek/air@latest
```

4. Setup environment:

```bash
cp .env.example .env
# Edit .env with your configuration
```

5. Run migrations:

```bash
go run cmd/migrate/main.go
```

6. Start development server:

```bash
air
```

### Code Organization

- **Domain-Driven Design**: Business logic separated from infrastructure
- **Dependency Injection**: Loose coupling between components
- **Interface Segregation**: Small, focused interfaces
- **Single Responsibility**: Each component has one clear purpose
