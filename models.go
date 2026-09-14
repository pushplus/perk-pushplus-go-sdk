package pushplus

/* ============================== 通用模型 ============================== */

// PageQuery 通用分页查询参数。
type PageQuery struct {
	// Current 页码，从 1 开始。
	Current int `json:"current,omitempty"`
	// PageSize 每页大小，默认 20，最大 50。
	PageSize int `json:"pageSize,omitempty"`
}

// NewPageQuery 创建分页查询参数。
func NewPageQuery(current, pageSize int) *PageQuery {
	return &PageQuery{Current: current, PageSize: pageSize}
}

// PageResult 通用分页结果。
type PageResult[T any] struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
	Pages    int `json:"pages"`
	List     []T `json:"list"`
}

/* ============================== 发送消息 ============================== */

// SendRequest 发送单条消息请求，对应 /send 接口。
type SendRequest struct {
	// Token 可省略，SDK 会自动注入配置中的 token。
	Token string `json:"token,omitempty"`
	// Title 消息标题。
	Title string `json:"title,omitempty"`
	// Content 消息内容，必填。
	Content string `json:"content,omitempty"`
	// Topic 群组编码；不填仅发送给自己。
	Topic string `json:"topic,omitempty"`
	// Template 消息模板，默认 html。
	Template Template `json:"template,omitempty"`
	// Channel 发送渠道，默认 wechat。
	Channel Channel `json:"channel,omitempty"`
	// Option 渠道配置参数。
	Option string `json:"option,omitempty"`
	// CallbackURL 发送结果异步回调地址。
	CallbackURL string `json:"callbackUrl,omitempty"`
	// Timestamp 毫秒时间戳；超过该时间不再发送（时效控制）。
	Timestamp int64 `json:"timestamp,omitempty"`
	// To 好友令牌，多个用逗号分隔。
	To string `json:"to,omitempty"`
	// Pre 预处理编码。
	Pre string `json:"pre,omitempty"`
	// PushID push 编码；Template 为 form/doc/excel 时必传。
	PushID string `json:"pushId,omitempty"`
}

// BatchSendRequest 多渠道发送消息请求，对应 /batchSend 接口。
// Channel/Option 为逗号分隔字符串且按顺序一一对应，可使用 AddChannel 累积追加。
type BatchSendRequest struct {
	Token    string   `json:"token,omitempty"`
	Title    string   `json:"title,omitempty"`
	Content  string   `json:"content,omitempty"`
	Topic    string   `json:"topic,omitempty"`
	Template Template `json:"template,omitempty"`
	// Channel 多渠道，逗号分隔，如 "wechat,webhook"。
	Channel string `json:"channel,omitempty"`
	// Option 多渠道 option，逗号分隔，与 Channel 一一对应。
	Option      string `json:"option,omitempty"`
	CallbackURL string `json:"callbackUrl,omitempty"`
	Timestamp   int64  `json:"timestamp,omitempty"`
	To          string `json:"to,omitempty"`
	Pre         string `json:"pre,omitempty"`
	// PushID push 编码；Template 为 form/doc/excel 时必传。
	PushID string `json:"pushId,omitempty"`
}

// AddChannel 追加一个渠道及其 option，内部自动以逗号拼接（与官方文档示例语义一致）。
func (r *BatchSendRequest) AddChannel(ch Channel, option string) *BatchSendRequest {
	if r.Channel == "" {
		r.Channel = string(ch)
		r.Option = option
		return r
	}
	r.Channel = r.Channel + "," + string(ch)
	r.Option = r.Option + "," + option
	return r
}

// BatchSendResult 多渠道发送的单渠道结果。
type BatchSendResult struct {
	ShortCode string  `json:"shortCode"`
	Message   string  `json:"message"`
	Code      int     `json:"code"`
	Channel   Channel `json:"channel"`
}

/* ============================== AccessKey ============================== */

// AccessKeyResult 获取 AccessKey 的结果。
type AccessKeyResult struct {
	AccessKey string `json:"accessKey"`
	// ExpiresIn 有效期（秒）。
	ExpiresIn int64 `json:"expiresIn"`
}

/* ============================== 开放接口 - 消息 ============================== */

// MessageItem 消息列表项。
type MessageItem struct {
	TopicName   string  `json:"topicName"`
	MessageType int     `json:"messageType"`
	Title       string  `json:"title"`
	ShortCode   string  `json:"shortCode"`
	Channel     Channel `json:"channel"`
	UpdateTime  string  `json:"updateTime"`
}

// SendMessageResult 消息发送结果查询。
type SendMessageResult struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"errorMessage"`
	UpdateTime   string `json:"updateTime"`
}

// StatusEnum 把整数状态转为 SendStatus 枚举。
func (r *SendMessageResult) StatusEnum() SendStatus {
	return SendStatus(r.Status)
}

/* ============================== 开放接口 - 用户 ============================== */

// VipInfo 会员信息。
type VipInfo struct {
	// IsVip 是否会员；0-否，1-是。
	IsVip int `json:"isVip"`
	// LastDay 会员到期日。
	LastDay string `json:"lastDay"`
}

// UserInfo 当前用户信息。
type UserInfo struct {
	OpenID      string `json:"openId"`
	UnionID     string `json:"unionId"`
	NickName    string `json:"nickName"`
	HeadImgURL  string `json:"headImgUrl"`
	UserSex     int    `json:"userSex"`
	Token       string `json:"token"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	EmailStatus int    `json:"emailStatus"`
	Birthday    string `json:"birthday"`
	Points      int    `json:"points"`
	// VipInfo 会员信息。
	VipInfo *VipInfo `json:"vipInfo,omitempty"`
	// VerifyStatus 实名认证状态；0-未实名，1-已实名。
	VerifyStatus int `json:"verifyStatus"`
}

