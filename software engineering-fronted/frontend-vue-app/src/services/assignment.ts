import request from '@/utils/request'

// 获取作业列表
export function getAssignments(params?: { page?: number; size?: number }) {
  return request.get('/assignments', { params })
}

// 获取作业详情（不含答案）
export function getAssignment(id: number) {
  return request.get(`/assignments/${id}`)
}

// 提交作业
export function submitAssignment(assignmentId: number, data: {
  answers: Array<{ question_id: number; answer: string }>
}) {
  return request.post(`/assignments/${assignmentId}/submit`, data)
}

// 查看成绩
export function getAssignmentResult(assignmentId: number) {
  return request.get(`/assignments/${assignmentId}/result`)
}
