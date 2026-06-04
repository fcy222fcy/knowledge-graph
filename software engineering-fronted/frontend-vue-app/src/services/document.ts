import request from '@/utils/request'

export interface DocumentItem {
  id: number
  title: string
  description: string
  filename: string
  file_size: number
  file_type: string
  status: string
  created_at: string
  updated_at?: string
}

export interface DocumentListParams {
  page: number
  size: number
  keyword?: string
  status?: string
}

export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  size: number
  total_page: number
}

// 上传文档
export function uploadDocument(file: File, title?: string, description?: string) {
  const formData = new FormData()
  formData.append('file', file)
  if (title) formData.append('title', title)
  if (description) formData.append('description', description)
  return request.post<DocumentItem>('/documents', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 文档列表
export function getDocumentList(params: DocumentListParams) {
  return request.get<PaginatedResponse<DocumentItem>>('/documents', { params })
}

// 文档详情
export function getDocumentDetail(id: number) {
  return request.get<DocumentItem>(`/documents/${id}`)
}

// 更新文档信息
export function updateDocument(id: number, data: { title?: string; description?: string }) {
  return request.put(`/documents/${id}`, data)
}

// 删除文档
export function deleteDocument(id: number) {
  return request.delete(`/documents/${id}`)
}

// 获取文档内容
export function getDocumentContent(id: number) {
  return request.get<{ id: number; title: string; content: string }>(`/documents/${id}/content`)
}
