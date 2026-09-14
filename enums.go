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
	// ChannelCmcc 新消息 ClawBot（中国移动 5G 消息）；仅支持中国移动用户。
	ChannelCmcc Channel = "cmcc"
	// ChannelQQ QQ 机器人；不带 Option 发给自己，Option 填配置编码则发到对应 QQ 群。
	ChannelQQ Channel = "qq"
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

// ForwardMode 消息规则总开关。
type ForwardMode int

const (
	// ForwardModeOff 关闭（推送与原来一致）。
	ForwardModeOff ForwardMode = 0
	// ForwardModeOnFallback 开启，未命中时仍按默认方式推送。
	ForwardModeOnFallback ForwardMode = 1
	// ForwardModeOnStrict 开启，未命中时不推送。
	ForwardModeOnStrict ForwardMode = 2
)

// ForwardSourceType 消息规则触发来源。
type ForwardSourceType int

const (
	// ForwardSourceTypeAll 全部。
	ForwardSourceTypeAll ForwardSourceType = 0
	// ForwardSourceTypeAPI 消息接口。
	ForwardSourceTypeAPI ForwardSourceType = 1
	// ForwardSourceTypeMail 邮件。
	ForwardSourceTypeMail ForwardSourceType = 2
)

// ForwardVarSourceType 模板变量来源。
type ForwardVarSourceType int

const (
	// ForwardVarSourceTypeHeader 请求头。
	ForwardVarSourceTypeHeader ForwardVarSourceType = 1
	// ForwardVarSourceTypeQuery Query 参数。
	ForwardVarSourceTypeQuery ForwardVarSourceType = 2
	// ForwardVarSourceTypeBody 请求体。
	ForwardVarSourceTypeBody ForwardVarSourceType = 3
	// ForwardVarSourceTypePath URL 路径。
	ForwardVarSourceTypePath ForwardVarSourceType = 4
	// ForwardVarSourceTypeSubject 主题（邮件）。
	ForwardVarSourceTypeSubject ForwardVarSourceType = 5
)

// ForwardExtractType 模板变量提取方式。
type ForwardExtractType int

const (
	// ForwardExtractTypeSerialized 序列化数据。
	ForwardExtractTypeSerialized ForwardExtractType = 1
	// ForwardExtractTypeRegex 正则表达式。
	ForwardExtractTypeRegex ForwardExtractType = 2
	// ForwardExtractTypeJSONPath JSONPath。
	ForwardExtractTypeJSONPath ForwardExtractType = 3
	// ForwardExtractTypeRaw 原始全文。
	ForwardExtractTypeRaw ForwardExtractType = 4
)

// ForwardMatchResult 触发记录匹配结果。
type ForwardMatchResult int

const (
	// ForwardMatchResultNotMatched 条件不满足。
	ForwardMatchResultNotMatched ForwardMatchResult = 0
	// ForwardMatchResultForwarded 已转发。
	ForwardMatchResultForwarded ForwardMatchResult = 1
	// ForwardMatchResultRateLimited 频率限制。
	ForwardMatchResultRateLimited ForwardMatchResult = 2
	// ForwardMatchResultOutOfTime 不在触发时间段。
	ForwardMatchResultOutOfTime ForwardMatchResult = 3
	// ForwardMatchResultError 执行异常。
	ForwardMatchResultError ForwardMatchResult = 4
)

// ForwardConditionOperator 图形化触发条件运算符。
type ForwardConditionOperator string

const (
	ForwardConditionOperatorEQ          ForwardConditionOperator = "eq"
	ForwardConditionOperatorNE          ForwardConditionOperator = "ne"
	ForwardConditionOperatorContains    ForwardConditionOperator = "contains"
	ForwardConditionOperatorNotContains ForwardConditionOperator = "notContains"
	ForwardConditionOperatorStartsWith  ForwardConditionOperator = "startsWith"
	ForwardConditionOperatorEndsWith    ForwardConditionOperator = "endsWith"
	ForwardConditionOperatorRegex       ForwardConditionOperator = "regex"
	ForwardConditionOperatorGT          ForwardConditionOperator = "gt"
	ForwardConditionOperatorGTE         ForwardConditionOperator = "gte"
	ForwardConditionOperatorLT          ForwardConditionOperator = "lt"
	ForwardConditionOperatorLTE         ForwardConditionOperator = "lte"
	ForwardConditionOperatorIn          ForwardConditionOperator = "in"
	ForwardConditionOperatorNotIn       ForwardConditionOperator = "notIn"
	ForwardConditionOperatorEmpty       ForwardConditionOperator = "empty"
	ForwardConditionOperatorNotEmpty    ForwardConditionOperator = "notEmpty"
)

// ForwardMessageType 消息规则发送目标的消息类型。也支持写成 {{变量名}}。
type ForwardMessageType string

const (
	ForwardMessageTypeOne    ForwardMessageType = "one"
	ForwardMessageTypeTopic  ForwardMessageType = "topic"
	ForwardMessageTypeFriend ForwardMessageType = "friend"
)
