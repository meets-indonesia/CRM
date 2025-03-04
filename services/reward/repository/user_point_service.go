package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinnaserwan/crm-be/services/reward/config"
)

// HTTPUserPointService implements UserPointService with HTTP requests
type HTTPUserPointService struct {
	userServiceURL string
	client         *http.Client
}

// NewHTTPUserPointService creates a new HTTPUserPointService
func NewHTTPUserPointService(cfg config.ServiceConfig) *HTTPUserPointService {
	return &HTTPUserPointService{
		userServiceURL: cfg.UserServiceURL,
		client:         &http.Client{},
	}
}

// CheckUserPoints checks the points of a user
func (s *HTTPUserPointService) CheckUserPoints(ctx context.Context, userID uint) (int, error) {
	url := fmt.Sprintf("%s/users/customer/%d/points", s.userServiceURL, userID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	// Add headers if needed (e.g., for service-to-service authentication)

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to check user points: status code %d", resp.StatusCode)
	}

	var response struct {
		UserID    uint   `json:"user_id"`
		Total     int    `json:"total"`
		Level     string `json:"level"`
		NextLevel string `json:"next_level,omitempty"`
		ToNext    int    `json:"to_next,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	return response.Total, nil
}
