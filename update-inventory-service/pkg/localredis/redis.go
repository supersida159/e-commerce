package localredis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	Timeout       = 10 * time.Second
	LockTimeout   = 5 * time.Second // Duration for how long the lock is held
	LockRetryTime = 100 * time.Millisecond
)

var ErrKeyLocked = fmt.Errorf("key is already locked")

// IRedis interface
//
//go:generate mockery --name=IRedis
// type IRedis interface {
// 	IsConnected() bool
// 	Get(key string, value interface{}) error
// 	Set(key string, value interface{}) error
// 	SetWithExpiration(key string, value interface{}, expiration time.Duration) error
// 	Remove(keys ...string) error
// 	Keys(pattern string) ([]string, error)
// 	RemovePattern(pattern string) error
// }ss

// Config redis
type Config struct {
	Address  string
	Password string
	Database int
}

type RedisWRealStore struct {
	Client    *goredis.Client
	RealStore RealStore
}

type RealStore interface {
}

// NewRedis Redis interface with config
func NewRedis(config Config, realStore RealStore) *RedisWRealStore {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	rdb := goredis.NewClient(&goredis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.Database,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logrus.Fatal("error at redis", pong, err)
		return nil
	}

	return &RedisWRealStore{
		Client:    rdb,
		RealStore: realStore,
	}
}
func (r *RedisWRealStore) GetClient() *goredis.Client {
	return r.Client
}

func (r *RedisWRealStore) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	if r.Client == nil {
		return false
	}

	_, err := r.Client.Ping(ctx).Result()
	if err != nil {
		return false
	}
	return true
}

func (r *RedisWRealStore) Get(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	strValue, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(strValue), value)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisWRealStore) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	bData, _ := json.Marshal(value)
	err := r.Client.Set(ctx, key, bData, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisWRealStore) Set(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	bData, _ := json.Marshal(value)
	err := r.Client.Set(ctx, key, bData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisWRealStore) Remove(keys ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	err := r.Client.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisWRealStore) Keys(pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	keys, err := r.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *RedisWRealStore) RemovePattern(pattern string) error {
	keys, err := r.Keys(pattern)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	err = r.Remove(keys...)
	if err != nil {
		return err
	}

	return nil
}

// LockKey locks a key with a specific timeout to avoid deadlocks
func (r *RedisWRealStore) LockKey(ctx context.Context, key string) error {
	for {
		ok, err := r.Client.SetNX(ctx, key+"_lock", "locked", LockTimeout).Result()
		if err != nil {
			return err
		}
		if ok {
			return nil // Lock acquired
		}

		// Wait for a short period before retrying
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled while trying to acquire lock on %s", key)
		case <-time.After(LockRetryTime):
		}
	}
}

// UnlockKey unlocks the specified key
func (r *RedisWRealStore) UnlockKey(ctx context.Context, key string) error {
	_, err := r.Client.Del(ctx, key+"_lock").Result()
	return err
}
func (r *RedisWRealStore) SetWithExpirationPreserve(ctx context.Context, key string, value interface{}) error {
	// Step 1: Get the current TTL (Time To Live) of the key
	ttl, err := r.Client.TTL(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("could not get TTL for key %s: %w", key, err)
	}

	// Step 2: Update the value of the key
	bData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("could not marshal value: %w", err)
	}

	// Set the new value (no expiration set here, meaning the key will persist)
	err = r.Client.Set(ctx, key, bData, 0).Err()
	if err != nil {
		return fmt.Errorf("could not set key %s: %w", key, err)
	}

	// Step 3: Reapply the original TTL (expiration time)
	if ttl > 0 {
		err = r.Client.Expire(ctx, key, ttl).Err()
		if err != nil {
			return fmt.Errorf("could not reset expiration for key %s: %w", key, err)
		}
	}

	return nil
}
