<template>
  <div class="card">
    <h3>学习趋势</h3>
    <v-chart v-if="data && data.daily_stats && data.daily_stats.length" :option="option" autoresize style="height: 240px" />
    <EmptyState v-else text="暂无趋势数据" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { TrendData } from '@/types/stats'
import EmptyState from '@/components/common/EmptyState.vue'

use([LineChart, GridComponent, TooltipComponent, CanvasRenderer])

const props = defineProps<{
  data: TrendData | null
}>()

const option = computed(() => {
  if (!props.data) return {}
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: props.data.daily_stats.map(i => i.date) },
    yAxis: { type: 'value', name: '小时' },
    series: [{
      type: 'line',
      data: props.data.daily_stats.map(i => i.learning_hours),
      smooth: true,
      areaStyle: { color: 'rgba(37, 99, 235, 0.1)' },
      itemStyle: { color: '#3b82f6' }
    }]
  }
})
</script>

<style scoped>
.card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.card h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
</style>
