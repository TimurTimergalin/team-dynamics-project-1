package repos

import (
	"context"
	"github.com/redis/go-redis/v9"
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
	set, err := r.rdb.SetNX(ctx, tokenId, 1, r.expiration).Result()
	if err != nil {
		return false, err
	}
	return set, nil
}