// UserLimitTime 用户限制发送时间。
type UserLimitTime struct {
	// SendLimit 1 无限制 / 2 短期限制 / 3 永久限制。
	SendLimit     int    `json:"sendLimit"`
	UserLimitTime string `json:"userLimitTime"`
}

// SendCount 用户当日发送次数统计。
type SendCount struct {
	WechatSendCount  int `json:"wechatSendCount"`
	CpSendCount      int `json:"cpSendCount"`
	WebhookSendCount int `json:"webhookSendCount"`
	MailSendCount    int `json:"mailSendCount"`
	QQBotSendCount   int `json:"qqBotSendCount"`
}

/* ============================== 开放接口 - 消息令牌 ============================== */

// MessageTokenItem 消息令牌列表项。
type MessageTokenItem struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	ExpireTime string `json:"expireTime"`
	Token      string `json:"token"`
}

// MessageTokenAddRequest 新增消息令牌请求。
type MessageTokenAddRequest struct {
	// Name 令牌名称，必填。
	Name string `json:"name,omitempty"`
	// ExpireTime 过期时间，如 "2030-01-01 00:00:00"。
	ExpireTime string `json:"expireTime,omitempty"`
}

// MessageTokenEditRequest 编辑消息令牌请求。
type MessageTokenEditRequest struct {
	ID         int64  `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	ExpireTime string `json:"expireTime,omitempty"`
}

// MessageTokenOption 消息令牌下拉选项。
type MessageTokenOption struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

/* ============================== 开放接口 - 群组 ============================== */

// TopicListQuery 群组列表查询参数。
type TopicListQuery struct {
	Current  int `json:"current,omitempty"`
	PageSize int `json:"pageSize,omitempty"`
	// Params 附加过滤参数，如 {"topicType": 0}。
	Params map[string]any `json:"params,omitempty"`
}

// NewTopicListQuery 创建群组列表查询参数；topicType 0 表示创建的群组。
func NewTopicListQuery(current, pageSize, topicType int) *TopicListQuery {
	return &TopicListQuery{
		Current:  current,
		PageSize: pageSize,
		Params:   map[string]any{"topicType": topicType},
	}
}

// TopicItem 群组列表项。
type TopicItem struct {
	Icon            string `json:"icon"`
	TopicID         int64  `json:"topicId"`
	TopicCode       string `json:"topicCode"`
	TopicName       string `json:"topicName"`
	NickName        string `json:"nickName"`
	CreateTime      string `json:"createTime"`
	TopicUserCount  int    `json:"topicUserCount"`
	TopicType       int    `json:"topicType"`
	IsApproved      int    `json:"isApproved"`
	FirstIsApproved int    `json:"firstIsApproved"`
	ApproveReason   string `json:"approveReason"`
	IsOpen          int    `json:"isOpen"`
}

// TopicDetail 群组详情。
type TopicDetail struct {
	TopicID         int64   `json:"topicId"`
	TopicName       string  `json:"topicName"`
	TopicCode       string  `json:"topicCode"`
	QrCodeImgURL    string  `json:"qrCodeImgUrl"`
	Contact         string  `json:"contact"`
	Introduction    string  `json:"introduction"`
	ReceiptMessage  string  `json:"receiptMessage"`
	NickName        string  `json:"nickName"`
	CreateTime      string  `json:"createTime"`
	TopicUserCount  int     `json:"topicUserCount"`
	Icon            string  `json:"icon"`
	AppID           string  `json:"appId"`
	TopicType       int     `json:"topicType"`
	Price           float64 `json:"price"`
	TopicDescribe   string  `json:"topicDescribe"`
	UserNickName    string  `json:"userNickName"`
	IsApproved      int     `json:"isApproved"`
	FirstIsApproved int     `json:"firstIsApproved"`
	ApproveReason   string  `json:"approveReason"`
	IsOpen          int     `json:"isOpen"`
}

// TopicAddRequest 新增群组请求。
type TopicAddRequest struct {
	TopicCode      string   `json:"topicCode,omitempty"`
	TopicName      string   `json:"topicName,omitempty"`
	Contact        string   `json:"contact,omitempty"`
	Introduction   string   `json:"introduction,omitempty"`
	ReceiptMessage string   `json:"receiptMessage,omitempty"`
	AppID          string   `json:"appId,omitempty"`
	Icon           string   `json:"icon,omitempty"`
	TopicType      int      `json:"topicType,omitempty"`
	Price          *float64 `json:"price,omitempty"`
	TopicDescribe  string   `json:"topicDescribe,omitempty"`
}

// TopicEditRequest 编辑群组请求。
type TopicEditRequest struct {
	// Topic 群组 ID，必填。
	Topic          int64    `json:"topic,omitempty"`
	TopicCode      string   `json:"topicCode,omitempty"`
	TopicName      string   `json:"topicName,omitempty"`
	Contact        string   `json:"contact,omitempty"`
	Introduction   string   `json:"introduction,omitempty"`
	ReceiptMessage string   `json:"receiptMessage,omitempty"`
	Icon           string   `json:"icon,omitempty"`
	Price          *float64 `json:"price,omitempty"`
	TopicDescribe  string   `json:"topicDescribe,omitempty"`
}

// TopicQrCode 群组二维码。
type TopicQrCode struct {
	QrCodeImgURL string `json:"qrCodeImgUrl"`
	// Forever 0 临时 / 1 永久。
	Forever int `json:"forever"`
}

/* ============================== 开放接口 - 群组用户 ============================== */

// TopicUserListQuery 群组用户列表查询参数。
type TopicUserListQuery struct {
	Current  int `json:"current,omitempty"`
	PageSize int `json:"pageSize,omitempty"`
	// Params 附加参数，如 {"topicId": 1}。
	Params map[string]any `json:"params,omitempty"`
}

// NewTopicUserListQuery 创建群组用户列表查询参数。
func NewTopicUserListQuery(current, pageSize int, topicID int64) *TopicUserListQuery {
	return &TopicUserListQuery{
		Current:  current,
		PageSize: pageSize,
		Params:   map[string]any{"topicId": topicID},
	}
}

// TopicUserItem 群组订阅用户。
type TopicUserItem struct {
	ID          int64  `json:"id"`
	NickName    string `json:"nickName"`
	OpenID      string `json:"openId"`
	HeadImgURL  string `json:"headImgUrl"`
	UserSex     int    `json:"userSex"`
	HavePhone   int    `json:"havePhone"`
	IsFollow    int    `json:"isFollow"`
	EmailStatus int    `json:"emailStatus"`
	FollowTime  string `json:"followTime"`
	Remark      string `json:"remark"`
}

// TopicUserBlacklistItem 群组订阅人黑名单列表项。
type TopicUserBlacklistItem struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"userId"`
	NickName   string `json:"nickName"`
	OpenID     string `json:"openId"`
	HeadImgURL string `json:"headImgUrl"`
	CreateTime string `json:"createTime"`
}

