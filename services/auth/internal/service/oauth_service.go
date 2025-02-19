// services/auth/internal/service/oauth_service.go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kevinnaserwan/crm-be/services/auth/internal/config"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
)

type OAuthService struct {
	cfg config.OAuthConfig
}

func NewOAuthService(cfg config.OAuthConfig) *OAuthService {
	return &OAuthService{
		cfg: cfg,
	}
}

// UserInfo contains profile data returned by OAuth providers
type UserInfo struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Picture   string
	Provider  model.OAuthProvider
}

// ExchangeGoogleToken exchanges authorization code for tokens and user info
func (s *OAuthService) ExchangeGoogleToken(ctx context.Context, code string, redirectURI string) (*UserInfo, string, string, time.Time, error) {
	// Exchange authorization code for access token
	tokenURL := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", s.cfg.GoogleClientID)
	data.Set("client_secret", s.cfg.GoogleClientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, "", "", time.Time{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", "", time.Time{}, fmt.Errorf("failed to exchange token: %s", body)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		IDToken      string `json:"id_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, "", "", time.Time{}, err
	}

	// Get user info using access token
	userInfo, err := s.getGoogleUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return nil, "", "", time.Time{}, err
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(tokenResp.ExpiresIn))
	return userInfo, tokenResp.AccessToken, tokenResp.RefreshToken, expiresAt, nil
}

func (s *OAuthService) getGoogleUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	userInfoURL := "https://www.googleapis.com/oauth2/v3/userinfo"
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info")
	}

	var userInfoResp struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfoResp); err != nil {
		return nil, err
	}

	if !userInfoResp.EmailVerified {
		return nil, errors.New("email not verified")
	}

	return &UserInfo{
		ID:        userInfoResp.Sub,
		Email:     userInfoResp.Email,
		FirstName: userInfoResp.GivenName,
		LastName:  userInfoResp.FamilyName,
		Picture:   userInfoResp.Picture,
		Provider:  model.GoogleProvider,
	}, nil
}

// VerifyAppleToken verifies and extracts information from an Apple ID token
// Implementation depends on Apple's documentation and requirements
func (s *OAuthService) VerifyAppleToken(ctx context.Context, idToken string) (*UserInfo, error) {
	// For Apple sign-in, you'd need to verify the JWT token
	// This is a placeholder - actual implementation requires Apple-specific JWT validation

	// TODO: Implement Apple token verification
	return nil, errors.New("apple sign-in not implemented yet")
}
