package errors

// 错误码定义
const (
	// 通用错误码
	CodeSuccess       = 200
	CodeBadRequest    = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeInternalError = 500

	// 用户相关错误码
	CodeUserNotFound      = 1001
	CodeUserAlreadyExists = 1002
	CodeInvalidPassword   = 1003
	CodeInvalidToken      = 1004

	// 资料相关错误码
	CodeDocumentNotFound      = 2001
	CodeDocumentUploadFailed  = 2002
	CodeDocumentInvalidFormat = 2003

	// 知识点相关错误码
	CodeKnowledgePointNotFound = 3001
	CodeRelationNotFound       = 3002

	// 题目相关错误码
	CodeQuestionNotFound = 4001

	// 答题相关错误码
	CodeQuizNotFound = 5001

	// 问答相关错误码
	CodeAskSessionNotFound = 6001
	CodeAskFailed          = 6002
)
