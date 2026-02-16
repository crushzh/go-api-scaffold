# Go API Scaffold

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)

> 开箱即用的 Go REST API 脚手架，集成 Gin、GORM、JWT 认证、Swagger 文档、gRPC 支持和代码生成器。

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md)

## 特性

- **Gin** HTTP 框架，内置 Recovery、CORS、请求 ID、日志、超时中间件
- **GORM** ORM，支持 SQLite / MySQL / PostgreSQL
- **JWT** 认证，支持角色权限控制
- **gRPC** 双协议支持（HTTP + gRPC）
- **Swagger** API 文档自动生成
- **代码生成器** — 一条命令生成完整 CRUD 模块
- **跨平台编译** — Linux (amd64/arm64/arm32)、Windows、macOS
- **Docker** 多阶段构建
- **前端嵌入** — `go:embed` 内嵌 SPA 前端
- **结构化日志** — Zap + Lumberjack 日志轮转
- **统一响应格式** — 标准错误码体系
- **优雅退出** — 信号处理
- **服务管理脚本** — systemd / 守护进程

## 快速开始

```bash
# 1. 克隆仓库
git clone https://github.com/crushzh/go-api-scaffold.git
cd go-api-scaffold

# 2. 安装依赖
go mod download

# 3. 运行
make run
# 服务启动: http://localhost:8080

# 4. 测试
curl http://localhost:8080/health

# 5. 登录（默认账号: admin / admin123）
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## 项目结构

```
go-api-scaffold/
├── cmd/
│   ├── server/           # 程序入口
│   └── gen/              # 代码生成器
├── internal/
│   ├── handler/          # 接口层（HTTP Handler + 中间件 + 路由）
│   ├── service/          # 业务逻辑层
│   ├── model/            # 数据模型（GORM）
│   ├── store/            # 数据访问层（仓储）
│   └── web/              # 嵌入式前端（go:embed）
├── pkg/
│   ├── config/           # 配置管理（Viper）
│   ├── logger/           # 日志（Zap + Lumberjack）
│   └── response/         # 统一 API 响应
├── api/proto/            # Protocol Buffer 定义
├── configs/              # 配置文件
├── templates/            # 代码生成模板
├── scripts/              # 部署脚本
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

**架构**: Handler -> Service -> Store（三层架构）

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌──────────┐
│ Handler  │───>│ Service │───>│  Store  │───>│ Database │
│ (HTTP)   │    │ (业务)   │    │ (仓储)  │    │ (GORM)   │
└─────────┘    └─────────┘    └─────────┘    └──────────┘
```

## 代码生成器

一条命令生成完整 CRUD 模块:

```bash
make gen name=order cn=订单
```

生成 4 个文件并自动注册路由:

| 文件 | 说明 |
|------|------|
| `internal/handler/order_handler.go` | HTTP CRUD 接口 + Swagger |
| `internal/service/order_service.go` | 业务逻辑 |
| `internal/model/order.go` | 数据模型 + DTO |
| `internal/store/order_repo.go` | 数据仓储 |

自动追加: `router.go` 路由注册、`store.go` 模型迁移。

## 配置

通过 `configs/config.yaml` 加载配置，支持 `APP_` 前缀的环境变量覆盖。

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
  expire: 24              # 小时
  refresh_hours: 168      # 7 天
```

## API 示例

```bash
# 健康检查
curl http://localhost:8080/health

# 登录
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.data.token')

# CRUD 操作
curl -X POST http://localhost:8080/api/v1/examples \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"test","description":"hello"}'

curl http://localhost:8080/api/v1/examples?page=1&page_size=10 \
  -H "Authorization: Bearer $TOKEN"
```

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

| 错误码范围 | 分类 | 说明 |
|-----------|------|------|
| 0 | 成功 | 操作成功 |
| 1001-1999 | 客户端 | 参数/校验错误 |
| 2001-2999 | 资源 | 不存在、冲突 |
| 3001-3999 | 业务 | 业务逻辑错误 |
| 4001-4999 | 认证 | 未认证、无权限、Token 过期 |
| 5001-5999 | 系统 | 内部错误、数据库错误、超时 |

## 部署

### Docker

```bash
docker-compose up -d
```

### 跨平台编译 + 手动部署

```bash
make build-all            # 全平台编译
make build-linux          # Linux amd64
make build-arm64          # Linux arm64
make build-arm32          # Linux arm32
make build-windows        # Windows

# 部署到服务器
scp -r build/linux-amd64/* user@server:/opt/myapp/
```

### 服务管理

```bash
./scripts/manage.sh start     # 启动（含守护进程）
./scripts/manage.sh stop      # 优雅停止
./scripts/manage.sh restart   # 重启
./scripts/manage.sh status    # 查看状态
```

## 自定义

重命名模块:

```bash
# macOS
find . -name "*.go" -exec sed -i '' 's|go-api-scaffold|my-project|g' {} +
sed -i '' 's|go-api-scaffold|my-project|g' go.mod Makefile configs/config.yaml

# Linux
find . -name "*.go" -exec sed -i 's|go-api-scaffold|my-project|g' {} +
sed -i 's|go-api-scaffold|my-project|g' go.mod Makefile configs/config.yaml

go mod tidy
```

## 贡献

参见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 许可证

[MIT](LICENSE)
