package pushplus

import (
	"context"
	"strconv"
)

// FormAPI 开放接口 - push 表单。
//
// 文档：https://www.pushplus.plus/doc/ecosystem/form/
// 基础路径：/push/api/open/form
//
// 表单开放接口不单独提供推送接口。发布后请通过 MessageAPI 推送填写页：
// template=form，pushId=formCode。
type FormAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newFormAPI(c *core, akm *AccessKeyManager) *FormAPI {
	return &FormAPI{core: c, akm: akm}
}

// List 我的表单分页。
func (a *FormAPI) List(ctx context.Context, query *FormListQuery) (*PageResult[FormListItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[FormListItem]](ctx, a.core, a.akm, "POST", "/push/api/open/form/list", body)
}

// Create 创建空白表单（草稿）。
func (a *FormAPI) Create(ctx context.Context, title string) (*FormListItem, error) {
	return executeOpen[*FormListItem](ctx, a.core, a.akm, "POST", "/push/api/open/form/create", map[string]any{"title": title})
}

// Copy 基于已有表单复制一份新草稿。
func (a *FormAPI) Copy(ctx context.Context, id int64) (*FormListItem, error) {
	path := appendQuery("/push/api/open/form/copy", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	return executeOpen[*FormListItem](ctx, a.core, a.akm, "POST", path, nil)
}

// Save 保存表单设计（仅更新草稿；已发布需再调用 Publish）。
func (a *FormAPI) Save(ctx context.Context, req *FormSaveRequest) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/push/api/open/form/save", req)
	return err
}

// Detail 表单详情（含草稿题目、主题、设置）。
func (a *FormAPI) Detail(ctx context.Context, id int64) (*FormDetail, error) {
	path := appendQuery("/push/api/open/form/detail", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	return executeOpen[*FormDetail](ctx, a.core, a.akm, "GET", path, nil)
}

// PublishDiff 草稿与发布快照的题目差异。
func (a *FormAPI) PublishDiff(ctx context.Context, id int64) (*FormPublishDiff, error) {
	path := appendQuery("/push/api/open/form/publishDiff", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	return executeOpen[*FormPublishDiff](ctx, a.core, a.akm, "GET", path, nil)
}

// Publish 发布表单，开始收集。
func (a *FormAPI) Publish(ctx context.Context, id int64) (*FormPublishResult, error) {
	path := appendQuery("/push/api/open/form/publish", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	return executeOpen[*FormPublishResult](ctx, a.core, a.akm, "POST", path, nil)
}

// Stop 停止收集。
func (a *FormAPI) Stop(ctx context.Context, id int64) error {
	path := appendQuery("/push/api/open/form/stop", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", path, nil)
	return err
}

// Delete 删除表单（不可恢复）。
func (a *FormAPI) Delete(ctx context.Context, id int64) error {
	path := appendQuery("/push/api/open/form/delete", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", path, nil)
	return err
}
