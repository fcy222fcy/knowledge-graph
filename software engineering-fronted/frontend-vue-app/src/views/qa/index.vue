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
import { getSessions, getSessionMessages, createSession, askQuestionStream } from '@/services/qa'
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

const handleNewSession = () => {
  // 不立即创建会话，等发送第一条消息时再创建
  currentSessionId.value = null
  messages.value = []
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

  // 先创建会话（如果需要）
  if (!currentSessionId.value) {
    try {
      const result = await createSession()
      currentSessionId.value = result.data.conversation_id
      await fetchSessions()
    } catch {
      ElMessage.error('创建会话失败')
      return
    }
  }

  // 添加用户消息
  const userMsg: ExtendedMessage = {
    message_id: Date.now(),
    role: 'user',
    content: question,
    created_at: new Date().toISOString()
  }
  messages.value.push(userMsg)

  // 添加空的 AI 消息占位
  const aiMsg: ExtendedMessage = {
    message_id: Date.now() + 1,
    role: 'assistant',
    content: '',
    created_at: new Date().toISOString()
  }
  messages.value.push(aiMsg)
  const aiMsgIndex = messages.value.length - 1

  isAsking.value = true

  try {
    const stream = askQuestionStream({
      question,
      conversation_id: currentSessionId.value!
    })

    for await (const event of stream) {
      if (event.type === 'session') {
        // 会话创建事件
        if (event.session_id) {
          currentSessionId.value = event.session_id
        }
      } else if (event.type === 'chunk') {
        // 流式内容块
        messages.value[aiMsgIndex].content += event.content
      } else if (event.type === 'done') {
        // 完成事件
        messages.value[aiMsgIndex].sources = event.sources
        messages.value[aiMsgIndex].relatedKnowledgePoints = event.related
        fetchSessions()
      } else if (event.type === 'error') {
        messages.value[aiMsgIndex].content = event.content || '提问失败，请重试'
      }
    }
  } catch (error) {
    console.error('流式提问失败:', error)
    messages.value[aiMsgIndex].content = '提问失败，请检查网络连接后重试'
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
