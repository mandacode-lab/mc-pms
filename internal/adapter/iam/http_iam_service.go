package iam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type HTTPIAMService struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPIAMService(baseURL string) (*HTTPIAMService, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL cannot be empty")
	}

	// Ping the IAM service to ensure it's reachable
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("%s/health", baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to reach IAM service at %s: %w", baseURL, err)
	}
	defer func() {
		_ = resp.Body.Close() // Best effort close
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IAM service health check failed with status code: %d", resp.StatusCode)
	}

	// If we reach here, the IAM service is reachable and healthy
	return &HTTPIAMService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

type PermissionCheckRequest struct {
	UserID   string         `json:"user_id"`
	Resource string         `json:"resource"`
	Action   string         `json:"action"`
	Context  map[string]any `json:"context,omitempty"`
}

type PermissionCheckResponse struct {
	Allowed bool   `json:"allowed"`
	Version string `json:"version"`
}

func (s *HTTPIAMService) HasPermission(ctx context.Context, userID string, permission out.Permission) (bool, error) {
	// Prepare request payload
	req := PermissionCheckRequest{
		UserID:   userID,
		Resource: string(permission.Resource),
		Action:   string(permission.Action),
		Context:  permission.Context,
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return false, fmt.Errorf("failed to marshal permission check request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/v1/permissions/check", s.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute HTTP request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return false, fmt.Errorf("failed to execute IAM request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close() // Best effort close
	}()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var permResp PermissionCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&permResp); err != nil {
		return false, fmt.Errorf("failed to decode permission check response: %w", err)
	}

	return permResp.Allowed, nil
}

func (s *HTTPIAMService) CheckHealth(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", s.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute health check request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close() // Best effort close
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("IAM service unhealthy, status code: %d", resp.StatusCode)
	}

	return nil
}
