package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-api-scaffold/pkg/logger"
	"go-api-scaffold/pkg/response"
)

// Recovery handles panics and returns 500
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("[Recovery] panic: %v", err)
				response.ServerError(c, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID, X-Total-Count")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RequestID adds a unique request ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("X-Request-ID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Logger logs errors and slow requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// Only log 4xx/5xx or slow requests (>1s)
		if status >= 400 || latency > time.Second {
			logger.Warnf("[%d] %s %s %s %v",
				status, c.Request.Method, c.Request.URL.Path, c.ClientIP(), latency)
		}
	}
}

// Timeout sets a deadline on the request context.
// Handlers should check ctx.Err() for long-running operations.
func Timeout(timeout time.Duration, skipPaths ...string) gin.HandlerFunc {
	skipMap := make(map[string]bool, len(skipPaths))
	for _, p := range skipPaths {
		skipMap[p] = true
	}

	return func(c *gin.Context) {
		for prefix := range skipMap {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				c.Next()
				return
			}
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// SecurityHeaders 添加安全响应头
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Cache-Control", "no-store")
		c.Next()
	}
}

// RateLimit IP 级别速率限制
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	type visitor struct {
		count   int
		resetAt time.Time
	}
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		v, exists := visitors[ip]
		now := time.Now()
		if !exists || now.After(v.resetAt) {
			visitors[ip] = &visitor{count: 1, resetAt: now.Add(window)}
			mu.Unlock()
			c.Next()
			return
		}
		v.count++
		if v.count > maxRequests {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    50004,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		mu.Unlock()
		c.Next()
	}
}

// DemoMode 演示模式中间件：拦截写操作
func DemoMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		// 允许登录
		if strings.HasSuffix(c.FullPath(), "/auth/login") {
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{
			"code":    30001,
			"message": "演示模式，不允许修改操作",
		})
		c.Abort()
	}
}

// HealthReady 就绪探针（含依赖检查）
func HealthReady(checks ...func() error) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, check := range checks {
			if err := check(); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "not ready",
					"error":  err.Error(),
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
}
