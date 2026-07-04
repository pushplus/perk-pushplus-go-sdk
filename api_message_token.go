package pushplus

import (
	"context"
	"strconv"
)

// MessageTokenAPI 开放接口 - 消息 token（文档「四. 消息 token 接口」）。
type MessageTokenAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newMessageTokenAPI(c *core, akm *AccessKeyManager) *MessageTokenAPI {
	return &MessageTokenAPI{core: c, akm: akm}
}

// List 消息令牌列表。
func (a *MessageTokenAPI) List(ctx context.Context, query *PageQuery) (*PageResult[MessageTokenItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[MessageTokenItem]](ctx, a.core, a.akm, "POST", "/api/open/token/list", body)
}

// Add 新增消息令牌，返回新令牌。
func (a *MessageTokenAPI) Add(ctx context.Context, req *MessageTokenAddRequest) (string, error) {
	return executeOpen[string](ctx, a.core, a.akm, "POST", "/api/open/token/add", req)
}

// Edit 编辑消息令牌。
func (a *MessageTokenAPI) Edit(ctx context.Context, req *MessageTokenEditRequest) (string, error) {
	return executeOpen[string](ctx, a.core, a.akm, "POST", "/api/open/token/edit", req)
}

// Delete 删除消息令牌。
func (a *MessageTokenAPI) Delete(ctx context.Context, id int64) (string, error) {
	path := appendQuery("/api/open/token/deleteToken", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	return executeOpen[string](ctx, a.core, a.akm, "DELETE", path, nil)
}

// SelectList 消息令牌下拉选项列表；tokenType 为 0 时表示默认类型。
func (a *MessageTokenAPI) SelectList(ctx context.Context, tokenType int) ([]MessageTokenOption, error) {
	path := appendQuery("/api/open/token/selectTokenList", []queryParam{{"type", strconv.Itoa(tokenType)}})
	return executeOpen[[]MessageTokenOption](ctx, a.core, a.akm, "GET", path, nil)
}
