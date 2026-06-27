<template>
  <div class="question-card">
    <!-- 题目头部：题号 + 类型标签 + 难度标签 -->
    <div class="card-header">
      <div class="title-row">
        <span class="question-index">{{ index }}.</span>
        <span class="question-title">{{ question.title }}</span>
      </div>
      <div class="tag-row">
        <el-tag
          :type="typeTagStyle.type"
          size="small"
          effect="dark"
          round
        >
          {{ typeLabel }}
        </el-tag>
        <el-tag
          :type="difficultyTagStyle.type"
          size="small"
          :effect="difficultyTagStyle.effect"
          round
        >
          {{ difficultyLabel }}
        </el-tag>
      </div>
    </div>

    <!-- 选项区域 -->
    <div class="options-area">
      <!-- 单选：el-radio-group -->
      <el-radio-group
        v-if="question.type === 'single'"
        :model-value="modelValue as string"
        @update:model-value="(val: string | number | boolean | undefined) => emit('update:modelValue', String(val))"
        class="options-group"
      >
        <el-radio
          v-for="option in question.options"
          :key="option.key"
          :value="option.key"
          class="option-item"
        >
          <span class="option-key">{{ option.key }}.</span>
          <span class="option-value">{{ option.value }}</span>
        </el-radio>
      </el-radio-group>

      <!-- 多选：el-checkbox-group -->
      <el-checkbox-group
        v-else-if="question.type === 'multiple'"
        :model-value="(Array.isArray(modelValue) ? modelValue : []) as string[]"
        @update:model-value="(val: (string | number)[]) => emit('update:modelValue', val.map(String))"
        class="options-group"
      >
        <el-checkbox
          v-for="option in question.options"
          :key="option.key"
          :value="option.key"
          class="option-item"
        >
          <span class="option-key">{{ option.key }}.</span>
          <span class="option-value">{{ option.value }}</span>
        </el-checkbox>
      </el-checkbox-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface QuestionOption {
  key: string
  value: string
}

interface Question {
  id: number
  title: string
  type: 'single' | 'multiple'
  difficulty: 'easy' | 'medium' | 'hard'
  options: QuestionOption[]
  answer: string
  explanation: string
}

const props = defineProps<{
  question: Question
  index: number
  modelValue: string | string[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string | string[]]
}>()

// 类型标签文案
const typeLabel = computed(() => {
  const map: Record<Question['type'], string> = {
    single: '单选',
    multiple: '多选',
  }
  return map[props.question.type]
})

// 类型标签样式
const typeTagStyle = computed(() => {
  const map: Record<Question['type'], { type: 'primary' | 'success' | 'warning' | 'info' | 'danger' }> = {
    single: { type: 'primary' },       // 蓝色
    multiple: { type: 'success' }, // 绿色
  }
  return map[props.question.type]
})

// 难度标签文案
const difficultyLabel = computed(() => {
  const map: Record<Question['difficulty'], string> = {
    easy: '简单',
    medium: '中等',
    hard: '困难',
  }
  return map[props.question.difficulty]
})

// 难度标签样式
const difficultyTagStyle = computed(() => {
  const map: Record<Question['difficulty'], { type: 'primary' | 'success' | 'warning' | 'info' | 'danger'; effect: 'dark' | 'light' | 'plain' }> = {
    easy: { type: 'success', effect: 'light' },
    medium: { type: 'warning', effect: 'light' },
    hard: { type: 'danger', effect: 'light' },
  }
  return map[props.question.difficulty]
})
</script>

<style scoped>
.question-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  transition: box-shadow 0.25s ease;
}

.question-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

/* 头部 */
.card-header {
  margin-bottom: 20px;
}

.title-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 12px;
}

.question-index {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
  flex-shrink: 0;
}

.question-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  line-height: 1.6;
}

.tag-row {
  display: flex;
  gap: 8px;
}

/* 选项区域 */
.options-group {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.option-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  margin: 0;
  border-radius: 8px;
  transition: background-color 0.2s ease;
  height: auto;
}

.option-item:hover {
  background-color: #f5f7fa;
}

.option-key {
  font-weight: 600;
  color: #606266;
  margin-right: 6px;
}

.option-value {
  color: #303133;
  line-height: 1.5;
}

/* Element Plus 样式覆盖 */
:deep(.el-radio),
:deep(.el-checkbox) {
  margin-right: 0;
  --el-radio-height: auto;
}

:deep(.el-radio__input),
:deep(.el-checkbox__input) {
  margin-top: 1px;
}

:deep(.el-radio__label),
:deep(.el-checkbox__label) {
  padding-left: 10px;
  font-size: 15px;
}
</style>
