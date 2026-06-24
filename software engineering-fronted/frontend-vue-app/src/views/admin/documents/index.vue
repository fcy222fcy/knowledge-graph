<template>
  <div class="admin-documents" v-loading="loading">
    <div class="page-header">
      <h2 class="page-title">资料管理</h2>
      <div class="header-actions">
        <el-select
          v-model="statusFilter"
          placeholder="审核状态"
          clearable
          style="width: 140px; margin-right: 12px;"
          @change="fetchDocuments"
        >
          <el-option label="全部" value="" />
          <el-option label="待审核" value="pending" />
          <el-option label="已通过" value="approved" />
          <el-option label="已拒绝" value="rejected" />
        </el-select>
        <el-input
          v-model="keyword"
          placeholder="搜索资料标题"
          clearable
          style="width: 240px; margin-right: 12px;"
          @clear="fetchDocuments"
          @keyup.enter="fetchDocuments"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="fetchDocuments">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 资料表格 -->
    <el-table :data="documents" stripe border style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
      <el-table-column prop="filename" label="文件名" width="150" show-overflow-tooltip />
      <el-table-column prop="file_type" label="类型" width="80" />
      <el-table-column label="大小" width="100">
        <template #default="{ row }">
          {{ formatFileSize(row.file_size) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="审核状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="审核意见" min-width="150" show-overflow-tooltip>
        <template #default="{ row }">
          <span v-if="row.review_comment">{{ row.review_comment }}</span>
          <span v-else class="text-muted">-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="上传时间" width="180" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleView(row)">查看</el-button>
          <el-button
            v-if="row.status === 'pending'"
            type="warning"
            link
            @click="handleReview(row)"
          >审核</el-button>
          <el-popconfirm
            title="确定要删除该资料吗？"
            @confirm="handleDelete(row)"
          >
            <template #reference>
              <el-button type="danger" link>删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      style="margin-top: 16px; justify-content: flex-end;"
      @size-change="fetchDocuments"
      @current-change="fetchDocuments"
    />

    <!-- 查看详情对话框 -->
    <el-dialog v-model="detailVisible" title="资料详情" width="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="标题">{{ detailDocument?.title }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ detailDocument?.description || '无' }}</el-descriptions-item>
        <el-descriptions-item label="文件名">{{ detailDocument?.filename }}</el-descriptions-item>
        <el-descriptions-item label="文件类型">{{ detailDocument?.file_type }}</el-descriptions-item>
        <el-descriptions-item label="文件大小">{{ formatFileSize(detailDocument?.file_size || 0) }}</el-descriptions-item>
        <el-descriptions-item label="审核状态">
          <el-tag :type="getStatusType(detailDocument?.status || '')">{{ getStatusLabel(detailDocument?.status || '') }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item v-if="detailDocument?.review_comment" label="审核意见">{{ detailDocument?.review_comment }}</el-descriptions-item>
        <el-descriptions-item label="上传时间">{{ detailDocument?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 审核对话框 -->
    <el-dialog v-model="reviewVisible" title="审核资料" width="500px">
      <el-form :model="reviewForm" label-width="80px">
        <el-form-item label="资料标题">
          <span>{{ reviewDocumentItem?.title }}</span>
        </el-form-item>
        <el-form-item label="审核结果" required>
          <el-radio-group v-model="reviewForm.status">
            <el-radio label="approved">通过</el-radio>
            <el-radio label="rejected">拒绝</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审核意见">
          <el-input
            v-model="reviewForm.comment"
            type="textarea"
            :rows="3"
            placeholder="请输入审核意见（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reviewVisible = false">取消</el-button>
        <el-button type="primary" :loading="reviewLoading" @click="submitReview">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import { getDocuments, deleteDocument, reviewDocument as reviewDocumentApi } from '@/services/admin'
import type { DocumentItem } from '@/services/admin'

const loading = ref(false)
const documents = ref<DocumentItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const keyword = ref('')
const statusFilter = ref('')

// 详情相关
const detailVisible = ref(false)
const detailDocument = ref<DocumentItem | null>(null)

// 审核相关
const reviewVisible = ref(false)
const reviewLoading = ref(false)
const reviewDocumentItem = ref<DocumentItem | null>(null)
const reviewForm = ref({
  status: 'approved' as 'approved' | 'rejected',
  comment: ''
})

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return types[status] || 'info'
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return labels[status] || '未知'
}

const fetchDocuments = async () => {
  loading.value = true
  try {
    const params: { page: number; size: number; keyword: string; status?: string } = {
      page: page.value,
      size: pageSize.value,
      keyword: keyword.value
    }
    if (statusFilter.value) {
      params.status = statusFilter.value
    }
    const res = await getDocuments(params)
    documents.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    ElMessage.error('获取资料列表失败')
  } finally {
    loading.value = false
  }
}

const handleView = (row: DocumentItem) => {
  detailDocument.value = row
  detailVisible.value = true
}

const handleReview = (row: DocumentItem) => {
  reviewDocumentItem.value = row
  reviewForm.value = { status: 'approved', comment: '' }
  reviewVisible.value = true
}

const submitReview = async () => {
  if (!reviewDocumentItem.value) return
  reviewLoading.value = true
  try {
    await reviewDocumentApi(reviewDocumentItem.value.id, {
      status: reviewForm.value.status,
      comment: reviewForm.value.comment || undefined
    })
    ElMessage.success('审核成功')
    reviewVisible.value = false
    fetchDocuments()
  } catch (error) {
    ElMessage.error('审核失败')
  } finally {
    reviewLoading.value = false
  }
}

const handleDelete = async (row: DocumentItem) => {
  try {
    await deleteDocument(row.id)
    ElMessage.success('删除成功')
    fetchDocuments()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  fetchDocuments()
})
</script>

<style scoped>
.admin-documents {
  background: #fff;
  padding: 24px;
  border-radius: 8px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.header-actions {
  display: flex;
  align-items: center;
}

.text-muted {
  color: #909399;
}
</style>
