import request from '@/utils/request'

export interface AskParams {
  question: string
  conversation_id?: number
}

export interface AskResponse {
  conversation_id: number
  question_id: number
  answer: string
  confidence: number
  sources: Array<{
    document_id: number
    document_title: string
    content: string
  }>
  related_knowledge_points: Array<{
    id: number
    name: string
    description: string
  }>
  created_at: string
}

export interface ConversationItem {
  conversation_id: number
  title: string
  last_question: string
  message_count: number
  updated_at: string
}

export interface MessageItem {
  message_id: number
  role: string
  content: string
  created_at: string
}

// 提问
export function askQuestion(data: AskParams) {
  return request.post<AskResponse>('/ask', data)
}

// 新建问答会话
export function createSession(title?: string) {
  return request.post('/ask/sessions', { title })
}

// 会话列表
export function getSessions(params?: { page?: number; size?: number }) {
  return request.get('/ask/sessions', { params })
}

// 会话消息列表
export function getSessionMessages(sessionId: number, params?: { page?: number; size?: number }) {
  return request.get(`/ask/sessions/${sessionId}/messages`, { params })
}

// 问答历史
export function getAskHistory(params?: { page?: number; size?: number; conversation_id?: number }) {
  return request.get('/ask/history', { params })
}
