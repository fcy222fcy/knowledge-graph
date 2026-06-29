<template>
  <div class="source-reference" v-if="sources.length || relatedPoints.length">
    <div class="ref-section" v-if="sources.length">
      <h4>📚 知识来源</h4>
      <div v-for="source in sources" :key="source.document_id" class="source-card">
        <div class="source-title">{{ source.document_title }}</div>
        <div class="source-content">{{ source.content.slice(0, 200) }}...</div>
      </div>
    </div>
    <div class="ref-section" v-if="relatedPoints.length">
      <h4>🔗 相关知识点</h4>
      <div class="tags">
        <el-tag v-for="point in relatedPoints" :key="point.id" size="small" type="info">
          {{ point.name }}
        </el-tag>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AskResponse } from '@/services/qa'

defineProps<{
  sources: AskResponse['sources']
  relatedPoints: AskResponse['related_knowledge_points']
}>()
</script>

<style scoped>
.source-reference {
  margin-top: 12px;
  padding: 12px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid var(--border-light);
}
.ref-section {
  margin-bottom: 12px;
}
.ref-section:last-child {
  margin-bottom: 0;
}
.ref-section h4 {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 8px;
}
.source-card {
  padding: 8px 12px;
  background: white;
  border-radius: 6px;
  margin-bottom: 6px;
}
.source-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--primary);
  margin-bottom: 4px;
}
.source-content {
  font-size: 12px;
  color: var(--text-muted);
  line-height: 1.5;
}
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.no-source {
  font-size: 12px;
  color: #999;
  padding: 8px 0;
}
</style>
