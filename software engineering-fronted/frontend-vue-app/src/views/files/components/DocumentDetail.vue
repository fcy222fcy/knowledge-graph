<template>
  <el-drawer v-model="visible" title="文档详情" size="500px">
    <template v-if="document">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="标题">{{ document.title }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ document.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="文件类型">{{ document.file_type }}</el-descriptions-item>
        <el-descriptions-item label="文件大小">{{ formatFileSize(document.file_size) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="(DOCUMENT_STATUS_MAP[document.status as keyof typeof DOCUMENT_STATUS_MAP]?.color || 'info') as any">
            {{ DOCUMENT_STATUS_MAP[document.status as keyof typeof DOCUMENT_STATUS_MAP]?.label || document.status }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="上传时间">{{ formatDate(document.created_at) }}</el-descriptions-item>
      </el-descriptions>

      <div class="content-preview" v-if="content">
        <h4>文档内容预览</h4>
        <pre class="content-text">{{ content }}</pre>
      </div>
    </template>
    <el-skeleton :rows="5" animated v-else-if="loading" />
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { DocumentItem } from '@/types/document'
import { DOCUMENT_STATUS_MAP } from '@/types/document'
import { getDocumentContent } from '@/services/document'
import { formatDate } from '@/utils'

const visible = defineModel<boolean>({ default: false })

const props = defineProps<{
  document: DocumentItem | null
}>()

const content = ref('')
const loading = ref(false)

const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

watch(() => props.document, async (doc) => {
  if (doc && visible.value) {
    loading.value = true
    try {
      const result = await getDocumentContent(doc.id)
      content.value = result.data.content
    } catch {
      content.value = ''
    } finally {
      loading.value = false
    }
  }
})

watch(visible, (val) => {
  if (!val) {
    content.value = ''
  }
})
</script>

<style scoped>
.content-preview {
  margin-top: 24px;
}
.content-preview h4 {
  margin-bottom: 12px;
  color: var(--text-primary);
}
.content-text {
  background: #f8fafc;
  padding: 16px;
  border-radius: 8px;
  font-size: 13px;
  line-height: 1.8;
  max-height: 400px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
