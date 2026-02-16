# Go API Scaffold

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/crushzh/go-api-scaffold?style=flat-square)](https://goreportcard.com/report/github.com/crushzh/go-api-scaffold)
[![Release](https://img.shields.io/github/v/release/crushzh/go-api-scaffold?style=flat-square)](https://github.com/crushzh/go-api-scaffold/releases)

> Production-ready Go REST API boilerplate with Gin, GORM, JWT auth, Swagger docs, gRPC support, and built-in code generator.

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md)

## Features

- **Gin** HTTP framework with recovery, CORS, request ID, logger, and timeout middleware
- **GORM** ORM with SQLite / MySQL / PostgreSQL support
- **JWT** authentication with role-based access control
- **gRPC** dual-protocol support (HTTP + gRPC)
- **Swagger** API documentation (via `swag`)
- **Code generator** — scaffold full CRUD modules in one command
- **Cross-platform build** — Linux (amd64/arm64/arm32), Windows, macOS
- **Docker** support with multi-stage build
- **Embedded frontend** — serve SPA via `go:embed`
- **Structured logging** with Zap + Lumberjack rotation
- **Unified response** format with standard error codes
- **Graceful shutdown** with signal handling
- **Service management** scripts (systemd / watchdog)

## Tech Stack

| Library | Version | Purpose |
|---------|---------|---------|
| [Go](https://go.dev/) | 1.21+ | Language runtime |
| [Gin](https://gin-gonic.com/) | v1.9 | HTTP framework |
| [GORM](https://gorm.io/) | v1.25 | ORM (SQLite/MySQL/PostgreSQL) |
| [Viper](https://github.com/spf13/viper) | v1.18 | Configuration management |
| [Zap](https://go.uber.org/zap) | v1.27 | Structured logging |
| [JWT](https://github.com/golang-jwt/jwt) | v5 | JSON Web Token authentication |
| [gRPC](https://grpc.io/) | v1.62 | RPC framework |
| [Lumberjack](https://github.com/natefinsh/lumberjack) | v2.2 | Log file rotation |

## Quick Start

```bash
# 1. Clone
git clone https://github.com/crushzh/go-api-scaffold.git
cd go-api-scaffold

# 2. Install dependencies
go mod download

# 3. Run
make run
# Server starts at http://localhost:8080

# 4. Test
curl http://localhost:8080/health

# 5. Login (default: admin / admin123)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## Project Structure

```
go-api-scaffold/
├── cmd/
│   ├── server/           # Application entry point
│   └── gen/              # Code generator
├── internal/
│   ├── handler/          # HTTP handlers + middleware + router
│   ├── service/          # Business logic layer
│   ├── model/            # Data models (GORM)
│   ├── store/            # Data access layer (repositories)
│   └── web/              # Embedded frontend (go:embed)
├── pkg/
│   ├── config/           # Configuration (Viper)
│   ├── logger/           # Logging (Zap + Lumberjack)
│   └── response/         # Unified API response
├── api/proto/            # Protocol Buffer definitions
├── configs/              # Configuration files
├── templates/            # Code generator templates
├── scripts/              # Deployment scripts
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

**Architecture**: Handler -> Service -> Store (3-layer)

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌──────────┐
│ Handler  │───>│ Service │───>│  Store  │───>│ Database │
│ (HTTP)   │    │ (Logic) │    │ (Repo)  │    │ (GORM)   │
└─────────┘    └─────────┘    └─────────┘    └──────────┘
```

## Code Generator

Generate a complete CRUD module with a single command:

```bash
make gen name=order cn=Order
```

This creates 4 files and auto-registers routes:

| File | Description |
|------|-------------|
| `internal/handler/order_handler.go` | HTTP CRUD endpoints + Swagger |
| `internal/service/order_service.go` | Business logic |
| `internal/model/order.go` | Data model + DTOs |
| `internal/store/order_repo.go` | Database repository |

Auto-appended: routes in `router.go`, migration in `store.go`.

## Configuration

Loaded from `configs/config.yaml`. Override with `APP_` environment variables.

```yaml
app:
  name: "myapp"
  mode: "debug"           # debug, release, test

server:
  host: "0.0.0.0"
  port: 8080

grpc:
  enabled: false
  port: 9090

database:
  type: "sqlite"          # sqlite, mysql, postgres
  path: "./data/app.db"

jwt:
  secret: "change-me-in-production"
  expire: 24              # hours
  refresh_hours: 168      # 7 days
```

Environment variable examples:

```bash
APP_SERVER_PORT=3000 APP_DATABASE_TYPE=mysql ./myapp
```

## API Examples

```bash
# Health check
curl http://localhost:8080/health

# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.data.token')

# CRUD operations
curl -X POST http://localhost:8080/api/v1/examples \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"test","description":"hello"}'

curl http://localhost:8080/api/v1/examples?page=1&page_size=10 \
  -H "Authorization: Bearer $TOKEN"
```

## Unified Response Format

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

| Code Range | Category | Description |
|------------|----------|-------------|
| 0 | Success | Operation successful |
| 1001-1999 | Client | Parameter / validation errors |
| 2001-2999 | Resource | Not found, conflict |
| 3001-3999 | Business | Business logic errors |
| 4001-4999 | Auth | Unauthorized, forbidden, expired |
| 5001-5999 | System | Internal, database, timeout |

## Deployment

### Docker

```bash
docker-compose up -d
```

### Cross-compile + Manual

```bash
make build-all            # Build all platforms
# Or specific target:
make build-linux          # Linux amd64
make build-arm64          # Linux arm64
make build-arm32          # Linux arm32
make build-windows        # Windows

# Deploy
scp -r build/linux-amd64/* user@server:/opt/myapp/
```

### Service Management

```bash
./scripts/manage.sh start     # Start with watchdog
./scripts/manage.sh stop      # Graceful stop
./scripts/manage.sh restart   # Restart
./scripts/manage.sh status    # Check status
```

## Customization

Rename the module for your project:

```bash
# macOS
find . -name "*.go" -exec sed -i '' 's|go-api-scaffold|my-project|g' {} +
sed -i '' 's|go-api-scaffold|my-project|g' go.mod Makefile configs/config.yaml

# Linux
find . -name "*.go" -exec sed -i 's|go-api-scaffold|my-project|g' {} +
sed -i 's|go-api-scaffold|my-project|g' go.mod Makefile configs/config.yaml

go mod tidy
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE)
