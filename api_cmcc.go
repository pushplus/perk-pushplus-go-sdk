package pushplus

import "context"

// CmccAPI 开放接口 - 新消息 ClawBot（文档「九. 新消息ClawBot接口」）。
// 需先在手机 5G 消息的「新消息ClawBot」应用号中获取 Channel API Key，再调用绑定接口。
// 发送消息时 channel 传 cmcc。仅支持中国移动用户。
type CmccAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newCmccAPI(c *core, akm *AccessKeyManager) *CmccAPI {
	return &CmccAPI{core: c, akm: akm}
}

// Bind 绑定新消息 ClawBot。apiKey 必须以 ak_ 或 app_ 开头。
func (a *CmccAPI) Bind(ctx context.Context, apiKey string) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/api/open/cmcc/bind", &CmccBindRequest{APIKey: apiKey})
	return err
}

// Info 查询绑定状态。
func (a *CmccAPI) Info(ctx context.Context) (*CmccInfo, error) {
	return executeOpen[*CmccInfo](ctx, a.core, a.akm, "GET", "/api/open/cmcc/info", nil)
}

// Unbind 解绑新消息 ClawBot。
func (a *CmccAPI) Unbind(ctx context.Context) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", "/api/open/cmcc/unbind", nil)
	return err
}

// SendTest 发送测试消息。未绑定会返回「未绑定新消息ClawBot」。
func (a *CmccAPI) SendTest(ctx context.Context) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", "/api/open/cmcc/test", nil)
	return err
}
