package pushplus

import (
	"context"
)

// DocAPI 开放接口 - push 文档。
//
// 文档：https://www.pushplus.plus/doc/ecosystem/doc/
// 基础路径：/push/api/open/doc
type DocAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newDocAPI(c *core, akm *AccessKeyManager) *DocAPI {
	return &DocAPI{core: c, akm: akm}
}

// List 我的文档分页。
func (a *DocAPI) List(ctx context.Context, query *DocListQuery) (*PageResult[DocListItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[DocListItem]](ctx, a.core, a.akm, "POST", "/push/api/open/doc/list", body)
}

// Create 创建空白文档。
func (a *DocAPI) Create(ctx context.Context, title string) (*DocVo, error) {
	return executeOpen[*DocVo](ctx, a.core, a.akm, "POST", "/push/api/open/doc/create", map[string]any{"title": title})
}

// Content 获取文档元信息与 HTML 草稿正文。
func (a *DocAPI) Content(ctx context.Context, docCode string) (*DocContent, error) {
	path := appendQuery("/push/api/open/doc/content", []queryParam{{"docCode", docCode}})
	return executeOpen[*DocContent](ctx, a.core, a.akm, "GET", path, nil)
}

// SaveContent 保存 HTML 草稿（不影响分享页，需再 Publish）。
func (a *DocAPI) SaveContent(ctx context.Context, docCode, content string) (*DocVo, error) {
	body := map[string]any{"docCode": docCode, "content": content}
	return executeOpen[*DocVo](ctx, a.core, a.akm, "POST", "/push/api/open/doc/saveContent", body)
}

// Publish 将草稿同步为分享页快照。
func (a *DocAPI) Publish(ctx context.Context, docCode string) (*DocVo, error) {
	path := appendQuery("/push/api/open/doc/publish", []queryParam{{"docCode", docCode}})
	return executeOpen[*DocVo](ctx, a.core, a.akm, "POST", path, nil)
}

// Rename 重命名。
func (a *DocAPI) Rename(ctx context.Context, docCode, title string) error {
	body := map[string]any{"docCode": docCode, "title": title}
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/push/api/open/doc/rename", body)
	return err
}

// Delete 删除文档。
func (a *DocAPI) Delete(ctx context.Context, docCode string) error {
	path := appendQuery("/push/api/open/doc/delete", []queryParam{{"docCode", docCode}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", path, nil)
	return err
}

// UpdateShare 更新分享设置。sharePerm：0 关闭 / 1 开启；shareLogin 为 nil 时沿用原值。
func (a *DocAPI) UpdateShare(ctx context.Context, docCode string, sharePerm int, shareLogin *int) (*DocVo, error) {
	body := map[string]any{"docCode": docCode, "sharePerm": sharePerm}
	if shareLogin != nil {
		body["shareLogin"] = *shareLogin
	}
	return executeOpen[*DocVo](ctx, a.core, a.akm, "POST", "/push/api/open/doc/updateShare", body)
}
