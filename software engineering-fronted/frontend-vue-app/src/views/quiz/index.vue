<template>
  <div class="quiz-container" v-loading="loading">
    <h2>答题中心</h2>

    <!-- 答题模式 -->
    <template v-if="!submitted">
      <!-- 统计区域 -->
      <div class="overview-cards">
        <div class="stat-card">
          <div class="stat-icon blue">
            <el-icon :size="20"><Document /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ questions.length }}</div>
            <div class="stat-label">总题数</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon green">
            <el-icon :size="20"><Check /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ answeredCount }}</div>
            <div class="stat-label">已答题数</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon orange">
            <el-icon :size="20"><TrendCharts /></el-icon>
          </div>
          <div class="stat-content progress-content">
            <div class="stat-label">完成进度</div>
            <el-progress
              :percentage="questions.length ? Math.round((answeredCount / questions.length) * 100) : 0"
              :stroke-width="8"
              :show-text="false"
            />
          </div>
        </div>
      </div>

      <!-- 并列题目卡片 -->
      <div class="question-grid">
        <div
          v-for="(q, index) in questions"
          :key="q.id"
          class="question-card"
          :class="'difficulty-' + q.difficulty"
        >
          <!-- 卡片头部 -->
          <div class="card-top">
            <span class="q-num">{{ index + 1 }}</span>
            <div class="q-tags">
              <el-tag :type="typeTag(q.type)" size="small" round>
                {{ isSingleChoice(q.type) ? '单选' : '多选' }}
              </el-tag>
              <el-tag :type="difficultyTag(q.difficulty)" size="small" round>
                {{ difficultyLabel(q.difficulty) }}
              </el-tag>
            </div>
          </div>

          <div class="q-title">{{ q.title }}</div>

          <!-- 选项 -->
          <div class="q-options" v-if="q.options && q.options.length > 0">
            <el-radio-group
              v-if="isSingleChoice(q.type)"
              v-model="answers[q.id]"
              class="options-list"
            >
              <el-radio
                v-for="opt in q.options"
                :key="opt.key"
                :value="opt.key"
                class="opt-item"
              >
                <span class="opt-letter">{{ opt.key }}</span>
                <span class="opt-text">{{ opt.value }}</span>
              </el-radio>
            </el-radio-group>

            <el-checkbox-group
              v-else
              v-model="answers[q.id]"
              class="options-list"
            >
              <el-checkbox
                v-for="opt in q.options"
                :key="opt.key"
                :value="opt.key"
                class="opt-item"
              >
                <span class="opt-letter">{{ opt.key }}</span>
                <span class="opt-text">{{ opt.value }}</span>
              </el-checkbox>
            </el-checkbox-group>
          </div>
        </div>
      </div>

      <!-- 底部提交 -->
      <div class="submit-bar">
        <div class="submit-info">
          已答 <strong>{{ answeredCount }}</strong> / {{ questions.length }} 题
        </div>
        <el-button
          type="primary"
          size="large"
          :disabled="answeredCount === 0"
          :loading="submitting"
          @click="handleSubmit"
        >
          提交答题
        </el-button>
      </div>
    </template>

    <!-- 结果模式 -->
    <template v-else>
      <div class="result-summary">
        <el-progress
          type="circle"
          :percentage="Math.round((correctCount / results.length) * 100)"
          :width="140"
          :stroke-width="10"
          :color="scoreColor"
        >
          <template #default="{ percentage }">
            <div class="circle-inner">
              <span class="circle-num">{{ percentage }}%</span>
              <span class="circle-label">{{ correctCount }}/{{ results.length }} 正确</span>
            </div>
          </template>
        </el-progress>
      </div>

      <div class="question-grid">
        <div
          v-for="(r, index) in results"
          :key="r.question_id"
          class="question-card result-card"
          :class="r.is_correct ? 'correct' : 'wrong'"
        >
          <div class="card-top">
            <span class="q-num" :class="r.is_correct ? 'num-correct' : 'num-wrong'">{{ index + 1 }}</span>
            <el-tag :type="r.is_correct ? 'success' : 'danger'" size="small" round>
              {{ r.is_correct ? '正确' : '错误' }}
            </el-tag>
          </div>

          <div class="q-title">{{ r.questionTitle }}</div>

          <div class="result-info">
            <div class="answer-line">
              <span class="answer-label">你的答案：</span>
              <span :class="r.is_correct ? 'text-success' : 'text-danger'">{{ r.user_answer }}</span>
            </div>
            <div v-if="!r.is_correct" class="answer-line">
              <span class="answer-label">正确答案：</span>
              <span class="text-success">{{ r.correctAnswer }}</span>
            </div>
          </div>

          <div v-if="r.explanation" class="explanation">
            {{ r.explanation }}
          </div>
        </div>
      </div>

      <div class="submit-bar">
        <el-button type="primary" size="large" @click="handleRetry">
          重新答题
        </el-button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Document, Check, TrendCharts } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getQuestions, submitQuiz } from '@/services/quiz'
