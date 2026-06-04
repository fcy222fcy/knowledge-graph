<template>
  <el-drawer v-model="visible" title="节点详情" size="400px">
    <template v-if="node">
      <div class="node-info">
        <h3>{{ node.name }}</h3>
        <p class="node-desc">{{ node.description || '暂无描述' }}</p>
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="ID">{{ node.id }}</el-descriptions-item>
          <el-descriptions-item label="分类">{{ node.category || '-' }}</el-descriptions-item>
          <el-descriptions-item label="关联文档ID">{{ node.document_id }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import type { GraphNode } from '@/types/graph'

const visible = defineModel<boolean>({ default: false })

defineProps<{
  node: GraphNode | null
}>()
</script>

<style scoped>
.node-info h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}
.node-desc {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 16px;
  line-height: 1.6;
}
</style>
