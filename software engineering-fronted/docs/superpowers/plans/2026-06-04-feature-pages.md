# SE智图问答 - 功能页面开发计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现5个核心功能页面：资料管理、分析统计、首页、知识图谱、问答中心

**Architecture:** Vue 3 Composition API + TypeScript + Element Plus，ECharts做图表可视化，D3.js做知识图谱力导向图

**Tech Stack:** Vue 3, TypeScript, Element Plus, ECharts, vue-echarts, D3.js

---

## 实现顺序（按依赖和复杂度递增）

| 阶段 | 页面 | 理由 |
|------|------|------|
| 1 | 资料管理 (Files) | 独立页面，纯CRUD+上传，是图谱构建的前置数据源 |
| 2 | 分析统计 (Stats) | 独立只读页面，引入图表库，为首页复用做铺垫 |
| 3 | 首页 (Dashboard) | 聚合页，复用Stats图表组件 |
| 4 | 知识图谱 (KnowledgeGraph) | 引入D3可视化库，与资料管理联动 |
| 5 | 问答中心 (QA) | 最复杂，涉及会话管理和聊天交互 |

---

## 前置准备：安装额外依赖

```bash
cd frontend-vue-app
npm install echarts vue-echarts d3 @types/d3
```

- **ECharts + vue-echarts**：图表库，`<v-chart>` 组件封装Vue 3生命周期
- **D3.js**：知识图谱力导向图可视化，`d3-force` 专为此设计

---

## 公共基础设施（先创建）

### Step 1: 类型定义文件

- [ ] **src/types/document.ts**

```typescript
import type { DocumentItem, DocumentListParams, PaginatedResponse } from '@/services/document'

export type { DocumentItem, DocumentListParams, PaginatedResponse }

export type DocumentStatus = 'pending' | 'processing' | 'completed' | 'failed'

export const DOCUMENT_STATUS_MAP: Record<DocumentStatus, { label: string; color: string }> = {
  pending: { label: '待处理', color: 'info' },
  processing: { label: '处理中', color: 'warning' },
  completed: { label: '已完成', color: 'success' },
  failed: { label: '失败', color: 'danger' },
}
```

- [ ] **src/types/stats.ts**

```typescript
export type { OverviewStats, HotKnowledgePoint, KnowledgeMastery, WeakPoint, TrendData } from '@/services/stats'
```

- [ ] **src/types/graph.ts**

```typescript
export type { GraphNode, GraphEdge, GraphData, GraphBuildResult } from '@/services/graph'
```

- [ ] **src/types/qa.ts**

```typescript
export type { AskParams, AskResponse, ConversationItem, MessageItem } from '@/services/qa'
```

### Step 2: 公共组件

- [ ] **src/components/common/EmptyState.vue**

```vue
<template>
  <div class="empty-state">
    <el-icon :size="48" color="#c0c4cc"><component :is="icon" /></el-icon>
    <p class="empty-text">{{ text }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
import { Document } from '@element-plus/icons-vue'

defineProps<{
  text?: string
  icon?: any
}>()

withDefaults(defineProps<{ text?: string; icon?: any }>(), {
  text: '暂无数据',
  icon: Document
})
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--text-muted);
}
.empty-text {
  margin-top: 12px;
  font-size: 14px;
}
</style>
```

- [ ] **src/components/common/StatCard.vue**

```vue
<template>
  <div class="stat-card">
    <div class="stat-icon" :style="{ background: iconBg }">
      <el-icon :size="20" :color="iconColor"><component :is="icon" /></el-icon>
    </div>
    <div class="stat-content">
      <div class="stat-value">
        {{ prefix }}{{ displayValue }}{{ suffix }}
      </div>
      <div class="stat-label">{{ title }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  title: string
  value: number | null | undefined
  unit?: string
  suffix?: string
  prefix?: string
  icon?: any
  iconBg?: string
  iconColor?: string
  precision?: number
}>(), {
  unit: '',
  suffix: '',
  prefix: '',
  iconBg: '#eff6ff',
  iconColor: '#2563eb',
  precision: 0
})

const displayValue = computed(() => {
  if (props.value == null) return '--'
  if (props.unit === '%') return (props.value * 100).toFixed(props.precision)
  return props.value.toFixed(props.precision)
})
</script>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
}
.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}
.stat-label {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 4px;
}
</style>
```

### Step 3: 组合式函数

- [ ] **src/composables/usePagination.ts**

```typescript
import { ref, reactive } from 'vue'

export function usePagination(defaultSize = 10) {
  const pagination = reactive({
    page: 1,
    size: defaultSize,
    total: 0
  })

  const handlePageChange = (page: number) => {
    pagination.page = page
  }

  const handleSizeChange = (size: number) => {
    pagination.size = size
    pagination.page = 1
  }

  const resetPage = () => {
    pagination.page = 1
  }

  return {
    pagination,
    handlePageChange,
    handleSizeChange,
    resetPage
  }
}
```

---

## Task 1: 资料管理页面

**Files:**
- Rewrite: `src/views/files/index.vue`
- Create: `src/views/files/components/DocumentTable.vue`
- Create: `src/views/files/components/UploadDialog.vue`
- Create: `src/views/files/components/DocumentDetail.vue`

### Step 1: 创建 DocumentTable.vue

```vue
<!-- src/views/files/components/DocumentTable.vue -->
<template>
  <el-table :data="data" v-loading="loading" stripe style="width: 100%">
    <el-table-column prop="title" label="文档名称" min-width="200" show-overflow-tooltip />
    <el-table-column prop="file_type" label="类型" width="100">
      <template #default="{ row }">
        <el-tag size="small" :type=" getFileTagType(row.file_type)">{{ row.file_type }}</el-tag>
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
          {{ DOCUMENT_STATUS_MAP[row.status as DocumentStatus]?.label || row.status }}
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
import type { DocumentStatus } from '@/types/document'
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
  return (DOCUMENT_STATUS_MAP[status as DocumentStatus]?.color || 'info') as any
}
</script>
```

### Step 2: 创建 UploadDialog.vue

