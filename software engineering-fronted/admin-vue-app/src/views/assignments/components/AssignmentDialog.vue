<template>
  <el-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="isEdit ? '编辑作业' : '新建作业'"
    width="850px"
    destroy-on-close
    @open="initForm"
  >
    <!-- 基本信息 -->
    <el-form :model="form" label-width="80px">
      <el-form-item label="作业名称" required>
        <el-input v-model="form.title" placeholder="请输入作业名称" />
      </el-form-item>
      <div style="display: flex; gap: 16px;">
        <el-form-item label="所属章节" style="flex: 1;">
          <el-input v-model="form.chapter" placeholder="如：第三章 软件需求" />
        </el-form-item>
        <el-form-item label="截止时间" required style="flex: 1;">
          <el-date-picker v-model="form.deadline" type="datetime" placeholder="选择截止时间" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%;" />
        </el-form-item>
      </div>
      <el-form-item label="作业说明">
        <el-input v-model="form.description" type="textarea" :rows="2" placeholder="可选：作业说明" />
      </el-form-item>
    </el-form>

    <el-divider />

    <!-- 题目列表 -->
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
      <h4 style="margin: 0;">📋 题目列表（共 {{ form.questions.length }} 题）</h4>
    </div>

    <div v-for="(q, idx) in form.questions" :key="idx" class="question-card">
      <div class="q-card-header">
        <span class="q-num">第 {{ idx + 1 }} 题</span>
        <div style="display: flex; gap: 8px; align-items: center;">
          <el-tag size="small">{{ q.type === 'single' ? '单选' : q.type === 'multiple' ? '多选' : '判断' }}</el-tag>
          <el-input-number v-model="q.score" :min="1" :max="100" size="small" controls-position="right" style="width: 80px;" />
          <span style="font-size: 12px; color: #999;">分</span>
          <el-button size="small" type="danger" text @click="removeQuestion(idx)">🗑️</el-button>
        </div>
      </div>

      <el-form label-width="60px" size="small">
        <el-form-item label="题目">
          <el-input v-model="q.title" type="textarea" :rows="2" placeholder="请输入题目内容" />
        </el-form-item>
        <el-form-item label="类型">
          <el-radio-group v-model="q.type" @change="onTypeChange(q)">
            <el-radio value="single">单选</el-radio>
            <el-radio value="multiple">多选</el-radio>
            <el-radio value="judge">判断</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="选项" v-if="q.type !== 'judge'">
          <div style="width: 100%;">
            <div v-for="(opt, oi) in q.options" :key="oi" style="display: flex; gap: 8px; margin-bottom: 6px; align-items: center;">
              <span class="opt-key">{{ opt.key }}</span>
              <el-input v-model="opt.value" placeholder="选项内容" style="flex: 1;" />
              <el-button size="small" text type="danger" @click="q.options.splice(oi, 1)" :disabled="q.options.length <= 2">✕</el-button>
            </div>
            <el-button size="small" text type="primary" @click="addOption(q)">＋ 添加选项</el-button>
          </div>
        </el-form-item>
        <el-form-item label="选项" v-else>
          <div style="width: 100%;">
            <div style="display: flex; gap: 8px; margin-bottom: 6px; align-items: center;">
              <span class="opt-key">✓</span>
              <el-input v-model="q.options[0].value" placeholder="正确（如：正确）" />
            </div>
            <div style="display: flex; gap: 8px; margin-bottom: 6px; align-items: center;">
              <span class="opt-key">✗</span>
              <el-input v-model="q.options[1].value" placeholder="错误（如：错误）" />
            </div>
          </div>
        </el-form-item>
        <el-form-item label="答案" required>
          <!-- 单选：下拉选择 -->
          <el-select
            v-if="q.type === 'single'"
            v-model="q.answer"
            placeholder="请选择正确答案"
            style="width: 200px;"
          >
            <el-option
              v-for="opt in q.options"
              :key="opt.key"
              :label="`${opt.key}. ${opt.value || '未填写'}`"
              :value="opt.key"
            />
          </el-select>
          <!-- 多选：复选框 -->
          <el-checkbox-group
            v-else-if="q.type === 'multiple'"
            v-model="q.multipleAnswer"
            @change="(val: string[]) => { q.answer = val.sort().join(',') }"
          >
            <el-checkbox v-for="opt in q.options" :key="opt.key" :value="opt.key" :label="`${opt.key}. ${opt.value || '未填写'}`" />
          </el-checkbox-group>
          <!-- 判断：下拉选择 ✓/✗ -->
          <el-select
            v-else
            v-model="q.answer"
            placeholder="请选择"
            style="width: 200px;"
          >
            <el-option label="✓ 正确" value="✓" />
            <el-option label="✗ 错误" value="✗" />
          </el-select>
        </el-form-item>
        <el-form-item label="解析">
          <el-input v-model="q.explanation" placeholder="可选：题目解析" />
        </el-form-item>
      </el-form>
    </div>

    <button class="add-btn" @click="addQuestion">＋ 添加新题目</button>

    <template #footer>
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <span style="font-size: 13px; color: #999;">共 {{ form.questions.length }} 题 · 总分 {{ totalScore }} 分</span>
        <div style="display: flex; gap: 8px;">
          <el-button @click="$emit('update:modelValue', false)">取消</el-button>
          <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { createAssignment, updateAssignment } from '@/services/assignment'

