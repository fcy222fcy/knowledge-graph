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

const fetchData = async () => {
  loading.value = true
  try {
    const results = await Promise.allSettled([
      getOverview(),
      getTrends(7),
      getAskHistory({ page: 1, size: 5 })
    ])
    if (results[0].status === 'fulfilled') overview.value = results[0].value.data
    if (results[1].status === 'fulfilled') trendData.value = results[1].value.data
    if (results[2].status === 'fulfilled') recentQA.value = results[2].value.data.list || []
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