```vue
<!-- src/views/files/components/UploadDialog.vue -->
<template>
  <el-dialog v-model="visible" title="上传文档" width="500px" destroy-on-close>
    <el-form label-width="80px">
      <el-form-item label="选择文件">
        <el-upload
          ref="uploadRef"
          drag
          :auto-upload="false"
          :limit="1"
          accept=".pdf,.docx,.pptx,.md,.txt"
          :on-change="handleFileChange"
          :on-exceed="handleExceed"
        >
          <el-icon class="el-icon--upload" :size="40"><UploadFilled /></el-icon>
          <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
          <template #tip>
            <div class="el-upload__tip">支持 .pdf, .docx, .pptx, .md, .txt，单文件最大 10MB</div>
          </template>
        </el-upload>
      </el-form-item>
      <el-form-item label="文档标题">
        <el-input v-model="title" placeholder="不填则使用文件名" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="description" type="textarea" :rows="2" placeholder="可选" />
      </el-form-item>
    </el-form>
    <el-progress v-if="uploading" :percentage="uploadProgress" :status="uploadStatus" />
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="uploading" @click="handleUpload">上传</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, genFileId } from 'element-plus'
import type { UploadFile, UploadInstance } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import { uploadDocument } from '@/services/document'

const visible = defineModel<boolean>({ default: false })
const emit = defineEmits<{ success: [] }>()

const uploadRef = ref<UploadInstance>()
const selectedFile = ref<File | null>(null)
const title = ref('')
const description = ref('')
const uploading = ref(false)
const uploadProgress = ref(0)

const uploadStatus = computed(() => uploadProgress.value >= 100 ? 'success' : undefined)

const handleFileChange = (file: UploadFile) => {
  selectedFile.value = file.raw || null
}

const handleExceed = (files: File[]) => {
  uploadRef.value?.clearFiles()
  const file = files[0]
  file.uid = genFileId()
  uploadRef.value?.handleStart(file)
}

const handleUpload = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择文件')
    return
  }

  const maxSize = 10 * 1024 * 1024
  if (selectedFile.value.size > maxSize) {
    ElMessage.error('文件大小不能超过 10MB')
    return
  }

  uploading.value = true
  uploadProgress.value = 0

  try {
    // 模拟进度
    const timer = setInterval(() => {
      if (uploadProgress.value < 90) uploadProgress.value += 10
    }, 200)

    await uploadDocument(selectedFile.value, title.value || undefined, description.value || undefined)

    clearInterval(timer)
    uploadProgress.value = 100
    ElMessage.success('上传成功')
    visible.value = false
    emit('success')
  } catch (error) {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}
</script>
```

### Step 3: 创建 DocumentDetail.vue

```vue
<!-- src/views/files/components/DocumentDetail.vue -->
<template>
  <el-drawer v-model="visible" title="文档详情" size="500px">
    <template v-if="document">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="标题">{{ document.title }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ document.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="文件类型">{{ document.file_type }}</el-descriptions-item>
        <el-descriptions-item label="文件大小">{{ formatFileSize(document.file_size) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="(DOCUMENT_STATUS_MAP[document.status as DocumentStatus]?.color || 'info') as any">
            {{ DOCUMENT_STATUS_MAP[document.status as DocumentStatus]?.label || document.status }}
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
import type { DocumentItem, DocumentStatus } from '@/types/document'
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
```

### Step 4: 重写 index.vue

```vue
<!-- src/views/files/index.vue -->
<template>
  <div class="files-container">
    <div class="page-header">
      <h2>资料管理</h2>
      <el-button type="primary" @click="uploadVisible = true">
        <el-icon><Upload /></el-icon> 上传文档
      </el-button>
    </div>

    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索文档..."
        prefix-icon="Search"
        clearable
        style="width: 300px"
        @keyup.enter="handleSearch"
        @clear="handleSearch"
      />
      <el-select v-model="filterStatus" placeholder="状态筛选" clearable style="width: 140px" @change="handleSearch">
        <el-option label="待处理" value="pending" />
        <el-option label="处理中" value="processing" />
        <el-option label="已完成" value="completed" />
        <el-option label="失败" value="failed" />
      </el-select>
    </div>

    <DocumentTable
      :data="documentList"
      :loading="loading"
      @view="handleViewDetail"
      @delete="handleDelete"
    />

    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @current-change="fetchDocuments"
        @size-change="handleSizeChange"
      />
    </div>

    <UploadDialog v-model="uploadVisible" @success="fetchDocuments" />
    <DocumentDetail v-model="detailVisible" :document="currentDocument" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import { getDocumentList, deleteDocument } from '@/services/document'
import type { DocumentItem } from '@/types/document'
import { usePagination } from '@/composables/usePagination'
import DocumentTable from './components/DocumentTable.vue'
import UploadDialog from './components/UploadDialog.vue'
import DocumentDetail from './components/DocumentDetail.vue'

const documentList = ref<DocumentItem[]>([])
const loading = ref(false)
const searchKeyword = ref('')
const filterStatus = ref('')
const uploadVisible = ref(false)
const detailVisible = ref(false)
const currentDocument = ref<DocumentItem | null>(null)

const { pagination, handleSizeChange, resetPage } = usePagination()

const fetchDocuments = async () => {
  loading.value = true
  try {
    const result = await getDocumentList({
      page: pagination.page,
      size: pagination.size,
      keyword: searchKeyword.value || undefined,
      status: filterStatus.value || undefined
    })
    documentList.value = result.data.list
    pagination.total = result.data.total
  } catch (error) {
    console.error('获取文档列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  resetPage()
  fetchDocuments()
}

const handleViewDetail = (row: DocumentItem) => {
  currentDocument.value = row
  detailVisible.value = true
}

const handleDelete = async (row: DocumentItem) => {
  try {
    await ElMessageBox.confirm(`确定要删除文档「${row.title}」吗？`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteDocument(row.id)
    ElMessage.success('删除成功')
    fetchDocuments()
  } catch {
    // 用户取消
  }
}

onMounted(() => {
  fetchDocuments()
})
</script>

<style scoped>
.files-container {
  padding: 0;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}
.search-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
```

### Step 5: 验证

```bash
cd frontend-vue-app
npx vue-tsc --noEmit
npx vite build
```

### Step 6: Commit

```bash
git add src/types/ src/components/common/ src/composables/ src/views/files/
git commit -m "feat: 实现资料管理页面"
```

---

## Task 2: 分析统计页面

**Files:**
- Rewrite: `src/views/stats/index.vue`
- Create: `src/views/stats/components/MasteryChart.vue`
- Create: `src/views/stats/components/TrendChart.vue`
- Create: `src/views/stats/components/HotPointsRank.vue`
- Create: `src/views/stats/components/WeakPointsList.vue`

