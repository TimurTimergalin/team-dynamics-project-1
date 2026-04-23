package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"team_dynamics/user_service/models"
	"time"
)

type SteamService interface {
	GetUserSummary(ctx context.Context, steamID string) (*models.SteamData, error)
}

type steamServiceImpl struct {
	apiKey     string
	httpClient *http.Client
}

// NewSteamService creates a new Steam service with the provided API key.
func NewSteamService(apiKey string) SteamService {
	return &steamServiceImpl{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetUserSummary fetches a user's Steam profile summary by their Steam ID.
func (s *steamServiceImpl) GetUserSummary(ctx context.Context, steamID string) (*models.SteamData, error) {
	url := fmt.Sprintf(
		"https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s",
		s.apiKey,
		steamID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Steam API returned status %d", resp.StatusCode)
	}

	var apiResponse struct {
		Response struct {
			Players []struct {
				SteamID     string `json:"steamid"`
				PersonaName string `json:"personaname"`
			} `json:"players"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(apiResponse.Response.Players) == 0 {
		return nil, fmt.Errorf("no player data found for steam ID %s", steamID)
	}

	player := apiResponse.Response.Players[0]

	return &models.SteamData{
		SteamId: player.SteamID,
		Name:    player.PersonaName,
	}, nil
}
