package pushplus

import "context"

// ClawBotAPI 开放接口 - 微信 ClawBot（文档「八. 微信 ClawBot 接口」）。
type ClawBotAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newClawBotAPI(c *core, akm *AccessKeyManager) *ClawBotAPI {
	return &ClawBotAPI{core: c, akm: akm}
}

// GetBotQrcode 获取 ClawBot 登录二维码。
func (a *ClawBotAPI) GetBotQrcode(ctx context.Context) (*ClawBotQrCode, error) {
	return executeOpen[*ClawBotQrCode](ctx, a.core, a.akm, "GET", "/api/open/clawBot/getBotQrcode", nil)
}

// GetQrcodeStatus 查询二维码扫码状态（注意：官方接口的 query 参数名即为 getQrcodeStatus）。
func (a *ClawBotAPI) GetQrcodeStatus(ctx context.Context, qrcode string) error {
	path := appendQuery("/api/open/clawBot/getQrcodeStatus", []queryParam{{"getQrcodeStatus", qrcode}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", path, nil)
	return err
}

// BotInfo 获取 ClawBot 绑定信息。
func (a *ClawBotAPI) BotInfo(ctx context.Context) (*ClawBotInfo, error) {
	return executeOpen[*ClawBotInfo](ctx, a.core, a.akm, "GET", "/api/open/clawBot/botInfo", nil)
}

// Unbind 解绑 ClawBot。
func (a *ClawBotAPI) Unbind(ctx context.Context) error {
	_, err := executeOpen[any](ctx, a.core, a.akm, "GET", "/api/open/clawBot/unbind", nil)
	return err
}

// GetMsg 获取 ClawBot 消息。
func (a *ClawBotAPI) GetMsg(ctx context.Context) ([]ClawBotMessage, error) {
	return executeOpen[[]ClawBotMessage](ctx, a.core, a.akm, "GET", "/api/open/clawBot/getMsg", nil)
}
