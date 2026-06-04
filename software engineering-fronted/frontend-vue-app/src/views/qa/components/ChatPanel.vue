<template>
  <div class="chat-panel">
    <div class="messages-area" ref="messagesRef">
      <EmptyChat v-if="!messages.length && !loading" @ask="(q) => $emit('ask', q)" />
      <template v-else>
        <MessageBubble v-for="msg in messages" :key="msg.message_id" :message="msg" />
        <SourceReference
          v-if="lastAssistantMsg"
          :sources="lastAssistantMsg.sources || []"
          :related-points="lastAssistantMsg.relatedKnowledgePoints || []"
        />
      </template>
      <div v-if="isAsking" class="typing-indicator">
        <span></span><span></span><span></span>
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
.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 12px 16px;
  background: var(--bg-card);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  width: fit-content;
  margin-left: 48px;
}
.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--text-muted);
  animation: bounce 1.4s infinite ease-in-out;
}
.typing-indicator span:nth-child(1) { animation-delay: -0.32s; }
.typing-indicator span:nth-child(2) { animation-delay: -0.16s; }
@keyframes bounce {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}
</style>
