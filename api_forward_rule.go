package pushplus

import (
	"context"
	"strconv"
)

// ForwardRuleAPI 开放接口 - 消息规则（文档「十四. 消息规则接口」，需会员）。
type ForwardRuleAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newForwardRuleAPI(c *core, akm *AccessKeyManager) *ForwardRuleAPI {
	return &ForwardRuleAPI{core: c, akm: akm}
}

// List 消息规则列表。
func (a *ForwardRuleAPI) List(ctx context.Context, query *PageQuery) (*PageResult[ForwardRuleItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[ForwardRuleItem]](ctx, a.core, a.akm, "POST", "/api/open/forwardRule/list", body)
}

// Detail 消息规则详情。
func (a *ForwardRuleAPI) Detail(ctx context.Context, ruleID int64) (*ForwardRuleDetail, error) {
	path := appendQuery("/api/open/forwardRule/detail", []queryParam{{"ruleId", strconv.FormatInt(ruleID, 10)}})
	return executeOpen[*ForwardRuleDetail](ctx, a.core, a.akm, "GET", path, nil)
}

// Add 新增消息规则。
func (a *ForwardRuleAPI) Add(ctx context.Context, req *ForwardRuleSaveRequest) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/api/open/forwardRule/add", req)
	return err
}

// Edit 修改消息规则；会整体覆盖模板变量和发送目标。
func (a *ForwardRuleAPI) Edit(ctx context.Context, req *ForwardRuleSaveRequest) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/api/open/forwardRule/edit", req)
	return err
}

// Delete 删除消息规则。
func (a *ForwardRuleAPI) Delete(ctx context.Context, ruleID int64) error {
	path := appendQuery("/api/open/forwardRule/delete", []queryParam{{"ruleId", strconv.FormatInt(ruleID, 10)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "DELETE", path, nil)
	return err
}

// ChangeStatus 启用 / 停用消息规则；1 启用，0 停用。
func (a *ForwardRuleAPI) ChangeStatus(ctx context.Context, ruleID int64, status int) error {
	path := appendQuery("/api/open/forwardRule/changeStatus", []queryParam{
		{"ruleId", strconv.FormatInt(ruleID, 10)},
		{"status", strconv.Itoa(status)},
	})
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", path, nil)
	return err
}

// Test 用模拟请求测试规则，不会真正发送消息。
func (a *ForwardRuleAPI) Test(ctx context.Context, req *ForwardRuleTestRequest) (*ForwardRuleTestResult, error) {
	return executeOpen[*ForwardRuleTestResult](ctx, a.core, a.akm, "POST", "/api/open/forwardRule/test", req)
}

// GetSetting 获取消息规则总开关模式。
func (a *ForwardRuleAPI) GetSetting(ctx context.Context) (*ForwardRuleSetting, error) {
	return executeOpen[*ForwardRuleSetting](ctx, a.core, a.akm, "GET", "/api/open/forwardRule/setting", nil)
}

// SaveSetting 设置总开关：0 关闭，1 开启且未命中仍推送，2 开启且未命中不推送。
func (a *ForwardRuleAPI) SaveSetting(ctx context.Context, mode int) error {
	path := appendQuery("/api/open/forwardRule/setting", []queryParam{{"mode", strconv.Itoa(mode)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", path, nil)
	return err
}
