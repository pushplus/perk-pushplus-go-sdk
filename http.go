package pushplus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// HTTPResponse 是 HTTPRequester 返回的原始响应。
type HTTPResponse struct {
	StatusCode int
	Body       string
}

// IsSuccessful HTTP 状态码是否在 2xx 区间。
func (r *HTTPResponse) IsSuccessful() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// HTTPRequester 抽象 HTTP 执行层，便于替换实现或在测试中注入 mock。
type HTTPRequester interface {
	// Execute 执行 JSON 文本请求；body 为空串表示无请求体。
	Execute(ctx context.Context, method, url string, headers map[string]string, body string) (*HTTPResponse, error)
	// ExecuteRaw 执行二进制请求体（如 multipart 上传）。
	ExecuteRaw(ctx context.Context, method, url string, headers map[string]string, body []byte) (*HTTPResponse, error)
}

// defaultHTTPRequester 基于标准库 net/http 的默认实现。
type defaultHTTPRequester struct {
	client     *http.Client
	logRequest bool
}

// NewDefaultHTTPRequester 创建基于 net/http 的默认 HTTPRequester。
func NewDefaultHTTPRequester(config *Config) HTTPRequester {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: config.ConnectTimeout,
		}).DialContext,
	}
	return &defaultHTTPRequester{
		client: &http.Client{
			Transport: transport,
			Timeout:   config.ReadTimeout,
		},
		logRequest: config.LogRequest,
	}
}

func (d *defaultHTTPRequester) Execute(ctx context.Context, method, reqURL string, headers map[string]string, body string) (*HTTPResponse, error) {
	var payload []byte
	if body != "" {
		payload = []byte(body)
	}
	return d.do(ctx, method, reqURL, headers, payload, "application/json;charset=UTF-8")
}

func (d *defaultHTTPRequester) ExecuteRaw(ctx context.Context, method, reqURL string, headers map[string]string, body []byte) (*HTTPResponse, error) {
	return d.do(ctx, method, reqURL, headers, body, "application/octet-stream")
}

func (d *defaultHTTPRequester) do(ctx context.Context, method, reqURL string, headers map[string]string, body []byte, defaultContentType string) (*HTTPResponse, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, reqURL, reader)
	if err != nil {
		return nil, newErrorWithCause(-1, "构造 HTTP 请求失败: "+err.Error(), err)
	}

	hasContentType := false
	for k, v := range headers {
		req.Header.Set(k, v)
		if strings.EqualFold(k, "Content-Type") {
			hasContentType = true
		}
	}
	if body != nil && !hasContentType {
		req.Header.Set("Content-Type", defaultContentType)
	}

	if d.logRequest {
		log.Printf("[pushplus] --> %s %s body=%s", method, reqURL, truncateForLog(body))
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, newErrorWithCause(-1, "PushPlus 接口 HTTP 调用失败: "+err.Error(), err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, newErrorWithCause(-1, "读取 PushPlus 响应失败: "+err.Error(), err)
	}

	if d.logRequest {
		log.Printf("[pushplus] <-- %d %s body=%s", resp.StatusCode, reqURL, truncateForLog(respBody))
	}

	return &HTTPResponse{StatusCode: resp.StatusCode, Body: string(respBody)}, nil
}

func truncateForLog(b []byte) string {
	const max = 2048
	if len(b) > max {
		return string(b[:max]) + "...(truncated)"
	}
	return string(b)
}

/* ============================== 请求执行核心 ============================== */

// core 持有配置与 HTTP 实现，供各 API 复用。
type core struct {
	config *Config
	http   HTTPRequester
}

// resolveURL 拼接绝对 URL。
func (c *core) resolveURL(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return c.config.ResolveBaseURL() + path
}

// queryParam 表示一个有序的 query 参数。
type queryParam struct {
	key   string
	value string
}

// appendQuery 在 path 上按顺序追加 query 参数（值已转为字符串，nil 参数由调用方跳过）。
func appendQuery(path string, params []queryParam) string {
	if len(params) == 0 {
		return path
	}
	var sb strings.Builder
	for _, p := range params {
		if sb.Len() > 0 {
			sb.WriteByte('&')
		}
		sb.WriteString(url.QueryEscape(p.key))
		sb.WriteByte('=')
		sb.WriteString(url.QueryEscape(p.value))
	}
	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}
	return path + sep + sb.String()
}

