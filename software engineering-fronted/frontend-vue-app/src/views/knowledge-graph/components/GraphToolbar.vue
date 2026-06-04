<template>
  <div class="graph-toolbar">
    <el-input
      :model-value="keyword"
      @update:model-value="$emit('update:keyword', $event)"
      placeholder="搜索知识点..."
      prefix-icon="Search"
      clearable
      style="width: 260px"
      @keyup.enter="$emit('search')"
      @clear="$emit('search')"
    />
    <el-select
      :model-value="relationType"
      @update:model-value="$emit('update:relationType', $event)"
      placeholder="关系类型"
      clearable
      style="width: 160px"
      @change="$emit('search')"
    >
      <el-option label="全部关系" value="" />
      <el-option label="RELATED" value="RELATED" />
      <el-option label="DEPENDS_ON" value="DEPENDS_ON" />
      <el-option label="PART_OF" value="PART_OF" />
    </el-select>
    <div class="toolbar-right">
      <el-button @click="$emit('refresh')">
        <el-icon><Refresh /></el-icon> 刷新
      </el-button>
      <el-button type="primary" @click="$emit('build')">
        <el-icon><Setting /></el-icon> 构建图谱
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Refresh, Setting } from '@element-plus/icons-vue'

defineProps<{
  keyword: string
  relationType: string
}>()

defineEmits<{
  'update:keyword': [value: string]
  'update:relationType': [value: string]
  search: []
  build: []
  refresh: []
}>()
</script>

<style scoped>
.graph-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}
.toolbar-right {
  margin-left: auto;
  display: flex;
  gap: 8px;
}
</style>
