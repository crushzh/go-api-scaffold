package service

import (
	"errors"
	"time"

	"go-api-scaffold/internal/model"
	"go-api-scaffold/internal/store"
	"go-api-scaffold/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims holds JWT custom claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenResponse is the token response
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// AuthService handles authentication logic
type AuthService struct {
	db           *store.Store
	jwtSecret    []byte
	expireHours  int
	refreshHours int
}

func NewAuthService(db *store.Store, secret string, expireHours, refreshHours int) *AuthService {
	svc := &AuthService{
		db:           db,
		jwtSecret:    []byte(secret),
		expireHours:  expireHours,
		refreshHours: refreshHours,
	}
	// Ensure default admin account exists
	svc.ensureDefaultAdmin()
	return svc
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(username, password string) (*TokenResponse, error) {
	var user model.User
	if err := s.db.DB().Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return s.generateToken(&user)
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// RefreshToken refreshes a JWT token
func (s *AuthService) RefreshToken(tokenStr string) (*TokenResponse, error) {
	claims, err := s.ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := s.db.DB().First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return s.generateToken(&user)
}

func (s *AuthService) generateToken(user *model.User) (*TokenResponse, error) {
	expiresAt := time.Now().Add(time.Duration(s.expireHours) * time.Hour)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		Token:     tokenStr,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

func (s *AuthService) ensureDefaultAdmin() {
	var count int64
	s.db.DB().Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("failed to create default admin: %v", err)
		return
	}

	admin := &model.User{
		Username: "admin",
		Password: string(hashed),
		Role:     "admin",
	}
	if err := s.db.DB().Create(admin).Error; err != nil {
		logger.Errorf("failed to create default admin: %v", err)
		return
	}
	logger.Info("default admin created: admin / admin123")
}
