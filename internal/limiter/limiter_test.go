package limiter

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	client.FlushDB(client.Context())
	return client
}

func TestNewRateLimiter(t *testing.T) {
	client := setupTestRedis()
	defer client.Close()

	limiter := NewRateLimiter("localhost:6379", "", 0, 10, time.Second*10, "basic")
	assert.NotNil(t, limiter)
	assert.Equal(t, 10, limiter.MaxRequests)
	assert.Equal(t, time.Second*10, limiter.BlockDuration)
	assert.Equal(t, "basic", limiter.LimiterType)
}

func TestRateLimiter_AllowWithinLimit(t *testing.T) {
	client := setupTestRedis()
	defer client.Close()

	limiter := &RateLimiter{
		Client:        client,
		MaxRequests:   5,
		BlockDuration: time.Second * 10,
		LimiterType:   "basic",
	}

	key := "test_key_within_limit"

	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(key, false)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}
}

func TestRateLimiter_AllowExceedLimit(t *testing.T) {
	client := setupTestRedis()
	defer client.Close()

	limiter := &RateLimiter{
		Client:        client,
		MaxRequests:   3,
		BlockDuration: time.Second * 5,
		LimiterType:   "basic",
	}

	key := "test_key_exceed_limit"

	for i := 0; i < 3; i++ {
		allowed, err := limiter.Allow(key, false)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := limiter.Allow(key, false)
	assert.NoError(t, err)
	assert.False(t, allowed)
}

func TestRateLimiter_BlockDuration(t *testing.T) {
	client := setupTestRedis()
	defer client.Close()

	limiter := &RateLimiter{
		Client:        client,
		MaxRequests:   2,
		BlockDuration: time.Second * 2,
		LimiterType:   "basic",
	}

	key := "test_key_block_duration"

	for i := 0; i < 2; i++ {
		allowed, err := limiter.Allow(key, false)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := limiter.Allow(key, false)
	assert.NoError(t, err)
	assert.False(t, allowed)

	time.Sleep(time.Second * 2)

	allowed, err = limiter.Allow(key, false)
	assert.NoError(t, err)
	assert.True(t, allowed)
}

func TestRateLimiter_MaxRequestsForToken(t *testing.T) {
	client := setupTestRedis()
	defer client.Close()

	limiter := &RateLimiter{
		Client:        client,
		MaxRequests:   5,
		BlockDuration: time.Second * 10,
		LimiterType:   "token",
	}

	token := "special-token"
	nonSpecialToken := "regular-token"

	assert.Equal(t, 100, limiter.MaxRequestsForToken(token))
	assert.Equal(t, 5, limiter.MaxRequestsForToken(nonSpecialToken))
}

func TestRateLimiter_TokenBasedLimit(t *testing.T) {
	client := setupTestRedis()
	defer client.Close()

	limiter := &RateLimiter{
		Client:        client,
		MaxRequests:   5,
		BlockDuration: time.Second * 10,
		LimiterType:   "token",
	}

	token := "special-token"

	for i := 0; i < 100; i++ {
		allowed, err := limiter.Allow(token, true)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := limiter.Allow(token, true)
	assert.NoError(t, err)
	assert.False(t, allowed)
}
