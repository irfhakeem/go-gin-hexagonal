# Go Gin Hexagonal Architecture

A complete Go web application implementing **Hexagonal Architecture** (Ports and Adapters) with proper separation of concerns, dependency inversion, and framework-independent core business logic.

## 🏗️ Architecture Overview

This project follows **Hexagonal Architecture** (also known as Ports and Adapters pattern) with strict adherence to the dependency rule. The domain is at the center, completely isolated from external dependencies:

```
┌──────────────────────────────────────────────────────────────┐
│                   PRIMARY ADAPTERS (Driving)                 │
│              HTTP Handlers, Routes, Middleware               │
└────────────────────────────┬─────────────────────────────────┘
                             │
                             ▼
┌──────────────────────────────────────────────────────────────┐
│                    PRIMARY PORTS (Driving)                   │
│          Interfaces: UserUseCase, EmailUseCase               │
└────────────────────────────┬─────────────────────────────────┘
                             │
                             ▼
┌──────────────────────────────────────────────────────────────┐
│                      DOMAIN (Core)                           │
│   ┌────────────────────────────────────────────────────┐    │
│   │  Domain Services (UserService, EmailService)       │    │
│   │  ↓ uses                                            │    │
│   │  Domain Models (User, RefreshToken, Gender)        │    │
│   │  ↓ may throw                                       │    │
│   │  Domain Errors (ErrUserNotFound, etc.)            │    │
│   └────────────────────────────────────────────────────┘    │
└────────────────────────────┬─────────────────────────────────┘
                             │
                             ▼
┌──────────────────────────────────────────────────────────────┐
│                   SECONDARY PORTS (Driven)                   │
│     Interfaces: Repositories, JWT, Crypto, SMTP, Storage     │
└────────────────────────────┬─────────────────────────────────┘
                             │
                             ▼
┌──────────────────────────────────────────────────────────────┐
│                  SECONDARY ADAPTERS (Driven)                 │
│        Postgres, Crypto Services, SMTP, LocalStorage         │
└──────────────────────────────────────────────────────────────┘
```

**Key Architectural Benefits:**
- **Pure Domain**: Business logic has ZERO dependencies on frameworks, databases, or external libraries
- **Dependency Inversion**: All dependencies point inward toward the domain
- **Port-Driven Design**: Domain defines interfaces (ports), adapters implement them
- **Testability**: Domain services can be tested with mock ports - no framework mocking required
- **Flexibility**: Swap any adapter (HTTP→gRPC, Postgres→MongoDB) without changing domain
- **Maintainability**: Clear boundaries, single responsibility, and explicit dependencies

## 📁 Project Structure

```
go-gin-clean/
├── cmd/                          # Application entrypoints
│   ├── server/main.go           # HTTP server
│   └── migrate/main.go          # Database migrations
├── internal/                    # Private application code
│   ├── domain/                  # Domain Layer (Hexagon Core)
│   │   ├── model/              # Domain entities
│   │   │   ├── user.go         # User entity
│   │   │   ├── refresh_token.go # RefreshToken entity
│   │   │   ├── audit.go        # Audit fields
│   │   │   └── gender.go       # Gender enum
│   │   ├── service/            # Domain services (business logic)
│   │   │   ├── user_service.go  # User business logic
│   │   │   └── email_service.go # Email business logic
│   │   └── error/              # Domain errors
│   │       └── errors.go       # Business errors
│   ├── ports/                  # Port Interfaces (Hexagon Boundary)
│   │   ├── primary/            # Driving Ports (what app offers)
│   │   │   ├── dto.go          # Shared DTOs for use cases
│   │   │   └── user_port.go    # UserUseCase, EmailUseCase interfaces
│   │   └── secondary/          # Driven Ports (what app needs)
│   │       ├── repository.go   # Repository interfaces
│   │       └── service.go      # External service interfaces
│   ├── adapters/               # Adapter Implementations
│   │   ├── primary/            # Driving Adapters (UI, API)
│   │   │   └── http/           # HTTP adapter (Gin framework)
│   │   │       ├── dto/        # HTTP-specific DTOs
│   │   │       ├── handlers/   # HTTP handlers
│   │   │       ├── mappers/    # DTO ↔ Domain mappers
│   │   │       ├── messages/   # Response messages
│   │   │       ├── response/   # Response utilities
│   │   │       ├── middleware.go
│   │   │       └── routes.go
│   │   └── secondary/          # Driven Adapters (Infrastructure)
│   │       ├── postgres/       # PostgreSQL implementation
│   │       │   ├── base_repository.go
│   │       │   ├── user_repository.go
│   │       │   └── refresh_token_repository.go
│   │       ├── crypto/         # Cryptography services
│   │       │   ├── jwt_service.go
│   │       │   ├── bcrypt_service.go
│   │       │   └── aes_service.go
│   │       ├── smtp/           # Email service
│   │       │   ├── smtp_service.go
│   │       │   └── templates/
│   │       └── localstorage/   # File storage service
│   │           └── localstorage_service.go
│   └── infrastructure/         # Infrastructure concerns
│       └── container.go        # Dependency injection
└── pkg/                        # Public libraries
    ├── config/                 # Configuration management
    └── utils/                  # Utility functions
```

