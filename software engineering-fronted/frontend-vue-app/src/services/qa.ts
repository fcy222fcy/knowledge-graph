import request, { USE_MOCK } from '@/utils/request'
import { mockQa } from './mock'

export interface AskParams {
  question: string
  conversation_id?: number
}

export interface AskResponse {
  conversation_id: number
  question_id: number
  answer: string
  confidence?: number
  sources: Array<{
    document_id: number
    document_title: string
    content: string
  }>
  related_knowledge_points: Array<{
    id: number
    name: string
    description?: string
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
export async function askQuestion(data: AskParams) {
  if (USE_MOCK) {
    return mockQa.askQuestion(data) as Promise<any>
  }
  return request.post<AskResponse>('/ask', data)
}

// 新建问答会话
export async function createSession(title?: string) {
  if (USE_MOCK) {
    return mockQa.createSession(title) as Promise<any>
  }
  return request.post('/ask/sessions', { title })
}

// 会话列表
export async function getSessions(params?: { page?: number; size?: number }) {
  if (USE_MOCK) {
    return mockQa.getSessions(params) as Promise<any>
  }
  return request.get('/ask/sessions', { params })
}

// 会话消息列表
export async function getSessionMessages(sessionId: number, params?: { page?: number; size?: number }) {
  if (USE_MOCK) {
    return mockQa.getSessionMessages(sessionId, params) as Promise<any>
  }
  return request.get(`/ask/sessions/${sessionId}/messages`, { params })
}

// 问答历史
export async function getAskHistory(params?: { page?: number; size?: number; conversation_id?: number }) {
  if (USE_MOCK) {
    return mockQa.getAskHistory(params) as Promise<any>
  }
  return request.get('/ask/history', { params })
}
