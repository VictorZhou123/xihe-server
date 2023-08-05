package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type dbRedis struct {
	Expiration time.Duration
}

func NewDBRedis(expiration int) dbRedis {
	return dbRedis{Expiration: time.Duration(expiration)}
}

func (r dbRedis) Create(
	ctx context.Context, key string, value interface{},
) *redis.StatusCmd {
	return client.Set(ctx, key, value, r.Expiration*time.Second)
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
	return client.Expire(ctx, key, 3*time.Second)
}

func (r dbRedis) GetKeysRegex(
	ctx context.Context, pattern string,
) (keys []string, err error) {
	var cursor uint64
	for {
		scanKeys, cursor, err := client.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			logrus.Warnf("getkeys error: %s", err.Error())

			return nil, errors.New("get redis keys error")
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	return
}

func (r dbRedis) DelKeys(
	ctx context.Context, keys []string,
) error {
	_, err := client.Del(ctx, keys...).Result()

	return err
}

func DB() *redis.Client {
	return client
}
