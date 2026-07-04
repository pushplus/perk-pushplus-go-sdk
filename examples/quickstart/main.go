// PushPlus Go SDK 快速开始示例。
//
// 运行前设置环境变量：
//
//	export PUSHPLUS_TOKEN=your_user_token
//	export PUSHPLUS_SECRET_KEY=your_secret_key  # 调用开放接口才需要
//
// 然后：go run ./examples/quickstart
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pushplus "github.com/pushplus/perk-pushplus-go-sdk"
)

func main() {
	token := os.Getenv("PUSHPLUS_TOKEN")
	if token == "" {
		log.Fatal("请先设置环境变量 PUSHPLUS_TOKEN")
	}

	client := pushplus.NewClient(
		pushplus.WithToken(token),
		pushplus.WithSecretKey(os.Getenv("PUSHPLUS_SECRET_KEY")),
	)
	ctx := context.Background()

	// 1. 发送一条最简单的消息
	shortCode, err := client.SendSimple(ctx, "Go SDK 测试", "Hello PushPlus")
	if err != nil {
		handleError(err)
		return
	}
	fmt.Println("发送成功，流水号:", shortCode)

	// 2. 使用 Markdown 模板发送
	_, err = client.Send(ctx, &pushplus.SendRequest{
		Title:     "部署完成",
		Content:   "# v1.0.0\n- env: prod",
		Template:  pushplus.TemplateMarkdown,
		Timestamp: time.Now().Add(5*time.Second).UnixMilli(), // 时效控制
	})
	if err != nil {
		handleError(err)
	}

	// 3. 多渠道发送
	req := (&pushplus.BatchSendRequest{
		Title:   "多渠道告警",
		Content: "CPU > 90%",
	}).
		AddChannel(pushplus.ChannelWechat, "").
		AddChannel(pushplus.ChannelWebhook, "bark")
	results, err := client.BatchSend(ctx, req)
	if err != nil {
		handleError(err)
	} else {
		for _, r := range results {
			fmt.Printf("渠道 %s: code=%d shortCode=%s\n", r.Channel, r.Code, r.ShortCode)
		}
	}

	// 4. 开放接口（需 secretKey；AccessKey 完全自动管理）
	if os.Getenv("PUSHPLUS_SECRET_KEY") != "" {
		me, err := client.User().MyInfo(ctx)
		if err != nil {
			handleError(err)
		} else {
			fmt.Println("当前用户:", me.NickName)
		}

		msgs, err := client.OpenMessage().List(ctx, pushplus.NewPageQuery(1, 10))
		if err != nil {
			handleError(err)
		} else {
			fmt.Printf("最近消息 %d 条，共 %d 条\n", len(msgs.List), msgs.Total)
		}

		result, err := client.OpenMessage().QueryResult(ctx, shortCode)
		if err != nil {
			handleError(err)
		} else {
			fmt.Println("发送结果状态:", result.StatusEnum())
		}
	}
}

func handleError(err error) {
	e, ok := pushplus.AsError(err)
	if !ok {
		log.Println("未知错误:", err)
		return
	}
	if e.IsRateLimited() {
		log.Println("PushPlus 限流，今天暂停推送:", e.Msg)
		return
	}
	switch e.ErrorCode() {
	case pushplus.ErrorCodeInvalidToken:
		log.Println("token 错误，立即排查配置")
	case pushplus.ErrorCodeNotVerified:
		log.Println("账号未实名认证")
	case pushplus.ErrorCodeInsufficientPoints:
		log.Println("积分不足")
	default:
		log.Printf("PushPlus 失败: code=%d, msg=%s", e.Code, e.Msg)
	}
}
