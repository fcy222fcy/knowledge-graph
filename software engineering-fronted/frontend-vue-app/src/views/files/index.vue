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
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import { getDocumentList, getDocumentDetail, deleteDocument } from '@/services/document'
import type { DocumentItem } from '@/types/document'
import { usePagination } from '@/composables/usePagination'
import DocumentTable from './components/DocumentTable.vue'
import UploadDialog from './components/UploadDialog.vue'
import DocumentDetail from './components/DocumentDetail.vue'

const route = useRoute()
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

// 处理从知识图谱跳转过来的文档ID参数
const handleDocumentParam = async () => {
  const docId = route.query.doc_id
  if (docId) {
    try {
      const result = await getDocumentDetail(Number(docId))
      if (result.data) {
        currentDocument.value = result.data
        detailVisible.value = true
      }
    } catch (error) {
      console.error('获取文档详情失败:', error)
    }
  }
}

onMounted(() => {
  fetchDocuments()
  handleDocumentParam()
})

// 监听路由参数变化
watch(() => route.query.doc_id, () => {
  handleDocumentParam()
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
