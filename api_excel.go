package pushplus

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExcelAPI 开放接口 - push 表格。
//
// 文档：https://www.pushplus.plus/doc/ecosystem/sheet/
// 基础路径：/push/api/open/excel
//
// 表格开放接口不单独提供推送接口。发布后请通过 MessageAPI 推送分享页：
// template=excel，pushId=docCode。
type ExcelAPI struct {
	core *core
	akm  *AccessKeyManager
}

func newExcelAPI(c *core, akm *AccessKeyManager) *ExcelAPI {
	return &ExcelAPI{core: c, akm: akm}
}

// List 我的表格分页。
func (a *ExcelAPI) List(ctx context.Context, query *DocListQuery) (*PageResult[DocListItem], error) {
	var body any = query
	if query == nil {
		body = struct{}{}
	}
	return executeOpen[*PageResult[DocListItem]](ctx, a.core, a.akm, "POST", "/push/api/open/excel/list", body)
}

// Create 创建空白表格。
func (a *ExcelAPI) Create(ctx context.Context, title string) (*ExcelVo, error) {
	return executeOpen[*ExcelVo](ctx, a.core, a.akm, "POST", "/push/api/open/excel/create", map[string]any{"title": title})
}

// ImportExcel 导入 Excel（.xlsx / .xls）创建表格。标题默认取文件名；创建后默认关闭分享，需再 Publish。
func (a *ExcelAPI) ImportExcel(ctx context.Context, fileBytes []byte, fileName string) (*ExcelVo, error) {
	if strings.TrimSpace(fileName) == "" {
		fileName = "workbook.xlsx"
	}
	return executeOpenMultipart[*ExcelVo](ctx, a.core, a.akm, "/push/api/open/excel/import", fileName, guessExcelContentType(fileName), fileBytes)
}

// ImportExcelFile 从本地路径导入 Excel 创建表格。
func (a *ExcelAPI) ImportExcelFile(ctx context.Context, filePath string) (*ExcelVo, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, newErrorWithCause(-1, "读取上传文件失败: "+err.Error(), err)
	}
	name := filepath.Base(filePath)
	if name == "." || name == "/" {
		name = "workbook.xlsx"
	}
	return a.ImportExcel(ctx, b, name)
}

// Content 获取表格元信息与整表 JSON 草稿。
func (a *ExcelAPI) Content(ctx context.Context, docCode string) (*ExcelContent, error) {
	path := appendQuery("/push/api/open/excel/content", []queryParam{{"docCode", docCode}})
	return executeOpen[*ExcelContent](ctx, a.core, a.akm, "GET", path, nil)
}

// SaveContent 整表覆盖保存草稿。content 可为 JSON 字符串或工作簿对象（SDK 会序列化）。
func (a *ExcelAPI) SaveContent(ctx context.Context, docCode string, content any) (*ExcelVo, error) {
	s, err := stringifyJSONContent(content)
	if err != nil {
		return nil, err
	}
	body := map[string]any{"docCode": docCode, "content": s}
	return executeOpen[*ExcelVo](ctx, a.core, a.akm, "POST", "/push/api/open/excel/saveContent", body)
}

// WriteCells 从指定起始单元格起，按二维数组向右向下写入（草稿）。
// rangeStart 如 A1；sheetName 为空则写入活动表 / 第一张表。
func (a *ExcelAPI) WriteCells(ctx context.Context, docCode, rangeStart string, values [][]any, sheetName string) (*ExcelVo, error) {
	body := map[string]any{"docCode": docCode, "range": rangeStart, "values": values}
	if sheetName != "" {
		body["sheetName"] = sheetName
	}
	return executeOpen[*ExcelVo](ctx, a.core, a.akm, "POST", "/push/api/open/excel/writeCells", body)
}

// Publish 将草稿同步为分享页快照。
func (a *ExcelAPI) Publish(ctx context.Context, docCode string) (*ExcelVo, error) {
	path := appendQuery("/push/api/open/excel/publish", []queryParam{{"docCode", docCode}})
	return executeOpen[*ExcelVo](ctx, a.core, a.akm, "POST", path, nil)
}

// Rename 重命名。
func (a *ExcelAPI) Rename(ctx context.Context, docCode, title string) error {
	body := map[string]any{"docCode": docCode, "title": title}
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", "/push/api/open/excel/rename", body)
	return err
}

// Delete 删除表格。
func (a *ExcelAPI) Delete(ctx context.Context, docCode string) error {
	path := appendQuery("/push/api/open/excel/delete", []queryParam{{"docCode", docCode}})
	_, err := executeOpen[any](ctx, a.core, a.akm, "POST", path, nil)
	return err
}

// UpdateShare 更新分享设置。sharePerm：0 关闭 / 1 开启；shareLogin 为 nil 时沿用原值。
func (a *ExcelAPI) UpdateShare(ctx context.Context, docCode string, sharePerm int, shareLogin *int) (*ExcelVo, error) {
	body := map[string]any{"docCode": docCode, "sharePerm": sharePerm}
	if shareLogin != nil {
		body["shareLogin"] = *shareLogin
	}
	return executeOpen[*ExcelVo](ctx, a.core, a.akm, "POST", "/push/api/open/excel/updateShare", body)
}

func guessExcelContentType(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".xls":
		return "application/vnd.ms-excel"
	default:
		return "application/octet-stream"
	}
}

func stringifyJSONContent(content any) (string, error) {
	switch v := content.(type) {
	case nil:
		return "", nil
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case json.RawMessage:
		return string(v), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("序列化表格内容失败: %w", err)
		}
		return string(b), nil
	}
}
