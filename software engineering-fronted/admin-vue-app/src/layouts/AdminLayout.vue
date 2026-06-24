<template>
  <div class="admin-layout">
    <!-- 左侧边栏 -->
    <aside class="sidebar">
      <!-- Logo 区域 -->
      <div class="sidebar-header">
        <div class="logo">
          <span class="logo-mark">SE</span>
          <div class="logo-text">
            <span class="title">SE智图问答</span>
            <span class="badge">教师端</span>
          </div>
        </div>
      </div>

      <!-- 导航菜单 -->
      <nav class="sidebar-nav">
        <router-link
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: isActive(item.path) }"
        >
          <el-icon :size="20"><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </router-link>
      </nav>

      <!-- 用户信息 -->
      <div class="sidebar-footer">
        <div class="user-card">
          <div class="avatar">
            {{ userStore.userInfo?.nickname?.charAt(0) || 'T' }}
          </div>
          <div class="user-info">
            <span class="name">{{ userStore.userInfo?.nickname || '教师用户' }}</span>
            <span class="role">教师</span>
          </div>
          <el-button
            type="danger"
            :icon="SwitchButton"
            circle
            size="small"
            @click="handleLogout"
          />
        </div>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import {
  DataAnalysis,
  Document,
  Share,
  Edit,
  User,
  SwitchButton,
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 侧边栏菜单项
const menuItems = [
  { path: '/admin', label: '仪表盘', icon: DataAnalysis },
  { path: '/admin/documents', label: '资料审核', icon: Document },
  { path: '/admin/knowledge', label: '知识点管理', icon: Share },
  { path: '/admin/questions', label: '题目管理', icon: Edit },
  { path: '/admin/students', label: '学生管理', icon: User },
]

// 判断菜单是否激活
function isActive(path: string) {
  if (path === '/admin') {
    return route.path === '/admin'
  }
  return route.path.startsWith(path)
}

// 退出登录
function handleLogout() {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    userStore.logout()
  }).catch(() => {})
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
}

/* 侧边栏 */
.sidebar {
  width: var(--sidebar-width);
  background: var(--bg-sidebar);
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 100;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-mark {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, var(--color-primary), var(--color-info));
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 16px;
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-text .title {
  color: white;
  font-size: 16px;
  font-weight: 600;
}

.logo-text .badge {
  color: var(--color-primary-light);
  font-size: 12px;
}

/* 导航菜单 */
.sidebar-nav {
  flex: 1;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: var(--radius-md);
  color: var(--text-sidebar);
  text-decoration: none;
  transition: all 0.2s;
}

.nav-item:hover {
  background: var(--bg-sidebar-hover);
  color: var(--text-sidebar-active);
}

.nav-item.active {
  background: var(--color-primary);
  color: var(--text-sidebar-active);
}

/* 用户卡片 */
.sidebar-footer {
  padding: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.user-card {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--color-primary), var(--color-info));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 14px;
}

.user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.user-info .name {
  color: white;
  font-size: 14px;
  font-weight: 500;
}

.user-info .role {
  color: var(--text-muted);
  font-size: 12px;
}

/* 主内容区 */
.main-content {
  flex: 1;
  margin-left: var(--sidebar-width);
  padding: 28px;
  background: var(--bg-main);
  min-height: 100vh;
}
</style>
