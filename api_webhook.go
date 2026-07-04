package pushplus

import (
	"context"
	"strconv"
)

// WebhookAPI 开放接口 - 渠道配置 webhook（文档「七. 渠道配置 - webhook」）。
type WebhookAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newWebhookAPI(c *core, akm *AccessKeyManager) *WebhookAPI {
	return &WebhookAPI{core: c, akm: akm}
}

// List webhook 列表。
func (a *WebhookAPI) List(ctx context.Context, query *PageQuery) (*PageResult[WebhookItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[WebhookItem]](ctx, a.core, a.akm, "POST", "/api/open/webhook/list", body)
}

// Detail webhook 详情。
func (a *WebhookAPI) Detail(ctx context.Context, webhookID int64) (*WebhookItem, error) {
	path := appendQuery("/api/open/webhook/detail", []queryParam{{"webhookId", strconv.FormatInt(webhookID, 10)}})
	return executeOpen[*WebhookItem](ctx, a.core, a.akm, "GET", path, nil)
}

// Add 新增 webhook，返回新配置 ID。
func (a *WebhookAPI) Add(ctx context.Context, req *WebhookSaveRequest) (int64, error) {
	return executeOpen[int64](ctx, a.core, a.akm, "POST", "/api/open/webhook/add", req)
}

// Edit 编辑 webhook。
func (a *WebhookAPI) Edit(ctx context.Context, req *WebhookSaveRequest) (string, error) {
	return executeOpen[string](ctx, a.core, a.akm, "POST", "/api/open/webhook/edit", req)
}
