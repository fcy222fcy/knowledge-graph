import type { DocumentItem, DocumentListParams, PaginatedResponse } from '@/services/document'

export type { DocumentItem, DocumentListParams, PaginatedResponse }

export type DocumentStatus = 'pending' | 'processing' | 'completed' | 'failed'

export const DOCUMENT_STATUS_MAP: Record<DocumentStatus, { label: string; color: string }> = {
  pending: { label: '待处理', color: 'info' },
  processing: { label: '处理中', color: 'warning' },
  completed: { label: '已完成', color: 'success' },
  failed: { label: '失败', color: 'danger' },
}
