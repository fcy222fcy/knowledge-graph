package document

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	pkgResponse "software_engineering/pkg/response"
)

const maxFileSizeIface = 50 << 20 // 50MB

var allowedExtsIface = map[string]bool{
	".md": true, ".txt": true, ".pdf": true,
	".docx": true, ".pptx": true,
}

// DocumentService 定义文档服务接口
type DocumentService interface {
	UploadDocument(userID uint, title, description, filename string, fileSize int64, fileType string, reader io.Reader) (*response.DocumentResponse, error) // 上传文档
	GetDocument(id uint) (*response.DocumentResponse, error)                                                                                                    // 获取文档详情
	GetDocumentContent(id uint) (*response.DocumentContentResponse, error)                                                                                      // 获取文档内容
	UpdateDocument(id uint, req request.UpdateDocumentRequest) error                                                                                             // 更新文档
	DeleteDocument(userID uint, id uint) error                                                                                                                  // 删除文档
	ListUserDocuments(userID uint, page, size int, keyword, status string) ([]response.DocumentResponse, int64, error)                                         // 列表文档
}

// DocumentController 文档控制器
type DocumentController struct {
	documentService DocumentService // 文档服务
}

// NewDocumentController 创建文档控制器实例
func NewDocumentController(documentService DocumentService) *DocumentController {
	return &DocumentController{documentService: documentService}
}

// UploadDocument 上传文档接口
func (ctrl *DocumentController) UploadDocument(c *gin.Context) {
	userID := c.GetUint("user_id")
	file, err := c.FormFile("file")
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "请上传文件")
		return
	}

	if file.Size > maxFileSizeIface {
		pkgResponse.Error(c, http.StatusBadRequest, "文件大小超出限制（最大 50MB）")
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	filename := fixFilenameEncoding(file.Filename)
	ext := filepath.Ext(filename)

	if !allowedExtsIface[ext] {
		pkgResponse.Error(c, http.StatusBadRequest, "不支持的文件类型，支持：.md .txt .pdf .docx .pptx")
		return
	}

	f, err := file.Open()
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, "文件读取失败")
		return
	}
	defer f.Close()

	resp, err := ctrl.documentService.UploadDocument(userID, title, description, filename, file.Size, ext, f)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// GetDocument 获取文档详情
func (ctrl *DocumentController) GetDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := ctrl.documentService.GetDocument(uint(id))
	if err != nil {
		pkgResponse.Error(c, http.StatusNotFound, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// GetDocumentContent 获取文档完整内容
func (ctrl *DocumentController) GetDocumentContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := ctrl.documentService.GetDocumentContent(uint(id))
	if err != nil {
		pkgResponse.Error(c, http.StatusNotFound, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// UpdateDocument 更新文档信息
func (ctrl *DocumentController) UpdateDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := ctrl.documentService.UpdateDocument(uint(id), req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, nil)
}

// DeleteDocument 删除文档
func (ctrl *DocumentController) DeleteDocument(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.documentService.DeleteDocument(userID, uint(id)); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, nil)
}

// ListDocuments 获取文档列表
func (ctrl *DocumentController) ListDocuments(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	list, total, err := ctrl.documentService.ListUserDocuments(userID, page, size, keyword, status)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// fixFilenameEncodingIface 修复 Windows 下 multipart 表单文件名的编码问题。
func fixFilenameEncodingIface(name string) string {
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

	raw := []byte(name)
	decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Bytes, _, err := transform.Bytes(decoder, raw)
	if err != nil {
		return name
	}
	return string(utf8Bytes)
}
