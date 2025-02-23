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
	} `json:"data"`
}

func NewAuthService(baseURL string) *AuthService {
	return &AuthService{
		baseURL: baseURL,
	}
}

func (s *AuthService) GetUserEmail(userID string) (string, error) {
	url := fmt.Sprintf("%s/api/users/%s", s.baseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
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

	if userResp.Data.Email == "" {
		return "", fmt.Errorf("email not found for user")
	}

	return userResp.Data.Email, nil
}
