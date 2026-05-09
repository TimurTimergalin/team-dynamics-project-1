package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SteamService interface {
	Validate(ctx context.Context, authToken string) (int64, error)
}

type steamServiceImpl struct {
	apiKey     string
	appId      string
	httpClient *http.Client
}

func NewSteamService(apiKey string, appId string) SteamService {
	return &steamServiceImpl{
		apiKey: apiKey,
		appId:  appId,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *steamServiceImpl) Validate(ctx context.Context, authToken string) (int64, error) {
	url := fmt.Sprintf(
		"https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v1/?key=%s&appid=%s&ticket=%s",
		s.apiKey, s.appId, authToken,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Steam API returned status %d", resp.StatusCode)
	}

	var apiResponse struct {
		Response struct {
			Params struct {
				Result          string `json:"result"`
				SteamId         string `json:"steamid"`
				OwnerSteamId    string `json:"ownersteamid"`
				VacBanned       bool   `json:"vacbanned"`
				PublisherBanned bool   `json:"publisherbanned"`
			} `json:"params"`
			Error *struct {
				ErrorCode int    `json:"errorcode"`
				ErrorDesc string `json:"errordesc"`
			} `json:"error"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResponse.Response.Error != nil {
		return 0, fmt.Errorf("Steam API error %d: %s",
			apiResponse.Response.Error.ErrorCode,
			apiResponse.Response.Error.ErrorDesc,
		)
	}

	if apiResponse.Response.Params.Result != "OK" {
		return 0, fmt.Errorf("Steam auth failed: result=%s", apiResponse.Response.Params.Result)
	}

	var steamId int64
	if _, err := fmt.Sscanf(apiResponse.Response.Params.SteamId, "%d", &steamId); err != nil {
		return 0, fmt.Errorf("failed to parse steam id: %w", err)
	}

	return steamId, nil
}
