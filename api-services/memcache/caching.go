package memcache

import (
	"errors"
	"sync"
	"time"
)

type cachedItem struct {
	value      interface{}
	expiration *time.Time // Use a pointer to time.Time to allow for nil (no expiration)
}

type caching struct {
	store  map[string]cachedItem
	locker *sync.RWMutex
}

type Caching interface {
	Read(key string) (interface{}, error)
	Write(key string, value interface{}, expiration ...time.Duration)
}

func NewCaching() *caching {
	return &caching{
		store:  make(map[string]cachedItem),
		locker: new(sync.RWMutex),
	}
}

func (c *caching) Read(key string) (interface{}, error) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	item, found := c.store[key]
	if !found {
		return nil, errors.New("key not found")
	}

	if item.expiration != nil && time.Now().After(*item.expiration) {
		// If the item has expired and expiration time is set, delete it and return an error
		c.delete(key)
		return item, errors.New("item expired")
	}

	return item.value, nil
}

func (c *caching) Write(key string, value interface{}, expiration ...time.Duration) {
	c.locker.Lock()
	defer c.locker.Unlock()

	var expirationTime *time.Time
	if len(expiration) > 0 {
		expTime := time.Now().Add(expiration[0])
		expirationTime = &expTime
	}

	c.store[key] = cachedItem{
		value:      value,
		expiration: expirationTime,
	}
}

// Helper function to delete expired entries
func (c *caching) delete(key string) {
	delete(c.store, key)
}
