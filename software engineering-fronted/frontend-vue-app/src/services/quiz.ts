import request from '@/utils/request'

// 题目类型
export interface Question {
  id: number
  title: string
  type: 'single' | 'multiple'
  difficulty: 'easy' | 'medium' | 'hard'
  knowledge_point_id: number
  options: Array<{ key: string; value: string }>
  answer: string
  explanation: string
  created_at: string
}

// 答题提交参数
export interface QuizSubmitParams {
  question_id: number
  user_answer: string
}

// 答题结果
export interface QuizResult {
  id: number
  question_id: number
  user_answer: string
  is_correct: boolean
  created_at: string
}

// 获取题目列表
export function getQuestions(params?: { page?: number; size?: number }) {
  return request.get<{ list: Question[]; total: number }>('/questions', { params })
}

// 获取单个题目
export function getQuestion(id: number) {
  return request.get<Question>(`/questions/${id}`)
}

// 提交答题
export function submitQuiz(data: QuizSubmitParams) {
  return request.post<QuizResult>('/quizzes/submit', data)
}

// 获取答题历史
export function getQuizHistory(params?: { page?: number; size?: number }) {
  return request.get<{ list: QuizResult[]; total: number }>('/quizzes/history', { params })
}