/* ============================== 开放接口 - Webhook ============================== */

// WebhookItem webhook 渠道配置。
type WebhookItem struct {
	ID              int64  `json:"id"`
	WebhookCode     string `json:"webhookCode"`
	WebhookName     string `json:"webhookName"`
	WebhookType     int    `json:"webhookType"`
	WebhookTypeName string `json:"webhookTypeName"`
	WebhookURL      string `json:"webhookUrl"`
	CreateTime      string `json:"createTime"`
	HTTPMethod      string `json:"httpMethod"`
	Headers         string `json:"headers"`
	Body            string `json:"body"`
}

// WebhookSaveRequest 新增/编辑 webhook 渠道配置请求。
type WebhookSaveRequest struct {
	ID          int64  `json:"id,omitempty"`
	WebhookCode string `json:"webhookCode,omitempty"`
	WebhookName string `json:"webhookName,omitempty"`
	WebhookType int    `json:"webhookType,omitempty"`
	WebhookURL  string `json:"webhookUrl,omitempty"`
	HTTPMethod  string `json:"httpMethod,omitempty"`
	Headers     string `json:"headers,omitempty"`
	Body        string `json:"body,omitempty"`
}

/* ============================== 开放接口 - 渠道配置 ============================== */

// MpItem 微信公众号渠道。
type MpItem struct {
	ID                 int64  `json:"id"`
	NickName           string `json:"nickName"`
	HeadImg            string `json:"headImg"`
	PrincipalName      string `json:"principalName"`
	AuthorizationAppid string `json:"authorizationAppid"`
	FuncInfo           string `json:"funcInfo"`
	ServiceType        int    `json:"serviceType"`
	VerifyType         int    `json:"verifyType"`
	Alias              string `json:"alias"`
	UpdateTime         string `json:"updateTime"`
}

// CpItem 企业微信渠道。
type CpItem struct {
	ID     int64  `json:"id"`
	CpName string `json:"cpName"`
	CpCode string `json:"cpCode"`
}

// MailItem 邮箱渠道。
type MailItem struct {
	ID       int64  `json:"id"`
	MailName string `json:"mailName"`
	MailCode string `json:"mailCode"`
}

// MailDetail 邮箱渠道详情。
type MailDetail struct {
	ID         int64  `json:"id"`
	MailName   string `json:"mailName"`
	MailCode   string `json:"mailCode"`
	Account    string `json:"account"`
	Password   string `json:"password"`
	SMTPServer string `json:"smtpServer"`
	SMTPSsl    string `json:"smtpSsl"`
	SMTPPort   string `json:"smtpPort"`
	CreateTime string `json:"createTime"`
}

/* ============================== 开放接口 - ClawBot ============================== */

// ClawBotQrCode ClawBot 登录二维码。
type ClawBotQrCode struct {
	URL    string `json:"url"`
	Qrcode string `json:"qrcode"`
}

// ClawBotInfo ClawBot 绑定信息。
type ClawBotInfo struct {
	CreateTime string `json:"createTime"`
	// HaveContextToken 是否已有上下文 token。
	HaveContextToken int `json:"haveContextToken"`
}

// ClawBotMessage ClawBot 消息。
type ClawBotMessage struct {
	// Type 1 文字 / 3 语音。
	Type int    `json:"type"`
	Text string `json:"text"`
}

/* ============================== 开放接口 - 新消息 ClawBot ============================== */

// CmccBindRequest 新消息 ClawBot 绑定请求。
type CmccBindRequest struct {
	// APIKey 中国移动新消息 Channel API Key，必须以 ak_ 或 app_ 开头。
	APIKey string `json:"apiKey"`
}

// CmccInfo 新消息 ClawBot 绑定状态。
type CmccInfo struct {
	// Bound 是否已绑定；0-未绑定，1-已绑定。
	Bound int `json:"bound"`
	// APIKeyMasked 脱敏后的 API Key。
	APIKeyMasked string `json:"apiKeyMasked"`
	CreateTime   string `json:"createTime"`
}

