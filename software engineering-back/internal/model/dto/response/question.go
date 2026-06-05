package response

type QuestionOption struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type QuestionResponse struct {
	ID                 uint             `json:"id"`
	Title              string           `json:"title"`
	Type               string           `json:"type"`
	Difficulty         string           `json:"difficulty"`
	KnowledgePointID   uint             `json:"knowledge_point_id"`
	KnowledgePointName string           `json:"knowledge_point_name,omitempty"`
	Options            []QuestionOption `json:"options"`
	Answer             string           `json:"answer,omitempty"`
	Explanation        string           `json:"explanation,omitempty"`
	CreatedAt          string           `json:"created_at"`
}
