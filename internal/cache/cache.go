package cache

import (
	"context"
	"sync"
	"time"

	"github.com/maypok86/otter"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any) bool
	Delete(key string)
}

type InMemoryCache struct {
	mu   sync.RWMutex
	data map[string]any
}

func NewCache(ctx context.Context) (Cache, error) {
	cache, err := otter.MustBuilder[string, any](1_000).
		CollectStats().
		Cost(func(key string, value any) uint32 {
			return 1
		}).
		WithTTL(time.Minute).
		Build()
	if err != nil {
		panic(err)
	}

	return cache, nil
}

func (c *InMemoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	return value, ok
}

func (c *InMemoryCache) Set(key string, value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	return true
}

func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
