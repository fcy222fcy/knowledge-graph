<template>
  <div class="empty-state">
    <el-icon :size="48" color="#c0c4cc"><component :is="icon" /></el-icon>
    <p class="empty-text">{{ text }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
import { Document } from '@element-plus/icons-vue'

withDefaults(defineProps<{ text?: string; icon?: any }>(), {
  text: '暂无数据',
  icon: Document
})
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--text-muted);
}
.empty-text {
  margin-top: 12px;
  font-size: 14px;
}
</style>
