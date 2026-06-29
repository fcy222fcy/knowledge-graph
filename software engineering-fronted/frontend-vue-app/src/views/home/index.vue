<template>
  <div class="home-container" v-loading="loading">
    <TodayStats :data="overview" />
    <QuickActions />

    <div class="home-grid">
      <RecentQAList :data="recentQA" />
      <LearningTrend :data="trendData" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { OverviewStats, TrendData } from '@/types/stats'
import { getOverview, getTrends } from '@/services/stats'
import { getAskHistory } from '@/services/qa'
import TodayStats from './components/TodayStats.vue'
import QuickActions from './components/QuickActions.vue'
import RecentQAList from './components/RecentQAList.vue'
import LearningTrend from './components/LearningTrend.vue'

const loading = ref(true)
const overview = ref<OverviewStats | null>(null)
const trendData = ref<TrendData | null>(null)
const recentQA = ref<any[]>([])

// 模拟数据
const mockOverview = {
  total_learning_hours: 86.3, total_questions_asked: 168, total_quizzes_taken: 42,
  average_correct_rate: 0.81, knowledge_points_mastered: 5, knowledge_points_total: 8,
  today_learning_hours: 2.5, today_questions_asked: 12
}
const mockTrend = {
  daily_stats: [
    { date: '06-21', questions_asked: 8, learning_hours: 1.5, correct_rate: 0.75 },
    { date: '06-22', questions_asked: 12, learning_hours: 2.3, correct_rate: 0.83 },
    { date: '06-23', questions_asked: 6, learning_hours: 1.0, correct_rate: 0.67 },
    { date: '06-24', questions_asked: 15, learning_hours: 3.2, correct_rate: 0.87 },
    { date: '06-25', questions_asked: 10, learning_hours: 2.0, correct_rate: 0.80 },
    { date: '06-26', questions_asked: 18, learning_hours: 3.5, correct_rate: 0.89 },
    { date: '06-27', questions_asked: 14, learning_hours: 2.8, correct_rate: 0.86 },
  ]
}

const fetchData = async () => {
  loading.value = true
  try {
    overview.value = mockOverview
    trendData.value = mockTrend
    const results = await Promise.allSettled([
      getAskHistory({ page: 1, size: 5 })
    ])
    if (results[0].status === 'fulfilled') recentQA.value = results[0].value.data.list || []
  } catch (error) {
    console.error('获取首页数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.home-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}
@media (max-width: 768px) {
  .home-grid {
    grid-template-columns: 1fr;
  }
}
</style>