### Step 1: 安装 echarts vue-echarts

```bash
cd frontend-vue-app
npm install echarts vue-echarts
```

### Step 2: 创建 MasteryChart.vue

```vue
<!-- src/views/stats/components/MasteryChart.vue -->
<template>
  <div class="chart-card">
    <h3>知识点掌握度</h3>
    <v-chart v-if="data.length" :option="option" autoresize style="height: 300px" />
    <EmptyState v-else text="暂无掌握度数据" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { KnowledgeMastery } from '@/types/stats'
import EmptyState from '@/components/common/EmptyState.vue'

use([BarChart, GridComponent, TooltipComponent, CanvasRenderer])

const props = defineProps<{
  data: KnowledgeMastery[]
}>()

const getLevelColor = (level: string) => {
  const map: Record<string, string> = {
    mastered: '#10b981',
    learning: '#3b82f6',
    weak: '#f43f5e'
  }
  return map[level] || '#94a3b8'
}

const option = computed(() => ({
  tooltip: {
    trigger: 'axis',
    formatter: (params: any) => {
      const item = props.data[params[0]?.dataIndex]
      if (!item) return ''
      return `${item.knowledge_point_name}<br/>掌握度: ${(item.mastery_rate * 100).toFixed(0)}%<br/>总题数: ${item.total_questions}<br/>正确: ${item.correct_answers}`
    }
  },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: {
    type: 'category',
    data: props.data.map(i => i.knowledge_point_name),
    axisLabel: { rotate: 30, fontSize: 11 }
  },
  yAxis: {
    type: 'value',
    max: 100,
    axisLabel: { formatter: '{value}%' }
  },
  series: [{
    type: 'bar',
    data: props.data.map(i => ({
      value: Math.round(i.mastery_rate * 100),
      itemStyle: { color: getLevelColor(i.level), borderRadius: [4, 4, 0, 0] }
    })),
    barMaxWidth: 40
  }]
}))
</script>

<style scoped>
.chart-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.chart-card h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
</style>
```

### Step 3: 创建 TrendChart.vue

```vue
<!-- src/views/stats/components/TrendChart.vue -->
<template>
  <div class="chart-card">
    <h3>学习趋势</h3>
    <v-chart v-if="data && data.daily_stats.length" :option="option" autoresize style="height: 300px" />
    <EmptyState v-else text="暂无趋势数据" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { TrendData } from '@/types/stats'
import EmptyState from '@/components/common/EmptyState.vue'

use([LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const props = defineProps<{
  data: TrendData | null
}>()

const option = computed(() => {
  if (!props.data) return {}
  const dates = props.data.daily_stats.map(i => i.date)
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['学习时长', '正确率'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: dates },
    yAxis: [
      { type: 'value', name: '时长(h)', axisLabel: { formatter: '{value}' } },
      { type: 'value', name: '正确率(%)', max: 100, axisLabel: { formatter: '{value}%' } }
    ],
    series: [
      {
        name: '学习时长',
        type: 'line',
        data: props.data.daily_stats.map(i => i.learning_hours),
        smooth: true,
        itemStyle: { color: '#3b82f6' }
      },
      {
        name: '正确率',
        type: 'line',
        yAxisIndex: 1,
        data: props.data.daily_stats.map(i => Math.round(i.correct_rate * 100)),
        smooth: true,
        itemStyle: { color: '#10b981' }
      }
    ]
  }
})
</script>

<style scoped>
.chart-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.chart-card h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
</style>
```

### Step 4: 创建 HotPointsRank.vue

```vue
<!-- src/views/stats/components/HotPointsRank.vue -->
<template>
  <div class="chart-card">
    <h3>热门知识点排行</h3>
    <div v-if="data.length" class="rank-list">
      <div v-for="(item, index) in data" :key="item.knowledge_point_id" class="rank-item">
        <span class="rank-num" :class="{ top3: index < 3 }">{{ index + 1 }}</span>
        <span class="rank-name">{{ item.knowledge_point_name }}</span>
        <span class="rank-heat">{{ item.heat }}</span>
      </div>
    </div>
    <EmptyState v-else text="暂无数据" />
  </div>
</template>

<script setup lang="ts">
import type { HotKnowledgePoint } from '@/types/stats'
import EmptyState from '@/components/common/EmptyState.vue'

defineProps<{
  data: HotKnowledgePoint[]
}>()
</script>

<style scoped>
.chart-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.chart-card h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
.rank-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.rank-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-radius: 8px;
  background: var(--bg-hover);
}
.rank-num {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  background: var(--border-light);
  margin-right: 12px;
  flex-shrink: 0;
}
.rank-num.top3 {
  color: white;
  background: linear-gradient(135deg, #f59e0b, #f97316);
}
.rank-name {
  flex: 1;
  font-size: 14px;
  color: var(--text-primary);
}
.rank-heat {
  font-size: 13px;
  color: var(--text-muted);
  font-weight: 500;
}
</style>
```

### Step 5: 创建 WeakPointsList.vue

```vue
<!-- src/views/stats/components/WeakPointsList.vue -->
<template>
  <div class="chart-card">
    <h3>薄弱知识点</h3>
    <div v-if="data.length" class="weak-list">
      <div v-for="item in data" :key="item.knowledge_point_id" class="weak-item">
        <div class="weak-header">
          <span class="weak-name">{{ item.knowledge_point_name }}</span>
          <span class="weak-rate">{{ (item.correct_rate * 100).toFixed(0) }}%</span>
        </div>
        <el-progress
          :percentage="item.correct_rate * 100"
          :color="item.correct_rate < 0.3 ? '#f43f5e' : '#f59e0b'"
          :show-text="false"
        />
        <div class="suggested" v-if="item.suggested_questions.length">
          <span class="suggested-label">建议练习：</span>
          <el-tag
            v-for="q in item.suggested_questions"
            :key="q.id"
            size="small"
            class="suggested-tag"
          >{{ q.title }}</el-tag>
        </div>
      </div>
    </div>
    <EmptyState v-else text="暂无薄弱知识点" />
  </div>
</template>

<script setup lang="ts">
import type { WeakPoint } from '@/types/stats'
import EmptyState from '@/components/common/EmptyState.vue'

defineProps<{
  data: WeakPoint[]
}>()
</script>

<style scoped>
.chart-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.chart-card h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
.weak-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.weak-item {
  padding: 16px;
  border-radius: 8px;
  background: var(--bg-hover);
}
.weak-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}
.weak-name {
  font-weight: 500;
  color: var(--text-primary);
}
.weak-rate {
  font-size: 13px;
  color: var(--rose);
  font-weight: 600;
}
.suggested {
  margin-top: 8px;
}
.suggested-label {
  font-size: 12px;
  color: var(--text-muted);
}
.suggested-tag {
  margin: 2px 4px 2px 0;
}
</style>
```

