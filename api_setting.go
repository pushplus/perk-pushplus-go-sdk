package pushplus

import (
	"context"
	"strconv"
)

// SettingAPI 开放接口 - 功能设置（文档「九. 功能设置接口」）。
type SettingAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newSettingAPI(c *core, akm *AccessKeyManager) *SettingAPI {
	return &SettingAPI{core: c, akm: akm}
}

// ListUserDefault 用户默认设置列表。
func (a *SettingAPI) ListUserDefault(ctx context.Context, query *PageQuery) (*PageResult[UserDefaultItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[UserDefaultItem]](ctx, a.core, a.akm, "POST", "/api/open/setting/listUserDefault", body)
}

// DetailUserDefault 用户默认设置详情。
func (a *SettingAPI) DetailUserDefault(ctx context.Context, id int64) (*UserDefaultDetail, error) {
	path := appendQuery("/api/open/setting/detailUserDefault", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	return executeOpen[*UserDefaultDetail](ctx, a.core, a.akm, "GET", path, nil)
}

// AddUserDefault 新增用户默认设置。
func (a *SettingAPI) AddUserDefault(ctx context.Context, req *UserDefaultSaveRequest) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/api/open/setting/addUserDefault", req)
	return err
}

// EditUserDefault 编辑用户默认设置。
func (a *SettingAPI) EditUserDefault(ctx context.Context, req *UserDefaultSaveRequest) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/api/open/setting/editUserDefault", req)
	return err
}

// DeleteUserDefault 删除用户默认设置。
func (a *SettingAPI) DeleteUserDefault(ctx context.Context, id int64) error {
	path := appendQuery("/api/open/setting/deleteUserDefault", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "DELETE", path, nil)
	return err
}

// ChangeReceiveLimit 修改接收限制；0 接收全部消息，1 不接收消息。
// 注意：官方接口的参数名即为 recevieLimit（原始拼写）。
func (a *SettingAPI) ChangeReceiveLimit(ctx context.Context, receiveLimit int) error {
	path := appendQuery("/api/open/setting/changeRecevieLimit",
		[]queryParam{{"recevieLimit", strconv.Itoa(receiveLimit)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", path, nil)
	return err
}

// ChangeIsSend 启用/禁用发送功能；0 禁用，1 启用。
func (a *SettingAPI) ChangeIsSend(ctx context.Context, isSend int) error {
	path := appendQuery("/api/open/setting/changeIsSend", []queryParam{{"isSend", strconv.Itoa(isSend)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", path, nil)
	return err
}

// ChangeOpenMessageType 修改消息打开方式；0 H5，1 小程序。
func (a *SettingAPI) ChangeOpenMessageType(ctx context.Context, openMessageType int) error {
	path := appendQuery("/api/open/setting/changeOpenMessageType",
		[]queryParam{{"openMessageType", strconv.Itoa(openMessageType)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", path, nil)
	return err
}

// ChangeExtensionForward 修改插件转发设置；0 否，1 是。
func (a *SettingAPI) ChangeExtensionForward(ctx context.Context, forward int) error {
	path := appendQuery("/api/open/setting/extension", []queryParam{{"forward", strconv.Itoa(forward)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", path, nil)
	return err
}
