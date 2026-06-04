<template>
  <div class="session-list">
    <div class="session-header">
      <h3>会话列表</h3>
      <el-button type="primary" size="small" @click="$emit('new')">
        <el-icon><Plus /></el-icon> 新建
      </el-button>
    </div>
    <div class="session-items" v-loading="loading">
      <div
        v-for="session in sessions"
        :key="session.conversation_id"
        class="session-item"
        :class="{ active: activeId === session.conversation_id }"
        @click="$emit('select', session.conversation_id)"
      >
        <div class="session-title">{{ session.title }}</div>
        <div class="session-meta">
          <span class="session-last">{{ session.last_question }}</span>
          <span class="session-count">{{ session.message_count }}条</span>
        </div>
      </div>
      <EmptyState v-if="!loading && !sessions.length" text="暂无会话" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import type { ConversationItem } from '@/types/qa'
import EmptyState from '@/components/common/EmptyState.vue'

defineProps<{
  sessions: ConversationItem[]
  activeId: number | null
  loading: boolean
}>()

defineEmits<{
  select: [id: number]
  new: []
}>()
</script>

<style scoped>
.session-list {
  width: 300px;
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  background: var(--bg-card);
}
.session-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--border-light);
}
.session-header h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}
.session-items {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}
.session-item {
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
  margin-bottom: 4px;
}
.session-item:hover {
  background: var(--bg-hover);
}
.session-item.active {
  background: var(--primary-light);
}
.session-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.session-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-muted);
}
.session-last {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  margin-right: 8px;
}
</style>
