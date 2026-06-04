<template>
  <div class="graph-container">
    <GraphToolbar
      v-model:keyword="searchKeyword"
      v-model:relationType="filterRelationType"
      @search="handleSearch"
      @build="buildVisible = true"
      @refresh="fetchGraphData"
    />

    <div class="graph-summary" v-if="graphData">
      <el-tag type="info">节点: {{ graphData.summary.node_count }}</el-tag>
      <el-tag type="info">关系: {{ graphData.summary.edge_count }}</el-tag>
    </div>

    <GraphCanvas
      :data="graphData"
      :loading="loading"
      :highlight="searchKeyword"
      @node-click="handleNodeClick"
    />

    <NodeDetail v-model="detailVisible" :node="selectedNode" />
    <BuildDialog v-model="buildVisible" :loading="buildLoading" @confirm="handleBuild" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { GraphData, GraphNode } from '@/types/graph'
import { getGraphData, buildGraph } from '@/services/graph'
import { debounce } from '@/utils'
import GraphToolbar from './components/GraphToolbar.vue'
import GraphCanvas from './components/GraphCanvas.vue'
import NodeDetail from './components/NodeDetail.vue'
import BuildDialog from './components/BuildDialog.vue'

const graphData = ref<GraphData | null>(null)
const loading = ref(false)
const searchKeyword = ref('')
const filterRelationType = ref('')
const selectedNode = ref<GraphNode | null>(null)
const detailVisible = ref(false)
const buildVisible = ref(false)
const buildLoading = ref(false)

const fetchGraphData = async () => {
  loading.value = true
  try {
    const result = await getGraphData({
      keyword: searchKeyword.value || undefined,
      relation_type: filterRelationType.value || undefined
    })
    graphData.value = result.data
  } catch (error) {
    console.error('获取图谱数据失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = debounce(() => {
  fetchGraphData()
}, 300)

const handleNodeClick = (node: GraphNode) => {
  selectedNode.value = node
  detailVisible.value = true
}

const handleBuild = async (documentIds: number[]) => {
  buildLoading.value = true
  try {
    await buildGraph(documentIds)
    ElMessage.success('图谱构建成功')
    buildVisible.value = false
    fetchGraphData()
  } catch (error) {
    console.error('构建图谱失败:', error)
  } finally {
    buildLoading.value = false
  }
}

onMounted(() => {
  fetchGraphData()
})
</script>

<style scoped>
.graph-container h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
.graph-summary {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
</style>
