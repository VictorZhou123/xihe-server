package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type dbRedis struct {
	Expiration time.Duration
}

func NewDBRedis(expiration time.Duration) dbRedis {
	return dbRedis{Expiration: expiration}
}

func (r dbRedis) Create(
	ctx context.Context, key string, value interface{},
) *redis.StatusCmd {
	return client.Set(ctx, key, value, r.Expiration)
}

func (r dbRedis) Get(
	ctx context.Context, key string,
) *redis.StringCmd {
	return client.Get(ctx, key)
}

func (r dbRedis) Delete(
	ctx context.Context, key string,
) *redis.IntCmd {
	return client.Del(ctx, key)
}

func (r dbRedis) Expire(
	ctx context.Context, key string, expire time.Duration,
) *redis.BoolCmd {
	return client.Expire(ctx, key, expire)
}

func DB() *redis.Client {
	return client
}
