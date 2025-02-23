package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"` // Pastikan order field sesuai
	Role   string    `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(userID uuid.UUID, email string, role string, secret string, duration time.Duration) (string, error) {
	fmt.Printf("GenerateJWT input - userID: %v, email: %v, role: %v\n", userID, email, role)

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	fmt.Printf("Claims created: %+v\n", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	fmt.Printf("Validated token claims: %+v\n", claims)

	return claims, nil
}
