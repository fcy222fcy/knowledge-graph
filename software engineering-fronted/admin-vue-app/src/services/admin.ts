import request from '@/utils/request'

// ========== 仪表盘统计 ==========
// 获取总览统计
export function getAnalyticsOverview() {
  return request.get('/teacher/analytics/overview')
}

// 获取用户统计
export function getUserStats() {
  return request.get('/teacher/analytics/user-stats')
}

// ========== 学生管理 ==========
// 获取学生列表
export function getStudents(params?: { page?: number; size?: number; keyword?: string }) {
  return request.get('/teacher/students', { params })
}

// 获取学生详情
export function getStudent(id: number) {
  return request.get(`/teacher/students/${id}`)
}

// 更新学生状态
export function updateStudentStatus(id: number, status: number) {
  return request.put(`/teacher/students/${id}/status`, { status })
}

// 删除学生
export function deleteStudent(id: number) {
  return request.delete(`/teacher/students/${id}`)
}

// ========== 资料审核 ==========
// 获取文档列表（教师审核）
export function getDocuments(params?: { page?: number; size?: number; status?: string }) {
  return request.get('/teacher/documents', { params })
}

// 获取文档详情
export function getDocument(id: number) {
  return request.get(`/teacher/documents/${id}`)
}

// 审核文档（通过/拒绝）
export function reviewDocument(id: number, data: { status: 'approved' | 'rejected'; remark?: string }) {
  return request.put(`/teacher/documents/${id}/review`, data)
}

// 删除文档
export function deleteDocument(id: number) {
  return request.delete(`/teacher/documents/${id}`)
}

// ========== 知识点管理 ==========
// 获取知识点列表
export function getKnowledgePoints(params?: { page?: number; size?: number; keyword?: string }) {
  return request.get('/teacher/knowledge', { params })
}

// 获取知识点详情
export function getKnowledgePoint(id: number) {
  return request.get(`/teacher/knowledge/${id}`)
}

// 创建知识点
export function createKnowledgePoint(data: { name: string; description?: string; category?: string }) {
  return request.post('/teacher/knowledge', data)
}

// 更新知识点
export function updateKnowledgePoint(id: number, data: { name?: string; description?: string; category?: string }) {
  return request.put(`/teacher/knowledge/${id}`, data)
}

// 删除知识点
export function deleteKnowledgePoint(id: number) {
  return request.delete(`/teacher/knowledge/${id}`)
}

// 获取知识点关系
export function getKnowledgeRelations(params?: { page?: number; size?: number }) {
  return request.get('/teacher/knowledge/relations', { params })
}

// ========== 题目管理 ==========
// 获取题目列表
export function getQuestions(params?: { page?: number; size?: number; keyword?: string; type?: string }) {
  return request.get('/teacher/questions', { params })
}

// 获取题目详情
export function getQuestion(id: number) {
  return request.get(`/teacher/questions/${id}`)
}

// 创建题目
export function createQuestion(data: {
  title: string
  type: string
  options?: string[]
  answer: string
  analysis?: string
  knowledge_point_id?: number
}) {
  return request.post('/teacher/questions', data)
}

// 更新题目
export function updateQuestion(id: number, data: Record<string, unknown>) {
  return request.put(`/teacher/questions/${id}`, data)
}

// 删除题目
export function deleteQuestion(id: number) {
  return request.delete(`/teacher/questions/${id}`)
}
