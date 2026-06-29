<template>
  <div class="card">
    <h3>最近问答</h3>
    <div v-if="data.length" class="qa-list">
      <router-link
        v-for="item in data"
        :key="item.id"
        to="/qa"
        class="qa-item"
      >
        <span class="qa-question">{{ item.title || item.last_question }}</span>
        <span class="qa-time">{{ formatDate(item.created_at) }}</span>
      </router-link>
    </div>
    <EmptyState v-else text="暂无问答记录" />
  </div>
</template>

<script setup lang="ts">
import { formatDate } from '@/utils'
import EmptyState from '@/components/common/EmptyState.vue'

defineProps<{
  data: any[]
}>()
</script>

<style scoped>
.card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.card h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
.qa-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.qa-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  background: var(--bg-hover);
  text-decoration: none;
  transition: background 0.15s;
}
.qa-item:hover {
  background: var(--border-light);
}
.qa-question {
  font-size: 13px;
  color: var(--text-primary);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-right: 12px;
}
.qa-time {
  font-size: 12px;
  color: var(--text-muted);
  flex-shrink: 0;
}
</style>
