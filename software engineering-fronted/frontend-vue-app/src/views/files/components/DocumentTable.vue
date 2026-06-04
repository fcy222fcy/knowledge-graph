<template>
  <el-table :data="data" v-loading="loading" stripe style="width: 100%">
    <el-table-column prop="title" label="文档名称" min-width="200" show-overflow-tooltip />
    <el-table-column prop="file_type" label="类型" width="100">
      <template #default="{ row }">
        <el-tag size="small" :type="getFileTagType(row.file_type)">{{ row.file_type }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="file_size" label="大小" width="100">
      <template #default="{ row }">
        {{ formatFileSize(row.file_size) }}
      </template>
    </el-table-column>
    <el-table-column prop="status" label="状态" width="100">
      <template #default="{ row }">
        <el-tag size="small" :type="getStatusType(row.status)">
          {{ DOCUMENT_STATUS_MAP[row.status as keyof typeof DOCUMENT_STATUS_MAP]?.label || row.status }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="created_at" label="上传时间" width="180">
      <template #default="{ row }">
        {{ formatDate(row.created_at) }}
      </template>
    </el-table-column>
    <el-table-column label="操作" width="160" fixed="right">
      <template #default="{ row }">
        <el-button type="primary" link size="small" @click="$emit('view', row)">详情</el-button>
        <el-button type="danger" link size="small" @click="$emit('delete', row)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import type { DocumentItem } from '@/types/document'
import { DOCUMENT_STATUS_MAP } from '@/types/document'
import { formatDate } from '@/utils'

defineProps<{
  data: DocumentItem[]
  loading: boolean
}>()

defineEmits<{
  view: [row: DocumentItem]
  delete: [row: DocumentItem]
}>()

const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const getFileTagType = (type: string) => {
  const map: Record<string, string> = {
    '.pdf': 'danger',
    '.docx': 'primary',
    '.pptx': 'warning',
    '.md': 'success',
    '.txt': 'info'
  }
  return map[type] || 'info'
}

const getStatusType = (status: string) => {
  return (DOCUMENT_STATUS_MAP[status as keyof typeof DOCUMENT_STATUS_MAP]?.color || 'info') as any
}
</script>
