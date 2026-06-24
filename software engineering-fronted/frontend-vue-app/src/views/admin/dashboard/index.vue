<template>
  <div class="admin-dashboard" v-loading="loading">
    <h2 class="page-title">系统概览</h2>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="4" v-for="item in statItems" :key="item.label">
        <el-card shadow="never" class="stat-card">
          <div class="stat-icon" :style="{ backgroundColor: item.color }">
            <el-icon :size="24"><component :is="item.icon" /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ item.value }}</div>
            <div class="stat-label">{{ item.label }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 用户角色分布 -->
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <span>用户角色分布</span>
          </template>
          <div class="role-stats">
            <div class="role-item" v-for="(count, role) in userStats" :key="role">
              <span class="role-name">{{ roleLabels[role as string] }}</span>
              <span class="role-count">{{ count }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getAnalyticsOverview, getUserStats } from '@/services/admin'
import type { AnalyticsOverview, UserStats } from '@/services/admin'
import {
  User,
  Document,
  Collection,
  ChatLineRound,
  Edit,
  DataAnalysis
} from '@element-plus/icons-vue'

const loading = ref(false)
const overview = ref<AnalyticsOverview>({
  user_count: 0,
  document_count: 0,
  knowledge_count: 0,
  question_count: 0,
  session_count: 0,
  quiz_count: 0
})
const userStats = ref<UserStats>({
  admin: 0,
  teacher: 0,
  student: 0
})

const roleLabels: Record<string, string> = {
  admin: '管理员',
  teacher: '教师',
  student: '学生'
}

const statItems = computed(() => [
  { label: '用户总数', value: overview.value.user_count, icon: User, color: '#409eff' },
  { label: '资料总数', value: overview.value.document_count, icon: Document, color: '#67c23a' },
  { label: '知识点', value: overview.value.knowledge_count, icon: Collection, color: '#e6a23c' },
  { label: '题目总数', value: overview.value.question_count, icon: Edit, color: '#f56c6c' },
  { label: '问答会话', value: overview.value.session_count, icon: ChatLineRound, color: '#909399' },
  { label: '答题记录', value: overview.value.quiz_count, icon: DataAnalysis, color: '#9b59b6' }
])

const fetchData = async () => {
  loading.value = true
  try {
    const [overviewRes, statsRes] = await Promise.all([
      getAnalyticsOverview(),
      getUserStats()
    ])
    overview.value = overviewRes.data
    userStats.value = statsRes.data
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.admin-dashboard {
  background: #fff;
  padding: 24px;
  border-radius: 8px;
}

.page-title {
  margin: 0 0 24px 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.stat-cards {
  display: flex;
}

.stat-card {
  text-align: center;
}

.stat-card :deep(.el-card__body) {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  margin-bottom: 12px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.role-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.role-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.role-name {
  color: #606266;
}

.role-count {
  font-size: 20px;
  font-weight: 600;
  color: #409eff;
}
</style>
