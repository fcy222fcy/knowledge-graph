<template>
  <div class="message-bubble" :class="message.role">
    <div class="bubble-avatar">
      <span v-if="message.role === 'user'">我</span>
      <span v-else>AI</span>
    </div>
    <div class="bubble-content">
      <div v-if="message.role === 'assistant'" class="bubble-text bubble-markdown" v-html="renderMarkdown(message.content)" />
      <div v-else class="bubble-text">{{ message.content }}</div>
      <div v-if="message.content" class="bubble-time">{{ formatTime(message.created_at) }}</div>
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

/** Markdown → HTML（标题、表格、代码块、水平线、列表、加粗、行内代码、换行） */
const renderMarkdown = (text: string): string => {
  if (!text) return ''

  // 1. 先提取代码块，用占位符替换，避免内部内容被其他规则干扰
  const codeBlocks: string[] = []
  let html = text.replace(/```(\w*)\n([\s\S]*?)```/g, (_m, _lang, code) => {
    const idx = codeBlocks.length
    codeBlocks.push(`<pre><code>${code}</code></pre>`)
    return `\x00CB${idx}\x00`
  })

  // 2. 表格：连续的 | 行
  html = html.replace(/((?:^\|.+\|$\n?)+)/gm, (block) => {
    const rows = block.trim().split('\n').filter(r => r.trim())
    if (rows.length < 2) return block
    const dataRows = rows.filter(r => !/^\|[\s\-:|]+\|$/.test(r.trim()))
    if (dataRows.length === 0) return block
    let table = '<table><thead><tr>'
    dataRows[0].split('|').filter(c => c.trim() !== '').forEach(c => {
      table += `<th>${c.trim()}</th>`
    })
    table += '</tr></thead><tbody>'
    for (let i = 1; i < dataRows.length; i++) {
      table += '<tr>'
      dataRows[i].split('|').filter(c => c.trim() !== '').forEach(c => {
        table += `<td>${c.trim()}</td>`
      })
      table += '</tr>'
    }
    table += '</tbody></table>'
    return table
  })

  // 3. 标题 h1-h4（从 h4 开始替换，避免 h4 被 h3 的正则误匹配）
  html = html.replace(/^####\s+(.+)$/gm, '<h4>$1</h4>')
  html = html.replace(/^###\s+(.+)$/gm, '<h3>$1</h3>')
  html = html.replace(/^##\s+(.+)$/gm, '<h2>$1</h2>')
  html = html.replace(/^#\s+(.+)$/gm, '<h1>$1</h1>')

  // 4. 水平线 --- 或 *** 或 ___
  html = html.replace(/^[-*_]{3,}\s*$/gm, '<hr>')

  // 5. 先把列表项转为 <li>，再包裹 <ul>（在换行转换之前）
  html = html.replace(/^[-*]\s+(.+)$/gm, '<li>$1</li>')
  html = html.replace(/^\d+\.\s+(.+)$/gm, '<li>$1</li>')
  html = html.replace(/((?:<li>.*?<\/li>\n?)+)/g, (match) => {
    return '<ul>' + match.replace(/\n/g, '') + '</ul>'
  })

  // 6. 行内元素
  html = html.replace(/`([^`]+)`/g, '<code>$1</code>')
  html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')

  // 7. 最后处理换行 → <br>（此时列表已包裹，不会破坏 <ul>）
  html = html.replace(/\n/g, '<br>')

  // 8. 还原代码块
  html = html.replace(/\x00CB(\d+)\x00/g, (_m, idx) => codeBlocks[Number(idx)])

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
  word-break: break-word;
}
.user .bubble-text {
  white-space: pre-wrap;
  background: #2563eb;
  color: white;
  border-bottom-right-radius: 4px;
}
.assistant .bubble-text {
  background: var(--bg-card);
  color: var(--text-primary);
  box-shadow: var(--shadow-sm);
  border-bottom-left-radius: 4px;
  white-space: normal;
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
.bubble-markdown :deep(h1) {
  font-size: 18px;
  font-weight: 700;
  margin: 12px 0 8px;
  padding-bottom: 6px;
  border-bottom: 1px solid #e2e8f0;
}
.bubble-markdown :deep(h2) {
  font-size: 16px;
  font-weight: 700;
  margin: 10px 0 6px;
  padding-bottom: 4px;
  border-bottom: 1px solid #f1f5f9;
}
.bubble-markdown :deep(h3) {
  font-size: 15px;
  font-weight: 600;
  margin: 8px 0 4px;
}
.bubble-markdown :deep(h4) {
  font-size: 14px;
  font-weight: 600;
  margin: 6px 0 4px;
}
.bubble-markdown :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 8px 0;
  font-size: 13px;
}
.bubble-markdown :deep(th),
.bubble-markdown :deep(td) {
  border: 1px solid #e2e8f0;
  padding: 6px 10px;
  text-align: left;
}
.bubble-markdown :deep(th) {
  background: #f8fafc;
  font-weight: 600;
}
.bubble-markdown :deep(td) {
  background: white;
}
.bubble-markdown :deep(hr) {
  border: none;
  border-top: 1px solid #e2e8f0;
  margin: 12px 0;
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
