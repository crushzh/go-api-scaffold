package handler

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go-api-scaffold/internal/service"
	"go-api-scaffold/internal/web"
	"go-api-scaffold/pkg/config"
	"go-api-scaffold/pkg/logger"
	"go-api-scaffold/pkg/response"

	_ "go-api-scaffold/docs/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter creates the HTTP router
func NewRouter(cfg *config.Config, authSvc *service.AuthService, exampleSvc *service.ExampleService) *gin.Engine {
	gin.SetMode(cfg.App.Mode)

	r := gin.New()

	// Global middleware
	r.Use(Recovery())
	r.Use(CORS())
	r.Use(RequestID())
	r.Use(Logger())
	r.Use(Timeout(30*time.Second, "/ws/", "/health", "/swagger/"))

	// ====== Base routes ======
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status":  "ok",
			"version": cfg.App.Version,
		})
	})

	// ====== API routes ======
	api := r.Group("/api/v1")
	{
		// Auth (no token required)
		authHandler := NewAuthHandler(authSvc)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Authenticated routes
		authorized := api.Group("")
		authorized.Use(AuthMiddleware(authSvc))
		{
			authorized.GET("/auth/profile", authHandler.GetProfile)

			// Example module CRUD
			exampleHandler := NewExampleHandler(exampleSvc)
			examples := authorized.Group("/examples")
			{
				examples.GET("", exampleHandler.List)
				examples.POST("", exampleHandler.Create)
				examples.GET("/:id", exampleHandler.Get)
				examples.PUT("/:id", exampleHandler.Update)
				examples.DELETE("/:id", exampleHandler.Delete)
			}

			// Admin-only routes example
			// admin := authorized.Group("")
			// admin.Use(RequireRole("admin"))
			// { ... }

			// GEN:ROUTE_REGISTER - Auto-appended by code generator, do not remove
		}
	}

	// ====== Swagger ======
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ====== Frontend static files ======
	registerFrontendRoutes(r)

	return r
}

// registerFrontendRoutes registers frontend static file routes (SPA support)
func registerFrontendRoutes(r *gin.Engine) {
	// Method 1: go:embed
	distFS, err := web.GetDistFS()
	if err == nil {
		registerEmbeddedFrontend(r, distFS)
		logger.Info("frontend files embedded in binary")
		return
	}

	// Method 2: external dist directory
	for _, dir := range []string{"dist", "web", "static", "public"} {
		absDir, _ := filepath.Abs(dir)
		if isDir(absDir) {
			r.Static("/assets", filepath.Join(absDir, "assets"))
			r.StaticFile("/", filepath.Join(absDir, "index.html"))
			r.NoRoute(spaFallback(absDir))
			logger.Infof("frontend directory: %s", absDir)
			return
		}
	}

	logger.Warn("frontend files not found (ignore in dev mode)")
}

func registerEmbeddedFrontend(r *gin.Engine, distFS fs.FS) {
	httpFS := http.FS(distFS)

	// Static assets (long-term cache)
	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		c.FileFromFS(c.Request.URL.Path, httpFS)
	})

	// SPA fallback
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API routes return 404
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/swagger/") {
			response.NotFound(c, "API not found")
			return
		}

		// Try to find the file
		if f, err := distFS.Open(strings.TrimPrefix(path, "/")); err == nil {
			f.Close()
			c.FileFromFS(path, httpFS)
			return
		}

		// Fall back to index.html
		if indexFile, err := fs.ReadFile(distFS, "index.html"); err == nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexFile)
			return
		}

		response.NotFound(c, "page not found")
	})
}

func spaFallback(staticDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/swagger/") {
			response.NotFound(c, "API not found")
			return
		}
		c.File(filepath.Join(staticDir, "index.html"))
	}
}

func isDir(path string) bool {
	info, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	_ = info
	// Check if directory exists
	entries, err := filepath.Glob(filepath.Join(path, "*"))
	return err == nil && len(entries) > 0
}
