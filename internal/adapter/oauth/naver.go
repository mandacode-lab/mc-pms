package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
)

const (
	NaverTokenURL    = "https://nid.naver.com/oauth2.0/token"
	NaverUserInfoURL = "https://openapi.naver.com/v1/nid/me"
	NaverAuthURL     = "https://nid.naver.com/oauth2.0/authorize"
)

type NaverOAuth struct {
	client *http.Client
}

func NewNaverOAuth() out.OAuthAPI {
	return &NaverOAuth{
		client: &http.Client{},
	}
}

func NewNaverOAuthWithClient(client *http.Client) out.OAuthAPI {
	return &NaverOAuth{
		client: client,
	}
}

type NaverTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

type NaverUserInfoResponse struct {
	ResultCode string `json:"resultcode"`
	Message    string `json:"message"`
	Response   struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		Name     string `json:"name"`
		Gender   string `json:"gender"`
		Age      string `json:"age"`
		Birthday string `json:"birthday"`
		Mobile   string `json:"mobile"`
	} `json:"response"`
}

func (n *NaverOAuth) GetAccessToken(ctx context.Context, code string, clientID string, clientSecret []byte, redirectURI string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", string(clientSecret))
	data.Set("code", code)
	data.Set("state", "STATE_STRING") // Naver requires state for token request

	req, err := http.NewRequestWithContext(ctx, "POST", NaverTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := n.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("making token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp NaverTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decoding token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

func (n *NaverOAuth) GetUserInfo(ctx context.Context, accessToken string) (*out.OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", NaverUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating user info request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making user info request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	var naverResp NaverUserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&naverResp); err != nil {
		return nil, fmt.Errorf("decoding user info response: %w", err)
	}

	if naverResp.ResultCode != "00" {
		return nil, fmt.Errorf("naver API error: %s", naverResp.Message)
	}

	rawData, err := json.Marshal(naverResp.Response)
	if err != nil {
		return nil, fmt.Errorf("marshaling raw data: %w", err)
	}

	return &out.OAuthUserInfo{
		ProviderID: naverResp.Response.ID,
		Email:      naverResp.Response.Email,
		Nickname:   naverResp.Response.Nickname,
		RawData:    rawData,
	}, nil
}

func (n *NaverOAuth) GetAuthURL(ctx context.Context, clientID string, scopes []string, redirectURI string, state string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("state", state)

	return NaverAuthURL + "?" + params.Encode()
}

