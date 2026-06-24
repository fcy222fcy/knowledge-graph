<template>
  <div class="dashboard-page">
    <div class="page-header">
      <h2>仪表盘</h2>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon" style="background: var(--color-primary-glow); color: var(--color-primary);">
          <el-icon :size="24"><User /></el-icon>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.studentCount }}</span>
          <span class="stat-label">学生总数</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: var(--color-success-light); color: var(--color-success);">
          <el-icon :size="24"><Document /></el-icon>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.documentCount }}</span>
          <span class="stat-label">资料总数</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: var(--color-warning-light); color: var(--color-warning);">
          <el-icon :size="24"><Clock /></el-icon>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.pendingReview }}</span>
          <span class="stat-label">待审核资料</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: var(--color-info-light); color: var(--color-info);">
          <el-icon :size="24"><Edit /></el-icon>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.questionCount }}</span>
          <span class="stat-label">题目总数</span>
        </div>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="quick-actions page-card">
      <h3>快捷操作</h3>
      <div class="action-grid">
        <router-link to="/admin/documents" class="action-item">
          <el-icon :size="32" color="var(--color-primary)"><Document /></el-icon>
          <span>审核资料</span>
        </router-link>
        <router-link to="/admin/questions" class="action-item">
          <el-icon :size="32" color="var(--color-success)"><Edit /></el-icon>
          <span>管理题目</span>
        </router-link>
        <router-link to="/admin/knowledge" class="action-item">
          <el-icon :size="32" color="var(--color-info)"><Share /></el-icon>
          <span>知识点管理</span>
        </router-link>
        <router-link to="/admin/students" class="action-item">
          <el-icon :size="32" color="var(--color-warning)"><User /></el-icon>
          <span>学生管理</span>
        </router-link>
      </div>
    </div>

    <!-- 最近活动 -->
    <div class="recent-activity page-card">
      <h3>最近活动</h3>
      <el-timeline>
        <el-timeline-item
          v-for="(activity, index) in recentActivities"
          :key="index"
          :timestamp="activity.time"
          :type="activity.type"
          placement="top"
        >
          {{ activity.content }}
        </el-timeline-item>
      </el-timeline>
      <el-empty v-if="recentActivities.length === 0" description="暂无活动记录" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { User, Document, Clock, Edit, Share } from '@element-plus/icons-vue'
import { getAnalyticsOverview } from '@/services/admin'

const stats = ref({
  studentCount: 0,
  documentCount: 0,
  pendingReview: 0,
  questionCount: 0,
})

const recentActivities = ref<Array<{
  content: string
  time: string
  type: '' | 'primary' | 'success' | 'warning' | 'danger'
}>>([])

onMounted(async () => {
  try {
    const data = await getAnalyticsOverview() as Record<string, unknown>
    stats.value = {
      studentCount: (data.student_count as number) || 0,
      documentCount: (data.document_count as number) || 0,
      pendingReview: (data.pending_review as number) || 0,
      questionCount: (data.question_count as number) || 0,
    }
    recentActivities.value = (data.recent_activities as typeof recentActivities.value) || []
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
})
</script>

<style scoped>
.dashboard-page {
  max-width: 1200px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-content {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-label {
  font-size: 14px;
  color: var(--text-muted);
}

.quick-actions {
  margin-bottom: 24px;
}

.quick-actions h3,
.recent-activity h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--text-primary);
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 24px;
  border-radius: var(--radius-md);
  background: var(--bg-hover);
  text-decoration: none;
  color: var(--text-primary);
  transition: all 0.2s;
}

.action-item:hover {
  background: var(--color-primary-glow);
  transform: translateY(-2px);
}

.recent-activity {
  min-height: 200px;
}
</style>
