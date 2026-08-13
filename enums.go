package pushplus

// Channel 发送渠道，JSON 序列化为字符串编码（与官方接口一致）。
type Channel string

const (
	// ChannelWechat 微信公众号。
	ChannelWechat Channel = "wechat"
	// ChannelWebhook 第三方 webhook。
	ChannelWebhook Channel = "webhook"
	// ChannelCp 企业微信应用。
	ChannelCp Channel = "cp"
	// ChannelMail 邮箱。
	ChannelMail Channel = "mail"
	// ChannelSms 短信。
	ChannelSms Channel = "sms"
	// ChannelVoice 语音。
	ChannelVoice Channel = "voice"
	// ChannelExtension 插件。
	ChannelExtension Channel = "extension"
	// ChannelApp App。
	ChannelApp Channel = "app"
	// ChannelClawBot 微信 ClawBot。
	ChannelClawBot Channel = "clawbot"
)

// Template 消息模板，JSON 序列化为字符串编码。
type Template string

const (
	TemplateHTML         Template = "html"
	TemplateTxt          Template = "txt"
	TemplateJSON         Template = "json"
	TemplateMarkdown     Template = "markdown"
	TemplateCloudMonitor Template = "cloudMonitor"
	TemplateJenkins      Template = "jenkins"
	TemplateRoute        Template = "route"
	TemplatePay          Template = "pay"
	// TemplateForm 表单格式模板；发送时需传 PushID（表单编码）。
	TemplateForm Template = "form"
	// TemplateDoc 文档格式模板（push 文档）；发送时需传 PushID。
	TemplateDoc Template = "doc"
	// TemplateExcel 表格格式模板（push 表格）；发送时需传 PushID。
	TemplateExcel Template = "excel"
)

// SendStatus 消息发送状态（整数编码）。
type SendStatus int

const (
	// SendStatusNotSent 未发送。
	SendStatusNotSent SendStatus = 0
	// SendStatusSending 发送中。
	SendStatusSending SendStatus = 1
	// SendStatusSuccess 发送成功。
	SendStatusSuccess SendStatus = 2
	// SendStatusFailed 发送失败。
	SendStatusFailed SendStatus = 3
)

// WebhookType webhook 渠道类型（整数编码）。
type WebhookType int

const (
	// WebhookTypeWorkWechatBot 企业微信机器人。
	WebhookTypeWorkWechatBot WebhookType = 1
	// WebhookTypeDingTalkBot 钉钉机器人。
	WebhookTypeDingTalkBot WebhookType = 2
	// WebhookTypeFeishuBot 飞书机器人。
	WebhookTypeFeishuBot WebhookType = 3
	// WebhookTypeServerChan Server酱。
	WebhookTypeServerChan WebhookType = 4
	// WebhookTypeWorkWechatApp 企业微信应用。
	WebhookTypeWorkWechatApp WebhookType = 6
	// WebhookTypeTencentLightLink 腾讯轻联。
	WebhookTypeTencentLightLink WebhookType = 7
	// WebhookTypeIFTTT IFTTT。
	WebhookTypeIFTTT WebhookType = 8
	// WebhookTypeJiJianYun 集简云。
	WebhookTypeJiJianYun WebhookType = 9
	// WebhookTypeGotify Gotify。
	WebhookTypeGotify WebhookType = 10
	// WebhookTypeWxPusher WxPusher。
	WebhookTypeWxPusher WebhookType = 11
	// WebhookTypeCustom 自定义。
	WebhookTypeCustom WebhookType = 12
	// WebhookTypeBark bark。
	WebhookTypeBark WebhookType = 50
)

// CallbackEvent 消息回调事件，JSON 序列化为字符串编码。
type CallbackEvent string

const (
	// CallbackEventMessageComplete 消息发送完成（官方拼写即为 message_complate）。
	CallbackEventMessageComplete CallbackEvent = "message_complate"
	// CallbackEventAddTopicUser 群组新增用户。
	CallbackEventAddTopicUser CallbackEvent = "add_topic_user"
	// CallbackEventAddFriend 新增好友。
	CallbackEventAddFriend CallbackEvent = "add_friend"
)

// ErrorCode PushPlus 业务返回码（整数编码），语义参考官方文档
// https://www.pushplus.plus/doc/guide/code.html。
type ErrorCode int

const (
	// ErrorCodeUnknown 未匹配任何已知业务码。
	ErrorCodeUnknown ErrorCode = -1
	// ErrorCodeOK 请求成功。
	ErrorCodeOK ErrorCode = 200
	// ErrorCodeNotLogin 未登录。
	ErrorCodeNotLogin ErrorCode = 302
	// ErrorCodeUnauthorized 请求未授权（AccessKey 无效或过期）。
	ErrorCodeUnauthorized ErrorCode = 401
	// ErrorCodeIPForbidden 禁止访问（IP 不在安全列表）。
	ErrorCodeIPForbidden ErrorCode = 403
	// ErrorCodeServerError 服务器内部错误。
	ErrorCodeServerError ErrorCode = 500
	// ErrorCodeDataError 数据处理异常。
	ErrorCodeDataError ErrorCode = 600
	// ErrorCodeForbiddenView 无权查看。
	ErrorCodeForbiddenView ErrorCode = 805
	// ErrorCodeInsufficientPoints 积分不足。
	ErrorCodeInsufficientPoints ErrorCode = 888
	// ErrorCodeRateLimited 请求次数过多。
	ErrorCodeRateLimited ErrorCode = 900
	// ErrorCodeInvalidToken token 无效。
	ErrorCodeInvalidToken ErrorCode = 903
	// ErrorCodeNotVerified 账号未实名认证。
	ErrorCodeNotVerified ErrorCode = 905
	// ErrorCodeValidationError 参数校验失败。
	ErrorCodeValidationError ErrorCode = 999
)

var knownErrorCodes = map[int]ErrorCode{
	200: ErrorCodeOK,
	302: ErrorCodeNotLogin,
	401: ErrorCodeUnauthorized,
	403: ErrorCodeIPForbidden,
	500: ErrorCodeServerError,
	600: ErrorCodeDataError,
	805: ErrorCodeForbiddenView,
	888: ErrorCodeInsufficientPoints,
	900: ErrorCodeRateLimited,
	903: ErrorCodeInvalidToken,
	905: ErrorCodeNotVerified,
	999: ErrorCodeValidationError,
}

// ErrorCodeOf 把整数业务码映射到 ErrorCode 枚举；未知返回 ErrorCodeUnknown。
func ErrorCodeOf(code int) ErrorCode {
	if ec, ok := knownErrorCodes[code]; ok {
		return ec
	}
	return ErrorCodeUnknown
}

// FormStatus push 表单状态。
type FormStatus int

const (
	// FormStatusDraft 草稿。
	FormStatusDraft FormStatus = 0
	// FormStatusCollecting 收集中。
	FormStatusCollecting FormStatus = 1
	// FormStatusStopped 已停止。
	FormStatusStopped FormStatus = 2
)

// SharePerm push 文档 / 表格分享权限。
type SharePerm int

const (
	// SharePermClosed 关闭分享。
	SharePermClosed SharePerm = 0
	// SharePermView 开启分享（仅可查看）。
	SharePermView SharePerm = 1
)

// ShareLogin push 文档 / 表格打开分享页是否需要登录。
type ShareLogin int

const (
	// ShareLoginAnonymous 免登录。
	ShareLoginAnonymous ShareLogin = 0
	// ShareLoginRequired 需登录。
	ShareLoginRequired ShareLogin = 1
)
