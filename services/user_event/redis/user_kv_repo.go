package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"team_dynamics/user_event/models"
)

type UserKvRepo interface {
	Register(ctx context.Context, player *models.PlayerUserData) (bool, error)
	Unregister(ctx context.Context, playerId int64) error
	Subscribe(ctx context.Context, playerId int64, otherPlayersId []int64) error
	NotifyBusy(ctx context.Context, playerId int64) error
	NotifyFree(ctx context.Context, playerId int64) error
	GetPlayerStatus(ctx context.Context, playerId int64) (models.PlayerStatus, error)
	GetSubscriptionsStatus(ctx context.Context, playerId int64) (map[int64]models.PlayerStatus, error)
	CreateChallenge(ctx context.Context, from, to int64) error
	AcceptChallenge(ctx context.Context, to, from int64) error
	DeclineChallenge(ctx context.Context, to, from int64) error
	ClearChallenge(ctx context.Context, from int64) error
}

type userKvRepoImpl struct {
	rdb *redis.Client
}

func MakeUserKvRepo(rdb *redis.Client) UserKvRepo {
	return &userKvRepoImpl{rdb}
}

func (r *userKvRepoImpl) Register(ctx context.Context, player *models.PlayerUserData) (bool, error) {
	keys := PlayerKeySet{player.Id}
	var registered bool
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		n, err := tx.Exists(ctx, keys.name()).Result()
		if err != nil {
			return err
		}
		if n == 1 {
			registered = false
			return nil
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.MSet(ctx, []interface{}{
				keys.name(), player.Name,
				keys.rating(), player.Rating,
				keys.status(), int64(models.Online),
			})
			return nil
		})
		if err != nil {
			return err
		}
		registered = true
		return nil
	}, keys.name())
	return registered, err
}

func (r *userKvRepoImpl) Subscribe(ctx context.Context, playerId int64, otherPlayersId []int64) error {
	keys := PlayerKeySet{playerId}
	members := make([]interface{}, len(otherPlayersId))
	for i, id := range otherPlayersId {
		members[i] = id
	}
	return r.rdb.SAdd(ctx, keys.subscriptions(), members...).Err()
}

func (r *userKvRepoImpl) GetPlayerStatus(ctx context.Context, playerId int64) (models.PlayerStatus, error) {
	keys := PlayerKeySet{playerId}
	val, err := r.rdb.Get(ctx, keys.status()).Int64()
	if err != nil {
		if err == redis.Nil {
			return models.Offline, nil
		}
		return 0, err
	}
	return models.PlayerStatus(val), nil
}



func (r *userKvRepoImpl) GetSubscriptionsStatus(ctx context.Context, playerId int64) (map[int64]models.PlayerStatus, error) {
	keys := PlayerKeySet{playerId}

	ids, err := r.rdb.SMembers(ctx, keys.subscriptions()).Result()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return map[int64]models.PlayerStatus{}, nil
	}

	statusKeys := make([]string, len(ids))
	parsedIds := make([]int64, len(ids))
	for i, idStr := range ids {
		var id int64
		fmt.Sscanf(idStr, "%d", &id)
		parsedIds[i] = id
		statusKeys[i] = PlayerKeySet{id}.status()
	}

	vals, err := r.rdb.MGet(ctx, statusKeys...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[int64]models.PlayerStatus, len(parsedIds))
	for i, val := range vals {
		if val == nil {
			result[parsedIds[i]] = models.Offline
			continue
		}
		var status int64
		fmt.Sscanf(val.(string), "%d", &status)
		result[parsedIds[i]] = models.PlayerStatus(status)
	}
	return result, nil
}

func (r *userKvRepoImpl) NotifyBusy(ctx context.Context, playerId int64) error {
	keys := PlayerKeySet{playerId}
	set, err := r.rdb.SetXX(ctx, keys.status(), int64(models.Busy), 0).Result()
	if err != nil {
		return err
	}
	if !set {
		return redis.Nil
	}
	return nil
}

func (r *userKvRepoImpl) NotifyFree(ctx context.Context, playerId int64) error {
	keys := PlayerKeySet{playerId}
	set, err := r.rdb.SetXX(ctx, keys.status(), int64(models.Online), 0).Result()
	if err != nil {
		return err
	}
	if !set {
		return redis.Nil
	}
	return nil
}

func (r *userKvRepoImpl) CreateChallenge(ctx context.Context, from, to int64) error {
	fromKeys := PlayerKeySet{from}
	set, err := r.rdb.SetNX(ctx, fromKeys.currentChallenge(), to, 0).Result()
	if err != nil {
		return err
	}
	if !set {
		return errors.New("challenge already exists")
	}
	return nil
}

func (r *userKvRepoImpl) AcceptChallenge(ctx context.Context, to, from int64) error {
	fromKeys := PlayerKeySet{from}
	toKeys := PlayerKeySet{to}
	return r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		val, err := tx.Get(ctx, fromKeys.currentChallenge()).Int64()
		if err != nil {
			return err
		}
		if val != to {
			return errors.New("challenge mismatch")
		}
		set, err := tx.SetNX(ctx, toKeys.currentChallenge(), from, 0).Result()
		if err != nil {
			return err
		}
		if !set {
			return errors.New("to player already has an outgoing challenge")
		}
		return nil
	}, fromKeys.currentChallenge())
}

func (r *userKvRepoImpl) DeclineChallenge(ctx context.Context, to, from int64) error {
	fromKeys := PlayerKeySet{from}
	return r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		val, err := tx.Get(ctx, fromKeys.currentChallenge()).Int64()
		if err != nil {
			return err
		}
		if val != to {
			return errors.New("challenge mismatch")
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, fromKeys.currentChallenge())
			return nil
		})
		return err
	}, fromKeys.currentChallenge())
}

func (r *userKvRepoImpl) ClearChallenge(ctx context.Context, from int64) error {
	fromKeys := PlayerKeySet{from}
	return r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		to, err := tx.Get(ctx, fromKeys.currentChallenge()).Int64()
		if err != nil {
			if err == redis.Nil {
				return errors.New("no active challenge")
			}
			return err
		}
		toKeys := PlayerKeySet{to}
		toChallenge, err := tx.Get(ctx, toKeys.currentChallenge()).Int64()
		if err != nil && err != redis.Nil {
			return err
		}
		if err == nil && toChallenge == from {
			return errors.New("challenge already accepted")
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, fromKeys.currentChallenge())
			return nil
		})
		return err
	}, fromKeys.currentChallenge())
}

func (r *userKvRepoImpl) Unregister(ctx context.Context, playerId int64) error {
	keys := PlayerKeySet{playerId}
	return r.rdb.Del(ctx, keys.keys()...).Err()
}
