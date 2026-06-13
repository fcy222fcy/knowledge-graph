<template>
  <div class="chart-card">
    <h3>知识点掌握度</h3>
    <v-chart v-if="data && data.length" :option="option" autoresize style="height: 300px" />
    <EmptyState v-else text="暂无掌握度数据" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { KnowledgeMastery } from '@/types/stats'
import EmptyState from '@/components/common/EmptyState.vue'

use([BarChart, GridComponent, TooltipComponent, CanvasRenderer])

const props = defineProps<{
  data: KnowledgeMastery[]
}>()

const getLevelColor = (level: string) => {
  const map: Record<string, string> = {
    mastered: '#10b981',
    learning: '#3b82f6',
    weak: '#f43f5e'
  }
  return map[level] || '#94a3b8'
}

const option = computed(() => ({
  tooltip: {
    trigger: 'axis',
    formatter: (params: any) => {
      const item = props.data[params[0]?.dataIndex]
      if (!item) return ''
      return `${item.knowledge_point_name}<br/>掌握度: ${(item.mastery_rate * 100).toFixed(0)}%<br/>总题数: ${item.total_questions}<br/>正确: ${item.correct_answers}`
    }
  },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: {
    type: 'category',
    data: props.data.map(i => i.knowledge_point_name),
    axisLabel: { rotate: 30, fontSize: 11 }
  },
  yAxis: {
    type: 'value',
    max: 100,
    axisLabel: { formatter: '{value}%' }
  },
  series: [{
    type: 'bar',
    data: props.data.map(i => ({
      value: Math.round(i.mastery_rate * 100),
      itemStyle: { color: getLevelColor(i.level), borderRadius: [4, 4, 0, 0] }
    })),
    barMaxWidth: 40
  }]
}))
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
