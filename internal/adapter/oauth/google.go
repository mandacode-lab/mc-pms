package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mandacode-com/serengeti/internal/port/out"
)

const (
	GoogleTokenURL    = "https://oauth2.googleapis.com/token"
	GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	GoogleAuthURL     = "https://accounts.google.com/o/oauth2/auth"
)

type GoogleOAuth struct {
	client *http.Client
}

func NewGoogleOAuth() out.OAuthAPI {
	return &GoogleOAuth{
		client: &http.Client{},
	}
}

func NewGoogleOAuthWithClient(client *http.Client) out.OAuthAPI {
	return &GoogleOAuth{
		client: client,
	}
}

type GoogleTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (g *GoogleOAuth) GetAccessToken(ctx context.Context, code string, clientID string, clientSecret []byte, redirectURI string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", string(clientSecret))
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", GoogleTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("making token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decoding token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

func (g *GoogleOAuth) GetUserInfo(ctx context.Context, accessToken string) (*out.OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", GoogleUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating user info request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making user info request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("decoding user info response: %w", err)
	}

	rawData, err := json.Marshal(googleUser)
	if err != nil {
		return nil, fmt.Errorf("marshaling raw data: %w", err)
	}

	return &out.OAuthUserInfo{
		ProviderID: googleUser.ID,
		Email:      googleUser.Email,
		Nickname:   googleUser.Name,
		RawData:    rawData,
	}, nil
}

func (g *GoogleOAuth) GetAuthURL(ctx context.Context, clientID string, scopes []string, redirectURI string, state string) string {
	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", strings.Join(scopes, " "))
	params.Set("state", state)

	return GoogleAuthURL + "?" + params.Encode()
}
