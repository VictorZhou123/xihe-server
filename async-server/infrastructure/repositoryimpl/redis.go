package repositoryimpl

import (
	"context"
)

type redisClient interface {
	GetKeysRegex(context.Context, string) ([]string, error)
	DelKeys(context.Context, []string) error
}
