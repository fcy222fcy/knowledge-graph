<template>
  <div class="empty-chat">
    <div class="empty-icon">💬</div>
    <h3>开始提问吧</h3>
    <p>基于知识图谱的智能问答，帮你理解软件工程知识</p>
    <div class="presets">
      <div
        v-for="q in presetQuestions"
        :key="q"
        class="preset-item"
        @click="$emit('ask', q)"
      >{{ q }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineEmits<{ ask: [question: string] }>()

const presetQuestions = [
  '什么是需求分析？',
  '软件测试有哪些方法？',
  '什么是软件生命周期？',
  '敏捷开发和瀑布模型的区别？'
]
</script>

<style scoped>
.empty-chat {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 40px;
  text-align: center;
}
.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}
h3 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}
p {
  font-size: 14px;
  color: var(--text-muted);
  margin-bottom: 32px;
}
.presets {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
  max-width: 400px;
}
.preset-item {
  padding: 12px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 10px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s;
}
.preset-item:hover {
  border-color: var(--primary);
  color: var(--primary);
  background: var(--primary-light);
}
</style>
