<template>
  <div class="stat-card">
    <div class="stat-icon" :style="{ background: iconBg }">
      <el-icon :size="20" :color="iconColor"><component :is="icon" /></el-icon>
    </div>
    <div class="stat-content">
      <div class="stat-value">
        {{ prefix }}{{ displayValue }}{{ suffix }}
      </div>
      <div class="stat-label">{{ title }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  title: string
  value: number | null | undefined
  unit?: string
  suffix?: string
  prefix?: string
  icon?: any
  iconBg?: string
  iconColor?: string
  precision?: number
}>(), {
  unit: '',
  suffix: '',
  prefix: '',
  iconBg: '#eff6ff',
  iconColor: '#2563eb',
  precision: 0
})

const displayValue = computed(() => {
  if (props.value == null) return '--'
  if (props.unit === '%') return (props.value * 100).toFixed(props.precision)
  return props.value.toFixed(props.precision)
})
</script>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
}
.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}
.stat-label {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 4px;
}
</style>
