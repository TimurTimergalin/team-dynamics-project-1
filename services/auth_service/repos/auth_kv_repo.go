package repos

import (
	"context"
	"github.com/redis/go-redis/v9"
	"team_dynamics/logging"
	"time"
)

type AuthKvRepo interface {
	MarkUsed(ctx context.Context, tokenId string) (bool, error)
}

type authKvRepoImpl struct {
	rdb        *redis.Client
	expiration time.Duration
}

func MakeAuthKvRepo(rdb *redis.Client, expiration time.Duration) AuthKvRepo {
	return &authKvRepoImpl{rdb: rdb, expiration: expiration}
}

func (r *authKvRepoImpl) MarkUsed(ctx context.Context, tokenId string) (bool, error) {
	logger := logging.GetLogger(ctx)
	logger.Debug("marking refresh token as used", "token_id", tokenId)
	set, err := r.rdb.SetNX(ctx, tokenId, 1, r.expiration).Result()
	if err != nil {
		logger.Error("failed to mark refresh token as used", "token_id", tokenId, "error", err)
		return false, err
	}
	if !set {
		logger.Debug("refresh token already marked as used", "token_id", tokenId)
	}
	return set, nil
}
