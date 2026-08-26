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

func TestQQBotBindAndGroupConfig(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/api/open/qqBot/getBindLink", 200,
		`{"code":200,"msg":"ok","data":{"url":"https://qun.qq.com/qunpro/robot/share?robot_appid=1","bindCode":"A1B2C3","expireSeconds":300,"botName":"pushplus"}}`)
	mock.enqueue("/api/open/qqBot/botInfo", 200,
		`{"code":200,"msg":"ok","data":{"isBind":1,"receiveStatus":1,"createTime":"2026-08-26 10:00:00","botInfo":{"appId":"1","username":"pushplus"}}}`)
	mock.enqueue("/api/open/qqBot/groupList", 200,
		`{"code":200,"msg":"ok","data":[{"id":9,"groupOpenId":"OPEN-1","status":1,"groupName":"运维告警群","groupTags":["运维"],"groupMemberNum":128}]}`)
	mock.enqueue("/api/open/qqBot/list", 200,
		`{"code":200,"msg":"ok","data":{"pageNum":1,"pageSize":20,"total":1,"pages":1,"list":[{"id":3,"qqName":"运维告警群","qqCode":"ops-group","sendType":2,"qqGroupId":9}]}}`)
	client := newTestClient(mock)
	ctx := context.Background()

	link, err := client.QQBot().GetBindLink(ctx, true)
	if err != nil {
		t.Fatalf("GetBindLink 失败: %v", err)
	}
	if link.BindCode != "A1B2C3" || link.ExpireSeconds != 300 {
		t.Fatalf("绑定链接解析错误: %+v", link)
	}
	if !strings.Contains(mock.requestsTo("getBindLink")[0].URL, "refresh=true") {
		t.Fatalf("refresh 参数未传递: %s", mock.requestsTo("getBindLink")[0].URL)
	}

	bind, err := client.QQBot().BotInfo(ctx)
	if err != nil {
		t.Fatalf("BotInfo 失败: %v", err)
	}
	if bind.IsBind != 1 || bind.BotInfo == nil || bind.BotInfo.Username != "pushplus" {
		t.Fatalf("绑定状态解析错误: %+v", bind)
	}

	groups, err := client.QQBot().GroupList(ctx)
	if err != nil {
		t.Fatalf("GroupList 失败: %v", err)
	}
	if len(groups) != 1 || groups[0].ID != 9 || len(groups[0].GroupTags) != 1 {
		t.Fatalf("QQ 群列表解析错误: %+v", groups)
	}

	if err := client.QQBot().Add(ctx, &QQBotSaveRequest{QQName: "运维告警群", QQCode: "ops-group", QQGroupID: 9}); err != nil {
		t.Fatalf("Add 失败: %v", err)
	}
	var addBody map[string]any
	if err := json.Unmarshal([]byte(mock.requestsTo("/api/open/qqBot/add")[0].Body), &addBody); err != nil {
		t.Fatalf("解析新增请求体失败: %v", err)
	}
	if addBody["sendType"] != float64(sendTypeQQGroup) || addBody["qqGroupId"] != float64(9) {
		t.Fatalf("新增请求体缺少默认 sendType 或群编号: %v", addBody)
	}

	page, err := client.QQBot().List(ctx, NewPageQuery(1, 20))
	if err != nil {
		t.Fatalf("List 失败: %v", err)
	}
	if len(page.List) != 1 || page.List[0].QQCode != "ops-group" {
		t.Fatalf("配置列表解析错误: %+v", page)
	}

	if err := client.QQBot().Delete(ctx, 3); err != nil {
		t.Fatalf("Delete 失败: %v", err)
	}
	del := mock.requestsTo("/api/open/qqBot/delete")[0]
	if del.Method != "DELETE" || !strings.Contains(del.URL, "id=3") {
		t.Fatalf("删除请求不正确: %s %s", del.Method, del.URL)
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

func TestFormCreateSavePublish(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/push/api/open/form/create", 200,
		`{"code":200,"msg":"ok","data":{"id":10001,"title":"用户满意度调查","status":0}}`)
	mock.enqueue("/push/api/open/form/save", 200, `{"code":200,"msg":"ok"}`)
	mock.enqueue("/push/api/open/form/publish", 200,
		`{"code":200,"msg":"ok","data":{"id":10001,"formCode":"a1b2c3d4","status":1}}`)
	client := newTestClient(mock)
	ctx := context.Background()

	created, err := client.Form().Create(ctx, "用户满意度调查")
	if err != nil {
		t.Fatalf("Create 失败: %v", err)
	}
	if created.ID != 10001 {
		t.Fatalf("期望 id=10001，实际 %d", created.ID)
	}
	err = client.Form().Save(ctx, &FormSaveRequest{
		ID:    10001,
		Title: "用户满意度调查",
		Items: []map[string]any{{"id": "q1", "type": "input", "label": "姓名"}},
	})
	if err != nil {
		t.Fatalf("Save 失败: %v", err)
	}
	published, err := client.Form().Publish(ctx, 10001)
	if err != nil {
		t.Fatalf("Publish 失败: %v", err)
	}
	if published.FormCode != "a1b2c3d4" {
		t.Fatalf("formCode 错误: %s", published.FormCode)
	}
	reqs := mock.requestsTo("/push/api/open/form/publish")
	if len(reqs) != 1 || !strings.Contains(reqs[0].URL, "id=10001") {
		t.Fatalf("publish 应带 query id: %+v", reqs)
	}
}

