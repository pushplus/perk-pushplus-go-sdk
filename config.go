package pushplus

import (
	"strings"
	"time"
)

// DefaultBaseURL 是 PushPlus 官方服务地址。
const DefaultBaseURL = "https://www.pushplus.plus"

// Config 是 SDK 的全部配置项，字段语义与默认值同 Java 版 PushPlusConfig 保持一致。
type Config struct {
	// Token 用户 token / 消息 token，发送消息使用。
	Token string
	// SecretKey 用户 secretKey，调用开放接口（获取 AccessKey）使用。
	SecretKey string
	// BaseURL 服务地址，默认 https://www.pushplus.plus。
	BaseURL string
	// ConnectTimeout 连接超时，默认 10s。
	ConnectTimeout time.Duration
	// ReadTimeout 读超时（整个请求超时），默认 30s。
	ReadTimeout time.Duration
	// AccessKeyRefreshAheadSeconds AccessKey 提前刷新秒数，默认 300。
	AccessKeyRefreshAheadSeconds int64
	// LogRequest 是否打印请求/响应日志，默认 false。
	LogRequest bool
	// RateLimitGuardEnabled 是否启用本地限流守卫（命中 code=900 后短路同 token 的发送请求），默认 true。
	RateLimitGuardEnabled bool
	// RateLimitCooldown 命中 code=900 后的本地禁推时长；为 0 表示禁推到本地时区"次日 0 点"。
	RateLimitCooldown time.Duration
}

// DefaultConfig 返回与 Java 版默认值一致的配置。
func DefaultConfig() *Config {
	return &Config{
		BaseURL:                      DefaultBaseURL,
		ConnectTimeout:               10 * time.Second,
		ReadTimeout:                  30 * time.Second,
		AccessKeyRefreshAheadSeconds: 300,
		RateLimitGuardEnabled:        true,
	}
}

// ResolveBaseURL 返回去掉尾部 "/" 的 BaseURL；为空时回退默认地址。
func (c *Config) ResolveBaseURL() string {
	base := strings.TrimSpace(c.BaseURL)
	if base == "" {
		base = DefaultBaseURL
	}
	return strings.TrimRight(base, "/")
}

// Option 是创建 Client 时的函数式选项。
type Option func(*clientOptions)

type clientOptions struct {
	config    *Config
	requester HTTPRequester
	clock     func() time.Time
}

// WithToken 设置用户 token / 消息 token。
func WithToken(token string) Option {
	return func(o *clientOptions) { o.config.Token = token }
}

// WithSecretKey 设置调用开放接口所需的 secretKey。
func WithSecretKey(secretKey string) Option {
	return func(o *clientOptions) { o.config.SecretKey = secretKey }
}

// WithBaseURL 设置服务地址。
func WithBaseURL(baseURL string) Option {
	return func(o *clientOptions) { o.config.BaseURL = baseURL }
}

// WithConnectTimeout 设置连接超时。
func WithConnectTimeout(d time.Duration) Option {
	return func(o *clientOptions) { o.config.ConnectTimeout = d }
}

// WithReadTimeout 设置读超时。
func WithReadTimeout(d time.Duration) Option {
	return func(o *clientOptions) { o.config.ReadTimeout = d }
}

// WithAccessKeyRefreshAheadSeconds 设置 AccessKey 提前刷新秒数。
func WithAccessKeyRefreshAheadSeconds(seconds int64) Option {
	return func(o *clientOptions) { o.config.AccessKeyRefreshAheadSeconds = seconds }
}

// WithLogRequest 开启/关闭请求日志。
func WithLogRequest(enabled bool) Option {
	return func(o *clientOptions) { o.config.LogRequest = enabled }
}

// WithRateLimitGuardEnabled 开启/关闭本地限流守卫。
func WithRateLimitGuardEnabled(enabled bool) Option {
	return func(o *clientOptions) { o.config.RateLimitGuardEnabled = enabled }
}

// WithRateLimitCooldown 设置命中 code=900 后的本地禁推时长；0 表示禁推到次日 0 点。
func WithRateLimitCooldown(d time.Duration) Option {
	return func(o *clientOptions) { o.config.RateLimitCooldown = d }
}

// WithConfig 直接使用完整配置（覆盖之前通过其它 Option 设置的全部配置字段）。
func WithConfig(config *Config) Option {
	return func(o *clientOptions) {
		if config != nil {
			cp := *config
			o.config = &cp
		}
	}
}

// WithHTTPRequester 自定义 HTTP 实现（例如注入 mock、或换成其它 HTTP 客户端）。
func WithHTTPRequester(requester HTTPRequester) Option {
	return func(o *clientOptions) { o.requester = requester }
}

// WithClock 注入时钟（主要供测试使用，影响限流守卫与 AccessKey 缓存）。
func WithClock(now func() time.Time) Option {
	return func(o *clientOptions) { o.clock = now }
}

func (c *Config) normalize() {
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.ConnectTimeout <= 0 {
		c.ConnectTimeout = 10 * time.Second
	}
	if c.ReadTimeout <= 0 {
		c.ReadTimeout = 30 * time.Second
	}
	if c.AccessKeyRefreshAheadSeconds < 0 {
		c.AccessKeyRefreshAheadSeconds = 0
	}
}
