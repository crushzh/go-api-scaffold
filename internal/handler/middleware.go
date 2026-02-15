package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go-api-scaffold/pkg/logger"
	"go-api-scaffold/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
