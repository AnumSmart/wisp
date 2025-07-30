package jwt_stuff

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Конфигурация JWT
type JWT struct {
	SecretAccKey    string
	SecretRefKey    string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

// Claims для JWT
type CustomClaims struct {
	Email     string `json:"email"`
	TokenType string `json:"type"` // "access" или "refresh"
	Role      string `json:"role"`
	UserId    string `json:"user_id"` // userID для извлечения из JWT токена
	IsActive  bool   `json:"is_active"`
	jwt.RegisteredClaims
}
