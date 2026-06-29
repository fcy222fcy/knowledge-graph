<template>
  <div class="students-page">
    <div class="page-header">
      <h2>学生管理</h2>
    </div>

    <!-- 搜索栏 -->
    <div class="filter-bar page-card">
      <el-input
        v-model="keyword"
        placeholder="搜索学生"
        :prefix-icon="Search"
        clearable
        style="width: 280px"
        @keyup.enter="fetchStudents"
      />
      <el-button type="primary" @click="fetchStudents">查询</el-button>
    </div>

    <!-- 学生列表 -->
    <div class="page-card">
      <el-table :data="students" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="140" />
        <el-table-column prop="nickname" label="昵称" width="140">
          <template #default="{ row }">
            <span v-if="row.nickname">{{ row.nickname }}</span>
            <span v-else class="text-muted">--</span>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button
              v-if="row.status === 1"
              type="warning"
              link
              size="small"
              @click="handleToggleStatus(row, 0)"
            >禁用</el-button>
            <el-button
              v-else
              type="success"
              link
              size="small"
              @click="handleToggleStatus(row, 1)"
            >启用</el-button>
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
          @size-change="fetchStudents"
          @current-change="fetchStudents"
        />
      </div>
    </div>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑学生信息"
      width="450px"
      destroy-on-close
    >
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="80px">
        <el-form-item label="用户名">
          <el-input :model-value="editForm.username" disabled />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="editForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="editForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="editForm.status" placeholder="请选择状态">
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="editSubmitLoading" @click="handleEditSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getStudents, updateStudent, updateStudentStatus, deleteStudent } from '@/services/admin'

function formatTime(time: string | null | undefined): string {
  if (!time) return '--'
  const date = new Date(time)
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const h = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  const s = String(date.getSeconds()).padStart(2, '0')
  return `${y}-${m}-${d} ${h}:${min}:${s}`
}

const loading = ref(false)
const students = ref<Record<string, unknown>[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const keyword = ref('')

const editDialogVisible = ref(false)
const editSubmitLoading = ref(false)
const editFormRef = ref<FormInstance>()
const editForm = reactive({
  id: 0,
  username: '',
  nickname: '',
  email: '',
  status: 1,
})

const editRules = {
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email' as const, message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
}

async function fetchStudents() {
  loading.value = true
  try {
    const data = await getStudents({
      page: currentPage.value,
      size: pageSize.value,
      keyword: keyword.value || undefined,
    }) as Record<string, unknown>
    students.value = (data.list as Record<string, unknown>[]) || []
    total.value = (data.total as number) || 0
  } catch (error) {
    console.error('获取学生列表失败:', error)
  } finally {
    loading.value = false
  }
}

async function handleToggleStatus(row: Record<string, unknown>, newStatus: number) {
  const action = newStatus === 1 ? '启用' : '禁用'
  try {
    await ElMessageBox.confirm(
      `确定要${action}学生「${row.username}」吗？`,
      `确认${action}`,
      { type: 'warning' }
    )
    await updateStudentStatus(row.id as number, newStatus)
    ElMessage.success(`已${action}`)
    fetchStudents()
  } catch {
    // 用户取消
  }
}

async function handleDelete(row: Record<string, unknown>) {
  try {
    await ElMessageBox.confirm(
      `确定要删除学生「${row.username}」吗？此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    )
    await deleteStudent(row.id as number)
    ElMessage.success('删除成功')
    fetchStudents()
  } catch {
    // 用户取消
  }
}

function handleEdit(row: Record<string, unknown>) {
  editForm.id = row.id as number
  editForm.username = (row.username as string) || ''
  editForm.nickname = (row.nickname as string) || ''
  editForm.email = (row.email as string) || ''
  editForm.status = (row.status as number) ?? 1
  editDialogVisible.value = true
}

async function handleEditSubmit() {
  if (!editFormRef.value) return

  await editFormRef.value.validate(async (valid) => {
    if (!valid) return

    editSubmitLoading.value = true
    try {
      await updateStudent(editForm.id, {
        nickname: editForm.nickname,
        email: editForm.email,
        status: editForm.status,
      })
      ElMessage.success('更新成功')
      editDialogVisible.value = false
      fetchStudents()
    } catch (error) {
      console.error('更新失败:', error)
    } finally {
      editSubmitLoading.value = false
    }
  })
}

onMounted(() => {
  fetchStudents()
})
</script>

<style scoped>
.students-page {
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

.text-muted {
  color: var(--text-muted);
}
</style>
