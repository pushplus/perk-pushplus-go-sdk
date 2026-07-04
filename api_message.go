package pushplus

import (
	"context"
	"strings"
)

// MessageAPI 发送消息接口，对应文档「二. 发送消息接口」与「三. 多渠道发送消息接口」。
//
// 内置本地限流守卫：当上游返回 code=900（请求次数过多）时，后续对同一 token 的发送调用
// 会在 SDK 内被直接短路返回错误，不再发起 HTTP，直到守卫到期自动解除。
// 可通过 Config.RateLimitGuardEnabled 关闭。
type MessageAPI struct {
	core  *core
	guard *RateLimitGuard
}

func newMessageAPI(c *core, guard *RateLimitGuard) *MessageAPI {
	return &MessageAPI{core: c, guard: guard}
}

// RateLimitGuard 暴露限流守卫，便于运维场景手动 Clear 或观察解禁时间。
func (a *MessageAPI) RateLimitGuard() *RateLimitGuard {
	return a.guard
}

// Send 发送单条消息；request.Token 可省略，SDK 会自动注入配置中的 token。
// 返回消息流水号 shortCode。
func (a *MessageAPI) Send(ctx context.Context, request *SendRequest) (string, error) {
	if request == nil {
		return "", newError(-1, "SendRequest 不能为空")
	}
	req := *request
	if strings.TrimSpace(req.Token) == "" {
		token, err := a.requireToken()
		if err != nil {
			return "", err
		}
		req.Token = token
	}
	if strings.TrimSpace(req.Content) == "" {
		return "", newError(-1, "发送消息 content 不能为空")
	}
	if err := a.guard.Check(req.Token); err != nil {
		return "", err
	}
	shortCode, err := executeForData[string](ctx, a.core, "POST", "/send", nil, &req)
	if err != nil {
		a.markIfRateLimited(req.Token, err)
		return "", err
	}
	return shortCode, nil
}

// SendSimple 便捷方法：以默认渠道、默认模板发送一条简单消息，返回消息流水号。
func (a *MessageAPI) SendSimple(ctx context.Context, title, content string) (string, error) {
	return a.Send(ctx, &SendRequest{Title: title, Content: content})
}

// BatchSend 多渠道发送消息，每个渠道返回一个 BatchSendResult。
func (a *MessageAPI) BatchSend(ctx context.Context, request *BatchSendRequest) ([]BatchSendResult, error) {
	if request == nil {
		return nil, newError(-1, "BatchSendRequest 不能为空")
	}
	req := *request
	if strings.TrimSpace(req.Token) == "" {
		token, err := a.requireToken()
		if err != nil {
			return nil, err
		}
		req.Token = token
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, newError(-1, "批量发送消息 content 不能为空")
	}
	if err := a.guard.Check(req.Token); err != nil {
		return nil, err
	}
	results, err := executeForData[[]BatchSendResult](ctx, a.core, "POST", "/batchSend", nil, &req)
	if err != nil {
		a.markIfRateLimited(req.Token, err)
		return nil, err
	}
	return results, nil
}

func (a *MessageAPI) requireToken() (string, error) {
	t := a.core.config.Token
	if strings.TrimSpace(t) == "" {
		return "", newError(-1, "发送消息需要 token，但 Config.Token 为空")
	}
	return t, nil
}

// markIfRateLimited 拦截 code=900 的业务错误并登记到本地限流守卫。
func (a *MessageAPI) markIfRateLimited(token string, err error) {
	if e, ok := AsError(err); ok && e.IsRateLimited() {
		a.guard.MarkBlocked(token)
	}
}