### Step 6: 重写 index.vue

```vue
<!-- src/views/stats/index.vue -->
<template>
  <div class="stats-container" v-loading="loading">
    <h2>分析统计</h2>

    <!-- 统计概览卡片 -->
    <div class="overview-cards" v-if="overview">
      <StatCard title="总学习时长" :value="overview.total_learning_hours" unit="小时" icon-bg="#eff6ff" icon-color="#2563eb" />
      <StatCard title="总提问数" :value="overview.total_questions_asked" icon-bg="#f0fdf4" icon-color="#10b981" />
      <StatCard title="总测验数" :value="overview.total_quizzes_taken" icon-bg="#fef3c7" icon-color="#f59e0b" />
      <StatCard title="平均正确率" :value="overview.average_correct_rate" unit="%" icon-bg="#ede9fe" icon-color="#8b5cf6" />
      <StatCard title="知识点掌握" :value="overview.knowledge_points_mastered" :suffix="'/' + overview.knowledge_points_total" icon-bg="#cffafe" icon-color="#06b6d4" />
    </div>

    <!-- 图表区 -->
    <div class="chart-grid">
      <MasteryChart :data="masteryList" />
      <HotPointsRank :data="hotPoints" />
    </div>

    <TrendChart :data="trendData" />

    <WeakPointsList :data="weakPoints" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { OverviewStats, KnowledgeMastery, HotKnowledgePoint, WeakPoint, TrendData } from '@/types/stats'
import { getOverview, getKnowledgeMastery, getWeakPoints, getHotKnowledgePoints, getTrends } from '@/services/stats'
import StatCard from '@/components/common/StatCard.vue'
import MasteryChart from './components/MasteryChart.vue'
import TrendChart from './components/TrendChart.vue'
import HotPointsRank from './components/HotPointsRank.vue'
import WeakPointsList from './components/WeakPointsList.vue'

const loading = ref(true)
const overview = ref<OverviewStats | null>(null)
const masteryList = ref<KnowledgeMastery[]>([])
const hotPoints = ref<HotKnowledgePoint[]>([])
const weakPoints = ref<WeakPoint[]>([])
const trendData = ref<TrendData | null>(null)

const fetchData = async () => {
  loading.value = true
  try {
    const results = await Promise.allSettled([
      getOverview(),
      getKnowledgeMastery(),
      getWeakPoints(10),
      getHotKnowledgePoints(10),
      getTrends(7)
    ])
    if (results[0].status === 'fulfilled') overview.value = results[0].value.data
    if (results[1].status === 'fulfilled') masteryList.value = results[1].value.data
    if (results[2].status === 'fulfilled') weakPoints.value = results[2].value.data
    if (results[3].status === 'fulfilled') hotPoints.value = results[3].value.data
    if (results[4].status === 'fulfilled') trendData.value = results[4].value.data
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.stats-container h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 20px;
}
.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}
.chart-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}
@media (max-width: 768px) {
  .chart-grid {
    grid-template-columns: 1fr;
  }
}
</style>
```

### Step 7: 验证

```bash
npx vue-tsc --noEmit
npx vite build
```

### Step 8: Commit

```bash
git add src/views/stats/
git commit -m "feat: 实现分析统计页面"
```

---

## Task 3: 首页 (Dashboard)

**Files:**
- Rewrite: `src/views/home/index.vue`
- Create: `src/views/home/components/TodayStats.vue`
- Create: `src/views/home/components/QuickActions.vue`
- Create: `src/views/home/components/RecentQAList.vue`
- Create: `src/views/home/components/LearningTrend.vue`

### Step 1: 创建 TodayStats.vue

```vue
<!-- src/views/home/components/TodayStats.vue -->
<template>
  <div class="overview-cards" v-if="data">
    <StatCard title="今日学习时长" :value="data.today_learning_hours" unit="小时"
      icon-bg="#eff6ff" icon-color="#2563eb" />
    <StatCard title="今日提问数" :value="data.today_questions_asked"
      icon-bg="#f0fdf4" icon-color="#10b981" />
    <StatCard title="平均正确率" :value="data.average_correct_rate" unit="%"
      icon-bg="#ede9fe" icon-color="#8b5cf6" />
    <StatCard title="知识点掌握" :value="data.knowledge_points_mastered"
      :suffix="'/' + data.knowledge_points_total" icon-bg="#cffafe" icon-color="#06b6d4" />
  </div>
</template>

<script setup lang="ts">
import type { OverviewStats } from '@/types/stats'
import StatCard from '@/components/common/StatCard.vue'

defineProps<{
  data: OverviewStats | null
}>()
</script>

<style scoped>
.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}
</style>
```

### Step 2: 创建 QuickActions.vue

```vue
<!-- src/views/home/components/QuickActions.vue -->
<template>
  <div class="quick-actions">
    <router-link to="/qa" class="action-card">
      <div class="action-icon" style="background: #eff6ff; color: #2563eb">
        <el-icon :size="24"><ChatDotRound /></el-icon>
      </div>
      <span class="action-title">智能问答</span>
      <span class="action-desc">向 AI 提问</span>
    </router-link>
    <router-link to="/files" class="action-card">
      <div class="action-icon" style="background: #f0fdf4; color: #10b981">
        <el-icon :size="24"><Upload /></el-icon>
      </div>
      <span class="action-title">上传资料</span>
      <span class="action-desc">管理学习文档</span>
    </router-link>
    <router-link to="/knowledge-graph" class="action-card">
      <div class="action-icon" style="background: #fef3c7; color: #f59e0b">
        <el-icon :size="24"><Share /></el-icon>
      </div>
      <span class="action-title">知识图谱</span>
      <span class="action-desc">可视化知识网络</span>
    </router-link>
    <router-link to="/stats" class="action-card">
      <div class="action-icon" style="background: #ede9fe; color: #8b5cf6">
        <el-icon :size="24"><DataAnalysis /></el-icon>
      </div>
      <span class="action-title">学习统计</span>
      <span class="action-desc">查看学习分析</span>
    </router-link>
  </div>
</template>

<script setup lang="ts">
import { ChatDotRound, Upload, Share, DataAnalysis } from '@element-plus/icons-vue'
</script>

<style scoped>
.quick-actions {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}
.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 24px 16px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  text-decoration: none;
  transition: all 0.2s;
  cursor: pointer;
}
.action-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}
.action-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.action-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}
.action-desc {
  font-size: 12px;
  color: var(--text-muted);
}
@media (max-width: 768px) {
  .quick-actions {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
```

