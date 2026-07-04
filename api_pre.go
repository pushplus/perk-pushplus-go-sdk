package pushplus

import (
	"context"
	"strconv"
)

// PreAPI 开放接口 - 预处理信息（文档「十一. 预处理信息接口」，需会员）。
type PreAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newPreAPI(c *core, akm *AccessKeyManager) *PreAPI {
	return &PreAPI{core: c, akm: akm}
}

// List 预处理列表。
func (a *PreAPI) List(ctx context.Context, query *PageQuery) (*PageResult[PreItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[PreItem]](ctx, a.core, a.akm, "POST", "/api/open/pre/list", body)
}

// Detail 预处理详情。
func (a *PreAPI) Detail(ctx context.Context, preID int64) (*PreDetail, error) {
	path := appendQuery("/api/open/pre/detail", []queryParam{{"preId", strconv.FormatInt(preID, 10)}})
	return executeOpen[*PreDetail](ctx, a.core, a.akm, "GET", path, nil)
}

// Add 新增预处理，返回新预处理 ID。
func (a *PreAPI) Add(ctx context.Context, req *PreSaveRequest) (int64, error) {
	return executeOpen[int64](ctx, a.core, a.akm, "POST", "/api/open/pre/add", req)
}

// Edit 编辑预处理。
func (a *PreAPI) Edit(ctx context.Context, req *PreSaveRequest) (string, error) {
	return executeOpen[string](ctx, a.core, a.akm, "POST", "/api/open/pre/edit", req)
}

// Delete 删除预处理。
func (a *PreAPI) Delete(ctx context.Context, preID int64) (string, error) {
	path := appendQuery("/api/open/pre/delete", []queryParam{{"preId", strconv.FormatInt(preID, 10)}})
	return executeOpen[string](ctx, a.core, a.akm, "DELETE", path, nil)
}

// Test 测试预处理，返回处理后的消息。
func (a *PreAPI) Test(ctx context.Context, req *PreTestRequest) (string, error) {
	return executeOpen[string](ctx, a.core, a.akm, "POST", "/api/open/pre/test", req)
}
