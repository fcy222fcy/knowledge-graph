<template>
  <div class="message-bubble" :class="message.role">
    <div class="bubble-avatar">
      <span v-if="message.role === 'user'">我</span>
      <span v-else>AI</span>
    </div>
    <div class="bubble-content">
      <div class="bubble-text">{{ message.content }}</div>
      <div class="bubble-time">{{ formatTime(message.created_at) }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { MessageItem } from '@/types/qa'

defineProps<{
  message: MessageItem
}>()

const formatTime = (dateStr: string) => {
  const d = new Date(dateStr)
  return d.getHours().toString().padStart(2, '0') + ':' + d.getMinutes().toString().padStart(2, '0')
}
</script>

<style scoped>
.message-bubble {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  max-width: 80%;
}
.message-bubble.user {
  margin-left: auto;
  flex-direction: row-reverse;
}
.bubble-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: white;
  flex-shrink: 0;
}
.user .bubble-avatar {
  background: linear-gradient(135deg, #2563eb, #3b82f6);
}
.assistant .bubble-avatar {
  background: linear-gradient(135deg, #8b5cf6, #a78bfa);
}
.bubble-content {
  flex: 1;
}
.bubble-text {
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}
.user .bubble-text {
  background: #2563eb;
  color: white;
  border-bottom-right-radius: 4px;
}
.assistant .bubble-text {
  background: var(--bg-card);
  color: var(--text-primary);
  box-shadow: var(--shadow-sm);
  border-bottom-left-radius: 4px;
}
.bubble-time {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 4px;
}
.user .bubble-time {
  text-align: right;
}
</style>
