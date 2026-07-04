package pushplus

import (
	"context"
	"strings"
	"sync"
)

// recordedRequest 记录一次 HTTP 调用。
type recordedRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Raw     bool
}

// mockHTTPRequester 按 URL 包含的 path 片段匹配响应队列；记录每次请求。
// 无匹配时返回 {"code":200,"msg":"ok"}。
type mockHTTPRequester struct {
	mu        sync.Mutex
	responses map[string][]*HTTPResponse
	errors    map[string]error
	Requests  []recordedRequest
}

func newMockHTTPRequester() *mockHTTPRequester {
	return &mockHTTPRequester{
		responses: make(map[string][]*HTTPResponse),
		errors:    make(map[string]error),
	}
}

// enqueue 为包含 pathFragment 的 URL 追加一个响应。
func (m *mockHTTPRequester) enqueue(pathFragment string, statusCode int, body string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses[pathFragment] = append(m.responses[pathFragment], &HTTPResponse{StatusCode: statusCode, Body: body})
}

func (m *mockHTTPRequester) Execute(_ context.Context, method, url string, headers map[string]string, body string) (*HTTPResponse, error) {
	return m.dispatch(method, url, headers, body, false)
}

func (m *mockHTTPRequester) ExecuteRaw(_ context.Context, method, url string, headers map[string]string, body []byte) (*HTTPResponse, error) {
	return m.dispatch(method, url, headers, string(body), true)
}

func (m *mockHTTPRequester) dispatch(method, url string, headers map[string]string, body string, raw bool) (*HTTPResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Requests = append(m.Requests, recordedRequest{Method: method, URL: url, Headers: headers, Body: body, Raw: raw})

	for fragment, err := range m.errors {
		if strings.Contains(url, fragment) {
			return nil, err
		}
	}
	for fragment, queue := range m.responses {
		if strings.Contains(url, fragment) && len(queue) > 0 {
			resp := queue[0]
			m.responses[fragment] = queue[1:]
			return resp, nil
		}
	}
	return &HTTPResponse{StatusCode: 200, Body: `{"code":200,"msg":"ok"}`}, nil
}

// requestsTo 返回 URL 包含指定片段的全部请求。
func (m *mockHTTPRequester) requestsTo(pathFragment string) []recordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []recordedRequest
	for _, r := range m.Requests {
		if strings.Contains(r.URL, pathFragment) {
			out = append(out, r)
		}
	}
	return out
}
