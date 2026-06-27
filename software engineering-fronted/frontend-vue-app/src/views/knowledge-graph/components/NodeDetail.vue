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
          <el-descriptions-item label="关联文档ID">{{ node.document_id }}</el-descriptions-item>
        </el-descriptions>
        <div class="node-actions">
          <el-button type="primary" class="action-btn" @click="$emit('edit', node)">
            <el-icon><Edit /></el-icon>
            编辑知识点
          </el-button>
          <el-button
            v-if="node.document_id"
            type="info"
            class="action-btn"
            @click="goToDocument"
          >
            <el-icon><Document /></el-icon>
            查看关联文档
          </el-button>
        </div>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Document, Edit } from '@element-plus/icons-vue'
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

const goToDocument = () => {
  if (props.node?.document_id) {
    // 跳转到文件列表页面，并传递文档ID参数以自动打开详情
    router.push({ path: '/files', query: { doc_id: String(props.node.document_id) } })
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
</style>
