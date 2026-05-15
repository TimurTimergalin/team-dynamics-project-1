package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type EosService interface {
	Validate(ctx context.Context, authToken string) (string, error)
}

type eosServiceImpl struct {
	clientId     string
	clientSecret string
	httpClient   *http.Client
}

func NewEosService(clientId string, clientSecret string) EosService {
	return &eosServiceImpl{
		clientId:     clientId,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (e *eosServiceImpl) Validate(ctx context.Context, authToken string) (string, error) {
	body := url.Values{}
	body.Set("grant_type", "external_auth")
	body.Set("external_auth_type", "epicgames_access_token")
	body.Set("external_auth_token", authToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.epicgames.dev/auth/v1/oauth/token",
		strings.NewReader(body.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(e.clientId, e.clientSecret)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("EOS API returned status %d", resp.StatusCode)
	}

	var apiResponse struct {
		AccountId string `json:"account_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	if apiResponse.AccountId == "" {
		return "", fmt.Errorf("EOS token validation returned empty account_id")
	}

	return apiResponse.AccountId, nil
}