### Step 3: 创建 RecentQAList.vue

```vue
<!-- src/views/home/components/RecentQAList.vue -->
<template>
  <div class="card">
    <h3>最近问答</h3>
    <div v-if="data.length" class="qa-list">
      <router-link
        v-for="item in data"
        :key="item.id"
        :to="'/qa'"
        class="qa-item"
      >
        <span class="qa-question">{{ item.question }}</span>
        <span class="qa-time">{{ formatDate(item.created_at) }}</span>
      </router-link>
    </div>
    <EmptyState v-else text="暂无问答记录" />
  </div>
</template>

<script setup lang="ts">
import { formatDate } from '@/utils'
import EmptyState from '@/components/common/EmptyState.vue'

defineProps<{
  data: any[]
}>()
</script>

<style scoped>
.card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.card h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
.qa-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.qa-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  background: var(--bg-hover);
  text-decoration: none;
  transition: background 0.15s;
}
.qa-item:hover {
  background: var(--border-light);
}
.qa-question {
  font-size: 13px;
  color: var(--text-primary);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-right: 12px;
}
.qa-time {
  font-size: 12px;
  color: var(--text-muted);
  flex-shrink: 0;
}
</style>
```

### Step 4: 创建 LearningTrend.vue

```vue
<!-- src/views/home/components/LearningTrend.vue -->
<template>
  <div class="card">
    <h3>学习趋势</h3>
    <v-chart v-if="data && data.daily_stats.length" :option="option" autoresize style="height: 240px" />
    <EmptyState v-else text="暂无趋势数据" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { TrendData } from '@/types/stats'
import EmptyState from '@/components/common/EmptyState.vue'

use([LineChart, GridComponent, TooltipComponent, CanvasRenderer])

const props = defineProps<{
  data: TrendData | null
}>()

const option = computed(() => {
  if (!props.data) return {}
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: props.data.daily_stats.map(i => i.date) },
    yAxis: { type: 'value', name: '小时' },
    series: [{
      type: 'line',
      data: props.data.daily_stats.map(i => i.learning_hours),
      smooth: true,
      areaStyle: { color: 'rgba(37, 99, 235, 0.1)' },
      itemStyle: { color: '#3b82f6' }
    }]
  }
})
</script>

<style scoped>
.card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.card h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}
</style>
```

### Step 5: 重写 index.vue

```vue
<!-- src/views/home/index.vue -->
<template>
  <div class="home-container" v-loading="loading">
    <TodayStats :data="overview" />
    <QuickActions />

    <div class="home-grid">
      <RecentQAList :data="recentQA" />
      <LearningTrend :data="trendData" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { OverviewStats, TrendData } from '@/types/stats'
import { getOverview, getTrends, getHotKnowledgePoints } from '@/services/stats'
import { getAskHistory } from '@/services/qa'
import TodayStats from './components/TodayStats.vue'
import QuickActions from './components/QuickActions.vue'
import RecentQAList from './components/RecentQAList.vue'
import LearningTrend from './components/LearningTrend.vue'

const loading = ref(true)
const overview = ref<OverviewStats | null>(null)
const trendData = ref<TrendData | null>(null)
const recentQA = ref<any[]>([])

const fetchData = async () => {
  loading.value = true
  try {
    const results = await Promise.allSettled([
      getOverview(),
      getTrends(7),
      getHotKnowledgePoints(5),
      getAskHistory({ page: 1, size: 5 })
    ])
    if (results[0].status === 'fulfilled') overview.value = results[0].value.data
    if (results[1].status === 'fulfilled') trendData.value = results[1].value.data
    if (results[3].status === 'fulfilled') recentQA.value = results[3].value.data.list || []
  } catch (error) {
    console.error('获取首页数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.home-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}
@media (max-width: 768px) {
  .home-grid {
    grid-template-columns: 1fr;
  }
}
</style>
```

### Step 6: 验证 + Commit

```bash
npx vue-tsc --noEmit && npx vite build
git add src/views/home/
git commit -m "feat: 实现首页 Dashboard"
```

---

## Task 4: 知识图谱页面

**Files:**
- Rewrite: `src/views/knowledge-graph/index.vue`
- Create: `src/views/knowledge-graph/components/GraphCanvas.vue`
- Create: `src/views/knowledge-graph/components/GraphToolbar.vue`
- Create: `src/views/knowledge-graph/components/NodeDetail.vue`
- Create: `src/views/knowledge-graph/components/BuildDialog.vue`

### Step 1: 安装 d3

```bash
cd frontend-vue-app
npm install d3 @types/d3
```

### Step 2: 创建 GraphCanvas.vue（核心复杂组件）

