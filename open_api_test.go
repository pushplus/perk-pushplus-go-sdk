package pushplus

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

const accessKeyResponse = `{"code":200,"msg":"ok","data":{"accessKey":"AK-1","expiresIn":7200}}`

func TestOpenMessageList(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/api/open/message/list", 200,
		`{"code":200,"msg":"ok","data":{"pageNum":1,"pageSize":20,"total":1,"pages":1,"list":[{"title":"t1","shortCode":"SC1","channel":"wechat","messageType":0}]}}`)
	client := newTestClient(mock)

	page, err := client.OpenMessage().List(context.Background(), NewPageQuery(1, 20))
	if err != nil {
		t.Fatalf("List 失败: %v", err)
	}
	if page.Total != 1 || len(page.List) != 1 || page.List[0].ShortCode != "SC1" {
		t.Fatalf("分页结果解析错误: %+v", page)
	}
	if page.List[0].Channel != ChannelWechat {
		t.Fatalf("channel 解析错误: %v", page.List[0].Channel)
	}
}

func TestOpenMessageQueryResult(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("sendMessageResult", 200,
		`{"code":200,"msg":"ok","data":{"status":2,"errorMessage":"","updateTime":"2026-07-04 10:00:00"}}`)
	client := newTestClient(mock)

	result, err := client.OpenMessage().QueryResult(context.Background(), "SC1")
	if err != nil {
		t.Fatalf("QueryResult 失败: %v", err)
	}
	if result.StatusEnum() != SendStatusSuccess {
		t.Fatalf("期望发送成功状态，实际 %v", result.StatusEnum())
	}

	reqs := mock.requestsTo("sendMessageResult")
	if !strings.Contains(reqs[0].URL, "shortCode=SC1") {
		t.Fatalf("URL 缺少 shortCode 参数: %s", reqs[0].URL)
	}
}

func TestTopicListParamsSerialization(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/api/open/topic/list", 200,
		`{"code":200,"msg":"ok","data":{"pageNum":1,"pageSize":20,"total":0,"pages":0,"list":[]}}`)
	client := newTestClient(mock)

	_, err := client.Topic().List(context.Background(), NewTopicListQuery(1, 20, 0))
	if err != nil {
		t.Fatalf("Topic List 失败: %v", err)
	}

	reqs := mock.requestsTo("/api/open/topic/list")
	var body map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatalf("请求体不是合法 JSON: %v", err)
	}
	params, ok := body["params"].(map[string]any)
	if !ok {
		t.Fatalf("params 字段缺失: %s", reqs[0].Body)
	}
	if params["topicType"] != float64(0) {
		t.Fatalf("topicType 序列化错误: %v", params["topicType"])
	}
}

func TestBusinessFailureWithStringData(t *testing.T) {
	// 业务失败时 data 可能是字符串（如"群组不存在"），不应触发 JSON 类型不匹配。
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/api/open/topic/detail", 200,
		`{"code":600,"msg":"数据处理异常","data":"群组不存在"}`)
	client := newTestClient(mock)

	_, err := client.Topic().Detail(context.Background(), 999)
	if err == nil {
		t.Fatal("期望业务失败返回错误")
	}
	e, _ := AsError(err)
	if e.Code != 600 {
		t.Fatalf("期望 code=600，实际 %d", e.Code)
	}
	if !strings.Contains(e.Msg, "群组不存在") {
		t.Fatalf("data 中的错误描述应追加到 msg: %s", e.Msg)
	}
}

func TestOpenMessageDetailURL(t *testing.T) {
	client := newTestClient(newMockHTTPRequester(), WithBaseURL("https://www.pushplus.plus/"))
	url := client.OpenMessage().DetailURL("SC9")
	if url != "https://www.pushplus.plus/shortMessage/SC9" {
		t.Fatalf("DetailURL 拼接错误: %s", url)
	}
}

func TestClawBotQrcodeStatusParamName(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	client := newTestClient(mock)

	if err := client.ClawBot().GetQrcodeStatus(context.Background(), "QR1"); err != nil {
		t.Fatalf("GetQrcodeStatus 失败: %v", err)
	}
	reqs := mock.requestsTo("getQrcodeStatus")
	if !strings.Contains(reqs[0].URL, "getQrcodeStatus=QR1") {
		t.Fatalf("query 参数名应保持官方拼写 getQrcodeStatus: %s", reqs[0].URL)
	}
}

func TestSettingReceiveLimitParamName(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	client := newTestClient(mock)

	if err := client.Setting().ChangeReceiveLimit(context.Background(), 1); err != nil {
		t.Fatalf("ChangeReceiveLimit 失败: %v", err)
	}
	reqs := mock.requestsTo("changeRecevieLimit")
	if len(reqs) != 1 || !strings.Contains(reqs[0].URL, "recevieLimit=1") {
		t.Fatalf("query 参数名应保持官方拼写 recevieLimit: %+v", reqs)
	}
}
