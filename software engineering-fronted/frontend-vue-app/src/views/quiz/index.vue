<template>
  <div class="quiz-container" v-loading="loading">
    <h2>答题中心</h2>

    <!-- 答题模式 -->
    <template v-if="!submitted">
      <!-- 统计区域 -->
      <div class="overview-cards">
        <div class="stat-card">
          <div class="stat-icon" style="background: #eff6ff">
            <el-icon :size="20" color="#2563eb"><Document /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ questions.length }}</div>
            <div class="stat-label">总题数</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon" style="background: #f0fdf4">
            <el-icon :size="20" color="#10b981"><Check /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ answeredCount }}</div>
            <div class="stat-label">已答题数</div>
          </div>
        </div>
        <div class="stat-card stat-card--progress">
          <div class="stat-icon" style="background: #fef3c7">
            <el-icon :size="20" color="#f59e0b"><TrendCharts /></el-icon>
          </div>
          <div class="stat-content stat-content--wide">
            <div class="stat-label" style="margin-top: 0; margin-bottom: 6px">完成进度</div>
            <el-progress
              :percentage="questions.length ? Math.round((answeredCount / questions.length) * 100) : 0"
              :stroke-width="8"
              :color="'#10b981'"
            />
          </div>
        </div>
      </div>

      <div class="question-list">
        <div
          v-for="(q, index) in questions"
          :key="q.id"
          class="question-card"
          :class="'difficulty-' + q.difficulty"
        >
          <div class="question-header">
            <span class="question-index">{{ index + 1 }}.</span>
            <span class="question-title">{{ q.title }}</span>
            <el-tag :type="typeTag(q.type)" size="small" class="question-type-tag">
              {{ isSingleChoice(q.type) ? '单选' : '多选' }}
            </el-tag>
            <el-tag :type="difficultyTag(q.difficulty)" size="small">
              {{ difficultyLabel(q.difficulty) }}
            </el-tag>
          </div>

          <!-- 单选 -->
          <el-radio-group
            v-if="isSingleChoice(q.type)"
            v-model="answers[q.id]"
            class="options-group"
          >
            <el-radio
              v-for="opt in q.options"
              :key="opt.key"
              :value="opt.key"
              class="option-item"
            >
              {{ opt.key }}. {{ opt.value }}
            </el-radio>
          </el-radio-group>

          <!-- 多选 -->
          <el-checkbox-group
            v-else-if="q.options && q.options.length > 0"
            v-model="answers[q.id]"
            class="options-group"
          >
            <el-checkbox
              v-for="opt in q.options"
              :key="opt.key"
              :value="opt.key"
              class="option-item"
            >
              {{ opt.key }}. {{ opt.value }}
            </el-checkbox>
          </el-checkbox-group>
        </div>
      </div>

      <div class="submit-bar">
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
      <!-- 圆形正确率 -->
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

      <div class="question-list">
        <div
          v-for="(r, index) in results"
          :key="r.question_id"
          class="question-card result-card"
          :class="r.is_correct ? 'correct' : 'wrong'"
        >
          <div class="question-header">
            <span class="question-index">{{ index + 1 }}.</span>
            <span class="question-title">{{ r.questionTitle }}</span>
            <el-tag :type="r.is_correct ? 'success' : 'danger'" size="small">
              {{ r.is_correct ? '正确' : '错误' }}
            </el-tag>
          </div>

          <div class="answer-line">
            <span>你的答案：</span>
            <span :class="r.is_correct ? 'text-success' : 'text-danger'">{{ r.user_answer }}</span>
          </div>
          <div v-if="!r.is_correct" class="answer-line">
            <span>正确答案：</span>
            <span class="text-success">{{ r.correctAnswer }}</span>
          </div>
          <div v-if="r.explanation" class="explanation">
            <strong>解析：</strong>{{ r.explanation }}
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

// answers[questionId] -> 单选为 string，多选为 string[]
const answers = ref<Record<number, string | string[]>>({})

// 提交后的结果（扩展了题目信息方便展示）
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

/** 单选答案格式化为 "A"，多选格式化为 "A,B,C" */
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
    // 初始化答案
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

      const result = await submitQuiz({
        question_id: q.id,
        user_answer: userAnswer
      })
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
  // 重置答案
  const init: Record<number, string | string[]> = {}
  for (const q of questions.value) {
    init[q.id] = q.type === 'multiple' ? [] : ''
  }
  answers.value = init
}

onMounted(() => {
  fetchQuestions()
})
</script>

<style scoped>
/* ===== 页面标题 ===== */
.quiz-container h2 {
  font-size: 20px;
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
  gap: 16px;
  padding: 20px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
}

.stat-card--progress {
  grid-column: span 1;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 4px;
}

.stat-content--wide {
  flex: 1;
  min-width: 0;
}

/* ===== 题目卡片 ===== */
.question-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.question-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px 24px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  border-left: 4px solid transparent;
  transition: box-shadow 0.25s ease;
}

.question-card:hover {
  box-shadow: var(--shadow-md);
}

/* 难度左边框 */
.question-card.difficulty-easy {
  border-left-color: #10b981;
}

.question-card.difficulty-medium {
  border-left-color: #f59e0b;
}

.question-card.difficulty-hard {
  border-left-color: #f43f5e;
}

/* 结果卡片 */
.question-card.result-card.correct {
  border-left: 4px solid #10b981;
}

.question-card.result-card.wrong {
  border-left: 4px solid #ef4444;
}

.question-header {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 14px;
  flex-wrap: wrap;
}

.question-index {
  font-weight: 600;
  color: var(--text-primary);
  flex-shrink: 0;
}

.question-title {
  font-weight: 500;
  color: var(--text-primary);
  flex: 1;
}

.question-type-tag {
  flex-shrink: 0;
}

.options-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-left: 22px;
}

.option-item {
  margin-right: 0 !important;
  color: var(--text-secondary);
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
  line-height: 1.2;
}

.circle-label {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}

.answer-line {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 6px;
  padding-left: 22px;
}

.text-success {
  color: #10b981;
  font-weight: 500;
}

.text-danger {
  color: #ef4444;
  font-weight: 500;
}

.explanation {
  margin-top: 10px;
  padding: 12px 16px;
  background: var(--bg);
  border-radius: 6px;
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
  padding-left: 22px;
}

/* ===== 底部提交栏（固定） ===== */
.submit-bar {
  position: sticky;
  bottom: 0;
  display: flex;
  justify-content: center;
  padding: 16px 0;
  background: var(--bg);
  z-index: 10;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .overview-cards {
    grid-template-columns: 1fr;
  }

  .question-card {
    padding: 16px;
  }

  .question-header {
    gap: 6px;
  }

  .result-summary {
    padding: 24px 16px;
  }
}
</style>
