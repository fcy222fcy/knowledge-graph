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
      <el-option label="关联" value="RELATED" />
      <el-option label="依赖" value="DEPENDS_ON" />
      <el-option label="组成部分" value="PART_OF" />
      <el-option label="是一种" value="IS_A" />
      <el-option label="示例" value="EXAMPLE_OF" />
      <el-option label="使用" value="USES" />
      <el-option label="实现" value="IMPLEMENTS" />
    </el-select>
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
