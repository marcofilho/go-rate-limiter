package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string, db int) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStorage{client: client}
}

func (r *RedisStorage) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisStorage) Expire(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Expire(ctx, key, duration).Err()
}

func (r *RedisStorage) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	return exists > 0, err
}
