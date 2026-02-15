# ==================== Build Stage ====================
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Dependency cache
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=1 go build -trimpath \
    -ldflags "-s -w -X 'main.Version=$(cat VERSION 2>/dev/null || echo dev)' -X 'main.BuildTime=$(date +%Y-%m-%dT%H:%M:%S)'" \
    -o /app/server ./cmd/server/

# ==================== Runtime Stage ====================
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# Copy artifacts from build stage
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/internal/web/dist ./dist

# Data directories
RUN mkdir -p data logs

EXPOSE 8080 9090

ENTRYPOINT ["./server"]
CMD ["-c", "configs/config.yaml"]
