package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"team_dynamics/match_history_service_v2/models"
	"time"
)

type PageKeyService interface {
	Deserialize(key string) (*models.PageKey, error)
	GetNewPageKey(matches []*models.AggregatedMatch) *string
}

type pageKeyServiceImpl struct{}

func MakePageKeyService() PageKeyService {
	return &pageKeyServiceImpl{}
}

func (s *pageKeyServiceImpl) serialize(pageKey models.PageKey) string {
	data, err := json.Marshal(pageKey)
	if err != nil {
		panic(fmt.Sprintf("Error while serializing page key: %v", err))
	}
	return base64.StdEncoding.EncodeToString(data)
}

func (s *pageKeyServiceImpl) Deserialize(key string) (*models.PageKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("cannot decode b64 of page key: %w", err)
	}
	var res models.PageKey
	if err := json.Unmarshal(bytes, &res); err != nil {
		return nil, fmt.Errorf("cannot parse json of page key: %w", err)
	}
	res.Before = res.Before.UTC()
	return &res, nil
}

func (s *pageKeyServiceImpl) GetNewPageKey(matches []*models.AggregatedMatch) *string {
	if len(matches) != 10 {
		return nil
	}
	var minTime time.Time
	for i, match := range matches {
		end := match.MatchObj.End
		if i == 0 || end.Before(minTime) {
			minTime = end
		}
	}
	res := s.serialize(models.PageKey{Before: minTime})
	return &res
}
