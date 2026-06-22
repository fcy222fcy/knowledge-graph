package response

// DocumentResponse 文档响应
type DocumentResponse struct {
	ID             uint   `json:"id"`               // 文档ID
	Title          string `json:"title"`            // 文档标题
	Description    string `json:"description"`      // 文档描述
	Filename       string `json:"filename"`         // 原始文件名
	FileSize       int64  `json:"file_size"`        // 文件大小（字节）
	FileType       string `json:"file_type"`        // 文件类型（pdf/docx/txt）
	Status         string `json:"status"`           // 处理状态
	ContentPreview string `json:"content_preview,omitempty"` // 内容预览（可选）
	CreatedAt      string `json:"created_at"`       // 创建时间
	UpdatedAt      string `json:"updated_at"`       // 更新时间
}

// DocumentContentResponse 文档内容响应
type DocumentContentResponse struct {
	ID      uint   `json:"id"`      // 文档ID
	Title   string `json:"title"`   // 文档标题
	Content string `json:"content"` // 文档完整内容
}

// DocumentDownloadResponse 文档下载响应
type DocumentDownloadResponse struct {
	ID      uint   `json:"id"`      // 文档ID
	Title   string `json:"title"`   // 文档标题
	URL     string `json:"url"`     // 预签名下载 URL
}
