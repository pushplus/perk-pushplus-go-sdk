package pushplus

import (
	"context"
	"strconv"
)

// ChannelAPI 开放接口 - 渠道配置：公众号/企业微信/邮箱（文档「七. 渠道配置」）。
type ChannelAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newChannelAPI(c *core, akm *AccessKeyManager) *ChannelAPI {
	return &ChannelAPI{core: c, akm: akm}
}

// MpList 微信公众号渠道列表。
func (a *ChannelAPI) MpList(ctx context.Context, query *PageQuery) (*PageResult[MpItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[MpItem]](ctx, a.core, a.akm, "POST", "/api/open/mp/list", body)
}

// CpList 企业微信渠道列表。
func (a *ChannelAPI) CpList(ctx context.Context, query *PageQuery) (*PageResult[CpItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[CpItem]](ctx, a.core, a.akm, "POST", "/api/open/cp/list", body)
}

// MailList 邮箱渠道列表。
func (a *ChannelAPI) MailList(ctx context.Context, query *PageQuery) (*PageResult[MailItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[MailItem]](ctx, a.core, a.akm, "POST", "/api/open/mail/list", body)
}

// MailDetail 邮箱渠道详情。
func (a *ChannelAPI) MailDetail(ctx context.Context, mailID int64) (*MailDetail, error) {
	path := appendQuery("/api/open/mail/detail", []queryParam{{"mailId", strconv.FormatInt(mailID, 10)}})
	return executeOpen[*MailDetail](ctx, a.core, a.akm, "GET", path, nil)
}
