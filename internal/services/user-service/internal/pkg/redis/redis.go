package redis

import (
	"context"
	"time"
)

// RedisService defines the interface for Redis operations
type RedisService interface {
	// Basic operations
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)

	// TTL operations
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Hash operations
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(ctx context.Context, key, field string, dest interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	// List operations
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPush(ctx context.Context, key string, values ...interface{}) error
	LPop(ctx context.Context, key string, dest interface{}) error
	RPop(ctx context.Context, key string, dest interface{}) error
	LLen(ctx context.Context, key string) (int64, error)

	// Set operations
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)

	// Connection
	Ping(ctx context.Context) error
	Close() error
}

// MockRedisService is a mock implementation for testing
type MockRedisService struct{}

func NewMockRedisService() RedisService {
	return &MockRedisService{}
}

func (m *MockRedisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}

func (m *MockRedisService) Get(ctx context.Context, key string, dest interface{}) error {
	return nil
}

func (m *MockRedisService) Del(ctx context.Context, keys ...string) error {
	return nil
}

func (m *MockRedisService) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (m *MockRedisService) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}

func (m *MockRedisService) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, nil
}

func (m *MockRedisService) HSet(ctx context.Context, key string, values ...interface{}) error {
	return nil
}

func (m *MockRedisService) HGet(ctx context.Context, key, field string, dest interface{}) error {
	return nil
}

func (m *MockRedisService) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return nil, nil
}

func (m *MockRedisService) HDel(ctx context.Context, key string, fields ...string) error {
	return nil
}

func (m *MockRedisService) LPush(ctx context.Context, key string, values ...interface{}) error {
	return nil
}

func (m *MockRedisService) RPush(ctx context.Context, key string, values ...interface{}) error {
	return nil
}

func (m *MockRedisService) LPop(ctx context.Context, key string, dest interface{}) error {
	return nil
}

func (m *MockRedisService) RPop(ctx context.Context, key string, dest interface{}) error {
	return nil
}

func (m *MockRedisService) LLen(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockRedisService) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return nil
}

func (m *MockRedisService) SRem(ctx context.Context, key string, members ...interface{}) error {
	return nil
}

func (m *MockRedisService) SMembers(ctx context.Context, key string) ([]string, error) {
	return nil, nil
}

func (m *MockRedisService) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return false, nil
}

func (m *MockRedisService) Ping(ctx context.Context) error {
	return nil
}

func (m *MockRedisService) Close() error {
	return nil
}
