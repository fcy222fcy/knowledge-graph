<template>
  <div class="chart-card">
    <h3>薄弱知识点</h3>
    <div v-if="data && data.length" class="weak-list">
      <div v-for="item in data" :key="item.knowledge_point_id" class="weak-item">
        <div class="weak-header">
          <span class="weak-name">{{ item.knowledge_point_name }}</span>
          <span class="weak-rate">{{ (item.correct_rate * 100).toFixed(0) }}%</span>
        </div>
        <el-progress
          :percentage="item.correct_rate * 100"
          :color="item.correct_rate < 0.3 ? '#f43f5e' : '#f59e0b'"
          :show-text="false"
        />
        <div class="suggested" v-if="item.suggested_questions.length">
          <span class="suggested-label">建议练习：</span>
          <el-tag
            v-for="q in item.suggested_questions"
            :key="q.id"
            size="small"
            class="suggested-tag"
          >{{ q.title }}</el-tag>
        </div>
      </div>
    </div>
    <EmptyState v-else text="暂无薄弱知识点" />
  </div>
</template>

<script setup lang="ts">
import type { WeakPoint } from '@/types/stats'
import EmptyState from '@/components/common/EmptyState.vue'

defineProps<{
  data: WeakPoint[]
}>()
</script>

<style scoped>
.chart-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.chart-card h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
.weak-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.weak-item {
  padding: 16px;
  border-radius: 8px;
  background: var(--bg-hover);
}
.weak-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}
.weak-name {
  font-weight: 500;
  color: var(--text-primary);
}
.weak-rate {
  font-size: 13px;
  color: #f43f5e;
  font-weight: 600;
}
.suggested {
  margin-top: 8px;
}
.suggested-label {
  font-size: 12px;
  color: var(--text-muted);
}
.suggested-tag {
  margin: 2px 4px 2px 0;
}
</style>
