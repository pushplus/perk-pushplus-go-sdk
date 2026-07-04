package pushplus

import (
	"context"
	"sync"
	"time"
)

// defaultAccessKeyTTLSeconds 服务端未返回 expiresIn 时的默认有效期（秒）。
const defaultAccessKeyTTLSeconds int64 = 7200

// AccessKeyManager 负责 AccessKey 的自动获取、缓存与过期前刷新。
// 并发安全；开放接口收到业务 code=401 时会调用 Invalidate 使缓存失效并触发重新获取。
type AccessKeyManager struct {
	config *Config
	api    *AccessKeyAPI
	now    func() time.Time

	mu        sync.Mutex
	cachedKey string
	expireAt  time.Time
}

// NewAccessKeyManager 创建 AccessKeyManager。
func NewAccessKeyManager(config *Config, api *AccessKeyAPI, now func() time.Time) *AccessKeyManager {
	if now == nil {
		now = time.Now
	}
	return &AccessKeyManager{config: config, api: api, now: now}
}

// GetAccessKey 返回有效的 AccessKey；缓存失效时自动刷新。
func (m *AccessKeyManager) GetAccessKey(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isValidLocked() {
		return m.cachedKey, nil
	}
	return m.refreshLocked(ctx)
}

// Invalidate 使缓存的 AccessKey 立即失效（下次调用时重新获取）。
func (m *AccessKeyManager) Invalidate() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cachedKey = ""
	m.expireAt = time.Time{}
}

func (m *AccessKeyManager) isValidLocked() bool {
	return m.cachedKey != "" && m.now().Before(m.expireAt)
}

func (m *AccessKeyManager) refreshLocked(ctx context.Context) (string, error) {
	result, err := m.api.GetAccessKey(ctx)
	if err != nil {
		return "", err
	}
	if result == nil || result.AccessKey == "" {
		return "", newError(-1, "获取 AccessKey 失败: 服务端未返回 accessKey")
	}

	expiresIn := result.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = defaultAccessKeyTTLSeconds
	}
	effectiveTTL := expiresIn - m.config.AccessKeyRefreshAheadSeconds
	if effectiveTTL < 1 {
		effectiveTTL = 1
	}

	m.cachedKey = result.AccessKey
	m.expireAt = m.now().Add(time.Duration(effectiveTTL) * time.Second)
	return m.cachedKey, nil
}
