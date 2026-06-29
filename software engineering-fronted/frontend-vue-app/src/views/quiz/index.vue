<template>
  <div class="quiz-container" v-loading="loading">
    <h2>答题中心</h2>

    <!-- 作业列表模式 -->
    <template v-if="mode === 'list'">
      <div v-if="assignments.length === 0" class="empty-state">
        <el-empty description="暂无作业" />
      </div>
      <div v-else class="assignment-grid">
        <div v-for="a in assignments" :key="a.id" class="assignment-card" :class="{ submitted: a.is_submitted }">
          <div class="a-header">
            <h3>{{ a.title }}</h3>
            <el-tag v-if="a.is_submitted" type="success" size="small">已完成</el-tag>
            <el-tag v-else-if="a.status === 'published'" type="warning" size="small">待完成</el-tag>
            <el-tag v-else type="info" size="small">{{ a.status }}</el-tag>
          </div>
          <div class="a-info">
            <span>📋 {{ a.question_num }} 题</span>
            <span>⭐ {{ a.total_score }} 分</span>
            <span>⏰ {{ a.deadline }}</span>
          </div>
          <div v-if="a.is_submitted && a.score !== undefined" class="a-score">
            得分：<strong>{{ a.score }}</strong> / {{ a.total_score }}
          </div>
          <div v-if="a.chapter" class="a-chapter">📚 {{ a.chapter }}</div>
          <div class="a-actions">
            <el-button v-if="a.is_submitted" size="small" @click="viewResult(a)">查看作业</el-button>
            <el-button v-else type="primary" size="small" @click="enterAssignment(a)">去答题</el-button>
          </div>
        </div>
      </div>
    </template>

    <!-- 答题模式 -->
    <template v-else-if="mode === 'answer'">
      <div class="answer-header">
        <el-button text @click="backToList">← 返回作业列表</el-button>
        <h3>{{ currentAssignment?.title }}</h3>
        <span class="answer-progress">已答 {{ answeredCount }} / {{ currentQuestions.length }} 题</span>
      </div>

      <div class="question-grid">
        <div
          v-for="(q, index) in currentQuestions"
          :key="q.id"
          class="question-card"
        >
          <div class="card-top">
            <span class="q-num">{{ index + 1 }}</span>
            <div class="q-tags">
              <el-tag :type="q.type === 'single' ? 'primary' : q.type === 'multiple' ? 'warning' : 'info'" size="small" round>
                {{ q.type === 'single' ? '单选' : q.type === 'multiple' ? '多选' : '判断' }}
              </el-tag>
              <el-tag size="small" round>{{ q.score }}分</el-tag>
            </div>
          </div>

          <div class="q-title">{{ q.title }}</div>

          <div class="q-options" v-if="q.options && q.options.length > 0">
            <!-- 单选 -->
            <el-radio-group
              v-if="q.type === 'single'"
              :model-value="(answers[q.id] as string)"
              @update:model-value="(val: string | number | boolean | undefined) => { answers[q.id] = String(val ?? '') }"
              class="options-list"
            >
              <el-radio v-for="opt in q.options" :key="opt.key" :value="opt.key" class="opt-item">
                <span class="opt-letter">{{ opt.key }}</span>
                <span class="opt-text">{{ opt.value }}</span>
              </el-radio>
            </el-radio-group>

            <!-- 多选 -->
            <el-checkbox-group
              v-else-if="q.type === 'multiple'"
              :model-value="(answers[q.id] as string[])"
              @update:model-value="(val: (string | number)[]) => { answers[q.id] = val.map(String) }"
              class="options-list"
            >
              <el-checkbox v-for="opt in q.options" :key="opt.key" :value="opt.key" class="opt-item">
                <span class="opt-letter">{{ opt.key }}</span>
                <span class="opt-text">{{ opt.value }}</span>
              </el-checkbox>
            </el-checkbox-group>

            <!-- 判断 -->
            <el-radio-group
              v-else
              :model-value="(answers[q.id] as string)"
              @update:model-value="(val: string | number | boolean | undefined) => { answers[q.id] = String(val ?? '') }"
              class="options-list"
            >
              <el-radio v-for="opt in q.options" :key="opt.key" :value="opt.key" class="opt-item">
                <span class="opt-letter">{{ opt.key }}</span>
                <span class="opt-text">{{ opt.value }}</span>
              </el-radio>
            </el-radio-group>
          </div>
        </div>
      </div>

      <div class="submit-bar">
        <div class="submit-info">
          已答 <strong>{{ answeredCount }}</strong> / {{ currentQuestions.length }} 题
        </div>
        <el-button type="primary" size="large" :disabled="answeredCount === 0" :loading="submitting" @click="handleSubmit">
          提交作业
        </el-button>
      </div>
    </template>

    <!-- 结果模式 -->
    <template v-else-if="mode === 'result'">
      <div class="result-summary">
        <div class="result-score">
          <div class="score-num">{{ submitResult?.score }}</div>
          <div class="score-total">/ {{ submitResult?.total_score }} 分</div>
        </div>
        <div class="result-status">
          <el-tag type="success" size="large">已提交</el-tag>
        </div>
      </div>

      <!-- 答题回顾 -->
      <div class="review-grid" v-if="submitResult?.questions?.length">
        <div
          v-for="(q, index) in submitResult.questions"
          :key="q.id"
          class="question-card"
          :class="{ correct: q.is_correct, wrong: !q.is_correct }"
        >
          <div class="card-top">
            <span class="q-num">{{ index + 1 }}</span>
            <div class="q-tags">
              <el-tag :type="q.is_correct ? 'success' : 'danger'" size="small" round>
                {{ q.is_correct ? '✓ 正确' : '✗ 错误' }}
              </el-tag>
              <el-tag :type="q.type === 'single' ? 'primary' : q.type === 'multiple' ? 'warning' : 'info'" size="small" round>
                {{ q.type === 'single' ? '单选' : q.type === 'multiple' ? '多选' : '判断' }}
              </el-tag>
              <el-tag size="small" round>{{ q.score }}分</el-tag>
            </div>
          </div>

          <div class="q-title">{{ q.title }}</div>

          <div class="q-options" v-if="q.options && q.options.length > 0">
            <div v-for="opt in q.options" :key="opt.key" class="q-option">
              <span
                class="opt-letter"
                :class="{
                  'opt-correct': q.answer.includes(opt.key),
                  'opt-wrong': q.my_answer.includes(opt.key) && !q.is_correct,
                  'opt-selected': q.my_answer.includes(opt.key)
                }"
              >{{ opt.key }}</span>
              <span class="opt-text">{{ opt.value }}</span>
              <span v-if="q.answer.includes(opt.key)" class="opt-tag correct-tag">正确答案</span>
              <span v-else-if="q.my_answer.includes(opt.key)" class="opt-tag wrong-tag">你的答案</span>
            </div>
          </div>

          <div class="q-meta">
            <span>你的答案：<strong :class="q.is_correct ? 'text-success' : 'text-danger'">{{ q.my_answer || '未作答' }}</strong></span>
            <span v-if="!q.is_correct">正确答案：<strong class="text-success">{{ q.answer }}</strong></span>
          </div>

          <div v-if="q.explanation" class="q-explanation">解析：{{ q.explanation }}</div>
        </div>
      </div>

      <div class="submit-bar">
        <el-button type="primary" size="large" @click="backToList">返回作业列表</el-button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAssignments, getAssignment, submitAssignment, getAssignmentResult } from '@/services/assignment'

