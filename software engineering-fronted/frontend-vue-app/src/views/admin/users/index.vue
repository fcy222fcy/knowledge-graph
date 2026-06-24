<template>
  <div class="admin-users" v-loading="loading">
    <div class="page-header">
      <h2 class="page-title">用户管理</h2>
      <div class="header-actions">
        <el-input
          v-model="keyword"
          placeholder="搜索用户名/昵称/邮箱"
          clearable
          style="width: 240px; margin-right: 12px;"
          @clear="fetchUsers"
          @keyup.enter="fetchUsers"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="filterRole"
          placeholder="全部角色"
          clearable
          style="width: 120px; margin-right: 12px;"
          @change="fetchUsers"
        >
          <el-option label="管理员" value="admin" />
          <el-option label="教师" value="teacher" />
          <el-option label="学生" value="student" />
        </el-select>
        <el-button type="primary" @click="fetchUsers">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 用户表格 -->
    <el-table :data="users" stripe border style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="nickname" label="昵称" width="120" />
      <el-table-column prop="email" label="邮箱" min-width="180" />
      <el-table-column prop="role" label="角色" width="100">
        <template #default="{ row }">
          <el-tag :type="getRoleType(row.role)">{{ getRoleLabel(row.role) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
          <el-button
            :type="row.status === 1 ? 'warning' : 'success'"
            link
            @click="handleToggleStatus(row)"
            :disabled="row.id === currentUserId"
          >
            {{ row.status === 1 ? '禁用' : '启用' }}
          </el-button>
          <el-dropdown trigger="click" @command="(cmd: string) => handleRoleChange(row, cmd)">
            <el-button type="info" link>
              角色
              <el-icon><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="admin" :disabled="row.role === 'admin'">管理员</el-dropdown-item>
                <el-dropdown-item command="teacher" :disabled="row.role === 'teacher'">教师</el-dropdown-item>
                <el-dropdown-item command="student" :disabled="row.role === 'student'">学生</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-popconfirm
            title="确定要删除该用户吗？"
            confirm-button-text="确定"
            cancel-button-text="取消"
            @confirm="handleDelete(row)"
          >
            <template #reference>
              <el-button type="danger" link :disabled="row.id === currentUserId">删除</el-button>
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
      @size-change="fetchUsers"
      @current-change="fetchUsers"
    />

    <!-- 编辑对话框 -->
    <el-dialog v-model="editVisible" title="编辑用户" width="400px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="昵称">
          <el-input v-model="editForm.nickname" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="editForm.email" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEdit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh, ArrowDown } from '@element-plus/icons-vue'
import { getUsers, updateUser, deleteUser, updateUserStatus, updateUserRole } from '@/services/admin'
import type { UserItem } from '@/services/admin'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const currentUserId = computed(() => userStore.userInfo?.id)

const loading = ref(false)
const submitting = ref(false)
const users = ref<UserItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const keyword = ref('')
const filterRole = ref('')

// 编辑相关
const editVisible = ref(false)
const editUserId = ref(0)
const editForm = ref({ nickname: '', email: '' })

const getRoleType = (role: string) => {
  const types: Record<string, string> = {
    admin: 'danger',
    teacher: 'warning',
    student: 'info'
  }
  return types[role] || 'info'
}

const getRoleLabel = (role: string) => {
  const labels: Record<string, string> = {
    admin: '管理员',
    teacher: '教师',
    student: '学生'
  }
  return labels[role] || '未知'
}

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await getUsers({
      page: page.value,
      size: pageSize.value,
      keyword: keyword.value,
      role: filterRole.value
    })
    users.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: UserItem) => {
  editUserId.value = row.id
  editForm.value = { nickname: row.nickname, email: row.email }
  editVisible.value = true
}

const submitEdit = async () => {
  submitting.value = true
  try {
    await updateUser(editUserId.value, editForm.value)
    ElMessage.success('更新成功')
    editVisible.value = false
    fetchUsers()
  } catch (error) {
    ElMessage.error('更新失败')
  } finally {
    submitting.value = false
  }
}

const handleToggleStatus = async (row: UserItem) => {
  try {
    await updateUserStatus(row.id, row.status === 1 ? 0 : 1)
    ElMessage.success('状态已更新')
    fetchUsers()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleRoleChange = async (row: UserItem, role: string) => {
  try {
    await updateUserRole(row.id, role)
    ElMessage.success('角色已更新')
    fetchUsers()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (row: UserItem) => {
  try {
    await deleteUser(row.id)
    ElMessage.success('删除成功')
    fetchUsers()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.admin-users {
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
</style>
