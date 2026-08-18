package pushplus

import (
	"context"
	"strconv"
)

// FriendAPI 开放接口 - 好友功能（文档「十. 好友功能接口」）。
type FriendAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newFriendAPI(c *core, akm *AccessKeyManager) *FriendAPI {
	return &FriendAPI{core: c, akm: akm}
}

// FriendQrCodeQuery 获取加好友二维码的可选参数；nil 字段不进入 query。
type FriendQrCodeQuery struct {
	AppID     string
	Content   string
	Second    *int
	ScanCount *int
}

// GetQrCode 获取加好友二维码；query 可为 nil。
func (a *FriendAPI) GetQrCode(ctx context.Context, query *FriendQrCodeQuery) (*FriendQrCode, error) {
	var params []queryParam
	if query != nil {
		if query.AppID != "" {
			params = append(params, queryParam{"appId", query.AppID})
		}
		if query.Content != "" {
			params = append(params, queryParam{"content", query.Content})
		}
		if query.Second != nil {
			params = append(params, queryParam{"second", strconv.Itoa(*query.Second)})
		}
		if query.ScanCount != nil {
			params = append(params, queryParam{"scanCount", strconv.Itoa(*query.ScanCount)})
		}
	}
	path := appendQuery("/api/open/friend/getQrCode", params)
	return executeOpen[*FriendQrCode](ctx, a.core, a.akm, "GET", path, nil)
}

// List 好友列表。
func (a *FriendAPI) List(ctx context.Context, query *PageQuery) (*PageResult[FriendItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[FriendItem]](ctx, a.core, a.akm, "POST", "/api/open/friend/list", body)
}

// Delete 删除好友。
func (a *FriendAPI) Delete(ctx context.Context, friendID int64) error {
	path := appendQuery("/api/open/friend/deleteFriend", []queryParam{{"friendId", strconv.FormatInt(friendID, 10)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", path, nil)
	return err
}

// EditRemark 编辑好友备注。
func (a *FriendAPI) EditRemark(ctx context.Context, id int64, remark string) error {
	body := map[string]any{"id": id, "remark": remark}
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/api/open/friend/editRemark", body)
	return err
}

// AddBlacklist 将好友加入黑名单。加入后将解除双方好友关系，对方无法再添加你。
func (a *FriendAPI) AddBlacklist(ctx context.Context, friendID int64) error {
	path := appendQuery("/api/open/friend/addBlacklist", []queryParam{{"friendId", strconv.FormatInt(friendID, 10)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", path, nil)
	return err
}

// BlacklistList 好友黑名单列表。
func (a *FriendAPI) BlacklistList(ctx context.Context, query *PageQuery) (*PageResult[FriendBlacklistItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[FriendBlacklistItem]](ctx, a.core, a.akm, "POST", "/api/open/friend/blacklistList", body)
}

// RemoveBlacklist 解除好友黑名单。解除后不会自动恢复好友关系，需重新扫码添加。
func (a *FriendAPI) RemoveBlacklist(ctx context.Context, id int64) error {
	path := appendQuery("/api/open/friend/removeBlacklist", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", path, nil)
	return err
}