```vue
<!-- src/views/knowledge-graph/components/GraphCanvas.vue -->
<template>
  <div ref="containerRef" class="graph-canvas" v-loading="loading">
    <svg ref="svgRef"></svg>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import * as d3 from 'd3'
import type { GraphData, GraphNode, GraphEdge } from '@/types/graph'

const props = defineProps<{
  data: GraphData | null
  loading: boolean
  highlight: string
}>()

const emit = defineEmits<{
  nodeClick: [node: GraphNode]
}>()

const containerRef = ref<HTMLDivElement>()
const svgRef = ref<SVGSVGElement>()
let simulation: d3.Simulation<GraphNode, GraphEdge> | null = null
let zoom: d3.ZoomBehavior<SVGSVGElement, unknown> | null = null

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#f43f5e', '#8b5cf6', '#06b6d4', '#ec4899']

const renderGraph = () => {
  if (!props.data || !svgRef.value || !containerRef.value) return

  // 清空旧内容
  d3.select(svgRef.value).selectAll('*').remove()

  const width = containerRef.value.clientWidth
  const height = containerRef.value.clientHeight

  const svg = d3.select(svgRef.value)
    .attr('width', width)
    .attr('height', height)

  // 缩放平移
  const g = svg.append('g')
  zoom = d3.zoom<SVGSVGElement, unknown>()
    .scaleExtent([0.3, 4])
    .on('zoom', (event) => g.attr('transform', event.transform))
  svg.call(zoom)

  // 准备数据
  const nodes = props.data.nodes.map(d => ({ ...d }))
  const edges = props.data.edges.map(d => ({ ...d }))
  const nodeMap = new Map(nodes.map(n => [n.id, n]))

  // 颜色映射
  const categories = [...new Set(nodes.map(n => n.category || '未知'))]
  const colorScale = d3.scaleOrdinal<string>().domain(categories).range(COLORS)

  // 力导向模拟
  simulation = d3.forceSimulation<GraphNode>(nodes)
    .force('link', d3.forceLink<GraphNode, GraphEdge>(edges).id(d => d.id).distance(120))
    .force('charge', d3.forceManyBody().strength(-300))
    .force('center', d3.forceCenter(width / 2, height / 2))
    .force('collide', d3.forceCollide().radius(30))

  // 绘制边
  const link = g.append('g')
    .selectAll('line')
    .data(edges)
    .join('line')
    .attr('stroke', '#cbd5e1')
    .attr('stroke-width', 1.5)
    .attr('stroke-opacity', 0.6)

  // 绘制节点
  const node = g.append('g')
    .selectAll('g')
    .data(nodes)
    .join('g')
    .call(d3.drag<SVGGElement, GraphNode>()
      .on('start', (event, d) => {
        if (!event.active) simulation!.alphaTarget(0.3).restart()
        d.fx = d.x
        d.fy = d.y
      })
      .on('drag', (event, d) => {
        d.fx = event.x
        d.fy = event.y
      })
      .on('end', (event, d) => {
        if (!event.active) simulation!.alphaTarget(0)
        d.fx = null
        d.fy = null
      })
    )

  node.append('circle')
    .attr('r', 16)
    .attr('fill', d => colorScale(d.category || '未知'))
    .attr('stroke', '#fff')
    .attr('stroke-width', 2)
    .style('cursor', 'pointer')
    .on('click', (_, d) => emit('nodeClick', d))
    .on('mouseenter', function(_, d) {
      d3.select(this).transition().attr('r', 20)
      // 高亮相邻
      link.attr('stroke', l => {
        const source = typeof l.source === 'object' ? l.source.id : l.source
        const target = typeof l.target === 'object' ? l.target.id : l.target
        return (source === d.id || target === d.id) ? '#3b82f6' : '#e2e8f0'
      }).attr('stroke-width', l => {
        const source = typeof l.source === 'object' ? l.source.id : l.source
        const target = typeof l.target === 'object' ? l.target.id : l.target
        return (source === d.id || target === d.id) ? 3 : 1
      })
    })
    .on('mouseleave', function() {
      d3.select(this).transition().attr('r', 16)
      link.attr('stroke', '#cbd5e1').attr('stroke-width', 1.5)
    })

  node.append('text')
    .text(d => d.name)
    .attr('dy', 30)
    .attr('text-anchor', 'middle')
    .attr('font-size', 12)
    .attr('fill', '#475569')

  // 更新位置
  simulation.on('tick', () => {
    link
      .attr('x1', d => (d.source as any).x)
      .attr('y1', d => (d.source as any).y)
      .attr('x2', d => (d.target as any).x)
      .attr('y2', d => (d.target as any).y)
    node.attr('transform', d => `translate(${d.x},${d.y})`)
  })
}

// 搜索高亮
const highlightNodes = (keyword: string) => {
  if (!svgRef.value || !keyword) return
  d3.select(svgRef.value).selectAll<SVGCircleElement, GraphNode>('circle')
    .transition()
    .attr('r', d => d.name.includes(keyword) ? 22 : 14)
    .attr('stroke', d => d.name.includes(keyword) ? '#f59e0b' : '#fff')
    .attr('stroke-width', d => d.name.includes(keyword) ? 3 : 2)
}

watch(() => props.data, async () => {
  await nextTick()
  renderGraph()
}, { deep: true })

watch(() => props.highlight, (val) => {
  highlightNodes(val)
})

onMounted(() => {
  renderGraph()
  // 监听容器尺寸变化
  if (containerRef.value) {
    const observer = new ResizeObserver(() => renderGraph())
    observer.observe(containerRef.value)
  }
})

onUnmounted(() => {
  simulation?.stop()
  simulation = null
})
</script>

<style scoped>
.graph-canvas {
  width: 100%;
  height: 500px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}
.graph-canvas svg {
  width: 100%;
  height: 100%;
}
</style>
```

### Step 3: 创建 GraphToolbar.vue

```vue
<!-- src/views/knowledge-graph/components/GraphToolbar.vue -->
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
      <el-option label="RELATED" value="RELATED" />
      <el-option label="DEPENDS_ON" value="DEPENDS_ON" />
      <el-option label="PART_OF" value="PART_OF" />
    </el-select>
    <div class="toolbar-right">
      <el-button @click="$emit('refresh')">
        <el-icon><Refresh /></el-icon> 刷新
      </el-button>
      <el-button type="primary" @click="$emit('build')">
        <el-icon><Setting /></el-icon> 构建图谱
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Refresh, Setting } from '@element-plus/icons-vue'

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
```

### Step 4: 创建 NodeDetail.vue

```vue
<!-- src/views/knowledge-graph/components/NodeDetail.vue -->
<template>
  <el-drawer v-model="visible" title="节点详情" size="400px">
    <template v-if="modelValue">
      <div class="node-info">
        <h3>{{ modelValue.name }}</h3>
        <p class="node-desc">{{ modelValue.description || '暂无描述' }}</p>
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="ID">{{ modelValue.id }}</el-descriptions-item>
          <el-descriptions-item label="分类">{{ modelValue.category || '-' }}</el-descriptions-item>
          <el-descriptions-item label="关联文档ID">{{ modelValue.document_id }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import type { GraphNode } from '@/types/graph'

const visible = defineModel<boolean>({ default: false })

defineProps<{
  modelValue: GraphNode | null
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
```

