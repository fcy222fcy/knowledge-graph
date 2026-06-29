import request from '@/utils/request'

// ========== 仪表盘统计 ==========
// 获取总览统计
export function getAnalyticsOverview() {
  return request.get('/admin/analytics/overview')
}

// 获取用户统计
export function getUserStats() {
  return request.get('/admin/analytics/users')
}

// ========== 学生管理 ==========
// 获取学生列表
export function getStudents(params?: { page?: number; size?: number; keyword?: string }) {
  return request.get('/admin/users', { params })
}

// 获取学生详情
export function getStudent(id: number) {
  return request.get(`/admin/users/${id}`)
}

// 更新学生信息
export function updateStudent(id: number, data: { nickname?: string; email?: string; status?: number }) {
  return request.put(`/admin/users/${id}`, data)
}

// 更新学生状态
export function updateStudentStatus(id: number, status: number) {
  return request.put(`/admin/users/${id}/status`, { status })
}

// 删除学生
export function deleteStudent(id: number) {
  return request.delete(`/admin/users/${id}`)
}

// ========== 资料审核 ==========
// 获取文档列表（教师审核）
export function getDocuments(params?: { page?: number; size?: number; status?: string }) {
  return request.get('/admin/documents', { params })
}

// 获取文档详情
export function getDocument(id: number) {
  return request.get(`/admin/documents/${id}`)
}

// 获取文档内容（用于预览）
export function getDocumentContent(id: number) {
  return request.get(`/documents/${id}/content`)
}

// 审核文档（通过/拒绝）
export function reviewDocument(id: number, data: { status: 'approved' | 'rejected'; comment?: string }) {
  return request.put(`/admin/documents/${id}/review`, data)
}

// 删除文档
export function deleteDocument(id: number) {
  return request.delete(`/admin/documents/${id}`)
}

// ========== 知识点管理 ==========
// 获取知识点列表
export function getKnowledgePoints(params?: { page?: number; size?: number; keyword?: string }) {
  return request.get('/admin/knowledge/points', { params })
}

// 获取知识点详情
export function getKnowledgePoint(id: number) {
  return request.get(`/admin/knowledge/points/${id}`)
}

// 创建知识点
export function createKnowledgePoint(data: { name: string; description?: string; category?: string }) {
  return request.post('/admin/knowledge/points', data)
}

// 更新知识点
export function updateKnowledgePoint(id: number, data: { name?: string; description?: string; category?: string }) {
  return request.put(`/admin/knowledge/points/${id}`, data)
}

// 删除知识点
export function deleteKnowledgePoint(id: number) {
  return request.delete(`/admin/knowledge/points/${id}`)
}

// 获取知识点关系
export function getKnowledgeRelations(params?: { page?: number; size?: number }) {
  return request.get('/admin/knowledge/relations', { params })
}

// ========== 题目管理 ==========
// 获取题目列表
export function getQuestions(params?: { page?: number; size?: number; keyword?: string; type?: string }) {
  return request.get('/admin/questions', { params })
}

// 获取题目详情
export function getQuestion(id: number) {
  return request.get(`/admin/questions/${id}`)
}

// 创建题目
export function createQuestion(data: {
  title: string
  type: string
  difficulty?: string
  options?: Array<{ key: string; value: string }>
  answer: string
  explanation?: string
  knowledge_point_id?: number
}) {
  return request.post('/admin/questions', data)
}

// 更新题目
export function updateQuestion(id: number, data: Record<string, unknown>) {
  return request.put(`/admin/questions/${id}`, data)
}

// 删除题目
export function deleteQuestion(id: number) {
  return request.delete(`/admin/questions/${id}`)
}
