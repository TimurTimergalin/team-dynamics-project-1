package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"team_dynamics/user_service/models"
)

type PageKeyService interface {
	Deserialize(key string) (*models.PageKey, error)
	GetNewPageKey(lastPage []*models.Friend) *string
}

type pageKeyServiceImpl struct {
}

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
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot parse json of page key: %w", err)
	}
	return &res, nil
}

func (s *pageKeyServiceImpl) GetNewPageKey(lastPage []*models.Friend) *string {
	if len(lastPage) != 20 {
		return nil
	}
	lastId := lastPage[len(lastPage)-1].Data.Id
	pageKey := models.PageKey{LastUserId: lastId}
	res := s.serialize(pageKey)
	return &res
}

func (s *pageKeyServiceImpl) GetInitialPageKey() *models.PageKey {
	return &models.PageKey{LastUserId: 0}
}
