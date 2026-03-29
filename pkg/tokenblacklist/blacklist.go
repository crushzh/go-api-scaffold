package tokenblacklist

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"go-api-scaffold/pkg/cache"
)

const KeyPrefix = "token:blacklist:"

// Blacklist Token 黑名单接口
type Blacklist interface {
	Add(ctx context.Context, token string, expiration time.Duration) error
	IsBlacklisted(ctx context.Context, token string) bool
}

// CacheBlacklist 基于缓存的 Token 黑名单
type CacheBlacklist struct {
	cache cache.Cache
}

// New 创建黑名单
func New(c cache.Cache) *CacheBlacklist {
	return &CacheBlacklist{cache: c}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Add 将 Token 加入黑名单
func (b *CacheBlacklist) Add(ctx context.Context, token string, expiration time.Duration) error {
	return b.cache.SetString(ctx, KeyPrefix+hashToken(token), "1", expiration)
}

// IsBlacklisted 检查 Token 是否在黑名单中（fail-closed：缓存出错时拒绝）
func (b *CacheBlacklist) IsBlacklisted(ctx context.Context, token string) bool {
	exists, err := b.cache.Exists(ctx, KeyPrefix+hashToken(token))
	if err != nil {
		return true // fail-closed
	}
	return exists
}

var _ Blacklist = (*CacheBlacklist)(nil)