/* ============================== 开放接口 - QQ 机器人 ============================== */

// QQBotBindLink QQ 机器人绑定链接与绑定码。
type QQBotBindLink struct {
	// URL 带参分享链接，用于生成扫码二维码；已绑定用户再次获取时可能为空。
	URL string `json:"url"`
	// BindCode 绑定码。已是好友时扫码收不到加好友事件，需私聊发送该码；认领 QQ 群也用此码。
	BindCode string `json:"bindCode"`
	// ExpireSeconds 有效期秒数，默认 300。
	ExpireSeconds int `json:"expireSeconds"`
	// BotAppID 为当前用户分配的官方机器人 appId。
	BotAppID  string `json:"botAppId"`
	BotName   string `json:"botName"`
	BotAvatar string `json:"botAvatar"`
}

// QQBotInfo QQ 机器人详情。
type QQBotInfo struct {
	BotID    string `json:"botId"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	AppID    string `json:"appId"`
	// ShareURL 官方分享链接，可用于拉机器人进群。
	ShareURL string `json:"shareUrl"`
}

// QQBotBindInfo QQ 机器人绑定状态。
type QQBotBindInfo struct {
	// IsBind 0 未绑定，1 已绑定。
	IsBind int `json:"isBind"`
	// ReceiveStatus 1 可接收，0 用户已关闭单聊接收。
	ReceiveStatus int    `json:"receiveStatus"`
	CreateTime    string `json:"createTime"`
	// BotInfo 机器人详情，取不到时为 nil。
	BotInfo *QQBotInfo `json:"botInfo"`
}

// QQGroupItem 机器人已加入的 QQ 群。
type QQGroupItem struct {
	// ID 群编号；新增渠道配置时作为 QQGroupID 使用。
	ID          int64  `json:"id"`
	GroupOpenID string `json:"groupOpenId"`
	GroupRemark string `json:"groupRemark"`
	// Status 1 在群，2 群消息接收关闭。
	Status int `json:"status"`
	// GroupName 群名称，接口未授权时为空。
	GroupName       string   `json:"groupName"`
	GroupFingerMemo string   `json:"groupFingerMemo"`
	GroupClassText  string   `json:"groupClassText"`
	GroupTags       []string `json:"groupTags"`
	GroupMemberNum  int      `json:"groupMemberNum"`
	CreateTime      string   `json:"createTime"`
}

// QQBotItem QQ 机器人渠道配置列表项。
type QQBotItem struct {
	ID     int64  `json:"id"`
	QQName string `json:"qqName"`
	// QQCode 配置编码；发送消息时作为 Option 传入。
	QQCode string `json:"qqCode"`
	// SendType 2 发到 QQ 群。
	SendType    int    `json:"sendType"`
	QQGroupID   int64  `json:"qqGroupId"`
	GroupRemark string `json:"groupRemark"`
	GroupOpenID string `json:"groupOpenId"`
	GroupName   string `json:"groupName"`
	UpdateTime  string `json:"updateTime"`
}

// QQBotSaveRequest 新增/修改 QQ 机器人渠道配置请求。
type QQBotSaveRequest struct {
	// ID 修改时必填。
	ID int64 `json:"id,omitempty"`
	// QQName 配置名称，必填，最多 64 个字符。
	QQName string `json:"qqName,omitempty"`
	// QQCode 配置编码，新增必填；仅支持字母、数字、下划线和中划线，创建后不可修改。
	QQCode string `json:"qqCode,omitempty"`
	// SendType 发送类型；留空时 SDK 自动填 2（发到 QQ 群）。
	SendType int `json:"sendType,omitempty"`
	// QQGroupID QQ 群编号，必填，取自 GroupList 返回的 ID。
	QQGroupID int64 `json:"qqGroupId,omitempty"`
}

/* ============================== 开放接口 - 功能设置 ============================== */

// UserDefaultItem 用户默认设置列表项。
type UserDefaultItem struct {
	ID         int64   `json:"id"`
	Channel    Channel `json:"channel"`
	ChannelTxt string  `json:"channelTxt"`
	UpdateTime string  `json:"updateTime"`
	Name       string  `json:"name"`
}

// UserDefaultDetail 用户默认设置详情。
type UserDefaultDetail struct {
	ID         int64   `json:"id"`
	Channel    Channel `json:"channel"`
	Option     string  `json:"option"`
	Pre        string  `json:"pre"`
	UpdateTime string  `json:"updateTime"`
	Name       string  `json:"name"`
	TokenID    int64   `json:"tokenId"`
}

// UserDefaultSaveRequest 新增/编辑用户默认设置请求。
type UserDefaultSaveRequest struct {
	ID      int64   `json:"id,omitempty"`
	Channel Channel `json:"channel,omitempty"`
	Option  string  `json:"option,omitempty"`
	Pre     string  `json:"pre,omitempty"`
	// TokenID 0 表示用户令牌。
	TokenID *int64 `json:"tokenId,omitempty"`
}

/* ============================== 开放接口 - 好友 ============================== */

// FriendItem 好友列表项。
type FriendItem struct {
	ID          int64  `json:"id"`
	FriendID    int64  `json:"friendId"`
	Token       string `json:"token"`
	HeadImgURL  string `json:"headImgUrl"`
	NickName    string `json:"nickName"`
	EmailStatus int    `json:"emailStatus"`
	HavePhone   int    `json:"havePhone"`
	IsFollow    int    `json:"isFollow"`
	Remark      string `json:"remark"`
	CreateTime  string `json:"createTime"`
}

// FriendQrCode 加好友二维码。
type FriendQrCode struct {
	QrCodeImgURL string `json:"qrCodeImgUrl"`
}

// FriendBlacklistItem 好友黑名单列表项。
type FriendBlacklistItem struct {
	ID         int64  `json:"id"`
	FriendID   int64  `json:"friendId"`
	NickName   string `json:"nickName"`
	HeadImgURL string `json:"headImgUrl"`
	CreateTime string `json:"createTime"`
}

/* ============================== 开放接口 - 预处理 ============================== */

// PreItem 预处理列表项。
type PreItem struct {
	ID      int64  `json:"id"`
	PreName string `json:"preName"`
	PreCode string `json:"preCode"`
	// ContentType 1 表示 JavaScript。
	ContentType int    `json:"contentType"`
	CreateTime  string `json:"createTime"`
}

// PreDetail 预处理详情。
type PreDetail struct {
	ID          int64  `json:"id"`
	PreName     string `json:"preName"`
	PreCode     string `json:"preCode"`
	ContentType int    `json:"contentType"`
	Content     string `json:"content"`
}

// PreSaveRequest 新增/编辑预处理请求。
type PreSaveRequest struct {
	ID          int64  `json:"id,omitempty"`
	Content     string `json:"content,omitempty"`
	PreName     string `json:"preName,omitempty"`
	PreCode     string `json:"preCode,omitempty"`
	ContentType int    `json:"contentType,omitempty"`
}

// PreTestRequest 预处理测试请求。
type PreTestRequest struct {
	Content     string `json:"content,omitempty"`
	ContentType int    `json:"contentType,omitempty"`
	Message     string `json:"message,omitempty"`
}

/* ============================== 开放接口 - 图片服务 ============================== */

// ImageUploadToken 图片上传凭证。
type ImageUploadToken struct {
	UploadToken string `json:"uploadToken"`
	UploadHost  string `json:"uploadHost"`
	UploadURL   string `json:"uploadUrl"`
	Bucket      string `json:"bucket"`
	ExpiresIn   int    `json:"expiresIn"`
}

// ImageUploadResult 七牛云上传结果（直接返回，不是 PushPlus 统一响应结构）。
type ImageUploadResult struct {
	// Errno 0 表示成功。
	Errno     int    `json:"errno"`
	Ext       string `json:"ext"`
	Fname     string `json:"fname"`
	Fsize     int64  `json:"fsize"`
	Hash      string `json:"hash"`
	Key       string `json:"key"`
	MimeType  string `json:"mimeType"`
	Msg       string `json:"msg"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
}

