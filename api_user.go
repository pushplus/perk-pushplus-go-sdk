package pushplus

import "context"

// UserAPI 开放接口 - 用户（文档「三. 用户接口」）。
type UserAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newUserAPI(c *core, akm *AccessKeyManager) *UserAPI {
	return &UserAPI{core: c, akm: akm}
}

// GetToken 获取当前用户 token。
func (a *UserAPI) GetToken(ctx context.Context) (string, error) {
	return executeOpen[string](ctx, a.core, a.akm, "GET", "/api/open/user/token", nil)
}

// MyInfo 获取当前用户信息。
func (a *UserAPI) MyInfo(ctx context.Context) (*UserInfo, error) {
	return executeOpen[*UserInfo](ctx, a.core, a.akm, "GET", "/api/open/user/myInfo", nil)
}

// GetLimitTime 获取用户限制发送时间。
func (a *UserAPI) GetLimitTime(ctx context.Context) (*UserLimitTime, error) {
	return executeOpen[*UserLimitTime](ctx, a.core, a.akm, "GET", "/api/open/user/userLimitTime", nil)
}

// GetSendCount 获取用户当日发送次数统计。
func (a *UserAPI) GetSendCount(ctx context.Context) (*SendCount, error) {
	return executeOpen[*SendCount](ctx, a.core, a.akm, "GET", "/api/open/user/sendCount", nil)
}