**Hexagonal Architecture Layers:**

1. **Domain (Core)**
   - **Models**: Pure business entities (User, RefreshToken)
   - **Services**: Business logic orchestrating entities
   - **Errors**: Domain-specific errors
   - ❗ **Zero external dependencies**

2. **Ports (Interfaces)**
   - **Primary Ports**: Interfaces for use cases (what application offers)
   - **Secondary Ports**: Interfaces for repositories/services (what application needs)
   - ❗ **Defined by domain, implemented by adapters**

3. **Adapters**
   - **Primary Adapters**: HTTP handlers, CLI, gRPC (driving the app)
   - **Secondary Adapters**: Database, SMTP, file storage (driven by the app)
   - ❗ **Implement port interfaces, depend on domain**

4. **Infrastructure**
   - Dependency injection and wiring
   - Application configuration

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL
- Git

### Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd go-gin-clean
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Environment setup**
   Copy `.env.example` to `.env` and configure:

   ```env
   # Server
   SERVER_HOST=localhost
   SERVER_PORT=8080
   ENVIRONMENT=development
   APP_FE_URL=
   TIMEOUT=30

   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USERNAME=postgres
   DB_PASSWORD=your_password
   DB_NAME=go_gin_clean
   DB_MAX_IDLE_CONNS=25
   DB_MAX_OPEN_CONNS=5

   # JWT
   JWT_ISSUER=go-gin-clean
   JWT_ACCESS_SECRET=your-access-secret-key
   JWT_REFRESH_SECRET=your-refresh-secret-key
   JWT_ACCESS_EXPIRY=1h
   JWT_REFRESH_EXPIRY=168h

   # AES Encryption
   AES_KEY=your-32-character-encryption-key
   AES_IV=your-16-character-iv-key

   # SMTP Email (optional)
   MAILER_HOST=smtp.gmail.com
   MAILER_PORT=587
   MAILER_SENDER=your-email@gmail.com
   MAILER_AUTH=your-email@gmail.com
   MAILER_PASSWORD=your-app-password
   ```

4. **Database migration**

   ```bash
   go run cmd/migrate/main.go migrate
   ```

5. **Start the server**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Health Check

- `GET /health` - Server health status

### Authentication (Public Routes)

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login (sets refresh token in cookie)
- `POST /api/v1/auth/refresh-token` - Refresh access token
- `POST /api/v1/auth/verify-email` - Email verification
- `POST /api/v1/auth/send-verify-email` - Send verification email
- `POST /api/v1/auth/send-reset-password` - Send reset password email
- `POST /api/v1/auth/reset-password` - Reset password with token

### Profile Management (Protected Routes)

- `GET /api/v1/profile` - Get current user profile
- `PUT /api/v1/profile` - Update current user profile
- `POST /api/v1/profile/change-password` - Change user password
- `POST /api/v1/profile/logout` - User logout

