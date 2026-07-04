package pushplus

import (
	"context"
	"fmt"
)

// OpenMessageAPI 开放接口 - 消息（文档「二. 消息接口」）。
type OpenMessageAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newOpenMessageAPI(c *core, akm *AccessKeyManager) *OpenMessageAPI {
	return &OpenMessageAPI{core: c, akm: akm}
}

// List 消息列表。
func (a *OpenMessageAPI) List(ctx context.Context, query *PageQuery) (*PageResult[MessageItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[MessageItem]](ctx, a.core, a.akm, "POST", "/api/open/message/list", body)
}

// QueryResult 查询消息发送结果。
func (a *OpenMessageAPI) QueryResult(ctx context.Context, shortCode string) (*SendMessageResult, error) {
	path := appendQuery("/api/open/message/sendMessageResult", []queryParam{{"shortCode", shortCode}})
	return executeOpen[*SendMessageResult](ctx, a.core, a.akm, "GET", path, nil)
}

// Delete 删除消息。
func (a *OpenMessageAPI) Delete(ctx context.Context, shortCode string) (string, error) {
	path := appendQuery("/api/open/message/deleteMessage", []queryParam{{"shortCode", shortCode}})
	return executeOpen[string](ctx, a.core, a.akm, "DELETE", path, nil)
}

// DetailURL 拼接消息详情页 URL（不发起 HTTP 请求）。
func (a *OpenMessageAPI) DetailURL(shortCode string) string {
	return fmt.Sprintf("%s/shortMessage/%s", a.core.config.ResolveBaseURL(), shortCode)
}
