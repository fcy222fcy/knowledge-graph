package request

type UpdateDocumentRequest struct {
	Title       string `json:"title" binding:"max=200"`
	Description string `json:"description" binding:"max=500"`
}
