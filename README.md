# Go Gin Hexagonal Architecture Template

> A production-ready Go REST API template following Hexagonal Architecture principles with comprehensive testing, clean code structure, and SonarQube quality gate compliance.

## 🏗️ Template Overview

This repository serves as a **clean architecture template** for building robust Go REST APIs:

- **Main Branch**: Basic template with core features
- **Feature/redis-rabbitmq Branch**: Advanced template with Redis caching and RabbitMQ messaging (COMING SOON!)
- **Quality Assured**: Passes SonarQube quality gate with 0 issues
- **Production Ready**: Complete with CI/CD, Docker, and comprehensive testing

## 📁 Project Structure

```
go-gin-hexagonal/
├── cmd/
│   ├── api/                # API server entry point
│   │   └── main.go
│   └── migrate/            # Database migration entry point
│       └── main.go
├── internal/
│   ├── adapter/            # Adapter Layer (infrastructure & delivery)
│   │   ├── database/
│   │   │   ├── gorm/       # GORM implementation & repositories
│   │   │   │   ├── base_repository.go
│   │   │   │   ├── gorm.go
│   │   │   │   ├── refresh_token_repository.go
│   │   │   │   ├── user_repository.go
│   │   │   │   ├── json/         # Seed data (user.json)
│   │   │   │   ├── schema/       # GORM models/schema
│   │   │   │   └── seeder/       # Seeder logic
│   │   ├── http/                 # HTTP layer (REST API)
│   │   │   ├── handlers/         # HTTP handlers (auth, user)
│   │   │   ├── message/          # Response messages (error, success)
│   │   │   ├── middleware/       # Middleware (CORS, CSRF, etc)
│   │   │   ├── routes/           # Route definitions
│   │   │   └── response.go       # Standardized response
│   │   ├── mailer/               # Email service adapter
│   │   │   ├── smtp.go
│   │   │   └── template/         # Email HTML templates
│   │   └── security/             # Security adapters (AES, bcrypt, JWT)
│   ├── application/              # Application Layer (use cases, business logic)
│   │   ├── dto/                  # Data Transfer Objects
│   │   ├── mapper/               # DTO <-> Entity mappers
│   │   └── service/              # Application services (auth, user, email)
│   └── domain/                   # Domain Layer (core business logic)
│       ├── entity/               # Domain entities
│       ├── ports/                # Interfaces (repositories, services, security, mailer)
├── pkg/                          # Shared utilities & configuration
│   ├── config/                   # App config loader
│   ├── errors/                   # Custom error types
│   └── utils/                    # Utility functions
├── test/                         # Test layer
│   ├── mock/                     # Mock implementations (external, repository)
│   └── user_test.go              # Test suites
├── tmp/                          # Temporary build/output files
├── .github/
│   └── workflows/                # CI/CD workflows
├── .air.toml                     # Air hot reload config
├── .env.example                  # Env template
├── Dockerfile                    # Docker config
├── go.mod
├── go.sum
├── makefile                      # Build automation
├── README.md
└── sonar-project.properties      # SonarQube config
```

## ✨ Key Features

### 🏛️ Architecture

- **Hexagonal Architecture**: Clean separation of concerns with ports and adapters
- **Domain-Driven Design**: Business logic isolated in domain layer
- **Dependency Inversion**: All dependencies point inward to domain
- **SOLID Principles**: Following all SOLID design principles

### 🔐 Security & Authentication

- **JWT Authentication**: Access and refresh token mechanism
- **Password Hashing**: Bcrypt for secure password storage
- **AES Encryption**: Data encryption for sensitive information
- **CORS Configuration**: Cross-origin resource sharing setup

### 📧 Email System

- **SMTP Integration**: Direct SMTP email sending capability
- **HTML Templates**: Beautiful email templates with data injection
- **Template Engine**: Dynamic content generation
- **Async Processing**: Non-blocking email delivery

### 🗄️ Database

- **PostgreSQL**: Production-grade database with GORM ORM
- **Migrations**: Database schema versioning and migration system
- **Seeding**: Initial data setup with JSON seed files
- **Generic Repository**: Base repository pattern with common operations

### 🧪 Testing & Quality

- **Comprehensive Testing**: Unit tests with testify framework
- **Mock Framework**: Complete mock implementations for all external dependencies
- **Test Suites**: Organized test suites with setup and teardown
- **SonarQube Integration**: Code quality assurance with 0 issues
- **GitHub Actions**: Automated CI/CD pipeline

