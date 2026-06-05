import request, { USE_MOCK } from '@/utils/request'
import { mockStats } from './mock'

export interface OverviewStats {
  today_learning_hours: number
  today_questions_asked: number
  total_learning_hours: number
  total_questions_asked: number
  total_quizzes_taken: number
  average_correct_rate: number
  knowledge_points_mastered: number
  knowledge_points_total: number
  mastery_rate?: number
}

export interface HotKnowledgePoint {
  knowledge_point_id: number
  knowledge_point_name: string
  heat: number
  question_count?: number
  quiz_count?: number
}

export interface KnowledgeMastery {
  knowledge_point_id: number
  knowledge_point_name: string
  total_questions: number
  correct_answers: number
  mastery_rate: number
  level: string
}

export interface WeakPoint {
  knowledge_point_id: number
  knowledge_point_name: string
  correct_rate: number
  suggested_questions: Array<{
    id: number
    title: string
  }>
}

export interface TrendData {
  daily_stats: Array<{
    date: string
    questions_asked: number
    learning_hours: number
    correct_rate: number
  }>
  weekly_trend?: Array<{
    week: string
    avg_correct_rate: number
    total_learning_hours: number
    total_questions_asked: number
  }>
}

// 总览统计
export async function getOverview() {
  if (USE_MOCK) {
    return mockStats.getOverview() as Promise<any>
  }
  return request.get<OverviewStats>('/analytics/overview')
}

// 热门知识点
export async function getHotKnowledgePoints(limit?: number) {
  if (USE_MOCK) {
    return mockStats.getHotKnowledgePoints(limit) as Promise<any>
  }
  return request.get<HotKnowledgePoint[]>('/analytics/hot-knowledge-points', { params: { limit } })
}

// 知识点掌握度
export async function getKnowledgeMastery() {
  if (USE_MOCK) {
    return mockStats.getKnowledgeMastery() as Promise<any>
  }
  return request.get<KnowledgeMastery[]>('/analytics/knowledge-mastery')
}

// 薄弱知识点
export async function getWeakPoints(limit?: number) {
  if (USE_MOCK) {
    return mockStats.getWeakPoints(limit) as Promise<any>
  }
  return request.get<WeakPoint[]>('/analytics/weak-points', { params: { limit } })
}

// 趋势数据
export async function getTrends(days?: number) {
  if (USE_MOCK) {
    return mockStats.getTrends(days) as Promise<any>
  }
  return request.get<TrendData>('/analytics/trends', { params: { days } })
}
