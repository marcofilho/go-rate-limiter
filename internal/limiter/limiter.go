package limiter

import (
	"context"
	"time"
)

type RateLimiter struct {
	Storage           Storage
	MaxRequests       int
	BlockDuration     time.Duration
	TokenRequestLimit int
}

func NewRateLimiter(storage Storage, maxRequests int, blockDuration time.Duration, tokenRequestLimit int) *RateLimiter {
	return &RateLimiter{
		Storage:           storage,
		MaxRequests:       maxRequests,
		BlockDuration:     blockDuration,
		TokenRequestLimit: tokenRequestLimit,
	}
}

func (r *RateLimiter) Allow(key string, isToken bool, tokenRequestLimit int) (bool, error) {
	ctx := context.Background()

	var maxRequests int
	if isToken {
		maxRequests = r.maxRequestsForToken(key, tokenRequestLimit)
	} else {
		maxRequests = r.MaxRequests
	}

	count, err := r.Storage.Incr(ctx, key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		err = r.Storage.Expire(ctx, key, time.Second)
		if err != nil {
			return false, err
		}
	}

	if count > int64(maxRequests) {
		return false, nil
	}

	return true, nil
}

func (r *RateLimiter) maxRequestsForToken(token string, tokenRequestLimit int) int {
	tokenLimits := map[string]int{
		"API-TOKEN": tokenRequestLimit,
	}

	if limit, exists := tokenLimits[token]; exists {
		return limit
	}

	return r.MaxRequests
}
