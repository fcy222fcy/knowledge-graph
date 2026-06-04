<template>
  <div class="qa-container">
    <SessionList
      :sessions="sessions"
      :active-id="currentSessionId"
      :loading="sessionsLoading"
      @select="handleSelectSession"
      @new="handleNewSession"
    />
    <ChatPanel
      :messages="messages"
      :loading="messagesLoading"
      :is-asking="isAsking"
      @ask="handleAskQuestion"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { ConversationItem } from '@/types/qa'
import type { AskResponse } from '@/services/qa'
import { getSessions, getSessionMessages, createSession, askQuestion } from '@/services/qa'
import { ElMessage } from 'element-plus'
import SessionList from './components/SessionList.vue'
import ChatPanel from './components/ChatPanel.vue'

interface ExtendedMessage {
  message_id: number
  role: 'user' | 'assistant'
  content: string
  created_at: string
  sources?: AskResponse['sources']
  relatedKnowledgePoints?: AskResponse['related_knowledge_points']
}

const sessions = ref<ConversationItem[]>([])
const currentSessionId = ref<number | null>(null)
const messages = ref<ExtendedMessage[]>([])
const isAsking = ref(false)
const sessionsLoading = ref(false)
const messagesLoading = ref(false)

const fetchSessions = async () => {
  sessionsLoading.value = true
  try {
    const result = await getSessions({ page: 1, size: 50 })
    sessions.value = result.data.list || []
  } catch (error) {
    console.error('获取会话列表失败:', error)
  } finally {
    sessionsLoading.value = false
  }
}

const handleNewSession = async () => {
  try {
    const result = await createSession()
    await fetchSessions()
    currentSessionId.value = result.data.conversation_id
    messages.value = []
  } catch (error) {
    console.error('创建会话失败:', error)
  }
}

const handleSelectSession = async (sessionId: number) => {
  currentSessionId.value = sessionId
  messagesLoading.value = true
  try {
    const result = await getSessionMessages(sessionId, { page: 1, size: 100 })
    messages.value = result.data.list || []
  } catch (error) {
    console.error('获取消息失败:', error)
  } finally {
    messagesLoading.value = false
  }
}

const handleAskQuestion = async (question: string) => {
  if (isAsking.value) return

  if (!currentSessionId.value) {
    try {
      const result = await createSession()
      currentSessionId.value = result.data.conversation_id
    } catch {
      ElMessage.error('创建会话失败')
      return
    }
  }

  const userMsg: ExtendedMessage = {
    message_id: Date.now(),
    role: 'user',
    content: question,
    created_at: new Date().toISOString()
  }
  messages.value.push(userMsg)

  isAsking.value = true
  try {
    const result = await askQuestion({
      question,
      conversation_id: currentSessionId.value!
    })
    const data = result.data

    const aiMsg: ExtendedMessage = {
      message_id: data.question_id,
      role: 'assistant',
      content: data.answer,
      created_at: data.created_at,
      sources: data.sources,
      relatedKnowledgePoints: data.related_knowledge_points
    }
    messages.value.push(aiMsg)

    fetchSessions()
  } catch (error) {
    ElMessage.error('提问失败')
  } finally {
    isAsking.value = false
  }
}

onMounted(() => {
  fetchSessions()
})
</script>

<style scoped>
.qa-container {
  display: flex;
  height: calc(100vh - 100px);
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}
</style>
