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

// 流式事件类型
export interface StreamEvent {
  type: 'chunk' | 'done' | 'session' | 'error'
  content: string
  session_id?: number
  confidence?: number
  sources?: Array<{
    document_id: number
    document_title: string
    content: string
  }>
  related?: Array<{
    id: number
    name: string
    description?: string
  }>
}

// 提问
export async function askQuestion(data: AskParams) {
  if (USE_MOCK) {
    return mockQa.askQuestion(data) as Promise<any>
  }
  return request.post<AskResponse>('/ask', data)
}

// 流式提问
export async function* askQuestionStream(data: AskParams) {
  if (USE_MOCK) {
    // Mock 模式下模拟流式输出
    const response = await mockQa.askQuestion(data) as any
    const answer = response.data?.answer || response.answer || '暂无回答'
    for (let i = 0; i < answer.length; i += 2) {
      yield {
        type: 'chunk',
        content: answer.slice(i, i + 2)
      }
      await new Promise(resolve => setTimeout(resolve, 30))
    }
    yield {
      type: 'done',
      content: '',
      confidence: response.data?.confidence || response.confidence || 0.8,
      sources: response.data?.sources || response.sources || [],
      related: response.data?.related_knowledge_points || response.related_knowledge_points || []
    }
    return
  }

  const baseURL = request.defaults?.baseURL || ''
  const token = localStorage.getItem('token') || ''

  const response = await fetch(`${baseURL}/ask/stream`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : ''
    },
    body: JSON.stringify(data)
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  const reader = response.body!.getReader()
  const decoder = new TextDecoder()
  let buffer = ''

  while (true) {
    const { done, value } = await reader.read()
    if (done) break

    buffer += decoder.decode(value, { stream: true })
    const lines = buffer.split('\n')
    buffer = lines.pop() || ''

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const dataStr = line.slice(6)
        if (dataStr === '[DONE]') {
          return
        }
        try {
          const event = JSON.parse(dataStr) as StreamEvent
          yield event
        } catch {
          // 忽略解析错误
        }
      }
    }
  }
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
