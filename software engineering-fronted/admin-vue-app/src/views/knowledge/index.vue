<template>
  <div class="knowledge-page">
    <div class="page-header">
      <h2>知识点管理</h2>
      <el-button type="primary" :icon="Plus" @click="handleCreate">新增知识点</el-button>
    </div>

    <!-- 搜索栏 -->
    <div class="filter-bar page-card">
      <el-input
        v-model="keyword"
        placeholder="搜索知识点"
        :prefix-icon="Search"
        clearable
        style="width: 280px"
        @keyup.enter="fetchKnowledgePoints"
      />
      <el-button type="primary" @click="fetchKnowledgePoints">查询</el-button>
    </div>

    <!-- 知识点列表 -->
    <div class="page-card">
      <el-table :data="knowledgePoints" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="知识点名称" min-width="200" />
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.category" size="small" type="info">{{ row.category }}</el-tag>
            <span v-else class="text-muted">--</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="250" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
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
          @size-change="fetchKnowledgePoints"
          @current-change="fetchKnowledgePoints"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑知识点' : '新增知识点'"
      width="500px"
    >
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入知识点名称" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-input v-model="formData.category" placeholder="请输入分类（可选）" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="4"
            placeholder="请输入知识点描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import {
  getKnowledgePoints,
  createKnowledgePoint,
  updateKnowledgePoint,
  deleteKnowledgePoint,
} from '@/services/admin'

const loading = ref(false)
const knowledgePoints = ref<Record<string, unknown>[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const keyword = ref('')

const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  name: '',
  category: '',
  description: '',
})

const formRules = {
  name: [{ required: true, message: '请输入知识点名称', trigger: 'blur' }],
}

async function fetchKnowledgePoints() {
  loading.value = true
  try {
    const data = await getKnowledgePoints({
      page: currentPage.value,
      size: pageSize.value,
      keyword: keyword.value || undefined,
    }) as Record<string, unknown>
    knowledgePoints.value = (data.list as Record<string, unknown>[]) || []
    total.value = (data.total as number) || 0
  } catch (error) {
    console.error('获取知识点列表失败:', error)
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  isEdit.value = false
  editId.value = null
  formData.name = ''
  formData.category = ''
  formData.description = ''
  dialogVisible.value = true
}

function handleEdit(row: Record<string, unknown>) {
  isEdit.value = true
  editId.value = row.id as number
  formData.name = row.name as string
  formData.category = (row.category as string) || ''
  formData.description = (row.description as string) || ''
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (isEdit.value && editId.value) {
        await updateKnowledgePoint(editId.value, { ...formData })
        ElMessage.success('更新成功')
      } else {
        await createKnowledgePoint({ ...formData })
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchKnowledgePoints()
    } catch (error) {
      console.error('操作失败:', error)
    } finally {
      submitLoading.value = false
    }
  })
}

async function handleDelete(row: Record<string, unknown>) {
  try {
    await ElMessageBox.confirm(
      `确定要删除知识点「${row.name}」吗？此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    )
    await deleteKnowledgePoint(row.id as number)
    ElMessage.success('删除成功')
    fetchKnowledgePoints()
  } catch {
    // 用户取消
  }
}

onMounted(() => {
  fetchKnowledgePoints()
})
</script>

<style scoped>
.knowledge-page {
  max-width: 1200px;
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

.text-muted {
  color: var(--text-muted);
}
</style>
