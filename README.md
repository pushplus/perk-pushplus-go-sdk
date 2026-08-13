# perk-pushplus-go-sdk

[PushPlus(推送加)](https://www.pushplus.plus) 官方接口的 Go SDK，覆盖**消息发送接口**与**全部开放接口**。

- 纯 Go SDK，标准库 `net/http`，无额外重依赖
- AccessKey **自动获取、缓存、过期前刷新、失效自动重试**，调用方无感知
- **本地限流守卫**：发送接口命中 `code=900`（请求次数过多）时自动短路同 token 的后续调用，避免无效请求与账号进一步受限（[官方建议](https://www.pushplus.plus/doc/guide/code.html)）
- 单条 `/send`、多渠道 `/batchSend`、消息回调（`message_complate` / `add_topic_user` / `add_friend`）类型化解析
- 全部开放接口：消息、用户、消息令牌、群组、群组用户、好友、Webhook、公众号/企业微信/邮箱渠道、ClawBot、功能设置、预处理、图片服务、push 表单、push 文档、push 表格
- 强类型枚举（`Channel`、`Template`、`SendStatus`、`WebhookType`、`CallbackEvent`、`ErrorCode`）

## 快速开始

### 1. 引入依赖

```bash
go get github.com/pushplus/perk-pushplus-go-sdk
```

### 2. 发送一条消息（最简）

```go
package main

import (
    "context"
    "log"

    pushplus "github.com/pushplus/perk-pushplus-go-sdk"
)

func main() {
    client := pushplus.NewClient(
        pushplus.WithToken("your_user_token"),     // 必填
        pushplus.WithSecretKey("your_secret_key"), // 调用开放接口才需要
    )

    shortCode, err := client.SendSimple(context.Background(), "标题", "Hello PushPlus")
    if err != nil {
        log.Fatal(err)
    }
    log.Println("流水号:", shortCode)
}
```

### 3. 使用 Markdown / HTML / JSON 消息

```go
_, err := client.Send(ctx, &pushplus.SendRequest{
    Title:       "部署完成",
    Content:     "# v1.0.0\n- env: prod",
    Template:    pushplus.TemplateMarkdown,
    Topic:       "ops",                              // 群组编码
    Channel:     pushplus.ChannelWechat,             // 默认就是 wechat，可省略
    CallbackURL: "https://you/cb",                   // 异步回调
    Timestamp:   time.Now().Add(5 * time.Second).UnixMilli(), // 时效控制
})
```

发送 push 表单消息时使用 `TemplateForm`，并传入表单编码 `PushID`：

```go
_, err := client.Send(ctx, &pushplus.SendRequest{
    Title:    "表单通知",
    Content:  "您有新的表单待填写",
    Template: pushplus.TemplateForm,
    PushID:   "表单编码",
})
```

### 4. 多渠道发送（`/batchSend`）

```go
req := (&pushplus.BatchSendRequest{
    Title:   "多渠道告警",
    Content: "CPU > 90%",
}).
    AddChannel(pushplus.ChannelWechat, "").
    AddChannel(pushplus.ChannelWebhook, "bark").
    AddChannel(pushplus.ChannelExtension, "")

results, err := client.BatchSend(ctx, req)
```

`AddChannel(ch, option)` 可链式调用，SDK 自动用逗号拼接 channel 与 option，与官方文档示例语义一致。

### 5. 异步发送

Go 中自行开 goroutine 即可：

```go
go func() {
    _, err := client.Send(ctx, &pushplus.SendRequest{Title: "t", Content: "c"})
    // ...
}()
```

## 开放接口（全量）

需要在 PushPlus 后台「开发设置」中：开启开放接口、配置 `secretKey`、把调用方所在公网 IP 加入安全 IP 列表。

```go
client := pushplus.NewClient(
    pushplus.WithToken("user_token"),
    pushplus.WithSecretKey("xxx"),
)

// AccessKey 完全自动管理 —— 直接调用就好
me, err := client.User().MyInfo(ctx)
msgs, err := client.OpenMessage().List(ctx, pushplus.NewPageQuery(1, 20))
result, err := client.OpenMessage().QueryResult(ctx, shortCode)

// 群组
topics, err := client.Topic().List(ctx, pushplus.NewTopicListQuery(1, 20, 0))
qr, err := client.Topic().QrCode(ctx, 1, nil, nil)

// Webhook 渠道配置
id, err := client.Webhook().Add(ctx, &pushplus.WebhookSaveRequest{
    WebhookCode: "bark",
    WebhookName: "我的 Bark",
    WebhookType: int(pushplus.WebhookTypeBark),
    WebhookURL:  "https://api.day.app/xxxx",
})

// 功能设置
err = client.Setting().ChangeIsSend(ctx, 1)

// 图片服务（一行上传到 PushPlus 图床）
img, err := client.Image().UploadFile(ctx, "/tmp/logo.png")
url := img.URL // 直接拿到可访问的图片地址
page, err := client.Image().List(ctx, pushplus.NewPageQuery(1, 10))
err = client.Image().Delete(ctx, page.List[0].ID)

// push 表单
form, err := client.Form().Create(ctx, "用户满意度调查")
err = client.Form().Save(ctx, &pushplus.FormSaveRequest{
    ID:    form.ID,
    Title: "用户满意度调查",
    Items: []map[string]any{{"id": "q_name", "type": "input", "label": "您的姓名", "required": true}},
})
published, err := client.Form().Publish(ctx, form.ID)

// push 文档
doc, err := client.Doc().Create(ctx, "本周工作同步")
_, err = client.Doc().SaveContent(ctx, doc.DocCode, "<h1>本周工作同步</h1><p>需求评审。</p>")
anon := 0
_, err = client.Doc().UpdateShare(ctx, doc.DocCode, 1, &anon)
_, err = client.Doc().Publish(ctx, doc.DocCode)

// push 表格
sheet, err := client.Excel().Create(ctx, "销售日报")
_, err = client.Excel().WriteCells(ctx, sheet.DocCode, "A1", [][]any{
    {"日期", "销售额"},
    {"2026-08-13", 12800},
}, "Sheet1")
_, err = client.Excel().Publish(ctx, sheet.DocCode)
```

各 API 一览：

| 方法 | 对应文档章节 |
| --- | --- |
| `client.Message()` | 二/三 发送消息接口 |
| `client.AccessKey()` | 一 获取 AccessKey（一般无需手动调用） |
| `client.OpenMessage()` | 二 消息接口（开放） |
| `client.User()` | 三 用户接口 |
| `client.MessageToken()` | 四 消息 token 接口 |
| `client.Topic()` | 五 群组接口 |
| `client.TopicUser()` | 六 群组用户接口 |
| `client.Webhook()` | 七 渠道配置 - webhook |
| `client.Channel()` | 七 渠道配置 - 公众号/企业微信/邮箱 |
| `client.ClawBot()` | 八 微信 ClawBot 接口 |
| `client.Setting()` | 九 功能设置接口 |
| `client.Friend()` | 十 好友功能接口 |
| `client.Pre()` | 十一 预处理信息接口 |
| `client.Image()` | 十二 图片服务接口 |
| `client.Form()` | push 表单开放接口 |
| `client.Doc()` | push 文档开放接口 |
| `client.Excel()` | push 表格开放接口 |

## 图片服务

PushPlus 提供基于七牛云的图片图床（30 天有效，可主动删除）。SDK 把"获取上传凭证 → 表单上传 → 解析返回 URL"封装成一步：

```go
// 1) 最常用：一行上传本地文件，得到可访问的图片 URL
r, err := client.Image().UploadFile(ctx, "/tmp/logo.png")
url := r.URL

// 2) 直接上传字节
client.Image().UploadBytes(ctx, bytes, "screenshot.png")

// 3) 已上传图片列表
page, err := client.Image().List(ctx, pushplus.NewPageQuery(1, 10))

// 4) 主动删除（未删除的图片默认 30 天后由系统自动清理）
client.Image().Delete(ctx, page.List[0].ID)
```

如需自己控制凭证的获取与上传过程（例如缓存 token、分布式上传），也可以拆开调用：

```go
token, err := client.Image().GetUploadToken(ctx)
r, err := client.Image().Upload(ctx, token, bytes, "a.png", "image/png")
```

> 注意：上传图片的真正请求会按七牛云规范以 `multipart/form-data` 提交到 `uploadUrl`，**不会**携带 PushPlus 的 `access-key`；其余三个接口（获取凭证 / 列表 / 删除）走 PushPlus 开放接口，自动带上 `access-key`。

## 消息回调解析

PushPlus 在消息发送完成、群组新增用户、新增好友时会回调你预置的 URL。SDK 提供类型安全的解析器：

```go
func callbackHandler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    p, err := pushplus.ParseCallback(body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    switch p.Event {
    case pushplus.CallbackEventMessageComplete:
        handle(p.MessageInfo) // SendStatus 枚举可直接拿
    case pushplus.CallbackEventAddTopicUser:
        handle(p.TopicUserInfo)
    case pushplus.CallbackEventAddFriend:
        handle(p.FriendInfo, p.QrCode)
    }
    w.Write([]byte("ok"))
}
```

## 配置项

| 字段 | 默认 | 说明 |
| --- | --- | --- |
| `Token` | – | 用户 token / 消息 token，发送消息使用 |
| `SecretKey` | – | 用户 secretKey，调用开放接口使用 |
| `BaseURL` | `https://www.pushplus.plus` | 服务地址 |
| `ConnectTimeout` | 10s | 连接超时 |
| `ReadTimeout` | 30s | 读超时 |
| `AccessKeyRefreshAheadSeconds` | 300 | AccessKey 提前刷新秒数 |
| `LogRequest` | false | 开启请求/响应日志 |
| `RateLimitGuardEnabled` | true | 是否启用本地限流守卫 |
| `RateLimitCooldown` | 0 | 命中 `code=900` 后的本地禁推时长；0 表示到"次日 0 点" |

```go
client := pushplus.NewClient(
    pushplus.WithToken("xxx"),
    pushplus.WithSecretKey("yyy"),
    pushplus.WithBaseURL("https://www.pushplus.plus"),
    pushplus.WithConnectTimeout(10*time.Second),
    pushplus.WithReadTimeout(30*time.Second),
    pushplus.WithAccessKeyRefreshAheadSeconds(300),
    pushplus.WithLogRequest(false),
    pushplus.WithRateLimitGuardEnabled(true),
    pushplus.WithRateLimitCooldown(48*time.Hour), // 留空（0）表示禁推到次日 0 点
)
```

如需自定义 HTTP 实现，通过 `WithHTTPRequester` 注入即可。

## 异常处理

所有错误都会包装成 `*pushplus.Error`：

- HTTP 调用失败：`Code` 为 HTTP 状态码
- 业务接口返回 `code != 200`：`Code` 为业务码、`Msg` 为业务消息
- SDK 参数校验失败：`Code = -1`
- **本地限流守卫命中**（不发起 HTTP）：`Code = 900`、`Msg` 包含 "本地限流守卫" 字样

`ErrorCode` 常量已经把官方文档的全部业务码语义化，无需再用魔法数字判断。

```go
shortCode, err := client.SendSimple(ctx, "t", "c")
if err != nil {
    e, ok := pushplus.AsError(err)
    if !ok {
        log.Fatal(err)
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
```

参考：[PushPlus 接口返回码说明](https://www.pushplus.plus/doc/guide/code.html)。

## 限流守卫（code=900 自动短路）

PushPlus 在请求次数过多时会返回 `code=900`，官方文档明确建议"根据返回值判断当天是否让程序继续调用发送消息接口，否则会让账号进一步受限"。SDK 默认替你做这件事：

- 任意一次 `client.Send(...)` / `client.BatchSend(...)` 命中 `code=900` 后，SDK 会按 token 维度记下"禁推至 X 时刻"。
- 同 token 后续发送调用不再发起 HTTP，直接返回 `Error(code=900, msg="本地限流守卫…")`。
- 默认禁推到**系统时区的次日 0 点**；通过 `WithRateLimitCooldown` 可改为固定时长（例如 2 天）。
- 仅作用于 `Message()`（`Send` / `BatchSend`），开放接口不受影响。
- 进程内单例，**不跨进程共享**——多实例部署时每个进程最多被命中一次。

可观察 / 可干预：

```go
guard := client.RateLimitGuard()

until := guard.BlockedUntil("user_token") // 零值 time.Time 表示未被限流
guard.Clear("user_token")               // 人工确认服务端已解禁后立即放行
```

需要完全关闭这个行为（不推荐）：

```go
client := pushplus.NewClient(
    pushplus.WithToken("xxx"),
    pushplus.WithRateLimitGuardEnabled(false),
)
```

或者按业务把禁推时长拉长 / 缩短：

```go
client := pushplus.NewClient(
    pushplus.WithToken("xxx"),
    pushplus.WithRateLimitCooldown(48*time.Hour), // 与文档示例的"2 天恢复正常"对齐
)
```

## 构建与测试

```bash
go build ./...
go vet ./...
go test ./...

# 运行示例（需设置环境变量）
export PUSHPLUS_TOKEN=your_token
export PUSHPLUS_SECRET_KEY=your_secret_key  # 开放接口示例需要
go run ./examples/quickstart
```

## License

Apache License 2.0