### User Management (Protected Routes)

- `GET /api/v1/users` - Get all users (paginated)
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Static Assets

- `GET /assets/*` - Serve static files from assets directory

## 🔧 Available Commands

### Database Commands

```bash
# Run migrations
go run cmd/migrate/main.go migrate

# Rollback migrations
go run cmd/migrate/main.go rollback

# Fresh migrations (rollback + migrate)
go run cmd/migrate/main.go fresh
```

### Development Commands

```bash
# Start server
go run cmd/server/main.go

# Build server
go build -o bin/server cmd/server/main.go

# Build migration tool
go build -o bin/migrate cmd/migrate/main.go

# Run tests
go test ./...

# Check for issues
go vet ./...
```

## 🏛️ Hexagonal Architecture Benefits

### 1. **Domain Independence**

- Domain has **ZERO dependencies** on frameworks, HTTP, databases, or external libraries
- Business logic is pure Go code with only domain model dependencies
- Easy to test domain services with simple mock interfaces
- **Port interfaces** define what domain needs, adapters provide implementations

### 2. **Dependency Inversion**

- All dependencies point **inward** toward the domain
- Domain defines interfaces (ports), infrastructure implements them
- Framework changes don't affect business logic
- Can swap Gin for Fiber/Echo without touching domain code

### 3. **Testability**

- Domain services tested with simple port mocks
- No need to mock framework-specific types (no `*gin.Context` in tests)
- Unit test business logic in isolation
- Integration test adapters independently

### 4. **Flexibility & Scalability**

- **Plug-and-play architecture**: Swap any adapter without domain changes
- Add new delivery mechanisms (GraphQL, gRPC, CLI) alongside HTTP
- Switch databases (Postgres → MongoDB) by implementing port interface
- Horizontal scaling through clear component boundaries

### 5. **Explicit Boundaries**

- **Primary Adapters** (HTTP, gRPC) → drive the application
- **Primary Ports** (UserUseCase) → what application offers
- **Domain Services** (UserService) → business logic
- **Secondary Ports** (UserRepository, JWTService) → what application needs
- **Secondary Adapters** (Postgres, Crypto) → infrastructure implementations

## 🧪 Testing

The Hexagonal Architecture makes testing straightforward with clear boundaries:

### Domain Service Testing (Pure Business Logic)
```go
// Test domain services with port mocks - no framework dependencies
func TestUserService_Login(t *testing.T) {
    // Arrange
    mockUserRepo := &mocks.UserRepository{}
    mockJWTService := &mocks.JWTService{}
    mockBcryptService := &mocks.BcryptService{}
    // ... other mocked ports

    userService := service.NewUserService(
        mockUserRepo,
        mockEmailService,
        mockRefreshTokenRepo,
        mockJWTService,
        mockBcryptService,
        mockAESService,
        mockMediaService,
    )

    loginReq := &primary.LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }

    // Act
    result, err := userService.Login(context.Background(), loginReq)

    // Assert - pure domain testing
    assert.NoError(t, err)
    assert.NotEmpty(t, result.AccessToken)
}
```

### Adapter Testing (Repository)
```go
// Test postgres adapter with real database or testcontainers
func TestUserRepository_FindByEmail(t *testing.T) {
    db := setupTestDB(t)
    userRepo := postgres.NewUserRepository(db)

    user := &model.User{
        Name:  "Test User",
        Email: "test@example.com",
    }

    // Test repository implementation
    savedUser, err := userRepo.Create(context.Background(), user)
    assert.NoError(t, err)

    foundUser, err := userRepo.FindByEmail(context.Background(), "test@example.com")
    assert.NoError(t, err)
    assert.Equal(t, savedUser.ID, foundUser.ID)
}
```

### HTTP Handler Testing
```go
// Test HTTP handlers with mocked use cases
func TestUserHandler_Login(t *testing.T) {
    mockUseCase := &mocks.UserUseCase{}
    mockMapper := mappers.NewUserMapper()

    handler := handlers.NewUserHandler(mockUseCase, mockMapper)

    // Setup Gin test context
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    // Test HTTP layer separately from business logic
    handler.Login(c)

    assert.Equal(t, http.StatusOK, w.Code)
}
```

