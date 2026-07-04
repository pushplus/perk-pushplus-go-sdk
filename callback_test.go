package pushplus

import "testing"

func TestParseMessageCompleteCallback(t *testing.T) {
	body := `{
		"event": "message_complate",
		"messageInfo": {"message": "内容", "shortCode": "SC1", "sendStatus": 2}
	}`
	p, err := ParseCallbackString(body)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if p.Event != CallbackEventMessageComplete {
		t.Fatalf("event 解析错误: %v", p.Event)
	}
	if p.MessageInfo == nil || p.MessageInfo.ShortCode != "SC1" {
		t.Fatalf("messageInfo 解析错误: %+v", p.MessageInfo)
	}
	if p.MessageInfo.SendStatusEnum() != SendStatusSuccess {
		t.Fatalf("sendStatus 枚举转换错误: %v", p.MessageInfo.SendStatusEnum())
	}
}

func TestParseAddTopicUserCallback(t *testing.T) {
	body := `{
		"event": "add_topic_user",
		"topicUserInfo": {
			"id": 10, "openId": "o1", "topicId": 3, "userSex": 1, "isFollow": 1,
			"nickName": "小明", "havePhone": 0, "topicCode": "ops", "topicName": "运维群",
			"headImgUrl": "https://img", "emailStatus": 0
		}
	}`
	p, err := ParseCallbackString(body)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if p.Event != CallbackEventAddTopicUser {
		t.Fatalf("event 解析错误: %v", p.Event)
	}
	if p.TopicUserInfo == nil || p.TopicUserInfo.TopicCode != "ops" || p.TopicUserInfo.ID != 10 {
		t.Fatalf("topicUserInfo 解析错误: %+v", p.TopicUserInfo)
	}
}

func TestParseAddFriendCallback(t *testing.T) {
	body := `{
		"event": "add_friend",
		"friendInfo": {"token": "ft", "friendId": 7, "isFollow": 1, "nickName": "小红",
			"havePhone": 1, "createTime": "2026-07-04 10:00:00", "emailStatus": 0},
		"qrCode": "QR-XYZ"
	}`
	p, err := ParseCallbackString(body)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if p.Event != CallbackEventAddFriend {
		t.Fatalf("event 解析错误: %v", p.Event)
	}
	if p.FriendInfo == nil || p.FriendInfo.FriendID != 7 {
		t.Fatalf("friendInfo 解析错误: %+v", p.FriendInfo)
	}
	if p.QrCode != "QR-XYZ" {
		t.Fatalf("qrCode 解析错误: %s", p.QrCode)
	}
}

func TestParseCallbackIgnoresUnknownFields(t *testing.T) {
	body := `{"event": "message_complate", "messageInfo": {"shortCode": "SC1"}, "unknownField": {"x": 1}}`
	p, err := ParseCallbackString(body)
	if err != nil {
		t.Fatalf("未知字段不应导致解析失败: %v", err)
	}
	if p.MessageInfo.ShortCode != "SC1" {
		t.Fatalf("解析错误: %+v", p)
	}
}

func TestParseCallbackInvalidJSON(t *testing.T) {
	_, err := ParseCallbackString("{not json")
	if err == nil {
		t.Fatal("非法 JSON 应返回错误")
	}
}
