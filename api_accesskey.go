package pushplus

import (
	"context"
	"strings"
)

// AccessKeyAPI 对应文档「一. 获取 AccessKey」。
// 一般无需手动调用：开放接口会通过 AccessKeyManager 自动获取、缓存与刷新。
type AccessKeyAPI struct {
	core *core
}

func newAccessKeyAPI(c *core) *AccessKeyAPI {
	return &AccessKeyAPI{core: c}
}

// GetAccessKey 使用配置中的 token + secretKey 获取 AccessKey。
func (a *AccessKeyAPI) GetAccessKey(ctx context.Context) (*AccessKeyResult, error) {
	return a.GetAccessKeyWith(ctx, a.core.config.Token, a.core.config.SecretKey)
}

// GetAccessKeyWith 使用指定 token + secretKey 获取 AccessKey。
// token 为用户 token（不支持消息 token）。
func (a *AccessKeyAPI) GetAccessKeyWith(ctx context.Context, token, secretKey string) (*AccessKeyResult, error) {
	if strings.TrimSpace(token) == "" {
		return nil, newError(-1, "获取 AccessKey 需要 token")
	}
	if strings.TrimSpace(secretKey) == "" {
		return nil, newError(-1, "获取 AccessKey 需要 secretKey")
	}
	body := map[string]string{"token": token, "secretKey": secretKey}
	result, err := executeForData[*AccessKeyResult](ctx, a.core, "POST", "/api/common/openApi/getAccessKey", nil, body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