### Port Interface Mocking
```go
// Easy to create mocks for port interfaces
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}
```

## 🔐 Security Features

### Password Security

- **Bcrypt Hashing**: Industry-standard password hashing
- **Salt Generation**: Automatic salt generation for each password
- **Cost Factor**: Configurable cost factor for security vs performance

### JWT Security

- **HMAC Signing**: Secure token signing with secret keys
- **Token Expiration**: Configurable expiration times
- **Refresh Rotation**: Secure refresh token rotation

### Data Encryption

- **AES Encryption**: Additional data encryption capabilities
- **Configurable Keys**: Environment-based encryption keys
- **PKCS7 Padding**: Standard padding for block cipher

## 🔒 Authentication

The application uses JWT-based authentication with the following features:

- **Access Token**: Short-lived token for API requests (1 hour)
- **Refresh Token**: Long-lived token stored in HTTP-only cookie (7 days)
- **Password Hashing**: Bcrypt for secure password storage
- **AES Encryption**: Additional data encryption capabilities

### Authentication Flow

1. **Login**: User provides email/password, receives access token and refresh token (in cookie)
2. **API Requests**: Include access token in Authorization header
3. **Token Refresh**: Automatic refresh using HTTP-only cookie
4. **Logout**: Clears refresh token and invalidates session

Include the access token in requests:

```
Authorization: Bearer <access_token>
```

## 📧 Email Features

The application includes comprehensive email functionality:

- **Email Verification**: Send verification emails to new users
- **Password Reset**: Send password reset emails with secure tokens
- **Template System**: HTML email templates with dynamic data
- **SMTP Integration**: Configurable SMTP service for email delivery

## 📁 File Upload

Local file storage implementation:

- **Avatar Upload**: Users can upload profile pictures
- **Local Storage**: Files stored in local filesystem
- **File Validation**: Type and size validation
- **Secure Paths**: Protected file path handling

## 🛠️ Development Guidelines

### Adding New Features (Hexagonal Architecture Flow)

1. **Define Domain Models**: Create or update entities in `internal/domain/model/`
   ```go
   // internal/domain/model/product.go
   type Product struct {
       ID    int64
       Name  string
       Price float64
   }
   ```

2. **Define Secondary Ports**: Create interfaces for what domain needs in `internal/ports/secondary/`
   ```go
   // internal/ports/secondary/repository.go
   type ProductRepository interface {
       FindByID(ctx context.Context, id int64) (*model.Product, error)
       Create(ctx context.Context, product *model.Product) (*model.Product, error)
   }
   ```

3. **Define Primary Ports**: Create use case interfaces in `internal/ports/primary/`
   ```go
   // internal/ports/primary/product_port.go
   type ProductUseCase interface {
       GetProduct(ctx context.Context, id int64) (*ProductInfo, error)
       CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductInfo, error)
   }
   ```

4. **Implement Domain Service**: Add business logic in `internal/domain/service/`
   ```go
   // internal/domain/service/product_service.go
   type ProductService struct {
       productRepo secondary.ProductRepository
   }

   func (s *ProductService) GetProduct(ctx context.Context, id int64) (*primary.ProductInfo, error) {
       product, err := s.productRepo.FindByID(ctx, id)
       // ... business logic
   }
   ```

5. **Implement Secondary Adapters**: Create infrastructure in `internal/adapters/secondary/`
   ```go
   // internal/adapters/secondary/postgres/product_repository.go
   type ProductRepository struct {
       db *gorm.DB
   }

   func (r *ProductRepository) FindByID(ctx context.Context, id int64) (*model.Product, error) {
       // Postgres implementation
   }
   ```

6. **Implement Primary Adapters**: Create HTTP handlers in `internal/adapters/primary/http/`
   ```go
   // internal/adapters/primary/http/handlers/product_handler.go
   type ProductHandler struct {
       productUseCase primary.ProductUseCase
   }
   ```

