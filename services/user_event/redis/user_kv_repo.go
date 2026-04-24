package redis

import (
	"github.com/redis/go-redis/v9"
	"team_dynamics/user_event/models"
)

type UserKvRepo interface {
	Register(player *models.PlayerUserData) error
	Unregister(playerId int64) error
	Subscribe(playerId int64, otherPlayersId []int64) error
	NotifyBusy(playerId int64) error
	NotifyFree(playerId int64) error
	GetPlayerStatus(playerId int64) error
	CreateChallenge(from, to int64) error
	AcceptChallenge(to, from int64) error
	DeclineChallenge(to, from int64) error
	ClearChallenge(from int64) error
}

type userKvRepoImpl struct {
	rdb *redis.Client
}
