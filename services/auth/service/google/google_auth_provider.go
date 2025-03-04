package service

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/auth/config"
	"google.golang.org/api/idtoken"
)

// GoogleAuthService implements GoogleAuthProvider
type GoogleAuthService struct {
	googleConfig config.GoogleAuthConfig
}

// NewGoogleAuthService creates a new GoogleAuthService
func NewGoogleAuthService(googleConfig config.GoogleAuthConfig) *GoogleAuthService {
	return &GoogleAuthService{
		googleConfig: googleConfig,
	}
}

// VerifyIDToken verifies a Google ID token
func (s *GoogleAuthService) VerifyIDToken(idToken string) (string, string, string, error) {
	// Verify the ID token
	payload, err := idtoken.Validate(context.Background(), idToken, s.googleConfig.ClientID)
	if err != nil {
		return "", "", "", err
	}

	// Extract user info from the payload
	googleID, ok := payload.Claims["sub"].(string)
	if !ok {
		return "", "", "", errors.New("invalid Google ID")
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return "", "", "", errors.New("invalid email")
	}

	name, ok := payload.Claims["name"].(string)
	if !ok {
		name = "Google User" // Default name if not provided
	}

	return googleID, email, name, nil
}
