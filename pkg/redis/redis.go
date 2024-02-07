package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/supersida159/e-commerce/src/users/entities"
)

const (
	Timeout = 10
)

// IRedis interface
//
//go:generate mockery --name=IRedis
type IRedis interface {
	IsConnected() bool
	Get(key string, value interface{}) error
	Set(key string, value interface{}) error
	SetWithExpiration(key string, value interface{}, expiration time.Duration) error
	Remove(keys ...string) error
	Keys(pattern string) ([]string, error)
	RemovePattern(pattern string) error
}

// Config redis
type Config struct {
	Address  string
	Password string
	Database int
}

type RedisWRealStore struct {
	cmd       goredis.Cmdable
	RealStore RealStore
}
type RealStore interface {
	FindUser(ctx context.Context, condition map[string]interface{}, moreInfores ...string) (*entities.User, error)
}

// NewRedis Redis interface with config
func NewRedis(config Config, realStore RealStore) *RedisWRealStore {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	rdb := goredis.NewClient(&goredis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.Database,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logrus.Fatal(pong, err)
		return nil
	}

	return &RedisWRealStore{
		cmd:       rdb,
		RealStore: realStore,
	}
}
func (r *RedisWRealStore) FindUser(ctx context.Context, condition map[string]interface{}, moreInfores ...string) (*entities.User, error) {
	userID := condition["id"].(int)
	var userInCache entities.User
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
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	if r.cmd == nil {
		return false
	}

	_, err := r.cmd.Ping(ctx).Result()
	if err != nil {
		return false
	}
	return true
}

func (r *RedisWRealStore) Get(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	strValue, err := r.cmd.Get(ctx, key).Result()
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
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	bData, _ := json.Marshal(value)
	err := r.cmd.Set(ctx, key, bData, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisWRealStore) Set(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	bData, _ := json.Marshal(value)
	err := r.cmd.Set(ctx, key, bData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisWRealStore) Remove(keys ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	err := r.cmd.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisWRealStore) Keys(pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	keys, err := r.cmd.Keys(ctx, pattern).Result()
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