const props = defineProps<{
  modelValue: boolean
  assignment?: any
}>()

const emit = defineEmits<{
  'update:modelValue': [val: boolean]
  saved: []
}>()

const isEdit = computed(() => !!props.assignment?.id)
const saving = ref(false)

const form = ref({
  title: '',
  description: '',
  chapter: '',
  deadline: '',
  questions: [] as any[]
})

const totalScore = computed(() => form.value.questions.reduce((s, q) => s + (q.score || 10), 0))

// 切换题目类型时重置答案
const onTypeChange = (q: any) => {
  q.answer = ''
  q.multipleAnswer = []
}

const initForm = () => {
  if (props.assignment?.id) {
    // 编辑模式 - 从详情接口获取
    form.value = {
      title: props.assignment.title || '',
      description: props.assignment.description || '',
      chapter: props.assignment.chapter || '',
      deadline: props.assignment.deadline || '',
      questions: (props.assignment.questions || []).map((q: any) => ({
        ...q,
        multipleAnswer: q.type === 'multiple' ? (q.answer || '').split(',').filter(Boolean) : []
      }))
    }
  } else {
    form.value = { title: '', description: '', chapter: '', deadline: '', questions: [] }
    addQuestion()
  }
}

const newQuestion = () => ({
  title: '',
  type: 'single',
  options: [
    { key: 'A', value: '' },
    { key: 'B', value: '' },
    { key: 'C', value: '' },
    { key: 'D', value: '' },
  ],
  answer: '',
  multipleAnswer: [] as string[],
  explanation: '',
  score: 10,
  sort_order: 0
})

const addQuestion = () => {
  form.value.questions.push(newQuestion())
}

const removeQuestion = (idx: number) => {
  form.value.questions.splice(idx, 1)
}

const addOption = (q: any) => {
  const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  const nextKey = letters[q.options.length] || String(q.options.length + 1)
  q.options.push({ key: nextKey, value: '' })
}

const handleSave = async () => {
  if (!form.value.title) return ElMessage.warning('请输入作业名称')
  if (!form.value.deadline) return ElMessage.warning('请选择截止时间')
  if (form.value.questions.length === 0) return ElMessage.warning('请至少添加一道题目')

  for (let i = 0; i < form.value.questions.length; i++) {
    const q = form.value.questions[i]
    if (!q.title) return ElMessage.warning(`第 ${i + 1} 题请输入题目内容`)
    if (!q.answer) return ElMessage.warning(`第 ${i + 1} 题请输入正确答案`)
    q.sort_order = i + 1
  }

  saving.value = true
  try {
    if (isEdit.value) {
      await updateAssignment(props.assignment.id, form.value)
    } else {
      await createAssignment(form.value)
    }
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    emit('update:modelValue', false)
    emit('saved')
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '操作失败')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.question-card {
  border: 1px solid #e0e0e0; border-radius: 8px; padding: 14px; margin-bottom: 12px;
}
.q-card-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;
  padding-bottom: 8px; border-bottom: 1px solid #f0f0f0;
}
.q-num { font-weight: 600; color: #1a73e8; font-size: 13px; background: #e8f0fe; padding: 2px 8px; border-radius: 4px; }
.opt-key {
  width: 24px; height: 24px; border-radius: 50%; background: #f0f0f0;
  display: inline-flex; align-items: center; justify-content: center;
  font-size: 12px; font-weight: 600; color: #666; flex-shrink: 0;
}
.add-btn {
  width: 100%; padding: 12px; border: 2px dashed #ccc; border-radius: 8px;
  background: transparent; color: #999; font-size: 14px; cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 6px;
}
.add-btn:hover { border-color: #1a73e8; color: #1a73e8; }
</style>