import type { Question, QuizResult } from '@/services/quiz'

// --- State ---
const loading = ref(true)
const submitting = ref(false)
const submitted = ref(false)
const questions = ref<Question[]>([])
const answers = ref<Record<number, string | string[]>>({})

const results = ref<(QuizResult & {
  questionTitle: string
  correctAnswer: string
  explanation: string
})[]>([])

// --- Computed ---
const answeredCount = computed(() => {
  return Object.values(answers.value).filter(v => {
    if (Array.isArray(v)) return v.length > 0
    return !!v
  }).length
})

const correctCount = computed(() => results.value.filter(r => r.is_correct).length)

const scoreColor = computed(() => {
  const pct = results.value.length ? (correctCount.value / results.value.length) * 100 : 0
  if (pct >= 80) return '#10b981'
  if (pct >= 60) return '#f59e0b'
  return '#f43f5e'
})

// --- Helpers ---
const typeTag = (type: string) => (isSingleChoice(type) ? 'primary' : 'warning')
const isSingleChoice = (type: string) => type === 'single' || type === 'single_choice'
const isMultipleChoice = (type: string) => type === 'multiple' || type === 'multiple_choice'
const difficultyTag = (d: string) => {
  const map: Record<string, string> = { easy: 'success', medium: 'warning', hard: 'danger' }
  return (map[d] || 'info') as any
}
const difficultyLabel = (d: string) => {
  const map: Record<string, string> = { easy: '简单', medium: '中等', hard: '困难' }
  return map[d] || d
}

const formatAnswer = (q: Question): string => {
  const v = answers.value[q.id]
  if (!v) return ''
  if (Array.isArray(v)) return v.slice().sort().join(',')
  return v
}

// --- Actions ---
const fetchQuestions = async () => {
  loading.value = true
  try {
    const result = await getQuestions({ page: 1, size: 100 })
    questions.value = result.data.list || []
    const init: Record<number, string | string[]> = {}
    for (const q of questions.value) {
      init[q.id] = isMultipleChoice(q.type) ? [] : ''
    }
    answers.value = init
  } catch (error) {
    console.error('获取题目列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    const quizResults: typeof results.value = []
    for (const q of questions.value) {
      const userAnswer = formatAnswer(q)
      if (!userAnswer) continue
      const result = await submitQuiz({ question_id: q.id, user_answer: userAnswer })
      quizResults.push({
        ...result.data,
        questionTitle: q.title,
        correctAnswer: q.answer,
        explanation: q.explanation
      })
    }
    results.value = quizResults
    submitted.value = true
  } catch (error) {
    console.error('提交答题失败:', error)
  } finally {
    submitting.value = false
  }
}

const handleRetry = () => {
  submitted.value = false
  results.value = []
  const init: Record<number, string | string[]> = {}
  for (const q of questions.value) {
    init[q.id] = isMultipleChoice(q.type) ? [] : ''
  }
  answers.value = init
}

onMounted(() => {
  fetchQuestions()
})
</script>

<style scoped>
.quiz-container {
  padding: 24px 32px;
  max-width: 1200px;
  margin: 0 auto;
}

.quiz-container h2 {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 20px;
}

/* ===== 统计卡片 ===== */
.overview-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: #fff;
}