const loading = ref(true)
const submitting = ref(false)
const mode = ref<'list' | 'answer' | 'result'>('list')

const assignments = ref<any[]>([])
const currentAssignment = ref<any>(null)
const currentQuestions = ref<any[]>([])
const answers = ref<Record<number, string | string[]>>({})
const submitResult = ref<any>(null)

const answeredCount = computed(() => {
  return Object.values(answers.value).filter(v => {
    if (Array.isArray(v)) return v.length > 0
    return !!v
  }).length
})

const fetchAssignments = async () => {
  loading.value = true
  try {
    const res = await getAssignments({ page: 1, size: 50 })
    assignments.value = res.data?.list || []
  } finally {
    loading.value = false
  }
}

const enterAssignment = async (a: any) => {
  loading.value = true
  try {
    const res = await getAssignment(a.id)
    currentAssignment.value = a
    currentQuestions.value = res.data?.questions || []
    const init: Record<number, string | string[]> = {}
    for (const q of currentQuestions.value) {
      init[q.id] = q.type === 'multiple' ? [] : ''
    }
    answers.value = init
    mode.value = 'answer'
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '获取作业详情失败')
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    const answerList = currentQuestions.value
      .filter(q => {
        const v = answers.value[q.id]
        if (Array.isArray(v)) return v.length > 0
        return !!v
      })
      .map(q => {
        const v = answers.value[q.id]
        const answer = Array.isArray(v) ? v.slice().sort().join(',') : v
        return { question_id: q.id, answer }
      })

    if (answerList.length === 0) {
      ElMessage.warning('请至少回答一道题')
      return
    }

    const res = await submitAssignment(currentAssignment.value.id, { answers: answerList })
    submitResult.value = res.data
    mode.value = 'result'
    ElMessage.success('作业提交成功！')
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

const backToList = () => {
  mode.value = 'list'
  currentAssignment.value = null
  currentQuestions.value = []
  answers.value = {}
  submitResult.value = null
  fetchAssignments()
}

const viewResult = async (a: any) => {
  loading.value = true
  try {
    const res = await getAssignmentResult(a.id)
    submitResult.value = res.data || res
    currentAssignment.value = a
    mode.value = 'result'
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '未提交该作业')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAssignments()
})
</script>

