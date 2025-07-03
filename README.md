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
│   ├── adapters/
│   │   ├── auth/
│   │   │   ├── bcrypt.go        # Password hashing adapter
│   │   │   └── jwt.go           # JWT token adapter
│   │   ├── database/
│   │   │   ├── refresh_token_repository.go
│   │   │   └── user_repository.go
│   │   └── http/
│   │       ├── handlers/
│   │       │   ├── auth_handler.go
│   │       │   └── user_handler.go
│   │       ├── middleware/
│   │       │   ├── cors.go
│   │       │   └── middleware.go
│   │       └── routes/
│   │           └── router.go
│   ├── application/
│   │   ├── dto/
│   │   │   ├── auth_dto.go
│   │   │   └── user_dto.go
│   │   └── service/
│   │       ├── auth_service.go
│   │       └── user_service.go
│   └── domain/
│       ├── entity/
│       │   ├── refresh_token.go
│       │   └── user.go
│       └── ports/
│           ├── repositories.go
│           └── services.go
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── postgres.go
│   ├── message/
│   │   ├── error.go
│   │   └── success.go
│   ├── response/
│   │   └── response.go
│   └── utils/
│       └── string_utils.go
├── go.mod
├── go.sum
└── README.md
```

## How to Run

### 1. Run API Server

Using Air (for hot reload):

```bash
air
```

Using Go run:

```bash
go run cmd/api/main.go
```

### 2. Run Database Migration

```bash
go run cmd/migrate/main.go
```
