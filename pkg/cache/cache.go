package cache

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// ErrNotFound 键不存在
var ErrNotFound = fmt.Errorf("cache: key not found")

// Cache 统一缓存接口
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	GetString(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	SetString(ctx context.Context, key, value string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string, expiration time.Duration) (int64, error)
	DeleteByPrefix(ctx context.Context, prefix string) error
	Close() error
	Name() string
	Ping(ctx context.Context) error
}

// MemoryCache 基于 go-cache 的内存缓存实现
type MemoryCache struct {
	cache             *gocache.Cache
	defaultExpiration time.Duration
	mu                sync.RWMutex
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(defaultExpiration, cleanupInterval time.Duration) *MemoryCache {
	if defaultExpiration == 0 {
		defaultExpiration = 15 * time.Minute
	}
	if cleanupInterval == 0 {
		cleanupInterval = 5 * time.Minute
	}
	return &MemoryCache{
		cache:             gocache.New(defaultExpiration, cleanupInterval),
		defaultExpiration: defaultExpiration,
	}
}

func (c *MemoryCache) Get(_ context.Context, key string) ([]byte, error) {
	val, found := c.cache.Get(key)
	if !found {
		return nil, ErrNotFound
	}
	switch v := val.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return nil, ErrNotFound
	}
}

func (c *MemoryCache) GetString(_ context.Context, key string) (string, error) {
	val, found := c.cache.Get(key)
	if !found {
		return "", ErrNotFound
	}
	switch v := val.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case int64:
		return fmt.Sprintf("%d", v), nil
	case int:
		return fmt.Sprintf("%d", v), nil
	default:
		return "", ErrNotFound
	}
}

func (c *MemoryCache) Set(_ context.Context, key string, value []byte, exp time.Duration) error {
	if exp == 0 {
		exp = c.defaultExpiration
	}
	c.cache.Set(key, value, exp)
	return nil
}

func (c *MemoryCache) SetString(_ context.Context, key, value string, exp time.Duration) error {
	if exp == 0 {
		exp = c.defaultExpiration
	}
	c.cache.Set(key, value, exp)
	return nil
}

func (c *MemoryCache) Delete(_ context.Context, key string) error {
	c.cache.Delete(key)
	return nil
}

func (c *MemoryCache) Exists(_ context.Context, key string) (bool, error) {
	_, found := c.cache.Get(key)
	return found, nil
}

// Increment 原子递增（互斥锁保证原子性）
func (c *MemoryCache) Increment(_ context.Context, key string, exp time.Duration) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if exp == 0 {
		exp = c.defaultExpiration
	}
	val, found := c.cache.Get(key)
	var count int64
	if found {
		switch v := val.(type) {
		case int64:
			count = v
		case string:
			fmt.Sscanf(v, "%d", &count)
		}
	}
	count++
	c.cache.Set(key, count, exp)
	return count, nil
}

func (c *MemoryCache) DeleteByPrefix(_ context.Context, prefix string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key := range c.cache.Items() {
		if strings.HasPrefix(key, prefix) {
			c.cache.Delete(key)
		}
	}
	return nil
}

func (c *MemoryCache) Close() error                 { c.cache.Flush(); return nil }
func (c *MemoryCache) Name() string                 { return "memory" }
func (c *MemoryCache) Ping(_ context.Context) error { return nil }

var _ Cache = (*MemoryCache)(nil)
