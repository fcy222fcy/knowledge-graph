import request from '@/utils/request'

export interface OverviewStats {
  today_learning_hours: number
  today_questions_asked: number
  total_learning_hours: number
  total_questions_asked: number
  total_quizzes_taken: number
  average_correct_rate: number
  knowledge_points_mastered: number
  knowledge_points_total: number
  mastery_rate: number
}

export interface HotKnowledgePoint {
  knowledge_point_id: number
  knowledge_point_name: string
  heat: number
  question_count: number
  quiz_count: number
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
  weekly_trend: Array<{
    week: string
    avg_correct_rate: number
    total_learning_hours: number
    total_questions_asked: number
  }>
}

// 总览统计
export function getOverview() {
  return request.get<OverviewStats>('/analytics/overview')
}

// 热门知识点
export function getHotKnowledgePoints(limit?: number) {
  return request.get<HotKnowledgePoint[]>('/analytics/hot-knowledge-points', { params: { limit } })
}

// 知识点掌握度
export function getKnowledgeMastery() {
  return request.get<KnowledgeMastery[]>('/analytics/knowledge-mastery')
}

// 薄弱知识点
export function getWeakPoints(limit?: number) {
  return request.get<WeakPoint[]>('/analytics/weak-points', { params: { limit } })
}

// 趋势数据
export function getTrends(days?: number) {
  return request.get<TrendData>('/analytics/trends', { params: { days } })
}