func TestExcelWriteCellsAndSaveObject(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/push/api/open/excel/writeCells", 200,
		`{"code":200,"msg":"ok","data":{"docCode":"Sh3xY7kP","publishDirty":true}}`)
	mock.enqueue("/push/api/open/excel/saveContent", 200,
		`{"code":200,"msg":"ok","data":{"docCode":"Sh3xY7kP","publishDirty":true}}`)
	client := newTestClient(mock)
	ctx := context.Background()

	_, err := client.Excel().WriteCells(ctx, "Sh3xY7kP", "A2", [][]any{{"2026-08-13", 12800}}, "Sheet1")
	if err != nil {
		t.Fatalf("WriteCells 失败: %v", err)
	}
	_, err = client.Excel().SaveContent(ctx, "Sh3xY7kP", map[string]any{"sheetOrder": []string{"sheet-1"}})
	if err != nil {
		t.Fatalf("SaveContent 失败: %v", err)
	}

	writeReqs := mock.requestsTo("/push/api/open/excel/writeCells")
	var writeBody map[string]any
	if err := json.Unmarshal([]byte(writeReqs[0].Body), &writeBody); err != nil {
		t.Fatalf("write body 不是 JSON: %v", err)
	}
	if writeBody["range"] != "A2" || writeBody["sheetName"] != "Sheet1" {
		t.Fatalf("writeCells 参数错误: %s", writeReqs[0].Body)
	}

	saveReqs := mock.requestsTo("/push/api/open/excel/saveContent")
	var saveBody map[string]any
	if err := json.Unmarshal([]byte(saveReqs[0].Body), &saveBody); err != nil {
		t.Fatalf("save body 不是 JSON: %v", err)
	}
	content, ok := saveBody["content"].(string)
	if !ok || !strings.Contains(content, "sheet-1") {
		t.Fatalf("content 应为 JSON 字符串: %s", saveReqs[0].Body)
	}
}

