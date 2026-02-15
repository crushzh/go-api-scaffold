# Go API Scaffold

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)

> 開箱即用的 Go REST API 腳手架，整合 Gin、GORM、JWT 認證、Swagger 文件、gRPC 支援和程式碼產生器。

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md)

## 特色

- **Gin** HTTP 框架，內建 Recovery、CORS、請求 ID、日誌、逾時中介軟體
- **GORM** ORM，支援 SQLite / MySQL / PostgreSQL
- **JWT** 認證，支援角色權限控制
- **gRPC** 雙協定支援（HTTP + gRPC）
- **Swagger** API 文件自動產生
- **程式碼產生器** — 一條指令產生完整 CRUD 模組
- **跨平台編譯** — Linux (amd64/arm64/arm32)、Windows、macOS
- **Docker** 多階段建置
- **前端嵌入** — `go:embed` 內嵌 SPA 前端
- **結構化日誌** — Zap + Lumberjack 日誌輪替
- **統一回應格式** — 標準錯誤碼體系
- **優雅退出** — 訊號處理
- **服務管理腳本** — systemd / 守護程序

## 快速開始

```bash
# 1. 複製儲存庫
git clone https://github.com/mrzhoong/go-api-scaffold.git
cd go-api-scaffold

# 2. 安裝依賴
go mod download

# 3. 執行
make run
# 服務啟動: http://localhost:8080

# 4. 測試
curl http://localhost:8080/health

# 5. 登入（預設帳號: admin / admin123）
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## 專案結構

```
go-api-scaffold/
├── cmd/
│   ├── server/           # 程式進入點
│   └── gen/              # 程式碼產生器
├── internal/
│   ├── handler/          # 介面層（HTTP Handler + 中介軟體 + 路由）
│   ├── service/          # 業務邏輯層
│   ├── model/            # 資料模型（GORM）
│   ├── store/            # 資料存取層（儲存庫）
│   └── web/              # 嵌入式前端（go:embed）
├── pkg/
│   ├── config/           # 設定管理（Viper）
│   ├── logger/           # 日誌（Zap + Lumberjack）
│   └── response/         # 統一 API 回應
├── api/proto/            # Protocol Buffer 定義
├── configs/              # 設定檔
├── templates/            # 程式碼產生範本
├── scripts/              # 部署腳本
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

**架構**: Handler -> Service -> Store（三層架構）

## 程式碼產生器

一條指令產生完整 CRUD 模組:

```bash
make gen name=order cn=訂單
```

產生 4 個檔案並自動註冊路由:

| 檔案 | 說明 |
|------|------|
| `internal/handler/order_handler.go` | HTTP CRUD 介面 + Swagger |
| `internal/service/order_service.go` | 業務邏輯 |
| `internal/model/order.go` | 資料模型 + DTO |
| `internal/store/order_repo.go` | 資料儲存庫 |

## 設定

透過 `configs/config.yaml` 載入設定，支援 `APP_` 前綴的環境變數覆寫。

## API 範例

```bash
# 健康檢查
curl http://localhost:8080/health

# 登入
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.data.token')

# CRUD 操作
curl -X POST http://localhost:8080/api/v1/examples \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"test","description":"hello"}'
```

## 統一回應格式

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

| 錯誤碼範圍 | 分類 | 說明 |
|-----------|------|------|
| 0 | 成功 | 操作成功 |
| 1001-1999 | 客戶端 | 參數/驗證錯誤 |
| 2001-2999 | 資源 | 不存在、衝突 |
| 4001-4999 | 認證 | 未認證、無權限 |
| 5001-5999 | 系統 | 內部錯誤、資料庫錯誤 |

## 部署

### Docker

```bash
docker-compose up -d
```

### 跨平台編譯

```bash
make build-all            # 全平台編譯
```

### 服務管理

```bash
./scripts/manage.sh start     # 啟動
./scripts/manage.sh stop      # 停止
./scripts/manage.sh restart   # 重新啟動
./scripts/manage.sh status    # 檢視狀態
```

## 自訂

重新命名模組:

```bash
# macOS
find . -name "*.go" -exec sed -i '' 's|go-api-scaffold|my-project|g' {} +
sed -i '' 's|go-api-scaffold|my-project|g' go.mod Makefile configs/config.yaml
go mod tidy
```

## 貢獻

參見 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 授權條款

[MIT](LICENSE)
