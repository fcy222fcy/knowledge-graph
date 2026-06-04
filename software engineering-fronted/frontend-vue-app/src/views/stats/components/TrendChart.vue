<template>
  <div class="chart-card">
    <h3>学习趋势</h3>
    <v-chart v-if="data && data.daily_stats.length" :option="option" autoresize style="height: 300px" />
    <EmptyState v-else text="暂无趋势数据" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { TrendData } from '@/types/stats'
import EmptyState from '@/components/common/EmptyState.vue'

use([LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const props = defineProps<{
  data: TrendData | null
}>()

const option = computed(() => {
  if (!props.data) return {}
  const dates = props.data.daily_stats.map(i => i.date)
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['学习时长', '正确率'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: dates },
    yAxis: [
      { type: 'value', name: '时长(h)', axisLabel: { formatter: '{value}' } },
      { type: 'value', name: '正确率(%)', max: 100, axisLabel: { formatter: '{value}%' } }
    ],
    series: [
      {
        name: '学习时长',
        type: 'line',
        data: props.data.daily_stats.map(i => i.learning_hours),
        smooth: true,
        itemStyle: { color: '#3b82f6' }
      },
      {
        name: '正确率',
        type: 'line',
        yAxisIndex: 1,
        data: props.data.daily_stats.map(i => Math.round(i.correct_rate * 100)),
        smooth: true,
        itemStyle: { color: '#10b981' }
      }
    ]
  }
})
</script>

<style scoped>
.chart-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.chart-card h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
</style>
