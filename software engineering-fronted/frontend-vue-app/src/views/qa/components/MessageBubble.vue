<template>
  <div class="message-bubble" :class="message.role">
    <div class="bubble-avatar">
      <span v-if="message.role === 'user'">我</span>
      <span v-else>AI</span>
    </div>
    <div class="bubble-content">
      <div v-if="message.role === 'assistant'" class="bubble-text bubble-markdown" v-html="renderMarkdown(message.content)" />
      <div v-else class="bubble-text">{{ message.content }}</div>
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

/** 简易 Markdown → HTML（加粗、代码块、行内代码、列表、换行） */
const renderMarkdown = (text: string): string => {
  if (!text) return ''
  let html = text
    // 代码块 ```...```
    .replace(/```(\w*)\n([\s\S]*?)```/g, '<pre><code>$2</code></pre>')
    // 行内代码 `...`
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    // 加粗 **...**
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    // 无序列表
    .replace(/^[•·\-]\s+(.+)$/gm, '<li>$1</li>')
    // 有序列表
    .replace(/^\d+\.\s+(.+)$/gm, '<li>$1</li>')
    // 换行
    .replace(/\n/g, '<br>')
  // 包裹连续 <li> 为 <ul>
  html = html.replace(/(<li>.*?<\/li>(?:<br>)?)+/g, (match) => {
    return '<ul>' + match.replace(/<br>/g, '') + '</ul>'
  })
  return html
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
.bubble-markdown :deep(pre) {
  background: #1e293b;
  color: #e2e8f0;
  padding: 10px 14px;
  border-radius: 6px;
  font-size: 13px;
  overflow-x: auto;
  margin: 8px 0;
}
.bubble-markdown :deep(code) {
  font-family: 'Consolas', 'Monaco', monospace;
}
.bubble-markdown :deep(:not(pre) > code) {
  background: rgba(37, 99, 235, 0.1);
  color: #2563eb;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
}
.bubble-markdown :deep(ul) {
  margin: 6px 0;
  padding-left: 20px;
}
.bubble-markdown :deep(li) {
  margin: 3px 0;
  line-height: 1.5;
}
.bubble-markdown :deep(strong) {
  font-weight: 600;
  color: #0f172a;
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
