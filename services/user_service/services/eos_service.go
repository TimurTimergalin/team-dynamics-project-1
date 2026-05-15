package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"team_dynamics/user_service/models"
	"time"
)

type EosService interface {
	GetUserSummary(ctx context.Context, eosId string) (*models.EosData, error)
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

func (e *eosServiceImpl) GetUserSummary(ctx context.Context, eosId string) (*models.EosData, error) {
	url := fmt.Sprintf("https://api.epicgames.dev/user/v1/accounts?accountId=%s", eosId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(e.clientId, e.clientSecret)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("EOS API returned status %d", resp.StatusCode)
	}

	var apiResponse []struct {
		AccountId   string `json:"accountId"`
		DisplayName string `json:"displayName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(apiResponse) == 0 {
		return nil, fmt.Errorf("no account data found for eos id %s", eosId)
	}

	return &models.EosData{
		EosId: apiResponse[0].AccountId,
		Name:  apiResponse[0].DisplayName,
	}, nil
}
