package pushplus

import (
	"context"
	"strings"
	"testing"
)

const uploadTokenResponse = `{"code":200,"msg":"ok","data":{
	"uploadToken":"QN-TOKEN","uploadHost":"https://upload.qiniup.com",
	"uploadUrl":"https://upload.qiniup.com","bucket":"pushplus","expiresIn":3600}}`

const qiniuOKResponse = `{"errno":0,"url":"https://img.pushplus.plus/a.png","key":"a.png","fsize":10}`

func TestImageGetUploadToken(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("uploadToken", 200, uploadTokenResponse)
	client := newTestClient(mock)

	token, err := client.Image().GetUploadToken(context.Background())
	if err != nil {
		t.Fatalf("GetUploadToken 失败: %v", err)
	}
	if token.UploadToken != "QN-TOKEN" || token.UploadURL != "https://upload.qiniup.com" {
		t.Fatalf("凭证解析错误: %+v", token)
	}
	reqs := mock.requestsTo("uploadToken")
	if reqs[0].Headers[headerAccessKey] == "" {
		t.Fatal("获取上传凭证应携带 access-key")
	}
}

func TestImageUploadMultipartWithoutAccessKey(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("uploadToken", 200, uploadTokenResponse)
	mock.enqueue("upload.qiniup.com", 200, qiniuOKResponse)
	client := newTestClient(mock)

	result, err := client.Image().UploadBytes(context.Background(), []byte("PNGDATA"), "logo.png")
	if err != nil {
		t.Fatalf("UploadBytes 失败: %v", err)
	}
	if result.URL != "https://img.pushplus.plus/a.png" {
		t.Fatalf("上传结果解析错误: %+v", result)
	}

	reqs := mock.requestsTo("upload.qiniup.com")
	if len(reqs) != 1 {
		t.Fatalf("期望 1 次七牛上传请求，实际 %d", len(reqs))
	}
	up := reqs[0]
	if !up.Raw {
		t.Fatal("上传应走 ExecuteRaw（二进制 multipart）")
	}
	if _, ok := up.Headers[headerAccessKey]; ok {
		t.Fatal("七牛上传请求不应携带 access-key")
	}
	if !strings.Contains(up.Headers["Content-Type"], "multipart/form-data; boundary=") {
		t.Fatalf("Content-Type 错误: %s", up.Headers["Content-Type"])
	}
	if !strings.Contains(up.Body, `name="token"`) || !strings.Contains(up.Body, "QN-TOKEN") {
		t.Fatal("multipart body 应包含 token 字段")
	}
	if !strings.Contains(up.Body, `name="file"; filename="logo.png"`) {
		t.Fatal("multipart body 应包含 file 字段与文件名")
	}
	if !strings.Contains(up.Body, "Content-Type: image/png") {
		t.Fatal("应按扩展名推断 image/png")
	}
	if !strings.Contains(up.Body, "PNGDATA") {
		t.Fatal("multipart body 应包含文件内容")
	}
}

func TestImageUploadQiniuError(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("uploadToken", 200, uploadTokenResponse)
	mock.enqueue("upload.qiniup.com", 200, `{"errno":400,"msg":"invalid token"}`)
	client := newTestClient(mock)

	_, err := client.Image().UploadBytes(context.Background(), []byte("X"), "a.png")
	if err == nil {
		t.Fatal("七牛 errno != 0 应返回错误")
	}
	e, _ := AsError(err)
	if e.Code != 400 || !strings.Contains(e.Msg, "invalid token") {
		t.Fatalf("错误信息不符: %v", err)
	}
}

func TestImageUploadValidation(t *testing.T) {
	mock := newMockHTTPRequester()
	client := newTestClient(mock)

	if _, err := client.Image().Upload(context.Background(), nil, []byte("x"), "a.png", ""); err == nil {
		t.Fatal("凭证为 nil 应报错")
	}
	if _, err := client.Image().UploadTo(context.Background(), "https://u", "tok", nil, "a.png", ""); err == nil {
		t.Fatal("文件内容为空应报错")
	}
	if len(mock.Requests) != 0 {
		t.Fatal("参数校验失败时不应发起 HTTP 请求")
	}
}

func TestImageListAndDelete(t *testing.T) {
	mock := newMockHTTPRequester()
	mock.enqueue("getAccessKey", 200, accessKeyResponse)
	mock.enqueue("/api/open/userImage/list", 200,
		`{"code":200,"msg":"ok","data":{"pageNum":1,"pageSize":10,"total":1,"pages":1,"list":[{"id":5,"imgUrl":"https://img/a.png"}]}}`)
	client := newTestClient(mock)

	page, err := client.Image().List(context.Background(), NewPageQuery(1, 10))
	if err != nil {
		t.Fatalf("List 失败: %v", err)
	}
	if len(page.List) != 1 || page.List[0].ID != 5 {
		t.Fatalf("列表解析错误: %+v", page)
	}

	if err := client.Image().Delete(context.Background(), page.List[0].ID); err != nil {
		t.Fatalf("Delete 失败: %v", err)
	}
	delReqs := mock.requestsTo("/api/open/userImage/delete")
	if len(delReqs) != 1 || !strings.Contains(delReqs[0].URL, "id=5") {
		t.Fatalf("删除请求错误: %+v", delReqs)
	}
	if delReqs[0].Method != "DELETE" {
		t.Fatalf("删除应使用 DELETE 方法，实际 %s", delReqs[0].Method)
	}
}
