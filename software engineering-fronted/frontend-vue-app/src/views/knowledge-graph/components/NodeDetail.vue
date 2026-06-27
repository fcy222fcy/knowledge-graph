<template>
  <el-drawer v-model="visible" title="节点详情" size="400px">
    <template v-if="node">
      <div class="node-info">
        <div class="node-header">
          <div class="node-dot" :style="{ background: getNodeColor(node.category) }"></div>
          <h3>{{ node.name }}</h3>
        </div>
        <p class="node-desc">{{ node.description || '暂无描述' }}</p>
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="分类">{{ node.category || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- 来源文档列表 -->
        <div class="sources-section" v-if="node.sources && node.sources.length > 0">
          <h4 class="sources-title">📚 来源文档 ({{ node.sources.length }})</h4>
          <div class="sources-list">
            <div
              v-for="(source, index) in node.sources"
              :key="source.document_id"
              class="source-item"
              @click="goToDocument(source.document_id)"
            >
              <span class="source-index">来源{{ index + 1 }}</span>
              <span class="source-name">{{ source.document_title || '未知文档' }}</span>
              <el-icon class="source-arrow"><ArrowRight /></el-icon>
            </div>
          </div>
        </div>

        <div class="node-actions">
          <el-button type="primary" class="action-btn" @click="$emit('edit', node)">
            <el-icon><Edit /></el-icon>
            编辑知识点
          </el-button>
        </div>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Document, Edit, ArrowRight } from '@element-plus/icons-vue'
import type { GraphNode } from '@/types/graph'

const visible = defineModel<boolean>({ default: false })
const router = useRouter()

const props = defineProps<{
  node: GraphNode | null
}>()

defineEmits<{
  (e: 'edit', node: GraphNode): void
}>()

const COLORS: Record<string, string> = {
  '软件工程': '#3b82f6',
  'SQL': '#10b981',
  'Redis': '#f59e0b',
  'AI': '#8b5cf6'
}

const getNodeColor = (category?: string) => {
  return COLORS[category || ''] || '#3b82f6'
}

const goToDocument = (docId?: number) => {
  const documentId = docId || props.node?.document_id
  if (documentId) {
    // 跳转到文件列表页面，并传递文档ID参数以自动打开详情
    router.push({ path: '/files', query: { doc_id: String(documentId) } })
  }
}
</script>

<style scoped>
.node-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}
.node-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}
.node-info h3 {
  font-size: 18px;
  font-weight: 600;
  color: #0f172a;
  margin: 0;
}
.node-desc {
  font-size: 14px;
  color: #64748b;
  margin-bottom: 16px;
  line-height: 1.6;
}
.node-actions {
  margin-top: 16px;
  display: flex;
  gap: 12px;
}
.action-btn {
  flex: 1;
}

/* 来源文档样式 */
.sources-section {
  margin: 16px 0;
  padding: 12px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.sources-title {
  font-size: 14px;
  font-weight: 600;
  color: #334155;
  margin: 0 0 12px 0;
}

.sources-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.source-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: white;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  cursor: pointer;
  transition: all 0.2s;
}

.source-item:hover {
  background: #eff6ff;
  border-color: #3b82f6;
}

.source-index {
  font-size: 12px;
  font-weight: 600;
  color: #3b82f6;
  background: #dbeafe;
  padding: 2px 8px;
  border-radius: 4px;
}

.source-name {
  flex: 1;
  font-size: 13px;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.source-arrow {
  color: #94a3b8;
  font-size: 14px;
}
</style>
