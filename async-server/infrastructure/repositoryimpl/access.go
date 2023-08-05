package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/async-server/domain/repository"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/common/infrastructure/redis"
)

func NewAccessRepoImpl(cfg *Config) repository.Access {
	return &accessRepoImpl{
		pgsql: pgsql.NewDBTable(cfg.Table.Access),
		redis: redis.NewDBRedis(3),
	}
}

type accessRepoImpl struct {
	pgsql pgsqlClient
	redis redisClient
}

func (impl accessRepoImpl) GetKeys() (keys []string, err error) {
	pattern := "ip_*"

	f := func(ctx context.Context) error {
		keys, err = impl.redis.GetKeysRegex(ctx, pattern)
		if err != nil {
			return err
		}

		return nil
	}

	if err = redis.WithContext(f); err != nil {
		return
	}

	return
}

func (impl accessRepoImpl) DelKeys(keys []string) error {
	f := func(ctx context.Context) error {
		return impl.redis.DelKeys(ctx, keys)
	}

	return redis.WithContext(f)
}
