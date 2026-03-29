package lockout

import (
	"context"
	"strconv"
	"time"

	"go-api-scaffold/pkg/cache"
)

const keyPrefix = "login_fail:"

// Manager 登录失败锁定管理器
type Manager struct {
	cache     cache.Cache
	threshold int
	duration  time.Duration
}

// New 创建锁定管理器
func New(c cache.Cache, threshold int, duration time.Duration) *Manager {
	if threshold <= 0 {
		threshold = 5
	}
	if duration <= 0 {
		duration = 15 * time.Minute
	}
	return &Manager{cache: c, threshold: threshold, duration: duration}
}

// RecordFailure 记录一次登录失败（原子操作），返回当前失败次数
func (m *Manager) RecordFailure(ctx context.Context, username string) int {
	count, err := m.cache.Increment(ctx, keyPrefix+username, m.duration)
	if err != nil {
		val, _ := m.cache.GetString(ctx, keyPrefix+username)
		c, _ := strconv.Atoi(val)
		c++
		_ = m.cache.SetString(ctx, keyPrefix+username, strconv.Itoa(c), m.duration)
		return c
	}
	return int(count)
}

// IsLocked 是否已被锁定
func (m *Manager) IsLocked(ctx context.Context, username string) bool {
	val, err := m.cache.GetString(ctx, keyPrefix+username)
	if err != nil {
		return false
	}
	count, _ := strconv.Atoi(val)
	return count >= m.threshold
}

// Clear 清除失败记录（登录成功后调用）
func (m *Manager) Clear(ctx context.Context, username string) {
	_ = m.cache.Delete(ctx, keyPrefix+username)
}

// Threshold 返回阈值
func (m *Manager) Threshold() int { return m.threshold }

// Duration 返回锁定时长
func (m *Manager) Duration() time.Duration { return m.duration }
