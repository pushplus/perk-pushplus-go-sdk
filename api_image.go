package pushplus

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ImageAPI 开放接口 - 图片服务（文档「十二. 图片服务接口」）。
//
// PushPlus 提供基于七牛云的图片图床（30 天有效，可主动删除）。
// UploadFile/UploadBytes 把"获取上传凭证 → 表单上传 → 解析返回 URL"封装成一步。
// 注意：上传图片的真正请求按七牛云规范以 multipart/form-data 提交到 uploadUrl，
// 不携带 PushPlus 的 access-key；其余接口走 PushPlus 开放接口，自动带上 access-key。
type ImageAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newImageAPI(c *core, akm *AccessKeyManager) *ImageAPI {
	return &ImageAPI{core: c, akm: akm}
}

/* ------------------------- 1. 获取上传凭证 ------------------------- */

// GetUploadToken 获取上传凭证。
func (a *ImageAPI) GetUploadToken(ctx context.Context) (*ImageUploadToken, error) {
	return executeOpen[*ImageUploadToken](ctx, a.core, a.akm, "GET", "/api/open/userImage/uploadToken", nil)
}

/* ------------------------- 2. 上传图片 ------------------------- */

// Upload 使用已获取的上传凭证把文件上传到七牛云。
// contentType 为空时退化为 application/octet-stream。
func (a *ImageAPI) Upload(ctx context.Context, token *ImageUploadToken, fileBytes []byte, fileName, contentType string) (*ImageUploadResult, error) {
	if token == nil {
		return nil, newError(-1, "上传凭证 token 不能为 nil")
	}
	if strings.TrimSpace(token.UploadToken) == "" {
		return nil, newError(-1, "上传凭证 uploadToken 不能为空")
	}
	uploadURL := token.UploadURL
	if strings.TrimSpace(uploadURL) == "" {
		uploadURL = token.UploadHost
	}
	if strings.TrimSpace(uploadURL) == "" {
		return nil, newError(-1, "上传凭证未返回 uploadUrl/uploadHost")
	}
	return a.UploadTo(ctx, uploadURL, token.UploadToken, fileBytes, fileName, contentType)
}

// UploadTo 低层方法：直接指定七牛上传地址与 token 上传。
// 该请求按七牛云规范提交 multipart/form-data，不携带 PushPlus 的 access-key。
func (a *ImageAPI) UploadTo(ctx context.Context, uploadURL, uploadToken string, fileBytes []byte, fileName, contentType string) (*ImageUploadResult, error) {
	if strings.TrimSpace(uploadURL) == "" {
		return nil, newError(-1, "uploadUrl 不能为空")
	}
	if strings.TrimSpace(uploadToken) == "" {
		return nil, newError(-1, "uploadToken 不能为空")
	}
	if len(fileBytes) == 0 {
		return nil, newError(-1, "上传文件内容不能为空")
	}
	if strings.TrimSpace(fileName) == "" {
		fileName = "file"
	}
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}

	boundary := "----PushPlusBoundary" + randomHex(16)
	body := buildMultipartBody(boundary, uploadToken, fileName, contentType, fileBytes)
	headers := map[string]string{
		"Content-Type": "multipart/form-data; boundary=" + boundary,
	}

	resp, err := a.core.http.ExecuteRaw(ctx, "POST", uploadURL, headers, body)
	if err != nil {
		if _, ok := AsError(err); ok {
			return nil, err
		}
		return nil, newErrorWithCause(-1, "上传图片到七牛云失败: "+err.Error(), err)
	}
	if !resp.IsSuccessful() {
		return nil, newError(resp.StatusCode,
			fmt.Sprintf("上传图片到七牛云失败: status=%d, body=%s", resp.StatusCode, resp.Body))
	}

	var result ImageUploadResult
	if err := json.Unmarshal([]byte(resp.Body), &result); err != nil {
		return nil, newErrorWithCause(-1, "解析七牛云上传响应失败: "+err.Error()+", payload="+resp.Body, err)
	}
	if result.Errno != 0 {
		return nil, newError(result.Errno,
			fmt.Sprintf("七牛云上传失败: errno=%d, msg=%s", result.Errno, result.Msg))
	}
	return &result, nil
}

