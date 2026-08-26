package pushplus

import (
	"context"
	"strconv"
)

// sendTypeQQGroup 发送到 QQ 群，目前渠道配置仅支持该类型。
const sendTypeQQGroup = 2

// QQBotAPI 开放接口 - QQ 机器人（文档「九. QQ机器人接口」）。
type QQBotAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newQQBotAPI(c *core, akm *AccessKeyManager) *QQBotAPI {
	return &QQBotAPI{core: c, akm: akm}
}

// GetBindLink 获取绑定链接与绑定码；refresh 为 true 时旧绑定码失效并重新生成。
func (a *QQBotAPI) GetBindLink(ctx context.Context, refresh bool) (*QQBotBindLink, error) {
	path := "/api/open/qqBot/getBindLink"
	if refresh {
		path = appendQuery(path, []queryParam{{"refresh", "true"}})
	}
	return executeOpen[*QQBotBindLink](ctx, a.core, a.akm, "GET", path, nil)
}

// BotInfo 查询绑定状态。
func (a *QQBotAPI) BotInfo(ctx context.Context) (*QQBotBindInfo, error) {
	return executeOpen[*QQBotBindInfo](ctx, a.core, a.akm, "GET", "/api/open/qqBot/botInfo", nil)
}

// Unbind 解绑 QQ 机器人。
func (a *QQBotAPI) Unbind(ctx context.Context) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", "/api/open/qqBot/unbind", nil)
	return err
}

// GroupList 获取机器人已加入的 QQ 群列表。
func (a *QQBotAPI) GroupList(ctx context.Context) ([]QQGroupItem, error) {
	return executeOpen[[]QQGroupItem](ctx, a.core, a.akm, "GET", "/api/open/qqBot/groupList", nil)
}

// List 获取 QQ 机器人渠道配置列表。
func (a *QQBotAPI) List(ctx context.Context, query *PageQuery) (*PageResult[QQBotItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[QQBotItem]](ctx, a.core, a.akm, "POST", "/api/open/qqBot/list", body)
}

// Add 新增 QQ 机器人渠道配置，用于把消息发送到指定 QQ 群；发给自己无需创建配置。
func (a *QQBotAPI) Add(ctx context.Context, req *QQBotSaveRequest) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/api/open/qqBot/add", withDefaultSendType(req))
	return err
}

// Edit 修改 QQ 机器人渠道配置；配置编码不可修改。
func (a *QQBotAPI) Edit(ctx context.Context, req *QQBotSaveRequest) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/api/open/qqBot/edit", withDefaultSendType(req))
	return err
}

// Delete 删除 QQ 机器人渠道配置。
func (a *QQBotAPI) Delete(ctx context.Context, id int64) error {
	path := appendQuery("/api/open/qqBot/delete", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "DELETE", path, nil)
	return err
}

// withDefaultSendType 复制一份请求并补上默认 sendType，避免修改调用方传入的结构体。
func withDefaultSendType(req *QQBotSaveRequest) *QQBotSaveRequest {
	if req == nil {
		return nil
	}
	copied := *req
	if copied.SendType == 0 {
		copied.SendType = sendTypeQQGroup
	}
	return &copied
}