7. **Wire Dependencies**: Update `internal/infrastructure/container.go`
   ```go
   productRepo := postgres.NewProductRepository(db)
   productUseCase := service.NewProductService(productRepo)
   ```

8. **Add Routes**: Register in `internal/adapters/primary/http/routes.go`

### Architecture Rules

1. **Dependency Rule**: Dependencies ALWAYS point inward
   - ✅ Domain → nothing (pure business logic)
   - ✅ Ports → Domain (interfaces reference domain models)
   - ✅ Services → Ports (services implement primary ports, use secondary ports)
   - ✅ Adapters → Ports + Domain (adapters implement ports)
   - ❌ Domain → Ports/Adapters (NEVER)

2. **Port Ownership**
   - Domain **defines** secondary ports (what it needs)
   - Domain **implements** primary ports (what it offers)
   - Adapters **implement** secondary ports
   - Adapters **use** primary ports

3. **Framework Isolation**
   - Gin, GORM, HTTP types → stay in adapters
   - Domain has no `*gin.Context`, `*gorm.DB`, etc.
   - Use plain Go types in domain

4. **Testing Strategy**
   - Domain services: Mock secondary ports
   - Adapters: Test with real dependencies or testcontainers
   - HTTP layer: Mock primary ports (use cases)

### Error Handling

- **Domain errors** in `internal/domain/error/` (business rules)
  ```go
  var ErrUserNotFound = errors.New("user not found")
  var ErrEmailAlreadyExists = errors.New("email already exists")
  ```
- **Adapter errors**: Wrap/translate infrastructure errors
- **HTTP errors**: Convert domain errors to HTTP responses

### Layer Communication Flow

```go
// ❌ WRONG - Domain importing from adapter
package service
import "go-gin-clean/internal/adapters/secondary/postgres"

// ✅ CORRECT - Domain importing only ports
package service
import "go-gin-clean/internal/ports/secondary"

// ❌ WRONG - Domain using framework types
func (s *UserService) Login(c *gin.Context) error

// ✅ CORRECT - Domain using domain/port types
func (s *UserService) Login(ctx context.Context, req *primary.LoginRequest) (*primary.LoginResponse, error)

// ❌ WRONG - Adapter bypassing ports
handler → domain service directly

// ✅ CORRECT - Adapter through ports
handler → primary port interface → domain service → secondary port interface → adapter
```

### Current Implementation Overview

**Domain Layer:**
- `model.User`, `model.RefreshToken`, `model.Gender` - Domain entities
- `service.UserService` - User business logic
- `service.EmailService` - Email business logic
- `error.Err*` - Business errors

**Ports:**
- `primary.UserUseCase`, `primary.EmailUseCase` - What app offers
- `secondary.UserRepository`, `secondary.JWTService`, etc. - What app needs

**Adapters:**
- `postgres.*Repository` - Database implementations
- `crypto.*Service` - Security service implementations
- `smtp.SMTPService` - Email service implementation
- `localstorage.LocalStorageService` - File storage implementation
- `http.Handler` - HTTP request handlers

## 📈 Performance & Production

### Build for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o bin/server cmd/server/main.go

# Set production environment
export ENVIRONMENT=production
```

### Docker (Optional)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
CMD ["./server"]
```

## 🌟 Features

### Core Features

- ✅ User Registration and Authentication
- ✅ JWT Token-based Authentication with Refresh Tokens
- ✅ Email Verification System
- ✅ Password Reset Functionality
- ✅ User Profile Management
- ✅ File Upload (Avatar)
- ✅ Pagination Support
- ✅ CORS Configuration
- ✅ Middleware Authentication

### Security Features

- ✅ Bcrypt Password Hashing
- ✅ JWT Token Security
- ✅ AES Data Encryption
- ✅ HTTP-Only Cookie for Refresh Tokens
- ✅ Input Validation and Sanitization

### Infrastructure Features

- ✅ Database Migration System
- ✅ Configuration Management
- ✅ SMTP Email Integration
- ✅ Local File Storage
- ✅ Structured Logging
- ✅ Graceful Shutdown
