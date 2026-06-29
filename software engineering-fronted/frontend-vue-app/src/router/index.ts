import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { title: '注册', requiresAuth: false }
  },
  // 前台路由（学生端）
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: 'home',
        name: 'Home',
        component: () => import('@/views/home/index.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'knowledge-graph',
        name: 'KnowledgeGraph',
        component: () => import('@/views/knowledge-graph/index.vue'),
        meta: { title: '知识图谱' }
      },
      {
        path: 'qa',
        name: 'QA',
        component: () => import('@/views/qa/index.vue'),
        meta: { title: '问答中心' }
      },
      {
        path: 'files',
        name: 'Files',
        component: () => import('@/views/files/index.vue'),
        meta: { title: '资料管理' }
      },
      {
        path: 'stats',
        name: 'Stats',
        component: () => import('@/views/stats/index.vue'),
        meta: { title: '分析统计' }
      },
      {
        path: 'quiz',
        name: 'Quiz',
        component: () => import('@/views/quiz/index.vue'),
        meta: { title: '答题' }
      }
    ]
  },
  // Admin 路由（教师权限）
  {
    path: '/admin',
    component: MainLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: 'knowledge-graph',
        name: 'AdminKnowledgeGraph',
        component: () => import('../views/admin/knowledge-graph/index.vue'),
        meta: { title: '知识图谱管理', requiresAuth: true, roles: ['teacher'] }
      }
    ]
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || ''} - 基于知识图谱的软件工程课程问答平台`

  const userStore = useUserStore()

  // 检查是否需要登录
  if (to.meta.requiresAuth && !userStore.token) {
    next('/login')
    return
  }

  // 已登录用户不允许访问登录/注册页，重定向到首页
  if ((to.path === '/login' || to.path === '/register') && userStore.token) {
    next('/home')
    return
  }

  next()
})

export default router
