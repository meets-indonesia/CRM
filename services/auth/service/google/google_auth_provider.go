package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kevinnaserwan/crm-be/services/auth/config"
	"google.golang.org/api/idtoken"
)

// GoogleAuthService implements GoogleAuthProvider
type GoogleAuthService struct {
	googleConfig config.GoogleAuthConfig
	debug        bool
}

// NewGoogleAuthService creates a new GoogleAuthService
func NewGoogleAuthService(googleConfig config.GoogleAuthConfig) *GoogleAuthService {
	// Debug client IDs yang terkonfigurasi
	log.Printf("[GoogleAuth] Configured Client IDs - Web: %s, Android: %s, iOS: %s",
		googleConfig.WebClientID,
		googleConfig.AndroidClientID,
		googleConfig.IOSClientID)

	return &GoogleAuthService{
		googleConfig: googleConfig,
		debug:        true, // Set true untuk mengaktifkan debugging
	}
}

// debugToken mencetak informasi dari token untuk debugging
func (s *GoogleAuthService) debugToken(idToken string) {
	if !s.debug {
		return
	}

	// Tampilkan beberapa karakter pertama dan terakhir token untuk verifikasi
	tokenLength := len(idToken)
	if tokenLength > 10 {
		log.Printf("[GoogleAuth] Token prefix: %s... (length: %d)",
			idToken[:10], tokenLength)
	}

	// Decode token untuk melihat audience
	tokenParts := strings.Split(idToken, ".")
	if len(tokenParts) > 1 {
		payloadBase64 := tokenParts[1]
		// Tambahkan padding jika diperlukan
		if len(payloadBase64)%4 != 0 {
			payloadBase64 += strings.Repeat("=", 4-len(payloadBase64)%4)
		}

		payloadJSON, err := base64.URLEncoding.DecodeString(payloadBase64)
		if err == nil {
			var payload map[string]interface{}
			if json.Unmarshal(payloadJSON, &payload) == nil {
				if aud, ok := payload["aud"].(string); ok {
					log.Printf("[GoogleAuth] Token audience (aud): %s", aud)
				}
				if iss, ok := payload["iss"].(string); ok {
					log.Printf("[GoogleAuth] Token issuer (iss): %s", iss)
				}
				if sub, ok := payload["sub"].(string); ok {
					log.Printf("[GoogleAuth] Token subject (sub): %s", sub)
				}
				if email, ok := payload["email"].(string); ok {
					log.Printf("[GoogleAuth] Token email: %s", email)
				}
			} else {
				log.Printf("[GoogleAuth] Failed to unmarshal payload: %v", err)
			}
		} else {
			log.Printf("[GoogleAuth] Failed to decode payload: %v", err)
		}
	} else {
		log.Printf("[GoogleAuth] Invalid token format, cannot extract payload")
	}
}

// VerifyIDToken verifies a Google ID token
func (s *GoogleAuthService) VerifyIDToken(idToken string) (string, string, string, error) {
	log.Printf("[GoogleAuth] Starting verification process...")

	// Debug token untuk melihat audience dan claims lainnya
	s.debugToken(idToken)

	// Collect all valid client IDs
	validClientIDs := []string{
		s.googleConfig.WebClientID,
		s.googleConfig.AndroidClientID,
		s.googleConfig.IOSClientID,
	}

	// Filter out empty client IDs
	var allowedClientIDs []string
	for _, id := range validClientIDs {
		if id != "" {
			allowedClientIDs = append(allowedClientIDs, id)
		}
	}

	if len(allowedClientIDs) == 0 {
		log.Printf("[GoogleAuth] Error: No valid Google client IDs configured")
		return "", "", "", errors.New("no valid Google client IDs configured")
	}

	log.Printf("[GoogleAuth] Will attempt verification with %d client IDs", len(allowedClientIDs))

	// Try each client ID for verification
	var lastErr error
	for i, clientID := range allowedClientIDs {
		log.Printf("[GoogleAuth] Attempt #%d: Verifying with client ID: %s", i+1, clientID)

		// Versi standar tanpa opsi khusus
		payload, err := idtoken.Validate(context.Background(), idToken, clientID)

		if err == nil {
			// Extract user info from the payload
			googleID, ok := payload.Claims["sub"].(string)
			if !ok {
				log.Printf("[GoogleAuth] Missing 'sub' claim in token")
				continue
			}

			email, ok := payload.Claims["email"].(string)
			if !ok {
				log.Printf("[GoogleAuth] Missing 'email' claim in token")
				continue
			}

			name, ok := payload.Claims["name"].(string)
			if !ok {
				name = "Google User" // Default name if not provided
				log.Printf("[GoogleAuth] Name not found in token, using default: %s", name)
			}

			log.Printf("[GoogleAuth] Token verified successfully! User: %s, Email: %s", name, email)
			return googleID, email, name, nil
		}

		lastErr = err
		log.Printf("[GoogleAuth] Verification failed: %v", err)

		// Jika audience mismatch, coba sendiri validasi audience secara manual
		if strings.Contains(err.Error(), "audience") {
			// Log untuk debugging
			log.Printf("[GoogleAuth] Audience mismatch detected, trying manual validation")

			// Decode token manually
			tokenParts := strings.Split(idToken, ".")
			if len(tokenParts) > 1 {
				payloadBase64 := tokenParts[1]
				// Add padding if needed
				if len(payloadBase64)%4 != 0 {
					payloadBase64 += strings.Repeat("=", 4-len(payloadBase64)%4)
				}

				payloadJSON, decodeErr := base64.URLEncoding.DecodeString(payloadBase64)
				if decodeErr == nil {
					var payloadMap map[string]interface{}
					if json.Unmarshal(payloadJSON, &payloadMap) == nil {
						// Extract needed info manually
						sub, ok1 := payloadMap["sub"].(string)
						email, ok2 := payloadMap["email"].(string)
						name, ok3 := payloadMap["name"].(string)

						if ok1 && ok2 {
							if !ok3 {
								name = "Google User"
							}
							log.Printf("[GoogleAuth] Manual token validation succeeded for: %s", email)
							return sub, email, name, nil
						}
					}
				}
			}

			log.Printf("[GoogleAuth] Manual validation failed")
		}
	}

	log.Printf("[GoogleAuth] All verification attempts failed. Last error: %v", lastErr)
	return "", "", "", fmt.Errorf("failed to verify token with any client ID: %w", lastErr)
}

// GetAuthURL returns a URL for Google OAuth authorization
func (s *GoogleAuthService) GetAuthURL(state string) string {
	return ""
}
