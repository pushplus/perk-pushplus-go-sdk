package pushplus

import (
	"context"
	"strconv"
)

// ForwardLogAPI 开放接口 - 消息规则触发记录（文档「十四. 消息规则接口」第 10、11 节）。
type ForwardLogAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newForwardLogAPI(c *core, akm *AccessKeyManager) *ForwardLogAPI {
	return &ForwardLogAPI{core: c, akm: akm}
}

// List 触发记录列表。
func (a *ForwardLogAPI) List(ctx context.Context, query *ForwardLogListQuery) (*PageResult[ForwardLogItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[ForwardLogItem]](ctx, a.core, a.akm, "POST", "/api/open/forwardLog/list", body)
}

// Detail 触发记录详情。
func (a *ForwardLogAPI) Detail(ctx context.Context, logID int64) (*ForwardLogDetail, error) {
	path := appendQuery("/api/open/forwardLog/detail", []queryParam{{"logId", strconv.FormatInt(logID, 10)}})
	return executeOpen[*ForwardLogDetail](ctx, a.core, a.akm, "GET", path, nil)
}
