<template>
  <div class="chart-card">
    <h3>热门知识点排行</h3>
    <div v-if="data.length" class="rank-list">
      <div v-for="(item, index) in data" :key="item.knowledge_point_id" class="rank-item">
        <span class="rank-num" :class="{ top3: index < 3 }">{{ index + 1 }}</span>
        <span class="rank-name">{{ item.knowledge_point_name }}</span>
        <span class="rank-heat">{{ item.heat }}</span>
      </div>
    </div>
    <EmptyState v-else text="暂无数据" />
  </div>
</template>

<script setup lang="ts">
import type { HotKnowledgePoint } from '@/types/stats'
import EmptyState from '@/components/common/EmptyState.vue'

defineProps<{
  data: HotKnowledgePoint[]
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
.rank-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.rank-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-radius: 8px;
  background: var(--bg-hover);
}
.rank-num {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  background: var(--border-light);
  margin-right: 12px;
  flex-shrink: 0;
}
.rank-num.top3 {
  color: white;
  background: linear-gradient(135deg, #f59e0b, #f97316);
}
.rank-name {
  flex: 1;
  font-size: 14px;
  color: var(--text-primary);
}
.rank-heat {
  font-size: 13px;
  color: var(--text-muted);
  font-weight: 500;
}
</style>
