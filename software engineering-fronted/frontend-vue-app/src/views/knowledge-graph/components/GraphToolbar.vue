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
    <div class="toolbar-right">
      <el-button @click="$emit('refresh')">
        <el-icon><Refresh /></el-icon> 刷新
      </el-button>
      <el-button type="primary" @click="$emit('upload')">
        <el-icon><Upload /></el-icon> 提交资料
      </el-button>
      <el-button type="primary" @click="$emit('build')">
        <el-icon><Setting /></el-icon> 构建图谱
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Refresh, Setting, Upload } from '@element-plus/icons-vue'

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
  upload: []
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
