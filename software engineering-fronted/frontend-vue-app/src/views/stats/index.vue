<template>
  <div class="stats-container" v-loading="loading">
    <h2>分析统计</h2>

    <div class="overview-cards" v-if="overview">
      <StatCard title="总学习时长" :value="overview.total_learning_hours" unit="小时" icon-bg="#eff6ff" icon-color="#2563eb" />
      <StatCard title="总提问数" :value="overview.total_questions_asked" icon-bg="#f0fdf4" icon-color="#10b981" />
      <StatCard title="总测验数" :value="overview.total_quizzes_taken" icon-bg="#fef3c7" icon-color="#f59e0b" />
      <StatCard title="平均正确率" :value="overview.average_correct_rate" unit="%" icon-bg="#ede9fe" icon-color="#8b5cf6" />
      <StatCard title="知识点掌握" :value="overview.knowledge_points_mastered" :suffix="'/' + overview.knowledge_points_total" icon-bg="#cffafe" icon-color="#06b6d4" />
    </div>

    <div class="chart-grid">
      <MasteryChart :data="masteryList" />
      <HotPointsRank :data="hotPoints" />
    </div>

    <TrendChart :data="trendData" />

    <WeakPointsList :data="weakPoints" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { OverviewStats, KnowledgeMastery, HotKnowledgePoint, WeakPoint, TrendData } from '@/types/stats'
import { getOverview, getKnowledgeMastery, getWeakPoints, getHotKnowledgePoints, getTrends } from '@/services/stats'
import StatCard from '@/components/common/StatCard.vue'
import MasteryChart from './components/MasteryChart.vue'
import TrendChart from './components/TrendChart.vue'
import HotPointsRank from './components/HotPointsRank.vue'
import WeakPointsList from './components/WeakPointsList.vue'

const loading = ref(true)
const overview = ref<OverviewStats | null>(null)
const masteryList = ref<KnowledgeMastery[]>([])
const hotPoints = ref<HotKnowledgePoint[]>([])
const weakPoints = ref<WeakPoint[]>([])
const trendData = ref<TrendData | null>(null)

const fetchData = async () => {
  loading.value = true
  try {
    const results = await Promise.allSettled([
      getOverview(),
      getKnowledgeMastery(),
      getWeakPoints(10),
      getHotKnowledgePoints(10),
      getTrends(7)
    ])
    if (results[0].status === 'fulfilled') overview.value = results[0].value.data
    if (results[1].status === 'fulfilled') masteryList.value = results[1].value.data
    if (results[2].status === 'fulfilled') weakPoints.value = results[2].value.data
    if (results[3].status === 'fulfilled') hotPoints.value = results[3].value.data
    if (results[4].status === 'fulfilled') trendData.value = results[4].value.data
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
.stats-container h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 20px;
}
.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}
.chart-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}
@media (max-width: 768px) {
  .chart-grid {
    grid-template-columns: 1fr;
  }
}
</style>
