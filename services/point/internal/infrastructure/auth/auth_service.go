package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthService struct {
	baseURL string
}

type UserResponse struct {
	Data struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"data"`
}

func NewAuthService(baseURL string) *AuthService {
	return &AuthService{
		baseURL: baseURL,
	}
}

func (s *AuthService) GetUserEmail(userID string) (string, error) {
	url := fmt.Sprintf("%s/api/users/%s", s.baseURL, userID)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to call auth service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return userResp.Data.Email, nil
}
