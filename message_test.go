package pushplus

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func newTestClient(mock *mockHTTPRequester, opts ...Option) *Client {
	base := []Option{
		WithToken("test-token"),
		WithSecretKey("test-secret"),
		WithHTTPRequester(mock),
	}
	return NewClient(append(base, opts...)...)
}

func TestSendInjectsDefaultToken(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("/send", 200, `{"code":200,"msg":"ok","data":"SC123"}`)
	client := newTestClient(mock)

	shortCode, err := client.SendSimple(context.Background(), "标题", "内容")
	if err != nil {
		t.Fatalf("SendSimple 失败: %v", err)
	}
	if shortCode != "SC123" {
		t.Fatalf("shortCode 期望 SC123，实际 %s", shortCode)
	}

	reqs := mock.requestsTo("/send")
	if len(reqs) != 1 {
		t.Fatalf("期望发起 1 次请求，实际 %d", len(reqs))
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatalf("请求体不是合法 JSON: %v", err)
	}
	if body["token"] != "test-token" {
		t.Fatalf("token 未自动注入，实际 %v", body["token"])
	}
}

func TestSendRequiresContent(t *testing.T) {
	mock := newMockHTTPRequester()
	client := newTestClient(mock)

	_, err := client.Send(context.Background(), &SendRequest{Title: "只有标题"})
	if err == nil {
		t.Fatal("期望 content 为空时报错")
	}
	e, ok := AsError(err)
	if !ok || e.Code != -1 {
		t.Fatalf("期望 code=-1 的参数校验错误，实际 %v", err)
	}
	if len(mock.Requests) != 0 {
		t.Fatal("参数校验失败时不应发起 HTTP 请求")
	}
}

func TestSendBusinessFailure903(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("/send", 200, `{"code":903,"msg":"token无效"}`)
	client := newTestClient(mock)

	_, err := client.SendSimple(context.Background(), "t", "c")
	if err == nil {
		t.Fatal("期望业务失败返回错误")
	}
	e, _ := AsError(err)
	if e.Code != 903 {
		t.Fatalf("期望 code=903，实际 %d", e.Code)
	}
	if e.ErrorCode() != ErrorCodeInvalidToken {
		t.Fatalf("期望 ErrorCodeInvalidToken，实际 %v", e.ErrorCode())
	}
	if !strings.Contains(e.Msg, "token无效") {
		t.Fatalf("错误消息应包含业务 msg，实际 %s", e.Msg)
	}
}

func TestBatchSendChannelCSV(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("/batchSend", 200,
		`{"code":200,"msg":"ok","data":[{"shortCode":"A","code":200,"channel":"wechat"},{"shortCode":"B","code":200,"channel":"webhook"}]}`)
	client := newTestClient(mock)

	req := (&BatchSendRequest{Title: "多渠道告警", Content: "CPU > 90%"}).
		AddChannel(ChannelWechat, "").
		AddChannel(ChannelWebhook, "bark").
		AddChannel(ChannelExtension, "")

	results, err := client.BatchSend(context.Background(), req)
	if err != nil {
		t.Fatalf("BatchSend 失败: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("期望 2 条结果，实际 %d", len(results))
	}
	if results[0].Channel != ChannelWechat || results[1].Channel != ChannelWebhook {
		t.Fatalf("channel 解析错误: %+v", results)
	}

	reqs := mock.requestsTo("/batchSend")
	var body map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatalf("请求体不是合法 JSON: %v", err)
	}
	if body["channel"] != "wechat,webhook,extension" {
		t.Fatalf("channel CSV 拼接错误: %v", body["channel"])
	}
	if body["option"] != ",bark," {
		t.Fatalf("option CSV 拼接错误: %v", body["option"])
	}
}

func TestSendHTTPFailure(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("/send", 500, "Internal Server Error")
	client := newTestClient(mock)

	_, err := client.SendSimple(context.Background(), "t", "c")
	if err == nil {
		t.Fatal("期望 HTTP 失败返回错误")
	}
	e, _ := AsError(err)
	if e.Code != 500 {
		t.Fatalf("期望 code=500（HTTP 状态码），实际 %d", e.Code)
	}
}