.stat-icon.blue { background: linear-gradient(135deg, #3b82f6, #2563eb); }
.stat-icon.green { background: linear-gradient(135deg, #34d399, #10b981); }
.stat-icon.orange { background: linear-gradient(135deg, #fbbf24, #f59e0b); }

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 2px;
}

.progress-content {
  flex: 1;
  min-width: 0;
}

/* ===== 并列题目网格 ===== */
.question-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.question-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  border-top: 3px solid transparent;
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: box-shadow 0.25s ease, transform 0.2s ease;
}

.question-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

/* 难度顶部边框 */
.question-card.difficulty-easy { border-top-color: #10b981; }
.question-card.difficulty-medium { border-top-color: #f59e0b; }
.question-card.difficulty-hard { border-top-color: #f43f5e; }

/* 结果卡片 */
.question-card.result-card.correct { border-top-color: #10b981; background: #f0fdf4; }
.question-card.result-card.wrong { border-top-color: #ef4444; background: #fef2f2; }

/* 卡片头部 */
.card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.q-num {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  background: var(--bg);
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.q-num.num-correct { background: #d1fae5; color: #059669; }
.q-num.num-wrong { background: #fee2e2; color: #dc2626; }

.q-tags {
  display: flex;
  gap: 6px;
}

.q-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  line-height: 1.6;
}

/* ===== 选项 ===== */
.q-options {
  margin-top: 4px;
}

.options-list {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.opt-item {
  margin-right: 0 !important;
  padding: 8px 10px;
  border-radius: 8px;
  border: 1px solid var(--border-light);
  transition: all 0.2s ease;
  height: auto !important;
  justify-content: flex-start;
}

.opt-item:hover {
  background: var(--bg-hover);
  border-color: var(--primary);
}

.opt-letter {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 6px;
  background: var(--bg);
  color: var(--text-secondary);
  font-size: 12px;
  font-weight: 600;
  margin-right: 8px;
  flex-shrink: 0;
}

.opt-text {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.4;
}

:deep(.el-radio__input),
:deep(.el-checkbox__input) {
  margin-top: 1px;
}

:deep(.el-radio__label),
:deep(.el-checkbox__label) {
  padding-left: 0;
  display: flex;
  align-items: center;
}

/* ===== 结果区域 ===== */
.result-summary {
  display: flex;
  justify-content: center;
  margin-bottom: 28px;
  padding: 32px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
}

.circle-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.circle-num {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
}

.circle-label {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}

.result-info {
  padding-left: 4px;
}

.answer-line {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.answer-label {
  color: var(--text-muted);
}

.text-success { color: #10b981; font-weight: 600; }
.text-danger { color: #ef4444; font-weight: 600; }

.explanation {
  padding: 10px 12px;
  background: var(--bg);
  border-radius: 8px;
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.6;
}

/* ===== 底部提交栏 ===== */
.submit-bar {
  position: sticky;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
  padding: 16px 24px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-light);
  z-index: 10;
}

.submit-info {
  font-size: 14px;
  color: var(--text-muted);
}

.submit-info strong {
  color: var(--primary);
  font-size: 16px;
}

/* ===== 响应式 ===== */
@media (max-width: 900px) {
  .question-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .quiz-container {
    padding: 16px;
  }

  .overview-cards {
    grid-template-columns: 1fr;
  }

  .question-card {
    padding: 16px;
  }
}
</style>
