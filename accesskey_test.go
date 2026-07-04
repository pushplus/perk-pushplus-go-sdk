package pushplus

import (
	"context"
	"testing"
	"time"
)

func TestAccessKeyCachedAcrossCalls(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/api/open/user/myInfo", 200, `{"code":200,"msg":"ok","data":{"nickName":"perk"}}`)
	mock.enqueue("/api/open/user/token", 200, `{"code":200,"msg":"ok","data":"tok"}`)
	client := newTestClient(mock)

	if _, err := client.User().MyInfo(context.Background()); err != nil {
		t.Fatalf("MyInfo 失败: %v", err)
	}
	if _, err := client.User().GetToken(context.Background()); err != nil {
		t.Fatalf("GetToken 失败: %v", err)
	}

	akReqs := mock.requestsTo("getAccessKey")
	if len(akReqs) != 1 {
		t.Fatalf("AccessKey 应只获取 1 次（缓存生效），实际 %d 次", len(akReqs))
	}
}

func TestOpenAPICarriesAccessKeyHeader(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/api/open/user/myInfo", 200, `{"code":200,"msg":"ok","data":{}}`)
	client := newTestClient(mock)

	if _, err := client.User().MyInfo(context.Background()); err != nil {
		t.Fatalf("MyInfo 失败: %v", err)
	}
	reqs := mock.requestsTo("/api/open/user/myInfo")
	if reqs[0].Headers[headerAccessKey] != "AK-1" {
		t.Fatalf("开放接口应携带 access-key header，实际 headers=%v", reqs[0].Headers)
	}
}

func TestAccessKeyRefreshRetryOn401(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, `{"code":200,"msg":"ok","data":{"accessKey":"AK-old","expiresIn":7200}}`)
	mock.enqueue("getAccessKey", 200, `{"code":200,"msg":"ok","data":{"accessKey":"AK-new","expiresIn":7200}}`)
	mock.enqueue("/api/open/user/myInfo", 200, `{"code":401,"msg":"请求未授权"}`)
	mock.enqueue("/api/open/user/myInfo", 200, `{"code":200,"msg":"ok","data":{"nickName":"perk"}}`)
	client := newTestClient(mock)

	info, err := client.User().MyInfo(context.Background())
	if err != nil {
		t.Fatalf("401 后应刷新 AccessKey 并重试成功，实际失败: %v", err)
	}
	if info.NickName != "perk" {
		t.Fatalf("重试结果解析错误: %+v", info)
	}

	akReqs := mock.requestsTo("getAccessKey")
	if len(akReqs) != 2 {
		t.Fatalf("401 后应重新获取 AccessKey，getAccessKey 期望 2 次，实际 %d", len(akReqs))
	}
	infoReqs := mock.requestsTo("/api/open/user/myInfo")
	if len(infoReqs) != 2 {
		t.Fatalf("业务请求期望重试 1 次（共 2 次），实际 %d", len(infoReqs))
	}
	if infoReqs[1].Headers[headerAccessKey] != "AK-new" {
		t.Fatalf("重试请求应携带新 AccessKey，实际 %v", infoReqs[1].Headers)
	}
}

func TestAccessKeyRetryStillFails(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/api/open/user/myInfo", 200, `{"code":401,"msg":"请求未授权"}`)
	mock.enqueue("/api/open/user/myInfo", 200, `{"code":401,"msg":"请求未授权"}`)
	client := newTestClient(mock)

	_, err := client.User().MyInfo(context.Background())
	if err == nil {
		t.Fatal("重试后仍 401 应返回错误")
	}
	e, _ := AsError(err)
	if e.Code != 401 {
		t.Fatalf("期望 code=401，实际 %d", e.Code)
	}
}

func TestAccessKeyExpiryTriggersRefresh(t *testing.T) {
	current := time.Date(2026, 7, 4, 10, 0, 0, 0, time.Local)
	mock := newMockHTTPRequester()
	// expiresIn=400，提前刷新 300 秒 => 有效期 100 秒
	mock.enqueue("getAccessKey", 200, `{"code":200,"msg":"ok","data":{"accessKey":"AK-1","expiresIn":400}}`)
	mock.enqueue("getAccessKey", 200, `{"code":200,"msg":"ok","data":{"accessKey":"AK-2","expiresIn":400}}`)
	mock.enqueue("/api/open/user/token", 200, `{"code":200,"msg":"ok","data":"t1"}`)
	mock.enqueue("/api/open/user/token", 200, `{"code":200,"msg":"ok","data":"t2"}`)
	client := newTestClient(mock, WithClock(func() time.Time { return current }))

	if _, err := client.User().GetToken(context.Background()); err != nil {
		t.Fatalf("首次调用失败: %v", err)
	}
	current = current.Add(101 * time.Second)
	if _, err := client.User().GetToken(context.Background()); err != nil {
		t.Fatalf("过期后调用失败: %v", err)
	}

	akReqs := mock.requestsTo("getAccessKey")
	if len(akReqs) != 2 {
		t.Fatalf("过期后应重新获取 AccessKey，期望 2 次，实际 %d", len(akReqs))
	}
	tokenReqs := mock.requestsTo("/api/open/user/token")
	if tokenReqs[1].Headers[headerAccessKey] != "AK-2" {
		t.Fatalf("过期后应使用新 AccessKey，实际 %v", tokenReqs[1].Headers)
	}
}

func TestGetAccessKeyRequiresSecretKey(t *testing.T) {
	mock := newMockHTTPRequester()
	client := NewClient(WithToken("only-token"), WithHTTPRequester(mock))

	_, err := client.User().MyInfo(context.Background())
	if err == nil {
		t.Fatal("缺少 secretKey 时开放接口应报错")
	}
	if len(mock.Requests) != 0 {
		t.Fatal("缺少 secretKey 时不应发起 HTTP 请求")
	}
}
