import request from '@/utils/request'

// 用户管理
export interface UserItem {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  role: 'admin' | 'teacher' | 'student'
  status: number
  created_at: string
  updated_at: string
}

export function getUsers(params?: { page?: number; size?: number; keyword?: string; role?: string }) {
  return request.get<{ list: UserItem[]; total: number }>('/admin/users', { params })
}

export function getUser(id: number) {
  return request.get<UserItem>(`/admin/users/${id}`)
}

export function updateUser(id: number, data: { nickname?: string; email?: string }) {
  return request.put(`/admin/users/${id}`, data)
}

export function deleteUser(id: number) {
  return request.delete(`/admin/users/${id}`)
}

export function updateUserStatus(id: number, status: number) {
  return request.put(`/admin/users/${id}/status`, { status })
}

export function updateUserRole(id: number, role: string) {
  return request.put(`/admin/users/${id}/role`, { role })
}

// 题目管理
export interface QuestionItem {
  id: number
  title: string
  type: string
  difficulty: string
  knowledge_point_id: number
  options: Array<{ key: string; value: string }>
  answer: string
  explanation: string
  created_at: string
}

export function getQuestions(params?: { page?: number; size?: number }) {
  return request.get<{ list: QuestionItem[]; total: number }>('/admin/questions', { params })
}

export function createQuestion(data: Omit<QuestionItem, 'id' | 'created_at'>) {
  return request.post('/admin/questions', data)
}

export function updateQuestion(id: number, data: Partial<QuestionItem>) {
  return request.put(`/admin/questions/${id}`, data)
}

export function deleteQuestion(id: number) {
  return request.delete(`/admin/questions/${id}`)
}

// 资料管理
export interface DocumentItem {
  id: number
  title: string
  description: string
  filename: string
  file_size: number
  file_type: string
  status: 'pending' | 'approved' | 'rejected'
  review_comment?: string
  user_id?: number
  created_at: string
  updated_at: string
}

export function getDocuments(params?: { page?: number; size?: number; keyword?: string; status?: string }) {
  return request.get<{ list: DocumentItem[]; total: number }>('/admin/documents', { params })
}

export function getDocument(id: number) {
  return request.get<DocumentItem>(`/admin/documents/${id}`)
}

export function deleteDocument(id: number) {
  return request.delete(`/admin/documents/${id}`)
}

export function reviewDocument(id: number, data: { status: 'approved' | 'rejected'; comment?: string }) {
  return request.put(`/admin/documents/${id}/review`, data)
}

// 知识点管理
export interface KnowledgePointItem {
  id: number
  name: string
  description: string
  document_id: number
  category: string
  created_at: string
}

export interface KnowledgeRelationItem {
  id: number
  source_id: number
  source_name: string
  target_id: number
  target_name: string
  relation_type: string
  description: string
  created_at: string
}

export function getKnowledgePoints(params?: { page?: number; size?: number }) {
  return request.get<{ list: KnowledgePointItem[]; total: number }>('/admin/knowledge/points', { params })
}

export function deleteKnowledgePoint(id: number) {
  return request.delete(`/admin/knowledge/points/${id}`)
}

export function getKnowledgeRelations(params?: { page?: number; size?: number }) {
  return request.get<{ list: KnowledgeRelationItem[]; total: number }>('/admin/knowledge/relations', { params })
}

export function deleteKnowledgeRelation(id: number) {
  return request.delete(`/admin/knowledge/relations/${id}`)
}

// 系统统计
export interface AnalyticsOverview {
  user_count: number
  document_count: number
  knowledge_count: number
  question_count: number
  session_count: number
  quiz_count: number
}

export interface UserStats {
  admin: number
  teacher: number
  student: number
}

export function getAnalyticsOverview() {
  return request.get<AnalyticsOverview>('/admin/analytics/overview')
}

export function getUserStats() {
  return request.get<UserStats>('/admin/analytics/users')
}

// 系统配置
export interface SystemConfig {
  ollama_url: string
  ollama_model: string
  minio_endpoint: string
  minio_bucket: string
  server_port: string
}

export function getSystemConfig() {
  return request.get<SystemConfig>('/admin/system/config')
}
