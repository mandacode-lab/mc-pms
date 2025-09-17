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
	baseURL      string
	httpClient   *http.Client
	clientAppID  string
	clientSecret string
}

func NewHTTPIAMService(baseURL, clientAppID, clientSecret string) *HTTPIAMService {
	return &HTTPIAMService{
		baseURL:      baseURL,
		clientAppID:  clientAppID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type PermissionCheckRequest struct {
	UserID   string         `json:"user_id"`
	Resource string         `json:"resource"`
	Action   string         `json:"action"`
	Context  map[string]any `json:"context,omitempty"`
}

type PermissionCheckResponse struct {
	Allowed bool   `json:"allowed"`
	Action  string `json:"action"`
	Version string `json:"version"`
}

func (s *HTTPIAMService) HasPermission(ctx context.Context, userID string, permission out.Permission) (out.PermissionAllowInfo, error) {
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
		return out.PermissionAllowInfo{}, fmt.Errorf("failed to marshal permission check request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/v1/permissions/check", s.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return out.PermissionAllowInfo{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Set Basic Authentication using ClientAppID and ClientSecret
	httpReq.SetBasicAuth(s.clientAppID, s.clientSecret)

	// Execute HTTP request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return out.PermissionAllowInfo{}, fmt.Errorf("failed to execute IAM request: %w", err)
	}
	defer resp.Body.Close()

	// Handle HTTP status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Permission check successful, parse response
		var permResp PermissionCheckResponse
		if err := json.NewDecoder(resp.Body).Decode(&permResp); err != nil {
			return out.PermissionAllowInfo{}, fmt.Errorf("failed to decode permission check response: %w", err)
		}

		return out.PermissionAllowInfo{
			Allowed: permResp.Allowed,
			Action:  out.Action(permResp.Action),
			Version: permResp.Version,
		}, nil

	case http.StatusUnauthorized:
		return out.PermissionAllowInfo{
			Allowed: false,
			Action:  permission.Action,
			Version: "",
		}, nil

	case http.StatusForbidden:
		return out.PermissionAllowInfo{
			Allowed: false,
			Action:  permission.Action,
			Version: "",
		}, nil

	case http.StatusBadRequest:
		return out.PermissionAllowInfo{}, fmt.Errorf("invalid permission check request")

	case http.StatusInternalServerError:
		return out.PermissionAllowInfo{}, fmt.Errorf("IAM service internal error")

	default:
		return out.PermissionAllowInfo{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

