<template>
  <div class="app-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <div class="logo">
          <div class="logo-mark">SE</div>
          <div class="logo-text">基于知识图谱的软件工程课程问答平台</div>
        </div>
      </div>

      <nav class="sidebar-nav">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: currentPath === item.path }"
        >
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.title }}</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="user-card" @click="handleLogout">
          <div class="user-avatar">{{ userAvatar }}</div>
          <div class="user-info">
            <div class="user-name">{{ userName }}</div>
            <div class="user-role">{{ roleLabel }}</div>
          </div>
        </div>
      </div>
    </aside>

    <!-- 主内容 -->
    <main class="main" :class="{ 'no-padding': currentPath === '/knowledge-graph' }">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import {
  HomeFilled,
  Share,
  ChatDotRound,
  DataAnalysis,
  Edit
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const currentPath = computed(() => route.path)
const userName = computed(() => userStore.userInfo?.nickname || userStore.userInfo?.username || '用户')
const userAvatar = computed(() => userName.value.charAt(0))

const roleLabel = computed(() => '软件工程专业')

const navItems = [
  { path: '/home', title: '首页', icon: HomeFilled },
  { path: '/knowledge-graph', title: '知识图谱', icon: Share },
  { path: '/qa', title: '问答中心', icon: ChatDotRound },
  { path: '/quiz', title: '答题', icon: Edit },
  { path: '/stats', title: '分析统计', icon: DataAnalysis }
]

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    userStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  }).catch(() => {})
}
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 220px;
  background: #ffffff;
  border-right: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  z-index: 100;
}

.sidebar-header {
  padding: 24px 20px;
  border-bottom: 1px solid #f1f5f9;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-mark {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #2563eb 0%, #3b82f6 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 14px;
}

.logo-text {
  font-size: 16px;
  font-weight: 600;
  color: #0f172a;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  color: #475569;
  cursor: pointer;
  transition: all 0.15s ease;
  font-size: 14px;
  font-weight: 450;
  margin-bottom: 4px;
  text-decoration: none;
}

.nav-item:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.nav-item.active {
  background: #eff6ff;
  color: #2563eb;
  font-weight: 500;
}

.sidebar-footer {
  padding: 16px 12px;
  border-top: 1px solid #f1f5f9;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}

.user-card:hover {
  background: #f1f5f9;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #2563eb 0%, #3b82f6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 15px;
  font-weight: 600;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #0f172a;
}

.user-role {
  font-size: 12px;
  color: #94a3b8;
}

.main {
  margin-left: 220px;
  flex: 1;
  min-height: 100vh;
  padding: 28px;
  background: #f8fafc;
}

.main.no-padding {
  padding: 0;
  display: flex;
  flex-direction: column;
}
</style>
