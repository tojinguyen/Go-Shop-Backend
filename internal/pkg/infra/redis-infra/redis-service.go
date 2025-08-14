package redis_infra

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisServiceInterface defines the contract for Redis operations
type RedisServiceInterface interface {
	// Connection & Management
	Ping() error
	Close() error

	// Basic Operations
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)

	// TTL Operations
	SetTTL(key string, expiration time.Duration) error
	GetTTL(key string) (time.Duration, error)

	// JSON Operations
	SetJSON(key string, value interface{}, expiration time.Duration) error
	GetJSON(key string, dest interface{}) error

	// Numeric Operations
	Increment(key string) (int64, error)
	Decrement(key string) (int64, error)

	// Hash Operations
	HSet(key, field string, value interface{}) error
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HDel(key, field string) error

	// List Operations
	LPush(key string, value interface{}) error
	RPush(key string, value interface{}) error
	LPop(key string) (string, error)
	RPop(key string) (string, error)
	LLen(key string) (int64, error)

	// Set Operations
	SAdd(key string, members ...interface{}) error
	SMembers(key string) ([]string, error)
	SIsMember(key string, member interface{}) (bool, error)
	SRem(key string, members ...interface{}) error

	// Utility
	FlushDB() error
}

type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisService creates a new Redis service instance
func NewRedisService(host, port string, password string, db int) *RedisService {
	redisAddr := fmt.Sprintf("%s:%s", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       db,
		PoolSize: 250,
		// Giữ một số kết nối "nhàn rỗi" để sẵn sàng phục vụ request mới ngay lập tức.
		MinIdleConns: 50,
		// Thời gian chờ tối đa để lấy một kết nối từ pool.
		PoolTimeout: 30 * time.Second,
		// Thời gian tối đa một kết nối có thể ở trạng thái nhàn rỗi trong pool.
		ConnMaxIdleTime: 10 * time.Minute,
	})

	return &RedisService{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Ping checks if Redis connection is alive
func (r *RedisService) Ping() error {
	_, err := r.client.Ping(r.ctx).Result()
	return err
}

// Set stores a key-value pair with optional expiration
func (r *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (r *RedisService) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// Delete removes a key from Redis
func (r *RedisService) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists
func (r *RedisService) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}

// SetJSON stores a JSON object
func (r *RedisService) SetJSON(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return r.client.Set(r.ctx, key, jsonData, expiration).Err()
}

// GetJSON retrieves and unmarshals a JSON object
func (r *RedisService) GetJSON(key string, dest interface{}) error {
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// SetTTL sets expiration time for an existing key
func (r *RedisService) SetTTL(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// GetTTL gets the time to live for a key
func (r *RedisService) GetTTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

// Increment increments a numeric value
func (r *RedisService) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// Decrement decrements a numeric value
func (r *RedisService) Decrement(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

// HSet sets a field in a hash
func (r *RedisService) HSet(key, field string, value interface{}) error {
	return r.client.HSet(r.ctx, key, field, value).Err()
}

// HGet gets a field from a hash
func (r *RedisService) HGet(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

// HGetAll gets all fields from a hash
func (r *RedisService) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

// HDel deletes a field from a hash
func (r *RedisService) HDel(key, field string) error {
	return r.client.HDel(r.ctx, key, field).Err()
}

// LPush pushes an element to the left of a list
func (r *RedisService) LPush(key string, value interface{}) error {
	return r.client.LPush(r.ctx, key, value).Err()
}

// RPush pushes an element to the right of a list
func (r *RedisService) RPush(key string, value interface{}) error {
	return r.client.RPush(r.ctx, key, value).Err()
}

// LPop pops an element from the left of a list
func (r *RedisService) LPop(key string) (string, error) {
	return r.client.LPop(r.ctx, key).Result()
}

// RPop pops an element from the right of a list
func (r *RedisService) RPop(key string) (string, error) {
	return r.client.RPop(r.ctx, key).Result()
}

// LLen gets the length of a list
func (r *RedisService) LLen(key string) (int64, error) {
	return r.client.LLen(r.ctx, key).Result()
}

// SAdd adds a member to a set
func (r *RedisService) SAdd(key string, members ...interface{}) error {
	return r.client.SAdd(r.ctx, key, members...).Err()
}

// SMembers gets all members of a set
func (r *RedisService) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}

// SIsMember checks if a member exists in a set
func (r *RedisService) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(r.ctx, key, member).Result()
}

// SRem removes a member from a set
func (r *RedisService) SRem(key string, members ...interface{}) error {
	return r.client.SRem(r.ctx, key, members...).Err()
}

// FlushDB clears the current database
func (r *RedisService) FlushDB() error {
	return r.client.FlushDB(r.ctx).Err()
}

// Close closes the Redis connection
func (r *RedisService) Close() error {
	return r.client.Close()
}
