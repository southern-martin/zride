package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/southern-martin/zride/backend/shared"
)

type ZaloOAuthService struct {
	appID     string
	appSecret string
	client    *http.Client
}

type ZaloUserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Birthday string `json:"birthday"`
	Gender   int    `json:"gender"`
}

type ZaloTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func NewZaloOAuthService(appID, appSecret string) *ZaloOAuthService {
	return &ZaloOAuthService{
		appID:     appID,
		appSecret: appSecret,
		client:    &http.Client{},
	}
}

func (z *ZaloOAuthService) ExchangeCodeForToken(code string) (*ZaloTokenResponse, error) {
	tokenURL := "https://oauth.zaloapp.com/v4/access_token"
	
	data := url.Values{}
	data.Set("app_id", z.appID)
	data.Set("app_secret", z.appSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")

	resp, err := z.client.PostForm(tokenURL, data)
	if err != nil {
		return nil, shared.NewExternalServiceError("failed to exchange code for token", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, shared.NewExternalServiceError("failed to read token response", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, shared.NewExternalServiceError(fmt.Sprintf("Zalo API error: %s", string(body)), nil)
	}

	var tokenResp ZaloTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, shared.NewExternalServiceError("failed to parse token response", err)
	}

	return &tokenResp, nil
}

func (z *ZaloOAuthService) GetUserInfo(accessToken string) (*ZaloUserInfo, error) {
	userInfoURL := "https://graph.zalo.me/v2.0/me"
	
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, shared.NewExternalServiceError("failed to create user info request", err)
	}

	// Add access token to header
	req.Header.Set("access_token", accessToken)

	resp, err := z.client.Do(req)
	if err != nil {
		return nil, shared.NewExternalServiceError("failed to get user info", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, shared.NewExternalServiceError("failed to read user info response", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, shared.NewExternalServiceError(fmt.Sprintf("Zalo API error: %s", string(body)), nil)
	}

	var userInfo ZaloUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, shared.NewExternalServiceError("failed to parse user info response", err)
	}

	return &userInfo, nil
}

func (z *ZaloOAuthService) RefreshAccessToken(refreshToken string) (*ZaloTokenResponse, error) {
	refreshURL := "https://oauth.zaloapp.com/v4/access_token"
	
	data := url.Values{}
	data.Set("app_id", z.appID)
	data.Set("app_secret", z.appSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	resp, err := z.client.PostForm(refreshURL, data)
	if err != nil {
		return nil, shared.NewExternalServiceError("failed to refresh access token", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, shared.NewExternalServiceError("failed to read refresh response", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, shared.NewExternalServiceError(fmt.Sprintf("Zalo API error: %s", string(body)), nil)
	}

	var tokenResp ZaloTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, shared.NewExternalServiceError("failed to parse refresh response", err)
	}

	return &tokenResp, nil
}