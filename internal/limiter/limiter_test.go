package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockStorage is a mock implementation of the Storage interface for testing
type MockStorage struct {
	data map[string]int64
	ttl  map[string]time.Time
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		data: make(map[string]int64),
		ttl:  make(map[string]time.Time),
	}
}

func (m *MockStorage) Incr(ctx context.Context, key string) (int64, error) {
	if _, exists := m.data[key]; !exists {
		m.data[key] = 0
	}
	m.data[key]++
	return m.data[key], nil
}

func (m *MockStorage) Expire(ctx context.Context, key string, duration time.Duration) error {
	m.ttl[key] = time.Now().Add(duration)
	return nil
}

func (m *MockStorage) Exists(ctx context.Context, key string) (bool, error) {
	expiration, exists := m.ttl[key]
	if !exists {
		return false, nil
	}
	if time.Now().After(expiration) {
		delete(m.data, key)
		delete(m.ttl, key)
		return false, nil
	}
	return true, nil
}

func TestRateLimiter_Allow(t *testing.T) {
	mockStorage := NewMockStorage()
	rateLimiter := NewRateLimiter(mockStorage, 5, 10*time.Second, 100)

	// Test IP-based rate limiting
	for i := 1; i <= 5; i++ {
		allowed, err := rateLimiter.Allow("192.168.1.1", false, 0)
		assert.NoError(t, err)
		assert.True(t, allowed, "Request %d should be allowed", i)
	}

	// Exceed the limit
	allowed, err := rateLimiter.Allow("192.168.1.1", false, 0)
	assert.NoError(t, err)
	assert.False(t, allowed, "Request should be blocked after exceeding the limit")

	// Test token-based rate limiting
	for i := 1; i <= 100; i++ {
		allowed, err := rateLimiter.Allow("API-TOKEN", true, 100)
		assert.NoError(t, err)
		assert.True(t, allowed, "Request %d should be allowed for token", i)
	}

	// Exceed the token limit
	allowed, err = rateLimiter.Allow("API-TOKEN", true, 100)
	assert.NoError(t, err)
	assert.False(t, allowed, "Request should be blocked after exceeding the token limit")
}
