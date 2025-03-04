package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kevinnaserwan/crm-be/services/auth/config"
	"github.com/kevinnaserwan/crm-be/services/auth/domain/entity"
)

// JWTService implements TokenGenerator
type JWTService struct {
	jwtConfig config.JWTConfig
}

// Claims adalah struktur untuk JWT claims
type Claims struct {
	UserID uint        `json:"user_id"`
	Email  string      `json:"email"`
	Role   entity.Role `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWTService
func NewJWTService(jwtConfig config.JWTConfig) *JWTService {
	return &JWTService{
		jwtConfig: jwtConfig,
	}
}

// GenerateToken generates a JWT token
func (s *JWTService) GenerateToken(userID uint, email string, role entity.Role) (string, time.Time, error) {
	// Set expiration time
	expirationTime := time.Now().Add(s.jwtConfig.Expiry)

	// Create claims
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "lrt-crm",
			Subject:   email,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ValidateToken validates a JWT token
func (s *JWTService) ValidateToken(tokenString string) (uint, string, entity.Role, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		return 0, "", "", err
	}

	// Validate claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, "", "", errors.New("invalid token")
	}

	return claims.UserID, claims.Email, claims.Role, nil
}