// ImageItem 已上传图片列表项。
type ImageItem struct {
	ID         int64  `json:"id"`
	ImgURL     string `json:"imgUrl"`
	Thumbnail  string `json:"thumbnail"`
	CreateTime string `json:"createTime"`
}

/* ============================== 回调模型 ============================== */

// CallbackPayload 消息回调统一载荷。
type CallbackPayload struct {
	Event CallbackEvent `json:"event"`
	// MessageInfo 事件为 message_complate 时有值。
	MessageInfo *MessageCompleteInfo `json:"messageInfo,omitempty"`
	// TopicUserInfo 事件为 add_topic_user 时有值。
	TopicUserInfo *TopicUserInfo `json:"topicUserInfo,omitempty"`
	// FriendInfo 事件为 add_friend 时有值。
	FriendInfo *FriendInfo `json:"friendInfo,omitempty"`
	// QrCode 事件为 add_friend 时的二维码标识。
	QrCode string `json:"qrCode,omitempty"`
}

// MessageCompleteInfo 消息发送完成回调数据。
type MessageCompleteInfo struct {
	Message    string `json:"message"`
	ShortCode  string `json:"shortCode"`
	SendStatus int    `json:"sendStatus"`
}

// SendStatusEnum 把整数状态转为 SendStatus 枚举。
func (m *MessageCompleteInfo) SendStatusEnum() SendStatus {
	return SendStatus(m.SendStatus)
}

// TopicUserInfo 群组新增用户回调数据。
type TopicUserInfo struct {
	ID          int64  `json:"id"`
	OpenID      string `json:"openId"`
	TopicID     int64  `json:"topicId"`
	UserSex     int    `json:"userSex"`
	IsFollow    int    `json:"isFollow"`
	NickName    string `json:"nickName"`
	HavePhone   int    `json:"havePhone"`
	TopicCode   string `json:"topicCode"`
	TopicName   string `json:"topicName"`
	HeadImgURL  string `json:"headImgUrl"`
	EmailStatus int    `json:"emailStatus"`
}

// FriendInfo 新增好友回调数据。
type FriendInfo struct {
	Token       string `json:"token"`
	FriendID    int64  `json:"friendId"`
	IsFollow    int    `json:"isFollow"`
	NickName    string `json:"nickName"`
	HavePhone   int    `json:"havePhone"`
	CreateTime  string `json:"createTime"`
	EmailStatus int    `json:"emailStatus"`
}

/* ============================== 开放接口 - form / doc / excel ============================== */

// FormListQuery 我的表单分页查询。官方结构是 {current, pageSize, params:{keyword, status}}。
type FormListQuery struct {
	Current  int            `json:"current,omitempty"`
	PageSize int            `json:"pageSize,omitempty"`
	Params   map[string]any `json:"params,omitempty"`
}

// NewFormListQuery 创建表单分页查询。
func NewFormListQuery(current, pageSize int) *FormListQuery {
	return &FormListQuery{Current: current, PageSize: pageSize}
}

