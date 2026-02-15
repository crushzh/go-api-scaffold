package handler

import (
	"strings"

	"go-api-scaffold/internal/service"
	"go-api-scaffold/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// LoginRequest is the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles user login
// @Summary  User login
// @Tags     Auth
// @Accept   json
// @Produce  json
// @Param    body body LoginRequest true "Login credentials"
// @Success  200  {object} response.Response{data=service.TokenResponse}
// @Router   /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "username and password required")
		return
	}

	token, err := h.authSvc.Login(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, token)
}

// RefreshToken refreshes an authentication token
// @Summary  Refresh token
// @Tags     Auth
// @Security Bearer
// @Produce  json
// @Success  200 {object} response.Response{data=service.TokenResponse}
// @Router   /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	tokenStr := extractToken(c)
	if tokenStr == "" {
		response.Unauthorized(c, "missing authentication token")
		return
	}

	token, err := h.authSvc.RefreshToken(tokenStr)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, token)
}

// GetProfile returns the current user's profile
// @Summary  Get current user profile
// @Tags     Auth
// @Security Bearer
// @Produce  json
// @Success  200 {object} response.Response
// @Router   /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	response.Success(c, gin.H{
		"user_id":  userID,
		"username": username,
		"role":     role,
	})
}

// ========================
// JWT Authentication Middleware
// ========================

// AuthMiddleware is the JWT authentication middleware
func AuthMiddleware(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			response.Unauthorized(c, "missing authentication token")
			c.Abort()
			return
		}

		claims, err := authSvc.ValidateToken(tokenStr)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRole is the role-based authorization middleware
func RequireRole(roles ...string) gin.HandlerFunc {
	roleMap := make(map[string]bool, len(roles))
	for _, r := range roles {
		roleMap[r] = true
	}

	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "user role not found")
			c.Abort()
			return
		}

		if !roleMap[role.(string)] {
			response.Forbidden(c, "insufficient permissions")
			c.Abort()
			return
		}
		c.Next()
	}
}

// extractToken extracts the token from the Authorization header
func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}
