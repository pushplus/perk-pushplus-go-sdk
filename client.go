// Package pushplus 是 PushPlus(推送加) 官方接口的 Go SDK，
// 覆盖消息发送接口与全部开放接口。功能与 perk-pushplus-sdk（Java 版）对齐：
//   - AccessKey 自动获取、缓存、过期前刷新、失效自动重试，调用方无感知
//   - 本地限流守卫：发送接口命中 code=900 时自动短路同 token 的后续调用
//   - 单条 /send、多渠道 /batchSend、消息回调类型化解析
//   - 全部开放接口：消息、用户、消息令牌、群组、群组用户、好友、Webhook、
//     公众号/企业微信/邮箱渠道、ClawBot、功能设置、预处理、图片服务、
//     push 表单、push 文档、push 表格
package pushplus

import (
	"context"
	"time"
)

// Client 是 PushPlus SDK 的入口，聚合全部 API；并发安全。
type Client struct {
	config           *Config
	httpRequester    HTTPRequester
	accessKeyManager *AccessKeyManager
	rateLimitGuard   *RateLimitGuard

	message      *MessageAPI
	accessKey    *AccessKeyAPI
	openMessage  *OpenMessageAPI
	user         *UserAPI
	messageToken *MessageTokenAPI
	topic        *TopicAPI
	topicUser    *TopicUserAPI
	friend       *FriendAPI
	webhook      *WebhookAPI
	channel      *ChannelAPI
	clawBot      *ClawBotAPI
	setting      *SettingAPI
	pre          *PreAPI
	image        *ImageAPI
	form         *FormAPI
	doc          *DocAPI
	excel        *ExcelAPI
}

// NewClient 创建 PushPlus 客户端。
//
//	client := pushplus.NewClient(
//	    pushplus.WithToken("your_user_token"),
//	    pushplus.WithSecretKey("your_secret_key"), // 调用开放接口才需要
//	)
func NewClient(opts ...Option) *Client {
	o := &clientOptions{config: DefaultConfig()}
	for _, opt := range opts {
		opt(o)
	}
	o.config.normalize()

	requester := o.requester
	if requester == nil {
		requester = NewDefaultHTTPRequester(o.config)
	}
	now := o.clock
	if now == nil {
		now = time.Now
	}

	c := &core{config: o.config, http: requester}
	accessKeyAPI := newAccessKeyAPI(c)
	akm := NewAccessKeyManager(o.config, accessKeyAPI, now)
	guard := NewRateLimitGuard(o.config, now)

	return &Client{
		config:           o.config,
		httpRequester:    requester,
		accessKeyManager: akm,
		rateLimitGuard:   guard,
		message:          newMessageAPI(c, guard),
		accessKey:        accessKeyAPI,
		openMessage:      newOpenMessageAPI(c, akm),
		user:             newUserAPI(c, akm),
		messageToken:     newMessageTokenAPI(c, akm),
		topic:            newTopicAPI(c, akm),
		topicUser:        newTopicUserAPI(c, akm),
		friend:           newFriendAPI(c, akm),
		webhook:          newWebhookAPI(c, akm),
		channel:          newChannelAPI(c, akm),
		clawBot:          newClawBotAPI(c, akm),
		setting:          newSettingAPI(c, akm),
		pre:              newPreAPI(c, akm),
		image:            newImageAPI(c, akm),
		form:             newFormAPI(c, akm),
		doc:              newDocAPI(c, akm),
		excel:            newExcelAPI(c, akm),
	}
}

/* ============================== 组件访问器 ============================== */

// Config 返回客户端配置。
func (c *Client) Config() *Config { return c.config }

// HTTPRequester 返回底层 HTTP 实现。
func (c *Client) HTTPRequester() HTTPRequester { return c.httpRequester }

// AccessKeyManager 返回 AccessKey 管理器。
func (c *Client) AccessKeyManager() *AccessKeyManager { return c.accessKeyManager }

// RateLimitGuard 返回本地限流守卫，便于手动 Clear 或观察解禁时间。
func (c *Client) RateLimitGuard() *RateLimitGuard { return c.rateLimitGuard }

/* ============================== API 访问器 ============================== */

// Message 发送消息接口（/send、/batchSend）。
func (c *Client) Message() *MessageAPI { return c.message }

// AccessKey 获取 AccessKey 接口（一般无需手动调用）。
func (c *Client) AccessKey() *AccessKeyAPI { return c.accessKey }

// OpenMessage 开放接口 - 消息。
func (c *Client) OpenMessage() *OpenMessageAPI { return c.openMessage }

// User 开放接口 - 用户。
func (c *Client) User() *UserAPI { return c.user }

// MessageToken 开放接口 - 消息 token。
func (c *Client) MessageToken() *MessageTokenAPI { return c.messageToken }

// Topic 开放接口 - 群组。
func (c *Client) Topic() *TopicAPI { return c.topic }

// TopicUser 开放接口 - 群组用户。
func (c *Client) TopicUser() *TopicUserAPI { return c.topicUser }

// Friend 开放接口 - 好友功能。
func (c *Client) Friend() *FriendAPI { return c.friend }

// Webhook 开放接口 - 渠道配置 webhook。
func (c *Client) Webhook() *WebhookAPI { return c.webhook }

// Channel 开放接口 - 渠道配置（公众号/企业微信/邮箱）。
func (c *Client) Channel() *ChannelAPI { return c.channel }

// ClawBot 开放接口 - 微信 ClawBot。
func (c *Client) ClawBot() *ClawBotAPI { return c.clawBot }

// Setting 开放接口 - 功能设置。
func (c *Client) Setting() *SettingAPI { return c.setting }

// Pre 开放接口 - 预处理信息（需会员）。
func (c *Client) Pre() *PreAPI { return c.pre }

// Image 开放接口 - 图片服务。
func (c *Client) Image() *ImageAPI { return c.image }

// Form 开放接口 - push 表单。
func (c *Client) Form() *FormAPI { return c.form }

// Doc 开放接口 - push 文档。
func (c *Client) Doc() *DocAPI { return c.doc }

// Excel 开放接口 - push 表格。
func (c *Client) Excel() *ExcelAPI { return c.excel }

/* ============================== 便捷转发方法 ============================== */

// SendSimple 以默认渠道、默认模板发送一条简单消息，返回消息流水号。
func (c *Client) SendSimple(ctx context.Context, title, content string) (string, error) {
	return c.message.SendSimple(ctx, title, content)
}

// Send 发送单条消息，返回消息流水号。
func (c *Client) Send(ctx context.Context, request *SendRequest) (string, error) {
	return c.message.Send(ctx, request)
}

// BatchSend 多渠道发送消息。
func (c *Client) BatchSend(ctx context.Context, request *BatchSendRequest) ([]BatchSendResult, error) {
	return c.message.BatchSend(ctx, request)
}