// NewFormListQueryFilter 创建带关键词 / 状态筛选的表单分页查询。
func NewFormListQueryFilter(current, pageSize int, keyword string, status *int) *FormListQuery {
	params := map[string]any{}
	if keyword != "" {
		params["keyword"] = keyword
	}
	if status != nil {
		params["status"] = *status
	}
	q := &FormListQuery{Current: current, PageSize: pageSize}
	if len(params) > 0 {
		q.Params = params
	}
	return q
}

// FormCover 表单封面页配置。
type FormCover struct {
	Enabled    *bool  `json:"enabled,omitempty"`
	Image      string `json:"image,omitempty"`
	ButtonText string `json:"buttonText,omitempty"`
}

// FormTheme 表单主题外观。
type FormTheme struct {
	PrimaryColor    string     `json:"primaryColor,omitempty"`
	BackgroundColor string     `json:"backgroundColor,omitempty"`
	HeaderImage     string     `json:"headerImage,omitempty"`
	BackgroundImage string     `json:"backgroundImage,omitempty"`
	Cover           *FormCover `json:"cover,omitempty"`
}

// FormSettings 表单收集 / 展示设置。
type FormSettings struct {
	EndTime            *string `json:"endTime,omitempty"`
	MaxResponses       *int    `json:"maxResponses,omitempty"`
	OncePerUser        *bool   `json:"oncePerUser,omitempty"`
	AllowAnonymous     *bool   `json:"allowAnonymous,omitempty"`
	Password           string  `json:"password,omitempty"`
	ShowQuestionNumber *bool   `json:"showQuestionNumber,omitempty"`
	OnePerPage         *bool   `json:"onePerPage,omitempty"`
	ShowPrevButton     *bool   `json:"showPrevButton,omitempty"`
	HideTitle          *bool   `json:"hideTitle,omitempty"`
	HideCopyright      *bool   `json:"hideCopyright,omitempty"`
	HideAd             *bool   `json:"hideAd,omitempty"`
	ShowOutline        *bool   `json:"showOutline,omitempty"`
	ThankText          string  `json:"thankText,omitempty"`
	RedirectEnabled    *bool   `json:"redirectEnabled,omitempty"`
	RedirectURL        string  `json:"redirectUrl,omitempty"`
	AllowEdit          *bool   `json:"allowEdit,omitempty"`
}

// FormListItem 表单列表项 / 创建、复制结果。
type FormListItem struct {
	ID            int64  `json:"id"`
	FormCode      string `json:"formCode,omitempty"`
	FillURL       string `json:"fillUrl,omitempty"`
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Status        int    `json:"status"`
	ResponseCount int    `json:"responseCount"`
	PublishTime   string `json:"publishTime,omitempty"`
	CreateTime    string `json:"createTime,omitempty"`
	UpdateTime    string `json:"updateTime,omitempty"`
}

// FormSaveRequest 保存表单设计。Items 为题目列表，每题至少含 id、type、label。
type FormSaveRequest struct {
	ID          int64            `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description,omitempty"`
	Items       []map[string]any `json:"items,omitempty"`
	Theme       *FormTheme       `json:"theme,omitempty"`
	Settings    *FormSettings    `json:"settings,omitempty"`
}

// FormDetail 表单详情。
type FormDetail struct {
	ID            int64            `json:"id"`
	FormCode      string           `json:"formCode,omitempty"`
	FillURL       string           `json:"fillUrl,omitempty"`
	Title         string           `json:"title,omitempty"`
	Description   string           `json:"description,omitempty"`
	Items         []map[string]any `json:"items,omitempty"`
	Theme         *FormTheme       `json:"theme,omitempty"`
	Settings      *FormSettings    `json:"settings,omitempty"`
	Status        int              `json:"status"`
	PublishDirty  bool             `json:"publishDirty"`
	ResponseCount int              `json:"responseCount"`
	PublishTime   string           `json:"publishTime,omitempty"`
	CreateTime    string           `json:"createTime,omitempty"`
	UpdateTime    string           `json:"updateTime,omitempty"`
}

// FormPublishDiff 草稿题目与发布快照差异。
type FormPublishDiff struct {
	Dirty         bool     `json:"dirty"`
	Breaking      bool     `json:"breaking"`
	ResponseCount int      `json:"responseCount"`
	Added         []string `json:"added"`
	Removed       []string `json:"removed"`
	TypeChanged   []string `json:"typeChanged"`
	OptionChanged []string `json:"optionChanged"`
}

// FormPublishResult 发布表单结果。
type FormPublishResult struct {
	ID             int64  `json:"id"`
	FormCode       string `json:"formCode,omitempty"`
	FillURL        string `json:"fillUrl,omitempty"`
	Title          string `json:"title,omitempty"`
	Status         int    `json:"status"`
	PreviousStatus int    `json:"previousStatus"`
	PublishDirty   bool   `json:"publishDirty"`
	PublishTime    string `json:"publishTime,omitempty"`
}

// DocListQuery 文档 / 表格分页查询。官方结构是 {current, pageSize, params:{keyword, shareEnabled}}。
type DocListQuery struct {
	Current  int            `json:"current,omitempty"`
	PageSize int            `json:"pageSize,omitempty"`
	Params   map[string]any `json:"params,omitempty"`
}

// NewDocListQuery 创建文档 / 表格分页查询。
func NewDocListQuery(current, pageSize int) *DocListQuery {
	return &DocListQuery{Current: current, PageSize: pageSize}
}

