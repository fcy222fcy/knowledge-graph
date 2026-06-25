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
        <el-table-column prop="created_at" label="注册时间" width="180" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getStudents, updateStudentStatus, deleteStudent } from '@/services/admin'

const loading = ref(false)
const students = ref<Record<string, unknown>[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const keyword = ref('')

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

onMounted(() => {
  fetchStudents()
})
</script>

<style scoped>
.students-page {
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
