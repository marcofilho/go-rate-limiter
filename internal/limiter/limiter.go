package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	Client        *redis.Client
	MaxRequests   int
	BlockDuration time.Duration
	LimiterType   string
}

func NewRateLimiter(redisAddr, redisPassword string, redisDB int, maxRequests int, blockDuration time.Duration, limiterType string) *RateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	return &RateLimiter{
		Client:        client,
		MaxRequests:   maxRequests,
		BlockDuration: blockDuration,
		LimiterType:   limiterType,
	}
}

func (r *RateLimiter) Allow(key string, isToken bool) (bool, error) {
	ctx := context.Background()

	var maxRequests int
	if isToken {
		maxRequests = r.MaxRequestsForToken(key)
	} else {
		maxRequests = r.MaxRequests
	}

	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		err = r.Client.Expire(ctx, key, time.Second).Err()
		if err != nil {
			return false, err
		}
	}

	if count > int64(maxRequests) {
		return false, nil
	}

	return true, nil
}

func (r *RateLimiter) MaxRequestsForToken(token string) int {
	if token == "special-token" {
		return 100
	}
	return r.MaxRequests
}