// DocListItem 文档 / 表格列表项。
type DocListItem struct {
	DocCode     string `json:"docCode,omitempty"`
	ShareURL    string `json:"shareUrl,omitempty"`
	Title       string `json:"title,omitempty"`
	SharePerm   int    `json:"sharePerm"`
	ShareLogin  int    `json:"shareLogin"`
	Perm        int    `json:"perm"`
	Published   bool   `json:"published"`
	PublishTime string `json:"publishTime,omitempty"`
	CreateTime  string `json:"createTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
}

// DocVo 文档信息（不含正文）。
type DocVo struct {
	DocCode      string `json:"docCode,omitempty"`
	ShareURL     string `json:"shareUrl,omitempty"`
	Title        string `json:"title,omitempty"`
	SharePerm    int    `json:"sharePerm"`
	ShareLogin   int    `json:"shareLogin"`
	Perm         int    `json:"perm"`
	Published    bool   `json:"published"`
	PublishDirty bool   `json:"publishDirty"`
	PublishTime  string `json:"publishTime,omitempty"`
	CreateTime   string `json:"createTime,omitempty"`
	UpdateTime   string `json:"updateTime,omitempty"`
}

// DocContent 文档内容（HTML 草稿）。
type DocContent struct {
	DocVo
	Content string `json:"content,omitempty"`
}

// ExcelVo 表格信息（不含正文）。
type ExcelVo struct {
	DocCode      string `json:"docCode,omitempty"`
	ShareURL     string `json:"shareUrl,omitempty"`
	Title        string `json:"title,omitempty"`
	SharePerm    int    `json:"sharePerm"`
	ShareLogin   int    `json:"shareLogin"`
	Perm         int    `json:"perm"`
	Published    bool   `json:"published"`
	PublishDirty bool   `json:"publishDirty"`
	PublishTime  string `json:"publishTime,omitempty"`
	CreateTime   string `json:"createTime,omitempty"`
	UpdateTime   string `json:"updateTime,omitempty"`
}

// ExcelContent 表格内容（整表 JSON 字符串草稿）。
type ExcelContent struct {
	ExcelVo
	Content string `json:"content,omitempty"`
}

/* ============================== 开放接口 - 消息规则 ============================== */

// ForwardConditionItem 图形化触发条件中的单条比较。
type ForwardConditionItem struct {
	VarName  string `json:"varName,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

// ForwardCondition 图形化触发条件。
type ForwardCondition struct {
	Logic string                 `json:"logic,omitempty"`
	Items []ForwardConditionItem `json:"items,omitempty"`
}

// ForwardVariable 模板变量。
type ForwardVariable struct {
	ID           int64  `json:"id,omitempty"`
	RuleID       int64  `json:"ruleId,omitempty"`
	VarName      string `json:"varName,omitempty"`
	SourceType   int    `json:"sourceType,omitempty"`
	ExtractType  int    `json:"extractType,omitempty"`
	ExtractKey   string `json:"extractKey,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty"`
	Sort         int    `json:"sort,omitempty"`
}

// ForwardTarget 发送目标。
type ForwardTarget struct {
	ID          int64  `json:"id,omitempty"`
	RuleID      int64  `json:"ruleId,omitempty"`
	Channel     string `json:"channel,omitempty"`
	Option      string `json:"option,omitempty"`
	MessageType string `json:"messageType,omitempty"`
	Topic       string `json:"topic,omitempty"`
	To          string `json:"to,omitempty"`
	Sort        int    `json:"sort,omitempty"`
}

// ForwardRuleItem 消息规则列表项。
type ForwardRuleItem struct {
	ID             int64  `json:"id"`
	TokenID        int64  `json:"tokenId"`
	TokenName      string `json:"tokenName"`
	RuleName       string `json:"ruleName"`
	SourceType     int    `json:"sourceType"`
	SourceTypeName string `json:"sourceTypeName"`
	Status         int    `json:"status"`
	Sort           int    `json:"sort"`
	ConditionExpr  string `json:"conditionExpr"`
	TargetCount    int    `json:"targetCount"`
	CreateTime     string `json:"createTime"`
}

// ForwardRuleDetail 消息规则详情。
type ForwardRuleDetail struct {
	ID              int64             `json:"id"`
	TokenID         int64             `json:"tokenId"`
	RuleName        string            `json:"ruleName"`
	SourceType      int               `json:"sourceType"`
	Status          int               `json:"status"`
	Sort            int               `json:"sort"`
	ConditionExpr   string            `json:"conditionExpr"`
	Condition       *ForwardCondition `json:"condition,omitempty"`
	TitleTemplate   string            `json:"titleTemplate"`
	ContentTemplate string            `json:"contentTemplate"`
	Template        string            `json:"template"`
	Pre             string            `json:"pre"`
	StopOnMatch     int               `json:"stopOnMatch"`
	LimitPeriod     int               `json:"limitPeriod"`
	LimitCount      int               `json:"limitCount"`
	ActiveStartTime string            `json:"activeStartTime"`
	ActiveEndTime   string            `json:"activeEndTime"`
	ActiveWeekdays  string            `json:"activeWeekdays"`
	Remark          string            `json:"remark"`
	Variables       []ForwardVariable `json:"variables"`
	Targets         []ForwardTarget   `json:"targets"`
	CreateTime      string            `json:"createTime"`
}

// ForwardRuleSaveRequest 新增 / 修改消息规则。修改时 ID 必填。
// TokenID / Status 等使用指针，以便正确发送 0（用户令牌 / 停用）。
type ForwardRuleSaveRequest struct {
	ID              int64             `json:"id,omitempty"`
	RuleName        string            `json:"ruleName,omitempty"`
	TokenID         *int64            `json:"tokenId,omitempty"`
	SourceType      *int              `json:"sourceType,omitempty"`
	Status          *int              `json:"status,omitempty"`
	Sort            *int              `json:"sort,omitempty"`
	Condition       *ForwardCondition `json:"condition,omitempty"`
	ConditionExpr   string            `json:"conditionExpr,omitempty"`
	TitleTemplate   string            `json:"titleTemplate,omitempty"`
	ContentTemplate string            `json:"contentTemplate,omitempty"`
	Template        string            `json:"template,omitempty"`
	Pre             string            `json:"pre,omitempty"`
	StopOnMatch     *int              `json:"stopOnMatch,omitempty"`
	LimitPeriod     *int              `json:"limitPeriod,omitempty"`
	LimitCount      *int              `json:"limitCount,omitempty"`
	ActiveStartTime string            `json:"activeStartTime,omitempty"`
	ActiveEndTime   string            `json:"activeEndTime,omitempty"`
	ActiveWeekdays  string            `json:"activeWeekdays,omitempty"`
	Remark          string            `json:"remark,omitempty"`
	Variables       []ForwardVariable `json:"variables,omitempty"`
	Targets         []ForwardTarget   `json:"targets,omitempty"`
}

// ForwardRuleTestRequest 测试消息规则请求。不会真正发送消息。
type ForwardRuleTestRequest struct {
	SourceType      *int              `json:"sourceType,omitempty"`
	ContentType     string            `json:"contentType,omitempty"`
	Headers         map[string]any    `json:"headers,omitempty"`
	Query           map[string]any    `json:"query,omitempty"`
	Body            string            `json:"body,omitempty"`
	Title           string            `json:"title,omitempty"`
	MailFrom        string            `json:"mailFrom,omitempty"`
	MailTo          string            `json:"mailTo,omitempty"`
	MailCc          string            `json:"mailCc,omitempty"`
	Condition       *ForwardCondition `json:"condition,omitempty"`
	ConditionExpr   string            `json:"conditionExpr,omitempty"`
	TitleTemplate   string            `json:"titleTemplate,omitempty"`
	ContentTemplate string            `json:"contentTemplate,omitempty"`
	Template        string            `json:"template,omitempty"`
	Pre             string            `json:"pre,omitempty"`
	Variables       []ForwardVariable `json:"variables,omitempty"`
}

// ForwardRuleTestResult 测试消息规则结果。
type ForwardRuleTestResult struct {
	Variables     map[string]any `json:"variables"`
	Matched       bool           `json:"matched"`
	ConditionExpr string         `json:"conditionExpr"`
	ErrorMessage  string         `json:"errorMessage"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	Template      string         `json:"template"`
}

// ForwardRuleSetting 消息规则总开关。
type ForwardRuleSetting struct {
	// Mode 0-关闭，1-开启且未命中仍推送，2-开启且未命中不推送。
	Mode int `json:"mode"`
}

// ForwardLogListQuery 触发记录分页查询。官方结构是 {current, pageSize, params:{ruleId, matchResult}}。
type ForwardLogListQuery struct {
	Current  int            `json:"current,omitempty"`
	PageSize int            `json:"pageSize,omitempty"`
	Params   map[string]any `json:"params,omitempty"`
}

// NewForwardLogListQuery 创建触发记录分页查询。
func NewForwardLogListQuery(current, pageSize int) *ForwardLogListQuery {
	return &ForwardLogListQuery{Current: current, PageSize: pageSize}
}

// NewForwardLogListQueryFilter 创建带规则编号 / 匹配结果筛选的触发记录分页查询。
func NewForwardLogListQueryFilter(current, pageSize int, ruleID *int64, matchResult *int) *ForwardLogListQuery {
	params := map[string]any{}
	if ruleID != nil {
		params["ruleId"] = *ruleID
	}
	if matchResult != nil {
		params["matchResult"] = *matchResult
	}
	q := &ForwardLogListQuery{Current: current, PageSize: pageSize}
	if len(params) > 0 {
		q.Params = params
	}
	return q
}

// ForwardLogItem 触发记录列表项。
type ForwardLogItem struct {
	ID              int64  `json:"id"`
	RuleID          int64  `json:"ruleId"`
	RuleName        string `json:"ruleName"`
	SourceType      int    `json:"sourceType"`
	SourceTypeName  string `json:"sourceTypeName"`
	RequestIP       string `json:"requestIp"`
	MatchResult     int    `json:"matchResult"`
	MatchResultName string `json:"matchResultName"`
	ShortCodes      string `json:"shortCodes"`
	ErrorMessage    string `json:"errorMessage"`
	CreateTime      string `json:"createTime"`
}

// ForwardLogDetail 触发记录详情。
type ForwardLogDetail struct {
	ID              int64  `json:"id"`
	RuleID          int64  `json:"ruleId"`
	RuleName        string `json:"ruleName"`
	SourceType      int    `json:"sourceType"`
	SourceTypeName  string `json:"sourceTypeName"`
	RequestIP       string `json:"requestIp"`
	RequestMethod   string `json:"requestMethod"`
	RequestHeaders  string `json:"requestHeaders"`
	RequestQuery    string `json:"requestQuery"`
	RequestBody     string `json:"requestBody"`
	Variables       string `json:"variables"`
	MatchResult     int    `json:"matchResult"`
	MatchResultName string `json:"matchResultName"`
	ShortCodes      string `json:"shortCodes"`
	ErrorMessage    string `json:"errorMessage"`
	CreateTime      string `json:"createTime"`
}
