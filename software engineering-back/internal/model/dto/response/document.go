package response

type DocumentResponse struct {
	ID             uint   `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Filename       string `json:"filename"`
	FileSize       int64  `json:"file_size"`
	FileType       string `json:"file_type"`
	Status         string `json:"status"`
	ContentPreview string `json:"content_preview,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type DocumentContentResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