<style scoped>
.quiz-container { padding: 24px 32px; max-width: 1200px; margin: 0 auto; }
.quiz-container h2 { font-size: 22px; font-weight: 600; color: var(--text-primary); margin-bottom: 20px; }

/* 作业卡片 */
.assignment-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 16px; }
.assignment-card {
  background: var(--bg-card); border-radius: var(--radius); padding: 20px;
  box-shadow: var(--shadow-sm); border: 1px solid var(--border-light);
  display: flex; flex-direction: column; gap: 10px;
  transition: all 0.2s ease;
}
.assignment-card.submitted {
  opacity: 0.85;
  border-color: #e1f5e1;
  background: #f9fdf9;
}
.a-header { display: flex; justify-content: space-between; align-items: center; }
.a-header h3 { margin: 0; font-size: 16px; font-weight: 600; }
.a-info { display: flex; gap: 16px; font-size: 13px; color: var(--text-muted); }
.a-chapter { font-size: 13px; color: var(--text-secondary); }
.a-actions { margin-top: auto; padding-top: 8px; }
.a-score { font-size: 14px; color: #4caf50; }
.a-score strong { font-size: 18px; }

/* 答题头部 */
.answer-header { display: flex; align-items: center; gap: 16px; margin-bottom: 20px; flex-wrap: wrap; }
.answer-header h3 { margin: 0; font-size: 18px; }
.answer-progress { font-size: 13px; color: var(--text-muted); margin-left: auto; }

/* 题目网格 */
.question-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; margin-bottom: 24px; }
.question-card {
  background: var(--bg-card); border-radius: var(--radius); padding: 20px;
  box-shadow: var(--shadow-sm); border: 1px solid var(--border-light);
  display: flex; flex-direction: column; gap: 12px;
}
.card-top { display: flex; align-items: center; justify-content: space-between; }
.q-num {
  width: 32px; height: 32px; border-radius: 10px; background: var(--bg);
  font-size: 14px; font-weight: 700; display: flex; align-items: center; justify-content: center;
}
.q-tags { display: flex; gap: 6px; }
.q-title { font-size: 14px; font-weight: 500; line-height: 1.6; }

/* 选项 */
.q-options { margin-top: 4px; }
.options-list { display: flex; flex-direction: column; align-items: flex-start; gap: 4px; }
.opt-item {
  margin-right: 0 !important; padding: 8px 10px; border-radius: 8px;
  border: 1px solid var(--border-light); transition: all 0.2s ease;
  height: auto !important; justify-content: flex-start;
}
.opt-item:hover { background: var(--bg-hover); border-color: var(--primary); }
.opt-letter {
  display: inline-flex; align-items: center; justify-content: center;
  width: 22px; height: 22px; border-radius: 6px; background: var(--bg);
  font-size: 12px; font-weight: 600; margin-right: 8px; flex-shrink: 0;
}
.opt-text { font-size: 13px; color: var(--text-secondary); line-height: 1.4; }
:deep(.el-radio__label), :deep(.el-checkbox__label) { padding-left: 0; display: flex; align-items: center; }

/* 结果 */
.result-summary {
  display: flex; flex-direction: column; align-items: center; gap: 16px;
  padding: 40px; background: var(--bg-card); border-radius: var(--radius);
  box-shadow: var(--shadow-sm); margin-bottom: 24px;
}
.score-num { font-size: 48px; font-weight: 700; color: var(--primary); }
.score-total { font-size: 16px; color: var(--text-muted); }

/* 答题回顾 */
.review-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; margin-bottom: 24px; }
.review-grid .question-card { border-left: 3px solid #eee; }
.review-grid .question-card.correct { border-left-color: #4caf50; background: #f9fdf9; }
.review-grid .question-card.wrong { border-left-color: #f56c6c; background: #fef5f5; }
.q-meta { font-size: 13px; display: flex; gap: 16px; color: var(--text-secondary); }
.text-success { color: #4caf50; }
.text-danger { color: #f56c6c; }
.opt-letter.opt-correct { background: #4caf50; color: #fff; }
.opt-letter.opt-wrong { background: #f56c6c; color: #fff; }
.opt-letter.opt-selected { border: 2px solid var(--primary); }
.opt-tag { font-size: 11px; padding: 1px 6px; border-radius: 4px; margin-left: 8px; }
.correct-tag { background: #e8f5e9; color: #4caf50; }
.wrong-tag { background: #fce4ec; color: #f56c6c; }

/* 底部提交栏 */
.submit-bar {
  position: sticky; bottom: 0; display: flex; align-items: center;
  justify-content: center; gap: 24px; padding: 16px 24px;
  background: var(--bg-card); border-radius: var(--radius);
  box-shadow: var(--shadow-md); border: 1px solid var(--border-light); z-index: 10;
}
.submit-info { font-size: 14px; color: var(--text-muted); }
.submit-info strong { color: var(--primary); font-size: 16px; }

@media (max-width: 900px) { .question-grid, .review-grid { grid-template-columns: 1fr; } }
</style>