/* ------------------------- 便捷上传方法 ------------------------- */

// UploadBytes 便捷方法：自动获取上传凭证后上传字节数组，按文件扩展名猜测 MIME 类型。
func (a *ImageAPI) UploadBytes(ctx context.Context, fileBytes []byte, fileName string) (*ImageUploadResult, error) {
	return a.UploadBytesWithType(ctx, fileBytes, fileName, guessContentTypeByName(fileName))
}

// UploadBytesWithType 便捷方法：自动获取上传凭证后上传字节数组，可指定 MIME 类型。
func (a *ImageAPI) UploadBytesWithType(ctx context.Context, fileBytes []byte, fileName, contentType string) (*ImageUploadResult, error) {
	token, err := a.GetUploadToken(ctx)
	if err != nil {
		return nil, err
	}
	return a.Upload(ctx, token, fileBytes, fileName, contentType)
}

// UploadFile 便捷方法：自动获取上传凭证后上传指定路径的文件。
func (a *ImageAPI) UploadFile(ctx context.Context, filePath string) (*ImageUploadResult, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, newError(-1, "上传文件路径不能为空")
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, newErrorWithCause(-1, "读取上传文件失败: "+err.Error(), err)
	}
	fileName := filepath.Base(filePath)
	return a.UploadBytesWithType(ctx, data, fileName, guessContentTypeByName(fileName))
}

// UploadStream 便捷方法：自动获取上传凭证后上传输入流（不关闭传入的流）。
func (a *ImageAPI) UploadStream(ctx context.Context, r io.Reader, fileName, contentType string) (*ImageUploadResult, error) {
	if r == nil {
		return nil, newError(-1, "输入流不能为 nil")
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, newErrorWithCause(-1, "读取上传输入流失败: "+err.Error(), err)
	}
	return a.UploadBytesWithType(ctx, data, fileName, contentType)
}

/* ------------------------- 3. 图片列表 / 4. 删除图片 ------------------------- */

// List 已上传图片列表。
func (a *ImageAPI) List(ctx context.Context, query *PageQuery) (*PageResult[ImageItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[ImageItem]](ctx, a.core, a.akm, "POST", "/api/open/userImage/list", body)
}

// Delete 主动删除图片；未删除的图片默认 30 天后由系统自动清理。
func (a *ImageAPI) Delete(ctx context.Context, id int64) error {
	path := appendQuery("/api/open/userImage/delete", []queryParam{{"id", strconv.FormatInt(id, 10)}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "DELETE", path, nil)
	return err
}

/* ------------------------- 内部辅助 ------------------------- */

func buildMultipartBody(boundary, token, fileName, contentType string, fileBytes []byte) []byte {
	var buf bytes.Buffer
	buf.Grow(len(fileBytes) + 512)

	buf.WriteString("--" + boundary + "\r\n")
	buf.WriteString("Content-Disposition: form-data; name=\"token\"\r\n\r\n")
	buf.WriteString(token + "\r\n")

	buf.WriteString("--" + boundary + "\r\n")
	buf.WriteString("Content-Disposition: form-data; name=\"file\"; filename=\"" + escapeFileName(fileName) + "\"\r\n")
	buf.WriteString("Content-Type: " + contentType + "\r\n\r\n")
	buf.Write(fileBytes)
	buf.WriteString("\r\n")

	buf.WriteString("--" + boundary + "--\r\n")
	return buf.Bytes()
}

func escapeFileName(name string) string {
	name = strings.ReplaceAll(name, "\"", "_")
	name = strings.ReplaceAll(name, "\r", " ")
	name = strings.ReplaceAll(name, "\n", " ")
	return name
}

func guessContentTypeByName(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand 在正常系统上不会失败；即便失败，boundary 也只需足够独特。
		for i := range b {
			b[i] = byte(i*31 + 7)
		}
	}
	return hex.EncodeToString(b)
}