### Step 5: 创建 BuildDialog.vue

```vue
<!-- src/views/knowledge-graph/components/BuildDialog.vue -->
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
```

### Step 6: 重写 index.vue

```vue
<!-- src/views/knowledge-graph/index.vue -->
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

    <NodeDetail v-model="detailVisible" :model-value="selectedNode" />
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
```

### Step 7: 验证 + Commit

```bash
npx vue-tsc --noEmit && npx vite build
git add src/views/knowledge-graph/
git commit -m "feat: 实现知识图谱页面 (D3 力导向图)"
```

---

## Task 5: 问答中心页面

**Files:**
- Rewrite: `src/views/qa/index.vue`
- Create: `src/views/qa/components/SessionList.vue`
- Create: `src/views/qa/components/ChatPanel.vue`
- Create: `src/views/qa/components/MessageBubble.vue`
- Create: `src/views/qa/components/SourceReference.vue`
- Create: `src/views/qa/components/EmptyChat.vue`

### Step 1: 创建 SessionList.vue

```vue
<!-- src/views/qa/components/SessionList.vue -->
<template>
  <div class="session-list">
    <div class="session-header">
      <h3>会话列表</h3>
      <el-button type="primary" size="small" @click="$emit('new')">
        <el-icon><Plus /></el-icon> 新建
      </el-button>
    </div>
    <div class="session-items" v-loading="loading">
      <div
        v-for="session in sessions"
        :key="session.conversation_id"
        class="session-item"
        :class="{ active: activeId === session.conversation_id }"
        @click="$emit('select', session.conversation_id)"
      >
        <div class="session-title">{{ session.title }}</div>
        <div class="session-meta">
          <span class="session-last">{{ session.last_question }}</span>
          <span class="session-count">{{ session.message_count }}条</span>
        </div>
      </div>
      <EmptyState v-if="!loading && !sessions.length" text="暂无会话" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import type { ConversationItem } from '@/types/qa'
import EmptyState from '@/components/common/EmptyState.vue'

defineProps<{
  sessions: ConversationItem[]
  activeId: number | null
  loading: boolean
}>()

defineEmits<{
  select: [id: number]
  new: []
}>()
</script>

<style scoped>
.session-list {
  width: 300px;
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  background: var(--bg-card);
}
.session-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--border-light);
}
.session-header h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}
.session-items {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}
.session-item {
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
  margin-bottom: 4px;
}
.session-item:hover {
  background: var(--bg-hover);
}
.session-item.active {
  background: var(--primary-light);
}
.session-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.session-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-muted);
}
.session-last {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  margin-right: 8px;
}
</style>
```

### Step 2: 创建 MessageBubble.vue

```vue
<!-- src/views/qa/components/MessageBubble.vue -->
<template>
  <div class="message-bubble" :class="message.role">
    <div class="bubble-avatar">
      <span v-if="message.role === 'user'">我</span>
      <span v-else>AI</span>
    </div>
    <div class="bubble-content">
      <div class="bubble-text">{{ message.content }}</div>
      <div class="bubble-time">{{ formatTime(message.created_at) }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { MessageItem } from '@/types/qa'

defineProps<{
  message: MessageItem
}>()

const formatTime = (dateStr: string) => {
  const d = new Date(dateStr)
  return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
}
</script>

<style scoped>
.message-bubble {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  max-width: 80%;
}
.message-bubble.user {
  margin-left: auto;
  flex-direction: row-reverse;
}
.bubble-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: white;
  flex-shrink: 0;
}
.user .bubble-avatar {
  background: linear-gradient(135deg, #2563eb, #3b82f6);
}
.assistant .bubble-avatar {
  background: linear-gradient(135deg, #8b5cf6, #a78bfa);
}
.bubble-content {
  flex: 1;
}
.bubble-text {
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}
.user .bubble-text {
  background: #2563eb;
  color: white;
  border-bottom-right-radius: 4px;
}
.assistant .bubble-text {
  background: var(--bg-card);
  color: var(--text-primary);
  box-shadow: var(--shadow-sm);
  border-bottom-left-radius: 4px;
}
.bubble-time {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 4px;
}
.user .bubble-time {
  text-align: right;
}
</style>
```

### Step 3: 创建 SourceReference.vue

```vue
<!-- src/views/qa/components/SourceReference.vue -->
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
</style>
```

### Step 4: 创建 EmptyChat.vue

```vue
<!-- src/views/qa/components/EmptyChat.vue -->
<template>
  <div class="empty-chat">
    <div class="empty-icon">💬</div>
    <h3>开始提问吧</h3>
    <p>基于知识图谱的智能问答，帮你理解软件工程知识</p>
    <div class="presets">
      <div
        v-for="q in presetQuestions"
        :key="q"
        class="preset-item"
        @click="$emit('ask', q)"
      >{{ q }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineEmits<{ ask: [question: string] }>()

const presetQuestions = [
  '什么是需求分析？',
  '软件测试有哪些方法？',
  '什么是软件生命周期？',
  '敏捷开发和瀑布模型的区别？'
]
</script>

<style scoped>
.empty-chat {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 40px;
  text-align: center;
}
.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}
h3 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}
p {
  font-size: 14px;
  color: var(--text-muted);
  margin-bottom: 32px;
}
.presets {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
  max-width: 400px;
}
.preset-item {
  padding: 12px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 10px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s;
}
.preset-item:hover {
  border-color: var(--primary);
  color: var(--primary);
  background: var(--primary-light);
}
</style>
```

### Step 5: 创建 ChatPanel.vue

