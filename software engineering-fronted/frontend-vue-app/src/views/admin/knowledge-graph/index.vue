<template>
  <div class="admin-graph-container">
    <div class="admin-header">
      <h2>知识图谱管理</h2>
      <div class="header-actions">
        <el-button @click="fetchGraphData" :loading="loading">
          <el-icon><Refresh /></el-icon> 刷新
        </el-button>
        <el-button type="primary" @click="showBuildPanel = !showBuildPanel">
          <el-icon><Setting /></el-icon> 重建图谱
        </el-button>
      </div>
    </div>

    <!-- 重建图谱面板 -->
    <el-collapse-transition>
      <div v-show="showBuildPanel" class="build-panel">
        <el-card shadow="never">
          <template #header>
            <div class="panel-header">
              <span>选择文档重建图谱</span>
              <el-button text @click="showBuildPanel = false">
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
          </template>

          <div v-loading="docsLoading">
            <el-empty v-if="!docsLoading && !documents.length" description="暂无可用文档" />
            <el-checkbox-group v-model="selectedDocIds" v-else>
              <div v-for="doc in documents" :key="doc.id" class="doc-item">
                <el-checkbox :label="doc.id">
                  <span class="doc-title">{{ doc.title }}</span>
                </el-checkbox>
                <el-tag size="small" :type="doc.status === 'completed' ? 'success' : 'info'">
                  {{ doc.status === 'completed' ? '已解析' : doc.status }}
                </el-tag>
              </div>
            </el-checkbox-group>
          </div>

          <div class="build-actions">
            <el-checkbox
              v-model="selectAll"
              :indeterminate="isIndeterminate"
              @change="handleSelectAll"
              v-if="documents.length"
            >
              全选
            </el-checkbox>
            <el-button
              type="primary"
              :loading="buildLoading"
              :disabled="!selectedDocIds.length"
              @click="handleBuild"
            >
              开始重建 (已选 {{ selectedDocIds.length }} 个文档)
            </el-button>
          </div>
        </el-card>
      </div>
    </el-collapse-transition>

    <!-- 变化统计 -->
    <el-collapse-transition>
      <div v-if="buildStats" class="stats-panel">
        <el-card shadow="never">
          <template #header>
            <div class="panel-header">
              <span>最近重建结果</span>
              <el-button text @click="buildStats = null">
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
          </template>
          <el-descriptions :column="4" border size="small">
            <el-descriptions-item label="新增知识点">
              <el-tag type="success">+{{ buildStats.created_points }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="新增关系">
              <el-tag type="success">+{{ buildStats.created_relations }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="文档分块数">
              {{ buildStats.chunk_count }}
            </el-descriptions-item>
            <el-descriptions-item label="向量数量">
              {{ buildStats.vector_count }}
            </el-descriptions-item>
          </el-descriptions>
          <div class="build-message" v-if="buildStats.message">
            <el-tag :type="buildStats.status === 'success' ? 'success' : 'danger'" size="small">
              {{ buildStats.status === 'success' ? '成功' : '失败' }}
            </el-tag>
            <span>{{ buildStats.message }}</span>
          </div>
        </el-card>
      </div>
    </el-collapse-transition>

    <!-- 图谱概览 -->
    <div class="graph-summary" v-if="graphData">
      <el-tag type="info">节点: {{ graphData.summary.node_count }}</el-tag>
      <el-tag type="info">关系: {{ graphData.summary.edge_count }}</el-tag>
    </div>

    <!-- 图谱画布 -->
    <div class="graph-canvas-wrapper">
      <GraphCanvas
        :data="graphData"
        :loading="loading"
        :highlight="searchKeyword"
        @node-click="handleNodeClick"
      />
    </div>

    <!-- 节点详情 -->
    <NodeDetail v-model="detailVisible" :node="selectedNode" />

    <!-- 节点编辑 -->
    <NodeEditDialog
      v-model="editVisible"
      :node="editingNode"
      @success="fetchGraphData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Setting, Close } from '@element-plus/icons-vue'
import type { GraphNode, GraphBuildResult } from '@/services/graph'
import { getGraphData, buildGraph } from '@/services/graph'
import { getDocumentList } from '@/services/document'
import type { DocumentItem } from '@/services/document'
import type { GraphData } from '@/types/graph'
import GraphCanvas from '@/views/knowledge-graph/components/GraphCanvas.vue'
import NodeDetail from '@/views/knowledge-graph/components/NodeDetail.vue'
import NodeEditDialog from '@/views/knowledge-graph/components/NodeEditDialog.vue'

const graphData = ref<GraphData | null>(null)
const loading = ref(false)
const searchKeyword = ref('')
const selectedNode = ref<GraphNode | null>(null)
const detailVisible = ref(false)
const editVisible = ref(false)
const editingNode = ref<GraphNode | null>(null)

// 重建相关
const showBuildPanel = ref(false)
const buildLoading = ref(false)
const docsLoading = ref(false)
const documents = ref<DocumentItem[]>([])
const selectedDocIds = ref<number[]>([])
const buildStats = ref<GraphBuildResult | null>(null)

const selectAll = computed({
  get: () => selectedDocIds.value.length === documents.value.length && documents.value.length > 0,
  set: () => {}
})

const isIndeterminate = computed(() => {
  return selectedDocIds.value.length > 0 && selectedDocIds.value.length < documents.value.length
})

const fetchGraphData = async () => {
  loading.value = true
  try {
    const result = await getGraphData()
    graphData.value = result.data
  } catch (error) {
    console.error('获取图谱数据失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchDocuments = async () => {
  docsLoading.value = true
  try {
    const result = await getDocumentList({ page: 1, size: 100 })
    documents.value = result.data.list
  } catch (error) {
    console.error('获取文档列表失败:', error)
    documents.value = []
  } finally {
    docsLoading.value = false
  }
}

const handleSelectAll = (val: boolean | string | number) => {
  if (val) {
    selectedDocIds.value = documents.value.map(doc => doc.id)
  } else {
    selectedDocIds.value = []
  }
}

const handleBuild = async () => {
  if (!selectedDocIds.value.length) return

  buildLoading.value = true
  try {
    const result = await buildGraph(selectedDocIds.value)
    buildStats.value = result.data as GraphBuildResult
    ElMessage.success('图谱重建完成')
    showBuildPanel.value = false
    fetchGraphData()
  } catch (error) {
    ElMessage.error('图谱重建失败')
    console.error('构建图谱失败:', error)
  } finally {
    buildLoading.value = false
  }
}

const handleNodeClick = (node: GraphNode) => {
  selectedNode.value = node
  detailVisible.value = true
}

watch(showBuildPanel, (val) => {
  if (val && !documents.value.length) {
    fetchDocuments()
  }
})

onMounted(() => {
  fetchGraphData()
})
</script>

<style scoped>
.admin-graph-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #f5f7fa;
}

.admin-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  background: #fff;
  border-bottom: 1px solid #e2e8f0;
}

.admin-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #0f172a;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.build-panel,
.stats-panel {
  padding: 0 24px;
  background: #f5f7fa;
}

.build-panel .el-card,
.stats-panel .el-card {
  margin-bottom: 0;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.doc-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f1f5f9;
}

.doc-item:last-child {
  border-bottom: none;
}

.doc-title {
  color: #334155;
  font-size: 14px;
}

.build-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
}

.build-message {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e2e8f0;
  font-size: 13px;
  color: #64748b;
}

.graph-summary {
  display: flex;
  gap: 8px;
  padding: 12px 24px;
  background: #fff;
  border-bottom: 1px solid #f1f5f9;
}

.graph-canvas-wrapper {
  flex: 1;
  overflow: hidden;
  background: #fff;
}
</style>
