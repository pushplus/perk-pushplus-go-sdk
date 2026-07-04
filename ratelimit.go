package pushplus

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// RateLimitGuard 本地限流守卫。
//
// PushPlus 在请求次数过多时返回 code=900，官方文档建议根据返回值判断当天是否继续调用发送接口。
// 守卫命中后会按 token 维度记录"禁推至 X 时刻"，同 token 的后续发送调用直接在本地短路，
// 不再发起 HTTP。仅作用于发送接口（/send、/batchSend），开放接口不受影响。
// 进程内单例，不跨进程共享。
type RateLimitGuard struct {
	config *Config
	now    func() time.Time

	mu      sync.Mutex
	blocked map[string]time.Time
}

// NewRateLimitGuard 创建限流守卫；now 为 nil 时使用系统时钟。
func NewRateLimitGuard(config *Config, now func() time.Time) *RateLimitGuard {
	if now == nil {
		now = time.Now
	}
	return &RateLimitGuard{
		config:  config,
		now:     now,
		blocked: make(map[string]time.Time),
	}
}

// Check 检查 token 是否处于本地禁推期；若是则返回 code=900 的 *Error（不发起 HTTP）。
func (g *RateLimitGuard) Check(token string) error {
	if !g.config.RateLimitGuardEnabled {
		return nil
	}
	key := strings.TrimSpace(token)
	if key == "" {
		return nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	until, ok := g.blocked[key]
	if !ok {
		return nil
	}
	if g.now().Before(until) {
		return newError(int(ErrorCodeRateLimited), fmt.Sprintf(
			"本地限流守卫：PushPlus 此前返回 code=900（请求次数过多），发送已被短路至 %s，可调用 RateLimitGuard.Clear 手动解除",
			until.Format(time.RFC3339)))
	}
	delete(g.blocked, key)
	return nil
}

// MarkBlocked 记录 token 命中限流，禁推到冷却期结束。
func (g *RateLimitGuard) MarkBlocked(token string) {
	if !g.config.RateLimitGuardEnabled {
		return
	}
	key := strings.TrimSpace(token)
	if key == "" {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.blocked[key] = g.cooldownUntil(g.now())
}

// Clear 手动清除 token 的禁推记录（例如人工确认服务端已解禁后立即放行）。
func (g *RateLimitGuard) Clear(token string) {
	key := strings.TrimSpace(token)
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.blocked, key)
}

// BlockedUntil 查询 token 的本地解禁时间；零值 time.Time 表示未被限流。
func (g *RateLimitGuard) BlockedUntil(token string) time.Time {
	key := strings.TrimSpace(token)
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.blocked[key]
}

// cooldownUntil 计算禁推截止时间：配置了 RateLimitCooldown 则为 now+cooldown，
// 否则为本地时区的次日 0 点。
func (g *RateLimitGuard) cooldownUntil(now time.Time) time.Time {
	if g.config.RateLimitCooldown > 0 {
		return now.Add(g.config.RateLimitCooldown)
	}
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
}
