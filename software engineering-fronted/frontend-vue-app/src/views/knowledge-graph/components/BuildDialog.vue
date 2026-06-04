<template>
  <el-dialog v-model="visible" title="构建知识图谱" width="500px">
    <p class="tip">选择要用于构建图谱的文档：</p>
    <el-checkbox-group v-model="selectedIds" v-loading="docsLoading">
      <div v-for="doc in documents" :key="doc.id" class="doc-item">
        <el-checkbox :label="doc.id">{{ doc.title }}</el-checkbox>
        <el-tag size="small" :type="doc.status === 'completed' ? 'success' : 'info'">
          {{ doc.status === 'completed' ? '已解析' : doc.status }}
        </el-tag>
      </div>
    </el-checkbox-group>
    <el-empty v-if="!docsLoading && !documents.length" description="暂无已解析的文档" />
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="loading" :disabled="!selectedIds.length" @click="handleConfirm">
        开始构建
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { DocumentItem } from '@/types/document'
import { getDocumentList } from '@/services/document'

const visible = defineModel<boolean>({ default: false })
const emit = defineEmits<{ confirm: [documentIds: number[]] }>()

defineProps<{ loading: boolean }>()

const documents = ref<DocumentItem[]>([])
const selectedIds = ref<number[]>([])
const docsLoading = ref(false)

watch(visible, async (val) => {
  if (val) {
    docsLoading.value = true
    try {
      const result = await getDocumentList({ page: 1, size: 100, status: 'completed' })
      documents.value = result.data.list
    } catch {
      documents.value = []
    } finally {
      docsLoading.value = false
    }
  } else {
    selectedIds.value = []
  }
})

const handleConfirm = () => {
  emit('confirm', [...selectedIds.value])
}
</script>

<style scoped>
.tip {
  margin-bottom: 16px;
  color: var(--text-secondary);
  font-size: 14px;
}
.doc-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-light);
}
</style>
