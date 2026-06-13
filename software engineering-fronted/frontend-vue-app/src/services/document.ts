import request, { USE_MOCK } from '@/utils/request'
import { mockDocument } from './mock'

export interface DocumentItem {
  id: number
  title: string
  description: string
  filename?: string
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
  page?: number
  size?: number
  total_page?: number
}

// 上传文档
export async function uploadDocument(
  file: File,
  title?: string,
  description?: string,
  onProgress?: (percent: number) => void
) {
  if (USE_MOCK) {
    return mockDocument.uploadDocument(file, title, description) as Promise<any>
  }
  const formData = new FormData()
  formData.append('file', file)
  if (title) formData.append('title', title)
  if (description) formData.append('description', description)
  return request.post<DocumentItem>('/documents', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: (e) => {
      if (e.total && onProgress) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    }
  })
}

// 文档列表
export async function getDocumentList(params: DocumentListParams) {
  if (USE_MOCK) {
    return mockDocument.getDocumentList(params) as Promise<any>
  }
  return request.get<PaginatedResponse<DocumentItem>>('/documents', { params })
}

// 文档详情
export async function getDocumentDetail(id: number) {
  if (USE_MOCK) {
    return mockDocument.getDocumentDetail(id) as Promise<any>
  }
  return request.get<DocumentItem>(`/documents/${id}`)
}

// 更新文档信息
export function updateDocument(id: number, data: { title?: string; description?: string }) {
  return request.put(`/documents/${id}`, data)
}

// 删除文档
export async function deleteDocument(id: number) {
  if (USE_MOCK) {
    return mockDocument.deleteDocument(id) as Promise<any>
  }
  return request.delete(`/documents/${id}`)
}

// 获取文档内容
export async function getDocumentContent(id: number) {
  if (USE_MOCK) {
    return mockDocument.getDocumentContent(id) as Promise<any>
  }
  return request.get<{ id: number; title: string; content: string }>(`/documents/${id}/content`)
}
