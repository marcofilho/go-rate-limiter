package limiter

import (
	"context"
	"time"
)

type Storage interface {
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, duration time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
}
