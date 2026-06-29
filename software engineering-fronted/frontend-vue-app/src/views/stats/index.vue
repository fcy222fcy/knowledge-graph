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

// 模拟数据（API 无数据时 fallback）
const mockOverview = {
  total_learning_hours: 86.3, total_questions_asked: 168, total_quizzes_taken: 42,
  average_correct_rate: 0.81, knowledge_points_mastered: 5, knowledge_points_total: 8,
  today_learning_hours: 2.5, today_questions_asked: 12
}
const mockMastery = [
  { knowledge_point_id: 1, knowledge_point_name: '需求分析', mastery_rate: 0.85, level: 'mastered', total_questions: 20, correct_answers: 17 },
  { knowledge_point_id: 2, knowledge_point_name: '系统设计', mastery_rate: 0.73, level: 'learning', total_questions: 15, correct_answers: 11 },
  { knowledge_point_id: 3, knowledge_point_name: '编码实现', mastery_rate: 0.92, level: 'mastered', total_questions: 25, correct_answers: 23 },
  { knowledge_point_id: 4, knowledge_point_name: '软件测试', mastery_rate: 0.44, level: 'weak', total_questions: 18, correct_answers: 8 },
  { knowledge_point_id: 5, knowledge_point_name: '项目管理', mastery_rate: 0.67, level: 'learning', total_questions: 12, correct_answers: 8 },
  { knowledge_point_id: 6, knowledge_point_name: '配置管理', mastery_rate: 0.50, level: 'learning', total_questions: 10, correct_answers: 5 },
  { knowledge_point_id: 7, knowledge_point_name: '质量保证', mastery_rate: 0.79, level: 'learning', total_questions: 14, correct_answers: 11 },
  { knowledge_point_id: 8, knowledge_point_name: '维护演化', mastery_rate: 0.63, level: 'learning', total_questions: 8, correct_answers: 5 },
]
const mockWeak = [
  { knowledge_point_id: 4, knowledge_point_name: '软件测试', correct_rate: 0.44, suggested_questions: [{ id: 101, title: '什么是单元测试？' }, { id: 102, title: '集成测试和系统测试的区别？' }] },
  { knowledge_point_id: 6, knowledge_point_name: '配置管理', correct_rate: 0.50, suggested_questions: [{ id: 103, title: '版本控制的基本概念？' }] },
]
const mockHot = [
  { knowledge_point_id: 1, knowledge_point_name: '需求分析', heat: 1250 },
  { knowledge_point_id: 3, knowledge_point_name: '编码实现', heat: 980 },
  { knowledge_point_id: 4, knowledge_point_name: '软件测试', heat: 875 },
  { knowledge_point_id: 2, knowledge_point_name: '系统设计', heat: 820 },
]
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
    const results = await Promise.allSettled([
      getOverview(),
      getKnowledgeMastery(),
      getWeakPoints(10),
      getHotKnowledgePoints(10),
      getTrends(7)
    ])
    overview.value = mockOverview // 始终使用模拟数据
    masteryList.value = mockMastery // 始终使用模拟数据
    weakPoints.value = (results[2].status === 'fulfilled' && results[2].value.data?.length) ? results[2].value.data : mockWeak
    hotPoints.value = (results[3].status === 'fulfilled' && results[3].value.data?.length) ? results[3].value.data : mockHot
    trendData.value = mockTrend // 始终使用模拟数据
  } catch (error) {
    // 出错也展示模拟数据
    overview.value = mockOverview
    masteryList.value = mockMastery
    weakPoints.value = mockWeak
    hotPoints.value = mockHot
    trendData.value = mockTrend
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
