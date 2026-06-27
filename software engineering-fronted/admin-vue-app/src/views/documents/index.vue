<template>
  <div class="documents-page">
    <div class="page-header">
      <h2>资料审核</h2>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar page-card">
      <el-select v-model="statusFilter" placeholder="审核状态" clearable style="width: 160px">
        <el-option label="待审核" value="pending" />
        <el-option label="已通过" value="approved" />
        <el-option label="已拒绝" value="rejected" />
      </el-select>
      <el-button type="primary" @click="fetchDocuments">查询</el-button>
    </div>

    <!-- 文档列表 -->
    <div class="page-card">
      <el-table :data="documents" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="资料名称" min-width="200" />
        <el-table-column prop="file_type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.file_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="uploader_name" label="上传者" width="120" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleView(row)">查看</el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              link
              size="small"
              @click="handleReview(row, 'approved')"
            >通过</el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="danger"
              link
              size="small"
              @click="handleReview(row, 'rejected')"
            >拒绝</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchDocuments"
          @current-change="fetchDocuments"
        />
      </div>
    </div>

    <!-- 审核对话框 -->
    <el-dialog
      v-model="reviewDialogVisible"
      title="审核资料"
      width="500px"
    >
      <el-form :model="reviewForm">
        <el-form-item label="审核结果">
          <el-radio-group v-model="reviewForm.status">
            <el-radio value="approved">通过</el-radio>
            <el-radio value="rejected">拒绝</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="reviewForm.comment"
            type="textarea"
            :rows="3"
            placeholder="请输入审核备注（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reviewDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="reviewLoading" @click="confirmReview">确定</el-button>
      </template>
    </el-dialog>

    <!-- 查看详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="资料详情"
      width="800px"
    >
      <el-descriptions :column="2" border v-if="currentViewDoc">
        <el-descriptions-item label="标题" :span="2">{{ currentViewDoc.title }}</el-descriptions-item>
        <el-descriptions-item label="文件名">{{ currentViewDoc.filename }}</el-descriptions-item>
        <el-descriptions-item label="文件类型">{{ currentViewDoc.file_type }}</el-descriptions-item>
        <el-descriptions-item label="审核状态">
          <el-tag :type="getStatusType(currentViewDoc.status as string)">
            {{ getStatusLabel(currentViewDoc.status as string) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="上传时间">{{ currentViewDoc.created_at }}</el-descriptions-item>
        <el-descriptions-item v-if="currentViewDoc.description" label="描述" :span="2">{{ currentViewDoc.description }}</el-descriptions-item>
        <el-descriptions-item v-if="currentViewDoc.remark" label="审核备注" :span="2">{{ currentViewDoc.remark }}</el-descriptions-item>
      </el-descriptions>

      <!-- 文档预览区域 -->
      <div class="document-preview">
        <div class="preview-header">
          <span class="preview-title">📄 文档预览</span>
        </div>
        <div class="preview-content" v-loading="contentLoading">
          <pre v-if="documentContent">{{ documentContent }}</pre>
          <el-empty v-else-if="!contentLoading" description="暂无文档内容" />
        </div>
      </div>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDocuments, getDocumentContent, reviewDocument } from '@/services/admin'

const loading = ref(false)
const documents = ref<Record<string, unknown>[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const statusFilter = ref('')

const reviewDialogVisible = ref(false)
const reviewLoading = ref(false)
const currentReviewDoc = ref<Record<string, unknown> | null>(null)
const reviewForm = reactive({
  status: 'approved' as 'approved' | 'rejected',
  comment: '',
})

const detailDialogVisible = ref(false)
const currentViewDoc = ref<Record<string, unknown> | null>(null)
const documentContent = ref('')
const contentLoading = ref(false)

function getStatusType(status: string) {
  const map: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
  }
  return map[status] || 'info'
}

function getStatusLabel(status: string) {
  const map: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
  }
  return map[status] || status
}

async function fetchDocuments() {
  loading.value = true
  try {
    const data = await getDocuments({
      page: currentPage.value,
      size: pageSize.value,
      status: statusFilter.value || undefined,
    }) as Record<string, unknown>
    documents.value = (data.list as Record<string, unknown>[]) || []
    total.value = (data.total as number) || 0
  } catch (error) {
    console.error('获取文档列表失败:', error)
  } finally {
    loading.value = false
  }
}

async function handleView(row: Record<string, unknown>) {
  currentViewDoc.value = row
  documentContent.value = ''
  detailDialogVisible.value = true

  // 获取文档内容用于预览
  contentLoading.value = true
  try {
    const res = await getDocumentContent(row.id as number) as Record<string, unknown>
    documentContent.value = (res.content as string) || ''
  } catch (error) {
    console.error('获取文档内容失败:', error)
    documentContent.value = '加载文档内容失败'
  } finally {
    contentLoading.value = false
  }
}

function handleReview(row: Record<string, unknown>, status: 'approved' | 'rejected') {
  currentReviewDoc.value = row
  reviewForm.status = status
  reviewForm.comment = ''
  reviewDialogVisible.value = true
}

async function confirmReview() {
  if (!currentReviewDoc.value) return

  const actionText = reviewForm.status === 'approved' ? '通过' : '拒绝'
  try {
    await ElMessageBox.confirm(
      `确定要${actionText}该资料吗？`,
      '确认审核',
      { type: 'warning' }
    )
  } catch {
    return
  }

  reviewLoading.value = true
  try {
    await reviewDocument(currentReviewDoc.value.id as number, {
      status: reviewForm.status,
      comment: reviewForm.comment,
    })
    ElMessage.success(`已${actionText}`)
    reviewDialogVisible.value = false
    fetchDocuments()
  } catch (error) {
    console.error('审核失败:', error)
  } finally {
    reviewLoading.value = false
  }
}

onMounted(() => {
  fetchDocuments()
})
</script>

<style scoped>
.documents-page {
  width: 100%;
}

.filter-bar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 20px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.document-preview {
  margin-top: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
}

.preview-header {
  background-color: #f5f7fa;
  padding: 12px 16px;
  border-bottom: 1px solid #e4e7ed;
}

.preview-title {
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

.preview-content {
  max-height: 400px;
  overflow-y: auto;
  padding: 16px;
  background-color: #fafafa;
}

.preview-content pre {
  margin: 0;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: #303133;
}
</style>
