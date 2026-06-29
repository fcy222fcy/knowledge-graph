import request from '@/utils/request'

// ========== 作业管理 ==========

// 获取作业列表
export function getAssignments(params?: { page?: number; size?: number }) {
  return request.get('/admin/assignments', { params })
}

// 获取作业详情（含答案）
export function getAssignment(id: number) {
  return request.get(`/admin/assignments/${id}`)
}

// 创建作业
export function createAssignment(data: {
  title: string
  description?: string
  chapter?: string
  deadline: string
  questions: Array<{
    title: string
    type: string
    options: Array<{ key: string; value: string }>
    answer: string
    explanation?: string
    score: number
    sort_order: number
  }>
}) {
  return request.post('/admin/assignments', data)
}

// 更新作业
export function updateAssignment(id: number, data: {
  title?: string
  description?: string
  chapter?: string
  deadline?: string
  questions?: Array<{
    id?: number
    title: string
    type: string
    options: Array<{ key: string; value: string }>
    answer: string
    explanation?: string
    score: number
    sort_order: number
  }>
}) {
  return request.put(`/admin/assignments/${id}`, data)
}

// 删除作业
export function deleteAssignment(id: number) {
  return request.delete(`/admin/assignments/${id}`)
}

// 发布作业
export function publishAssignment(id: number) {
  return request.put(`/admin/assignments/${id}/publish`)
}

// 关闭作业
export function closeAssignment(id: number) {
  return request.put(`/admin/assignments/${id}/close`)
}

// 获取提交列表
export function getSubmissions(assignmentId: number, params?: { page?: number; size?: number }) {
  return request.get(`/admin/assignments/${assignmentId}/submissions`, { params })
}
