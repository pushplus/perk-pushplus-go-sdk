package pushplus

import (
	"context"
	"strconv"
)

// TopicUserAPI 开放接口 - 群组用户（文档「六. 群组用户接口」）。
type TopicUserAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newTopicUserAPI(c *core, akm *AccessKeyManager) *TopicUserAPI {
	return &TopicUserAPI{core: c, akm: akm}
}

// SubscriberList 群组订阅用户列表；query 不可为 nil（需携带 topicId）。
func (a *TopicUserAPI) SubscriberList(ctx context.Context, query *TopicUserListQuery) (*PageResult[TopicUserItem], error) {
	if query == nil {
		return nil, newError(-1, "TopicUserListQuery 不能为空")
	}
	return executeOpen[*PageResult[TopicUserItem]](ctx, a.core, a.akm, "POST", "/api/open/topicUser/subscriberList", query)
}

// DeleteUser 删除群组用户。
func (a *TopicUserAPI) DeleteUser(ctx context.Context, topicRelationID int64) (string, error) {
	path := appendQuery("/api/open/topicUser/deleteTopicUser",
		[]queryParam{{"topicRelationId", strconv.FormatInt(topicRelationID, 10)}})
	return executeOpen[string](ctx, a.core, a.akm, "POST", path, nil)
}

// EditRemark 编辑群组用户备注。
func (a *TopicUserAPI) EditRemark(ctx context.Context, id int64, remark string) error {
	body := map[string]any{"id": id, "remark": remark}
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/api/open/topicUser/editRemark", body)
	return err
}
