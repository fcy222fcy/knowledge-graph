<template>
  <div class="quick-actions">
    <router-link to="/qa" class="action-card">
      <div class="action-icon" style="background: #eff6ff; color: #2563eb">
        <el-icon :size="24"><ChatDotRound /></el-icon>
      </div>
      <span class="action-title">智能问答</span>
      <span class="action-desc">向 AI 提问</span>
    </router-link>
    <router-link to="/files" class="action-card">
      <div class="action-icon" style="background: #f0fdf4; color: #10b981">
        <el-icon :size="24"><Upload /></el-icon>
      </div>
      <span class="action-title">上传资料</span>
      <span class="action-desc">管理学习文档</span>
    </router-link>
    <router-link to="/knowledge-graph" class="action-card">
      <div class="action-icon" style="background: #fef3c7; color: #f59e0b">
        <el-icon :size="24"><Share /></el-icon>
      </div>
      <span class="action-title">知识图谱</span>
      <span class="action-desc">可视化知识网络</span>
    </router-link>
    <router-link to="/stats" class="action-card">
      <div class="action-icon" style="background: #ede9fe; color: #8b5cf6">
        <el-icon :size="24"><DataAnalysis /></el-icon>
      </div>
      <span class="action-title">学习统计</span>
      <span class="action-desc">查看学习分析</span>
    </router-link>
  </div>
</template>

<script setup lang="ts">
import { ChatDotRound, Upload, Share, DataAnalysis } from '@element-plus/icons-vue'
</script>

<style scoped>
.quick-actions {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}
.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 24px 16px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  text-decoration: none;
  transition: all 0.2s;
  cursor: pointer;
}
.action-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}
.action-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.action-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}
.action-desc {
  font-size: 12px;
  color: var(--text-muted);
}
@media (max-width: 768px) {
  .quick-actions {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
