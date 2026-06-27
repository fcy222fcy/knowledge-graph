<template>
  <div class="chat-panel">
    <div class="messages-area" ref="messagesRef">
      <EmptyChat v-if="!messages.length && !loading" @ask="(q) => $emit('ask', q)" />
      <template v-else>
        <MessageBubble v-for="msg in visibleMessages" :key="msg.message_id" :message="msg" />
        <SourceReference
          v-if="lastAssistantMsg"
          :sources="lastAssistantMsg.sources || []"
          :related-points="lastAssistantMsg.relatedKnowledgePoints || []"
        />
      </template>
      <div v-if="isAsking" class="thinking-message">
        <div class="thinking-avatar">
          <el-icon :size="20"><ChatDotRound /></el-icon>
        </div>
        <div class="thinking-content">
          <div class="thinking-text">
            <span class="thinking-icon">🔍</span>
            正在为你寻找相关知识...
          </div>
          <div class="thinking-hint">AI 正在分析问题并检索知识库</div>
        </div>
      </div>
    </div>

    <div class="input-area">
      <el-input
        v-model="inputText"
        type="textarea"
        :autosize="{ minRows: 1, maxRows: 4 }"
        placeholder="输入你的问题... (Enter 发送, Shift+Enter 换行)"
        :disabled="isAsking"
        @keydown.enter.exact.prevent="handleSend"
      />
      <el-button
        type="primary"
        :disabled="!inputText.trim() || isAsking"
        :loading="isAsking"
        @click="handleSend"
      >发送</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { ChatDotRound } from '@element-plus/icons-vue'
import type { MessageItem } from '@/types/qa'
import type { AskResponse } from '@/services/qa'
import MessageBubble from './MessageBubble.vue'
import SourceReference from './SourceReference.vue'
import EmptyChat from './EmptyChat.vue'

interface ExtendedMessage extends MessageItem {
  sources?: AskResponse['sources']
  relatedKnowledgePoints?: AskResponse['related_knowledge_points']
}

const props = defineProps<{
  messages: ExtendedMessage[]
  loading: boolean
  isAsking: boolean
}>()

const emit = defineEmits<{
  ask: [question: string]
}>()

const inputText = ref('')
const messagesRef = ref<HTMLDivElement>()

const lastAssistantMsg = computed(() => {
  const assistantMsgs = props.messages.filter(m => m.role === 'assistant')
  return assistantMsgs[assistantMsgs.length - 1] || null
})

// 过滤掉内容为空的 assistant 消息
const visibleMessages = computed(() => {
  return props.messages.filter(m => m.role !== 'assistant' || m.content)
})

const handleSend = () => {
  const text = inputText.value.trim()
  if (!text || props.isAsking) return
  emit('ask', text)
  inputText.value = ''
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesRef.value) {
    messagesRef.value.scrollTop = messagesRef.value.scrollHeight
  }
}

watch(() => props.messages.length, scrollToBottom)
watch(() => props.isAsking, scrollToBottom)
</script>

<style scoped>
.chat-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}
.input-area {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--border-light);
  background: var(--bg-card);
  align-items: flex-end;
}
.input-area .el-textarea {
  flex: 1;
}
.thinking-message {
  display: flex;
  gap: 12px;
  padding: 16px;
  margin-bottom: 16px;
  animation: fadeIn 0.3s ease;
}
.thinking-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary-color), var(--primary-light));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}
.thinking-content {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 12px 16px;
  box-shadow: var(--shadow-sm);
}
.thinking-text {
  font-size: 14px;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}
.thinking-icon {
  animation: pulse 1.5s infinite;
}
.thinking-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
