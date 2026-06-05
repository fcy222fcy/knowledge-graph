package request

type BuildGraphRequest struct {
	DocumentIDs []uint `json:"document_ids" binding:"required,min=1"`
}
