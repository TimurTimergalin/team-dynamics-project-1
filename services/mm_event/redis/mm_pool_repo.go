package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"team_dynamics/mm_event/models"
	"time"
)

type MMPoolRepo interface {
	AddConnection(ctx context.Context, playerId int64, ttl time.Duration) (bool, error)
	RemoveConnection(ctx context.Context, playerId int64) error
	AddToPool(ctx context.Context, player *models.Player) error
	RemoveFromPool(ctx context.Context, playerId int64) error
	CheckInPool(ctx context.Context, playerId int64) (bool, error)
}

const acquireScriptCode = `
local LOCK_KEY = KEYS[1] -- "mmlock"
local REMOVERS_COUNT = KEYS[2] -- "removers_count"
local MMEVENT_OWNER_ID = "mmevent"

local lock_ttl_millis = tonumber(ARGV[1])

local prev_owner_id = redis.call("GET", LOCK_KEY)
if prev_owner_id and prev_owner_id != MMEVENT_OWNER_ID then
    return "0"
end
redis.call("SET", LOCK_KEY, MMEVENT_OWNER_ID, "PX", lock_ttl_millis)
redis.call("INCR", REMOVERS_COUNT)
return "1"
`

const releaseScriptCode = `
local LOCK_KEY = KEYS[1] -- "mmlock"
local REMOVERS_COUNT = KEYS[2] -- "removers_count"
local MMEVENT_OWNER_ID = "mmevent"

if redis.call("GET", LOCK_KEY) ~= MMEVENT_OWNER_ID then
    return "0"
end

local removers_count = redis.call("GET", REMOVERS_COUNT)
if not removers_count or tonumber(removers_count) == 0 then
    return "0"
end

removers_count = redis.call("DECR", REMOVERS_COUNT)
if removers_count == 0 then
    redis.call("DEL", LOCK_KEY)
end
return "1"
`

const (
	LockKey          = "mmlock"
	PoolKey          = "mmpool"
	RemoversCountKey = "removers_count"
)

type MMPoolRepoConfig struct {
	LockTTL time.Duration
}

type mmPoolRepoImpl struct {
	rdb           *redis.Client
	acquireScript *redis.Script
	releaseScript *redis.Script
	config        *MMPoolRepoConfig
}

func MakeMMPoolRepo(rdb *redis.Client, config *MMPoolRepoConfig) MMPoolRepo {
	return &mmPoolRepoImpl{
		rdb:           rdb,
		acquireScript: redis.NewScript(acquireScriptCode),
		releaseScript: redis.NewScript(releaseScriptCode),
		config:        config,
	}
}

func (m *mmPoolRepoImpl) AddConnection(ctx context.Context, playerId int64, ttl time.Duration) (bool, error) {
	pKeys := playerKeys{playerId}
	added, err := m.rdb.SetNX(ctx, pKeys.connection(), "1", ttl).Result()
	if err != nil {
		return false, err
	}
	if !added {
		_ = m.rdb.Expire(ctx, pKeys.connection(), ttl)
	}
	return added, nil
}

func (m *mmPoolRepoImpl) RemoveConnection(ctx context.Context, playerId int64) error {
	pKeys := playerKeys{playerId}
	removed, err := m.rdb.Del(ctx, pKeys.connection()).Result()
	if err != nil {
		return err
	}
	if removed == 0 {
		return errors.New("PlayerId was removed already")
	}
	return nil
}

func (m *mmPoolRepoImpl) AddToPool(ctx context.Context, player *models.Player) error {
	pKeys := playerKeys{player.Id}
	return m.rdb.Watch(ctx, func(tx *redis.Tx) error {
		n, err := m.rdb.Exists(ctx, pKeys.rating()).Result()
		if err != nil {
			return err
		}
		if n == 1 {
			return nil
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.MSet(ctx, []interface{}{
				pKeys.rating(), player.Rating,
				pKeys.fleet(), player.Fleet,
				pKeys.name(), player.Name,
				pKeys.regId(), player.RegId,
				pKeys.displayedRating(), player.DisplayedRating,
			})
			pipe.RPush(ctx, PoolKey, player.Id)
			return nil
		})
		return err
	}, pKeys.rating())
}

func (m *mmPoolRepoImpl) RemoveFromPool(ctx context.Context, playerId int64) (resErr error) {
	keys := []string{LockKey, RemoversCountKey}
	argv := []interface{}{m.config.LockTTL / time.Millisecond}
	aq, err := m.acquireScript.Run(ctx, m.rdb, keys, argv...).Result()
	if err != nil {
		return err
	}
	if aq == "0" {
		return redis.TxFailedErr
	}
	defer func(errOut *error) {
		keys := []string{LockKey, RemoversCountKey}
		rel, err := m.releaseScript.Run(ctx, m.rdb, keys).Result()
		if err != nil {
			*errOut = err
		}
		if rel == "0" {
			*errOut = errors.New("cannot release lock")
		}
	}(&resErr)
	pKeys := playerKeys{playerId}
	_, err = m.rdb.Del(ctx, pKeys.keys()...).Result()
	return err
}

func (m *mmPoolRepoImpl) CheckInPool(ctx context.Context, playerId int64) (bool, error) {
	pKeys := playerKeys{playerId}
	res, err := m.rdb.Exists(ctx, pKeys.rating()).Result()
	return res == 1, err
}
