package pushplus

import (
	"context"
	"testing"
	"time"
)

const rateLimitedResponse = `{"code":900,"msg":"请求次数过多"}`

func TestRateLimitGuardBlocksAfter900(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("/send", 200, rateLimitedResponse)
	client := newTestClient(mock)

	_, err := client.SendSimple(context.Background(), "t", "c")
	e, _ := AsError(err)
	if e == nil || !e.IsRateLimited() {
		t.Fatalf("期望限流错误，实际 %v", err)
	}

	// 第二次调用应被本地短路，不再发起 HTTP。
	_, err = client.SendSimple(context.Background(), "t", "c")
	e, _ = AsError(err)
	if e == nil || !e.IsRateLimited() {
		t.Fatalf("期望本地限流守卫短路，实际 %v", err)
	}
	if len(mock.requestsTo("/send")) != 1 {
		t.Fatalf("命中 900 后同 token 不应再发起 HTTP，实际请求 %d 次", len(mock.requestsTo("/send")))
	}

	// BlockedUntil 应返回非零解禁时间。
	if client.RateLimitGuard().BlockedUntil("test-token").IsZero() {
		t.Fatal("BlockedUntil 应返回解禁时间")
	}

	// Clear 后恢复放行。
	client.RateLimitGuard().Clear("test-token")
	mock.enqueue("/send", 200, `{"code":200,"msg":"ok","data":"SC1"}`)
	if _, err := client.SendSimple(context.Background(), "t", "c"); err != nil {
		t.Fatalf("Clear 后应恢复发送，实际 %v", err)
	}
}

func TestRateLimitGuardDisabled(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("/send", 200, rateLimitedResponse)
	mock.enqueue("/send", 200, rateLimitedResponse)
	client := newTestClient(mock, WithRateLimitGuardEnabled(false))

	_, _ = client.SendSimple(context.Background(), "t", "c")
	_, _ = client.SendSimple(context.Background(), "t", "c")

	if len(mock.requestsTo("/send")) != 2 {
		t.Fatalf("关闭守卫后每次调用都应发起 HTTP，实际 %d 次", len(mock.requestsTo("/send")))
	}
}

func TestRateLimitCustomCooldown(t *testing.T) {
	current := time.Date(2026, 7, 4, 10, 0, 0, 0, time.Local)
	mock := newMockHTTPRequester()
	mock.enqueue("/send", 200, rateLimitedResponse)
	client := newTestClient(mock,
		WithRateLimitCooldown(48*time.Hour),
		WithClock(func() time.Time { return current }))

	_, _ = client.SendSimple(context.Background(), "t", "c")

	until := client.RateLimitGuard().BlockedUntil("test-token")
	expect := current.Add(48 * time.Hour)
	if !until.Equal(expect) {
		t.Fatalf("自定义 cooldown 解禁时间期望 %v，实际 %v", expect, until)
	}

	// 冷却期内仍被短路。
	current = current.Add(47 * time.Hour)
	_, err := client.SendSimple(context.Background(), "t", "c")
	e, _ := AsError(err)
	if e == nil || !e.IsRateLimited() {
		t.Fatalf("冷却期内应被短路，实际 %v", err)
	}

	// 冷却期结束后自动解除。
	current = current.Add(2 * time.Hour)
	mock.enqueue("/send", 200, `{"code":200,"msg":"ok","data":"SC1"}`)
	if _, err := client.SendSimple(context.Background(), "t", "c"); err != nil {
		t.Fatalf("冷却期结束后应恢复发送，实际 %v", err)
	}
}

func TestRateLimitDefaultCooldownIsNextMidnight(t *testing.T) {
	current := time.Date(2026, 7, 4, 15, 30, 0, 0, time.Local)
	mock := newMockHTTPRequester()
	mock.enqueue("/send", 200, rateLimitedResponse)
	client := newTestClient(mock, WithClock(func() time.Time { return current }))

	_, _ = client.SendSimple(context.Background(), "t", "c")

	until := client.RateLimitGuard().BlockedUntil("test-token")
	expect := time.Date(2026, 7, 5, 0, 0, 0, 0, time.Local)
	if !until.Equal(expect) {
		t.Fatalf("默认应禁推至次日 0 点 %v，实际 %v", expect, until)
	}
}

func TestRateLimitGuardIsPerToken(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("/send", 200, rateLimitedResponse)
	mock.enqueue("/send", 200, `{"code":200,"msg":"ok","data":"SC2"}`)
	client := newTestClient(mock)

	_, _ = client.SendSimple(context.Background(), "t", "c")

	// 其它 token 不受影响。
	if _, err := client.Send(context.Background(), &SendRequest{Token: "another-token", Content: "c"}); err != nil {
		t.Fatalf("其它 token 不应被短路，实际 %v", err)
	}
}
