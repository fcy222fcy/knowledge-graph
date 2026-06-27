package document

import (
	"net/http"
	"path/filepath"
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// maxFileSize 上传文件大小上限 50MB
const maxFileSize = 50 << 20 // 50MB

// allowedExts 允许上传的文件类型白名单
var allowedExts = map[string]bool{
	".md": true, ".txt": true, ".pdf": true,
	".docx": true, ".pptx": true,
}

// UploadDocument 上传文档接口，支持 .md .txt .pdf .docx .pptx 格式
func UploadDocument(c *gin.Context) {
	userID := c.GetUint("user_id")
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请上传文件")
		return
	}

	// 文件大小限制
	if file.Size > maxFileSize {
		response.Error(c, http.StatusBadRequest, "文件大小超出限制（最大 50MB）")
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	filename := fixFilenameEncoding(file.Filename)
	ext := filepath.Ext(filename)

	// 文件类型白名单
	if !allowedExts[ext] {
		response.Error(c, http.StatusBadRequest, "不支持的文件类型，支持：.md .txt .pdf .docx .pptx")
		return
	}

	f, err := file.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "文件读取失败")
		return
	}
	defer f.Close()

	resp, err := service.UploadDocument(userID, title, description, filename, file.Size, ext, f)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetDocument 获取文档详情
func GetDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetDocument(uint(id))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, resp)
}

// GetDocumentContent 获取文档完整内容
func GetDocumentContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetDocumentContent(uint(id))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, resp)
}

// UpdateDocument 更新文档信息
func UpdateDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateDocument(uint(id), req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// DeleteDocument 删除文档
func DeleteDocument(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteDocument(userID, uint(id)); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// ListDocuments 获取当前用户的文档列表
func ListDocuments(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	list, total, err := service.ListUserDocuments(userID, page, size, keyword, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// fixFilenameEncoding 修复 Windows 下 multipart 表单文件名的编码问题。
// curl 在 Windows 上发送 multipart 时，文件名使用系统默认编码（GBK），
// 而 Go 的 mime/multipart 直接将 GBK 字节当作字符串返回。
// 此函数检测并将其转为 UTF-8。
func fixFilenameEncoding(name string) string {
	// 如果全是 ASCII 或已合法的 Unicode，直接返回
	hasNonASCII := false
	for _, r := range name {
		if r > unicode.MaxASCII {
			hasNonASCII = true
			break
		}
	}
	if !hasNonASCII {
		return name
	}

	// 尝试将字符串的底层字节用 GBK 解码，再转为 UTF-8 字符串
	raw := []byte(name)
	decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Bytes, _, err := transform.Bytes(decoder, raw)
	if err != nil {
		// 解码失败，说明不是 GBK，原样返回
		return name
	}
	return string(utf8Bytes)
}