func TestFriendAndTopicUserBlacklist(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/api/open/friend/addBlacklist", 200, `{"code":200,"msg":"ok"}`)
	mock.enqueue("/api/open/friend/blacklistList", 200,
		`{"code":200,"msg":"ok","data":{"pageNum":1,"pageSize":20,"total":1,"pages":1,"list":[{"id":4,"friendId":1322,"nickName":"昵称"}]}}`)
	mock.enqueue("/api/open/friend/removeBlacklist", 200, `{"code":200,"msg":"ok"}`)
	mock.enqueue("/api/open/topicUser/addBlacklist", 200, `{"code":200,"msg":"ok"}`)
	mock.enqueue("/api/open/topicUser/blacklistList", 200,
		`{"code":200,"msg":"ok","data":{"pageNum":1,"list":[{"id":1,"userId":1322}]}}`)
	mock.enqueue("/api/open/topicUser/removeBlacklist", 200, `{"code":200,"msg":"ok"}`)
	client := newTestClient(mock)
	ctx := context.Background()

	if err := client.Friend().AddBlacklist(ctx, 1322); err != nil {
		t.Fatalf("AddBlacklist 失败: %v", err)
	}
	page, err := client.Friend().BlacklistList(ctx, NewPageQuery(1, 20))
	if err != nil {
		t.Fatalf("BlacklistList 失败: %v", err)
	}
	if page.List[0].FriendID != 1322 {
		t.Fatalf("friendId 错误: %+v", page.List[0])
	}
	if err := client.Friend().RemoveBlacklist(ctx, 4); err != nil {
		t.Fatalf("RemoveBlacklist 失败: %v", err)
	}
	if err := client.TopicUser().AddBlacklist(ctx, 10); err != nil {
		t.Fatalf("topic AddBlacklist 失败: %v", err)
	}
	users, err := client.TopicUser().BlacklistList(ctx, NewTopicUserListQuery(1, 20, 100))
	if err != nil {
		t.Fatalf("topic BlacklistList 失败: %v", err)
	}
	if users.List[0].ID != 1 {
		t.Fatalf("id 错误: %+v", users.List[0])
	}
	if err := client.TopicUser().RemoveBlacklist(ctx, 1); err != nil {
		t.Fatalf("topic RemoveBlacklist 失败: %v", err)
	}
	reqs := mock.requestsTo("/api/open/friend/addBlacklist")
	if len(reqs) != 1 || !strings.Contains(reqs[0].URL, "friendId=1322") {
		t.Fatalf("friend addBlacklist query 错误: %+v", reqs)
	}
	topicList := mock.requestsTo("/api/open/topicUser/blacklistList")
	var body map[string]any
	if err := json.Unmarshal([]byte(topicList[0].Body), &body); err != nil {
		t.Fatalf("body 不是 JSON: %v", err)
	}
	params := body["params"].(map[string]any)
	if params["topicId"] != float64(100) {
		t.Fatalf("topicId 错误: %v", params["topicId"])
	}
}

func TestFormListUsesCurrentAndParams(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/push/api/open/form/list", 200,
		`{"code":200,"msg":"ok","data":{"pageNum":1,"pageSize":20,"total":0,"pages":0,"list":[]}}`)
	client := newTestClient(mock)
	status := 1
	_, err := client.Form().List(context.Background(), NewFormListQueryFilter(1, 20, "满意度", &status))
	if err != nil {
		t.Fatalf("Form List 失败: %v", err)
	}
	reqs := mock.requestsTo("/push/api/open/form/list")
	var body map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatalf("body 不是 JSON: %v", err)
	}
	if body["current"] != float64(1) {
		t.Fatalf("current 错误: %v", body["current"])
	}
	params := body["params"].(map[string]any)
	if params["keyword"] != "满意度" || params["status"] != float64(1) {
		t.Fatalf("params 错误: %v", params)
	}
}

func TestDocImport(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/push/api/open/doc/import", 200,
		`{"code":200,"msg":"ok","data":{"docCode":"Ab3xY7kP","title":"本周工作同步"}}`)
	client := newTestClient(mock)
	ctx := context.Background()

	imported, err := client.Doc().ImportWord(ctx, []byte("hello"), "本周工作同步.docx")
	if err != nil {
		t.Fatalf("ImportWord 失败: %v", err)
	}
	if imported.DocCode != "Ab3xY7kP" {
		t.Fatalf("docCode 错误: %s", imported.DocCode)
	}
	reqs := mock.requestsTo("/push/api/open/doc/import")
	if len(reqs) != 1 || !reqs[0].Raw {
		t.Fatalf("import 应走 ExecuteRaw: %+v", reqs)
	}
	if !strings.Contains(reqs[0].Headers["Content-Type"], "multipart/form-data; boundary=") {
		t.Fatalf("Content-Type 错误: %s", reqs[0].Headers["Content-Type"])
	}
}

func TestExcelImport(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/push/api/open/excel/import", 200,
		`{"code":200,"msg":"ok","data":{"docCode":"Sh3xY7kP","title":"销售日报"}}`)
	client := newTestClient(mock)
	ctx := context.Background()

	imported, err := client.Excel().ImportExcel(ctx, []byte("xlsx"), "销售日报.xlsx")
	if err != nil {
		t.Fatalf("ImportExcel 失败: %v", err)
	}
	if imported.DocCode != "Sh3xY7kP" {
		t.Fatalf("docCode 错误: %s", imported.DocCode)
	}
}
