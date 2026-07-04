package pushplus

import (
	"context"
	"strconv"
)

// TopicAPI 开放接口 - 群组（文档「五. 群组接口」）。
type TopicAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newTopicAPI(c *core, akm *AccessKeyManager) *TopicAPI {
	return &TopicAPI{core: c, akm: akm}
}

// List 群组列表。
func (a *TopicAPI) List(ctx context.Context, query *TopicListQuery) (*PageResult[TopicItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[TopicItem]](ctx, a.core, a.akm, "POST", "/api/open/topic/list", body)
}

// Detail 群组详情（自己创建的群组）。
func (a *TopicAPI) Detail(ctx context.Context, topicID int64) (*TopicDetail, error) {
	path := appendQuery("/api/open/topic/detail", []queryParam{{"topicId", strconv.FormatInt(topicID, 10)}})
	return executeOpen[*TopicDetail](ctx, a.core, a.akm, "GET", path, nil)
}

// JoinDetail 加入的群组详情。
func (a *TopicAPI) JoinDetail(ctx context.Context, topicID int64) (*TopicDetail, error) {
	path := appendQuery("/api/open/topic/joinTopicDetail", []queryParam{{"topicId", strconv.FormatInt(topicID, 10)}})
	return executeOpen[*TopicDetail](ctx, a.core, a.akm, "GET", path, nil)
}

// Add 新增群组，返回新群组 ID。
func (a *TopicAPI) Add(ctx context.Context, req *TopicAddRequest) (int64, error) {
	return executeOpen[int64](ctx, a.core, a.akm, "POST", "/api/open/topic/add", req)
}

// Edit 编辑群组。
func (a *TopicAPI) Edit(ctx context.Context, req *TopicEditRequest) (string, error) {
	return executeOpen[string](ctx, a.core, a.akm, "POST", "/api/open/topic/editTopic", req)
}

// QrCode 获取群组二维码。second/scanCount 传 nil 表示使用服务端默认值；
// scanCount 为 -1 表示不限扫码次数。
func (a *TopicAPI) QrCode(ctx context.Context, topicID int64, second, scanCount *int) (*TopicQrCode, error) {
	params := []queryParam{{"topicId", strconv.FormatInt(topicID, 10)}}
	if second != nil {
		params = append(params, queryParam{"second", strconv.Itoa(*second)})
	}
	if scanCount != nil {
		params = append(params, queryParam{"scanCount", strconv.Itoa(*scanCount)})
	}
	path := appendQuery("/api/open/topic/qrCode", params)
	return executeOpen[*TopicQrCode](ctx, a.core, a.akm, "GET", path, nil)
}

// Exit 退出群组。
func (a *TopicAPI) Exit(ctx context.Context, topicID int64) (string, error) {
	path := appendQuery("/api/open/topic/exitTopic", []queryParam{{"topicId", strconv.FormatInt(topicID, 10)}})
	return executeOpen[string](ctx, a.core, a.akm, "GET", path, nil)
}

// Delete 删除群组。
func (a *TopicAPI) Delete(ctx context.Context, topicID int64) (string, error) {
	path := appendQuery("/api/open/topic/delete", []queryParam{{"topicId", strconv.FormatInt(topicID, 10)}})
	return executeOpen[string](ctx, a.core, a.akm, "GET", path, nil)
}

// SetOpen 上架/下架群组；isOpen 1 上架、0 下架。
func (a *TopicAPI) SetOpen(ctx context.Context, topicID int64, isOpen int) (string, error) {
	body := map[string]any{"topic": topicID, "isOpen": isOpen}
	return executeOpen[string](ctx, a.core, a.akm, "POST", "/api/open/topic/isOpen", body)
}
