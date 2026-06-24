<template>
  <div class="admin-layout">
    <!-- 左侧菜单 -->
    <aside class="admin-sidebar">
      <div class="sidebar-header">
        <h1 class="logo">SE智图问答</h1>
        <span class="admin-badge">后台管理</span>
      </div>

      <nav class="sidebar-nav">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: isActive(item.path) }"
        >
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="user-card">
          <el-avatar :size="36" :src="userInfo?.avatar">
            {{ userInfo?.nickname?.charAt(0) || userInfo?.username?.charAt(0) }}
          </el-avatar>
          <div class="user-info">
            <div class="user-name">{{ userInfo?.nickname || userInfo?.username }}</div>
            <div class="user-role">{{ roleLabel }}</div>
          </div>
        </div>
        <el-button type="info" text @click="handleLogout">
          <el-icon><SwitchButton /></el-icon>
          退出
        </el-button>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="admin-main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  DataAnalysis,
  User,
  Document,
  Files,
  Collection,
  SwitchButton
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const userInfo = computed(() => userStore.userInfo)

const roleLabel = computed(() => {
  const roles: Record<string, string> = {
    admin: '系统管理员',
    teacher: '教师',
    student: '学生'
  }
  return roles[userInfo.value?.role || 'student'] || '未知'
})

// 导航菜单项
const navItems = computed(() => {
  const items = [
    { path: '/admin', label: '仪表盘', icon: DataAnalysis, roles: ['admin', 'teacher'] },
    { path: '/admin/users', label: '用户管理', icon: User, roles: ['admin'] },
    { path: '/admin/questions', label: '题目管理', icon: Document, roles: ['admin', 'teacher'] },
    { path: '/admin/documents', label: '资料管理', icon: Files, roles: ['admin', 'teacher'] },
    { path: '/admin/knowledge', label: '知识点管理', icon: Collection, roles: ['admin', 'teacher'] }
  ]

  // 根据用户角色过滤菜单项
  return items.filter(item => item.roles.includes(userInfo.value?.role || 'student'))
})

const isActive = (path: string) => {
  return route.path === path
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background-color: #f5f7fa;
}

.admin-sidebar {
  width: 240px;
  background-color: #fff;
  border-right: 1px solid #e4e7ed;
  display: flex;
  flex-direction: column;
  position: fixed;
  height: 100vh;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e4e7ed;
  text-align: center;
}

.logo {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.admin-badge {
  display: inline-block;
  padding: 2px 8px;
  background-color: #409eff;
  color: #fff;
  font-size: 12px;
  border-radius: 4px;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 0;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  color: #606266;
  text-decoration: none;
  transition: all 0.3s;
  gap: 10px;
}

.nav-item:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.nav-item.active {
  background-color: #ecf5ff;
  color: #409eff;
  border-right: 3px solid #409eff;
}

.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-info {
  flex: 1;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.user-role {
  font-size: 12px;
  color: #909399;
}

.admin-main {
  flex: 1;
  margin-left: 240px;
  padding: 24px;
  min-height: 100vh;
}
</style>
