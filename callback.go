package pushplus

import "encoding/json"

// ParseCallback 解析 PushPlus 消息回调（消息发送完成 / 群组新增用户 / 新增好友）。
//
//	payload, err := pushplus.ParseCallback(body)
//	switch payload.Event {
//	case pushplus.CallbackEventMessageComplete:
//	    handle(payload.MessageInfo)
//	case pushplus.CallbackEventAddTopicUser:
//	    handle(payload.TopicUserInfo)
//	case pushplus.CallbackEventAddFriend:
//	    handle(payload.FriendInfo, payload.QrCode)
//	}
func ParseCallback(body []byte) (*CallbackPayload, error) {
	var payload CallbackPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, newErrorWithCause(-1, "解析 PushPlus 回调失败: "+err.Error(), err)
	}
	return &payload, nil
}

// ParseCallbackString 同 ParseCallback，入参为字符串。
func ParseCallbackString(body string) (*CallbackPayload, error) {
	return ParseCallback([]byte(body))
}