// apiResponse 是 PushPlus 统一响应结构 { code, msg, data }。
type apiResponse[T any] struct {
	Code *int   `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func (r *apiResponse[T]) isSuccess() bool {
	return r.Code != nil && *r.Code == 200
}

func (r *apiResponse[T]) codeOr(def int) int {
	if r.Code == nil {
		return def
	}
	return *r.Code
}

// execute 执行请求并按两阶段方式解析统一响应（不校验业务 code）：
//  1. 先解析 code/msg 与 data 的原始 JSON，避免业务失败时 data 为字符串导致类型不匹配；
//  2. 仅在 code==200 时把 data 反序列化为目标类型；
//  3. code!=200 时把 data 的文本追加到 msg 上，便于上层组装错误信息。
func execute[T any](ctx context.Context, c *core, method, path string, headers map[string]string, body any) (*apiResponse[T], error) {
	reqURL := c.resolveURL(path)
	var bodyJSON string
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, newErrorWithCause(-1, "序列化请求体失败: "+err.Error(), err)
		}
		bodyJSON = string(b)
	}

	resp, err := c.http.Execute(ctx, method, reqURL, headers, bodyJSON)
	if err != nil {
		if _, ok := AsError(err); ok {
			return nil, err
		}
		return nil, newErrorWithCause(-1, "PushPlus 接口 HTTP 调用失败: "+err.Error(), err)
	}
	if !resp.IsSuccessful() {
		return nil, newError(resp.StatusCode,
			fmt.Sprintf("PushPlus 接口 HTTP 调用失败: status=%d, body=%s", resp.StatusCode, resp.Body))
	}
	return parseAPIResponse[T](resp.Body)
}

func parseAPIResponse[T any](responseBody string) (*apiResponse[T], error) {
	var raw apiResponse[json.RawMessage]
	if err := json.Unmarshal([]byte(responseBody), &raw); err != nil {
		return nil, newErrorWithCause(-1, "解析 PushPlus 响应失败: "+err.Error()+", payload="+responseBody, err)
	}

	result := &apiResponse[T]{Code: raw.Code, Msg: raw.Msg}
	if raw.isSuccess() {
		if len(raw.Data) == 0 || string(raw.Data) == "null" {
			return result, nil
		}
		if err := json.Unmarshal(raw.Data, &result.Data); err != nil {
			return nil, newErrorWithCause(-1,
				"解析 PushPlus 响应 data 字段失败: "+err.Error()+", payload="+responseBody, err)
		}
		return result, nil
	}

	// 业务失败：把 data 中的字符串/简单值附加到 msg，方便使用者直接拿到错误描述。
	dataText := extractDataText(raw.Data)
	if dataText != "" {
		if result.Msg == "" {
			result.Msg = dataText
		} else if !strings.Contains(result.Msg, dataText) {
			result.Msg = result.Msg + ": " + dataText
		}
	}
	return result, nil
}

func extractDataText(data json.RawMessage) string {
	if len(data) == 0 {
		return ""
	}
	s := strings.TrimSpace(string(data))
	if s == "null" {
		return ""
	}
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		return str
	}
	if s[0] == '{' || s[0] == '[' {
		return s
	}
	return s
}

// executeForData 执行请求并直接返回 data；业务 code != 200 时返回 *Error。
func executeForData[T any](ctx context.Context, c *core, method, path string, headers map[string]string, body any) (T, error) {
	var zero T
	resp, err := execute[T](ctx, c, method, path, headers, body)
	if err != nil {
		return zero, err
	}
	if !resp.isSuccess() {
		return zero, newError(resp.codeOr(-1),
			fmt.Sprintf("PushPlus 接口业务失败: code=%d, msg=%s", resp.codeOr(-1), resp.Msg))
	}
	return resp.Data, nil
}

/* ============================== 开放接口执行 ============================== */

// headerAccessKey 开放接口的鉴权 header 名。
const headerAccessKey = "access-key"

// codeAccessKeyInvalid AccessKey 失效相关业务码（触发自动刷新重试）。
const codeAccessKeyInvalid = 401

// executeOpen 执行开放接口请求：自动携带 access-key；当业务 code=401 时
// invalidate 缓存的 AccessKey 并刷新后重试一次。
func executeOpen[T any](ctx context.Context, c *core, akm *AccessKeyManager, method, path string, body any) (T, error) {
	var zero T

	key, err := akm.GetAccessKey(ctx)
	if err != nil {
		return zero, err
	}
	resp, err := execute[T](ctx, c, method, path, map[string]string{headerAccessKey: key}, body)
	if err != nil {
		return zero, err
	}
	if resp.isSuccess() {
		return resp.Data, nil
	}

	if resp.Code != nil && *resp.Code == codeAccessKeyInvalid {
		akm.Invalidate()
		retryKey, err := akm.GetAccessKey(ctx)
		if err != nil {
			return zero, err
		}
		retry, err := execute[T](ctx, c, method, path, map[string]string{headerAccessKey: retryKey}, body)
		if err != nil {
			return zero, err
		}
		if retry.isSuccess() {
			return retry.Data, nil
		}
		return zero, newError(retry.codeOr(-1),
			fmt.Sprintf("PushPlus 开放接口业务失败(重试后): code=%d, msg=%s", retry.codeOr(-1), retry.Msg))
	}

	return zero, newError(resp.codeOr(-1),
		fmt.Sprintf("PushPlus 开放接口业务失败: code=%d, msg=%s", resp.codeOr(-1), resp.Msg))
}
