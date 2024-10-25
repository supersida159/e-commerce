package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	// Updated import

	goredis "github.com/redis/go-redis/v9" // Updated import path

	"github.com/sirupsen/logrus"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
)

const (
	Timeout       = 10 * time.Second
	LockTimeout   = 5 * time.Second // Duration for how long the lock is held
	LockRetryTime = 100 * time.Millisecond
)

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
// }

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

type Client interface {
	RunExpireOrder(ctx context.Context)
}
type RealStore interface {
	FindUser(ctx context.Context, condition map[string]interface{}, moreInfores ...string) (*entities_user.User, error)
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
func (r *RedisWRealStore) FindUser(ctx context.Context, condition map[string]interface{}, moreInfores ...string) (*entities_user.User, error) {
	userID := condition["id"].(int)
	var userInCache entities_user.User
	err := r.Get(fmt.Sprintf("user-%d", userID), &userInCache)
	if err != nil {
		return &userInCache, nil
	}
	userInRealStore, err := r.RealStore.FindUser(ctx, condition, moreInfores...)
	if err != nil {
		return nil, err
	}
	go func() {
		r.Set(fmt.Sprintf("user-%d", userID), userInRealStore)
	}()
	return userInRealStore, nil
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