### 🔧 Development Experience

- **Hot Reload**: Air for automatic restart during development
- **Environment Configuration**: Flexible env-based configuration
- **Docker Support**: Ready for containerized deployment
- **Make Commands**: Simplified build and deployment commands

## 🌟 Template Branches

### Main Branch (Basic Template)

- Core hexagonal architecture implementation
- JWT authentication system
- Email service with SMTP
- PostgreSQL with GORM
- Comprehensive testing suite
- SonarQube quality gate compliance

### Feature/redis-rabbitmq Branch (Advanced Template - COMING SOON!)

- Everything from main branch
- Redis caching layer
- RabbitMQ message queuing
- Advanced async processing
- Distributed system patterns

## 📊 Quality Assurance

This template has been rigorously tested and passes all quality gates:

- ✅ **SonarQube Quality Gate**: 0 issues
- ✅ **Code Coverage**: Comprehensive test coverage
- ✅ **Security Scan**: No security vulnerabilities
- ✅ **Code Smells**: Clean, maintainable code
- ✅ **Duplication**: No code duplication
- ✅ **Maintainability**: High maintainability rating

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 12+
- SMTP server (Gmail, SendGrid, etc.)
- Docker (optional)

### 1. Clone the Template

```bash
# Clone the repository
git clone <repository-url>
cd go-gin-hexagonal

# For advanced features (Redis + RabbitMQ)
git checkout feature/redis-rabbitmq
```

### 2. Environment Setup

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your configuration
```

### 3. Environment Configuration

Update `.env` with your settings:

```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=3000
APP_ENV=development
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_here
JWT_ACCESS_EXPIRY=1h
JWT_REFRESH_EXPIRY=168h

# AES Encryption
AES_KEY=your_32_character_aes_key_here
AES_IV=your_16_character_initialization_vector

# Email Configuration
MAILER_HOST=smtp.gmail.com
MAILER_PORT=587
MAILER_SENDER=Your App <no-reply@yourapp.com>
MAILER_AUTH=your_email@gmail.com
MAILER_PASSWORD=your_app_password
```

### 4. Database Setup

```bash
# Install dependencies
go mod download

# Run migrations only
go run cmd/migrate/main.go --migrate

# Run migrations with seed data
go run cmd/migrate/main.go --seed

# Fresh migration (drop and recreate)
go run cmd/migrate/main.go --fresh
```

### 5. Run the Application

**Development (with hot reload):**

```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest

# Start development server
air
```

**Production:**

```bash
# Build and run
go build -o api cmd/api/main.go
./api
```

**Docker:**

```bash
# Build Docker image
docker build -t go-gin-hexagonal .

# Run with Docker
docker run -p 3000:3000 go-gin-hexagonal
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific test suite
go test -v ./test/
```

### Test Structure

- **Unit Tests**: Complete test coverage for all services
- **Mock Testing**: All external dependencies mocked
- **Integration Tests**: End-to-end testing scenarios
- **Test Suites**: Organized with testify suite framework

## 🏗️ Architecture Deep Dive

### Hexagonal Architecture Mapping

1. **Domain Layer** (`internal/domain/`)

   - Pure business logic, entities, and interfaces (ports)
   - No external dependencies

2. **Application Layer** (`internal/application/`)

   - Use case orchestration, service logic, DTOs, and mappers
   - Manages business workflow and communication between layers

3. **Adapter Layer** (`internal/adapter/`)

   - Implementation of ports (repositories, mailer, security, etc.)
   - External integrations: database, HTTP, email, security
   - Delivery (HTTP handlers, middleware, response)

4. **Infrastructure/Utility Layer** (`pkg/`)

   - Application configuration, error handling, general utilities
   - Does not contain business logic

5. **Test Layer** (`test/`)
   - Mocks, test suites, and comprehensive testing

## 🔒 Security Features

- **JWT Authentication**: Stateless authentication with access/refresh tokens
- **Password Security**: Bcrypt hashing with proper salt rounds
- **Data Encryption**: AES encryption for sensitive data
- **CORS Protection**: Configurable cross-origin policies
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: GORM ORM protection

### Code Organization

- **Domain-Driven Design**: Business logic separated from infrastructure
- **Dependency Injection**: Loose coupling between components
- **Interface Segregation**: Small, focused interfaces
- **Single Responsibility**: Each component has one clear purpose
