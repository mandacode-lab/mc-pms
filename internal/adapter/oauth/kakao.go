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
	KakaoTokenURL    = "https://kauth.kakao.com/oauth/token"
	KakaoUserInfoURL = "https://kapi.kakao.com/v2/user/me"
	KakaoAuthURL     = "https://kauth.kakao.com/oauth/authorize"
)

type KakaoOAuth struct {
	client *http.Client
}

func NewKakaoOAuth() out.OAuthAPI {
	return &KakaoOAuth{
		client: &http.Client{},
	}
}

func NewKakaoOAuthWithClient(client *http.Client) out.OAuthAPI {
	return &KakaoOAuth{
		client: client,
	}
}

type KakaoTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type KakaoUserInfo struct {
	ID           int64  `json:"id"`
	ConnectedAt  string `json:"connected_at"`
	KakaoAccount struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"is_email_verified"`
		EmailValid    bool   `json:"is_email_valid"`
		Profile       struct {
			Nickname        string `json:"nickname"`
			ProfileImageURL string `json:"profile_image_url"`
			ThumbnailURL    string `json:"thumbnail_image_url"`
		} `json:"profile"`
	} `json:"kakao_account"`
}

func (k *KakaoOAuth) GetAccessToken(ctx context.Context, code string, clientID string, clientSecret []byte, redirectURI string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", string(clientSecret))
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)

	req, err := http.NewRequestWithContext(ctx, "POST", KakaoTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := k.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("making token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp KakaoTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decoding token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

func (k *KakaoOAuth) GetUserInfo(ctx context.Context, accessToken string) (*out.OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", KakaoUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating user info request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making user info request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	var kakaoUser KakaoUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&kakaoUser); err != nil {
		return nil, fmt.Errorf("decoding user info response: %w", err)
	}

	rawData, err := json.Marshal(kakaoUser)
	if err != nil {
		return nil, fmt.Errorf("marshaling raw data: %w", err)
	}

	return &out.OAuthUserInfo{
		ProviderID: fmt.Sprintf("%d", kakaoUser.ID),
		Email:      kakaoUser.KakaoAccount.Email,
		Nickname:   kakaoUser.KakaoAccount.Profile.Nickname,
		RawData:    rawData,
	}, nil
}

func (k *KakaoOAuth) GetAuthURL(ctx context.Context, clientID string, scopes []string, redirectURI string, state string) string {
	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", strings.Join(scopes, ","))
	params.Set("state", state)

	return KakaoAuthURL + "?" + params.Encode()
}