```vue
<!-- src/views/qa/components/ChatPanel.vue -->
<template>
  <div class="chat-panel">
    <div class="messages-area" ref="messagesRef">
      <EmptyChat v-if="!messages.length && !loading" @ask="(q) => $emit('ask', q)" />
      <template v-else>
        <MessageBubble v-for="msg in messages" :key="msg.message_id" :message="msg" />
        <SourceReference
          v-if="lastAssistantMsg"
          :sources="lastAssistantMsg.sources || []"
          :related-points="lastAssistantMsg.relatedKnowledgePoints || []"
        />
      </template>
      <div v-if="isAsking" class="typing-indicator">
        <span></span><span></span><span></span>
      </div>
    </div>

    <div class="input-area">
      <el-input
        v-model="inputText"
        type="textarea"
        :autosize="{ minRows: 1, maxRows: 4 }"
        placeholder="输入你的问题... (Enter 发送, Shift+Enter 换行)"
        :disabled="isAsking"
        @keydown.enter.exact.prevent="handleSend"
      />
      <el-button
        type="primary"
        :disabled="!inputText.trim() || isAsking"
        :loading="isAsking"
        @click="handleSend"
      >发送</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import type { MessageItem } from '@/types/qa'
import type { AskResponse } from '@/services/qa'
import MessageBubble from './MessageBubble.vue'
import SourceReference from './SourceReference.vue'
import EmptyChat from './EmptyChat.vue'

interface ExtendedMessage extends MessageItem {
  sources?: AskResponse['sources']
  relatedKnowledgePoints?: AskResponse['related_knowledge_points']
}

const props = defineProps<{
  messages: ExtendedMessage[]
  loading: boolean
  isAsking: boolean
}>()

const emit = defineEmits<{
  ask: [question: string]
}>()

const inputText = ref('')
const messagesRef = ref<HTMLDivElement>()

const lastAssistantMsg = computed(() => {
  const assistantMsgs = props.messages.filter(m => m.role === 'assistant')
  return assistantMsgs[assistantMsgs.length - 1] || null
})

const handleSend = () => {
  const text = inputText.value.trim()
  if (!text || props.isAsking) return
  emit('ask', text)
  inputText.value = ''
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesRef.value) {
    messagesRef.value.scrollTop = messagesRef.value.scrollHeight
  }
}

watch(() => props.messages.length, scrollToBottom)
watch(() => props.isAsking, scrollToBottom)
</script>

<style scoped>
.chat-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}
.input-area {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--border-light);
  background: var(--bg-card);
  align-items: flex-end;
}
.input-area .el-textarea {
  flex: 1;
}
.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 12px 16px;
  background: var(--bg-card);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  width: fit-content;
  margin-left: 48px;
}
.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--text-muted);
  animation: bounce 1.4s infinite ease-in-out;
}
.typing-indicator span:nth-child(1) { animation-delay: -0.32s; }
.typing-indicator span:nth-child(2) { animation-delay: -0.16s; }
@keyframes bounce {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}
</style>
```

### Step 6: 重写 index.vue

```vue
<!-- src/views/qa/index.vue -->
<template>
  <div class="qa-container">
    <SessionList
      :sessions="sessions"
      :active-id="currentSessionId"
      :loading="sessionsLoading"
      @select="handleSelectSession"
      @new="handleNewSession"
    />
    <ChatPanel
      :messages="messages"
      :loading="messagesLoading"
      :is-asking="isAsking"
      @ask="handleAskQuestion"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { ConversationItem, MessageItem } from '@/types/qa'
import type { AskResponse } from '@/services/qa'
import { getSessions, getSessionMessages, createSession, askQuestion } from '@/services/qa'
import { ElMessage } from 'element-plus'
import SessionList from './components/SessionList.vue'
import ChatPanel from './components/ChatPanel.vue'

interface ExtendedMessage extends MessageItem {
  sources?: AskResponse['sources']
  relatedKnowledgePoints?: AskResponse['related_knowledge_points']
}

const sessions = ref<ConversationItem[]>([])
const currentSessionId = ref<number | null>(null)
const messages = ref<ExtendedMessage[]>([])
const isAsking = ref(false)
const sessionsLoading = ref(false)
const messagesLoading = ref(false)

const fetchSessions = async () => {
  sessionsLoading.value = true
  try {
    const result = await getSessions({ page: 1, size: 50 })
    sessions.value = result.data.list || []
  } catch (error) {
    console.error('获取会话列表失败:', error)
  } finally {
    sessionsLoading.value = false
  }
}

const handleNewSession = async () => {
  try {
    const result = await createSession()
    await fetchSessions()
    currentSessionId.value = result.data.conversation_id
    messages.value = []
  } catch (error) {
    console.error('创建会话失败:', error)
  }
}

const handleSelectSession = async (sessionId: number) => {
  currentSessionId.value = sessionId
  messagesLoading.value = true
  try {
    const result = await getSessionMessages(sessionId, { page: 1, size: 100 })
    messages.value = result.data.list || []
  } catch (error) {
    console.error('获取消息失败:', error)
  } finally {
    messagesLoading.value = false
  }
}

const handleAskQuestion = async (question: string) => {
  if (isAsking.value) return

  // 如果没有当前会话，先创建
  if (!currentSessionId.value) {
    try {
      const result = await createSession()
      currentSessionId.value = result.data.conversation_id
    } catch {
      ElMessage.error('创建会话失败')
      return
    }
  }

  // 添加用户消息
  const userMsg: ExtendedMessage = {
    message_id: Date.now(),
    role: 'user',
    content: question,
    created_at: new Date().toISOString()
  }
  messages.value.push(userMsg)

  isAsking.value = true
  try {
    const result = await askQuestion({
      question,
      conversation_id: currentSessionId.value
    })
    const data = result.data

    // 添加 AI 回复
    const aiMsg: ExtendedMessage = {
      message_id: data.question_id,
      role: 'assistant',
      content: data.answer,
      created_at: data.created_at,
      sources: data.sources,
      relatedKnowledgePoints: data.related_knowledge_points
    }
    messages.value.push(aiMsg)

    // 刷新会话列表
    fetchSessions()
  } catch (error) {
    ElMessage.error('提问失败')
  } finally {
    isAsking.value = false
  }
}

onMounted(() => {
  fetchSessions()
})
</script>

<style scoped>
.qa-container {
  display: flex;
  height: calc(100vh - 100px);
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}
</style>
```

### Step 7: 验证 + Commit

```bash
npx vue-tsc --noEmit && npx vite build
git add src/views/qa/
git commit -m "feat: 实现问答中心页面"
```

---

## 最终验证

```bash
cd frontend-vue-app
npm run dev
```

手动验证所有页面：
1. 登录 → 首页展示统计卡片和快捷入口
2. 资料管理 → 上传文件、搜索、查看详情、删除
3. 分析统计 → 图表正确渲染、数据展示
4. 知识图谱 → 力导向图渲染、节点交互、构建图谱
5. 问答中心 → 新建会话、提问、查看回答和来源

```bash
# 最终提交
git add -A
git commit -m "feat: 完成所有功能页面开发"
```
