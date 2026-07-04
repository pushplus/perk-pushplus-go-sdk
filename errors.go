package pushplus

import "fmt"

// Error 是 SDK 的统一错误类型：
//   - HTTP 调用失败：Code 为 HTTP 状态码
//   - 业务接口返回 code != 200：Code 为业务码、Msg 为业务消息
//   - SDK 参数校验失败：Code = -1
//   - 本地限流守卫命中（不发起 HTTP）：Code = 900、Msg 包含 "本地限流守卫" 字样
type Error struct {
	Code int
	Msg  string
	// Cause 底层错误（如网络错误、JSON 解析错误），可能为 nil。
	Cause error
}

// Error 实现 error 接口。
func (e *Error) Error() string {
	return fmt.Sprintf("pushplus: code=%d, msg=%s", e.Code, e.Msg)
}

// Unwrap 支持 errors.Is / errors.As 链式匹配底层错误。
func (e *Error) Unwrap() error {
	return e.Cause
}

// ErrorCode 把 Code 映射到 ErrorCode 枚举；未知返回 ErrorCodeUnknown。
func (e *Error) ErrorCode() ErrorCode {
	return ErrorCodeOf(e.Code)
}

// IsRateLimited 是否为限流错误（code=900）。
func (e *Error) IsRateLimited() bool {
	return e.Code == int(ErrorCodeRateLimited)
}

func newError(code int, msg string) *Error {
	return &Error{Code: code, Msg: msg}
}

func newErrorWithCause(code int, msg string, cause error) *Error {
	return &Error{Code: code, Msg: msg, Cause: cause}
}

// AsError 尝试把任意 error 断言为 *pushplus.Error。
func AsError(err error) (*Error, bool) {
	e, ok := err.(*Error)
	return e, ok
}
