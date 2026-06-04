<template>
  <div class="overview-cards" v-if="data">
    <StatCard title="今日学习时长" :value="data.today_learning_hours" unit="小时"
      icon-bg="#eff6ff" icon-color="#2563eb" />
    <StatCard title="今日提问数" :value="data.today_questions_asked"
      icon-bg="#f0fdf4" icon-color="#10b981" />
    <StatCard title="平均正确率" :value="data.average_correct_rate" unit="%"
      icon-bg="#ede9fe" icon-color="#8b5cf6" />
    <StatCard title="知识点掌握" :value="data.knowledge_points_mastered"
      :suffix="'/' + data.knowledge_points_total" icon-bg="#cffafe" icon-color="#06b6d4" />
  </div>
</template>

<script setup lang="ts">
import type { OverviewStats } from '@/types/stats'
import StatCard from '@/components/common/StatCard.vue'

defineProps<{
  data: OverviewStats | null
}>()
</script>

<style scoped>
.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}
</style>
