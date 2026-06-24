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
    meta: { requiresAuth: true, role: 'student' },
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
  // 后台路由（管理员/老师端）
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/dashboard/index.vue'),
        meta: { title: '仪表盘' }
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/users/index.vue'),
        meta: { title: '用户管理', adminOnly: true }
      },
      {
        path: 'questions',
        name: 'AdminQuestions',
        component: () => import('@/views/admin/questions/index.vue'),
        meta: { title: '题目管理' }
      },
      {
        path: 'documents',
        name: 'AdminDocuments',
        component: () => import('@/views/admin/documents/index.vue'),
        meta: { title: '资料管理' }
      },
      {
        path: 'knowledge',
        name: 'AdminKnowledge',
        component: () => import('@/views/admin/knowledge/index.vue'),
        meta: { title: '知识点管理' }
      },
      {
        path: 'system',
        name: 'AdminSystem',
        component: () => import('@/views/admin/system/index.vue'),
        meta: { title: '系统设置', adminOnly: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || ''} - SE智图问答`

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

  const role = userStore.userInfo?.role

  // 需要管理员权限的路由（/admin/*）
  if (to.meta.requiresAdmin) {
    if (role !== 'admin' && role !== 'teacher') {
      next('/home')
      return
    }
  }

  // 学生不能访问 /admin/* 路由
  if (to.path.startsWith('/admin') && role === 'student') {
    next('/home')
    return
  }

  // 管理员/教师默认跳转到后台
  if (to.path === '/' && (role === 'admin' || role === 'teacher')) {
    next('/admin')
    return
  }

  next()
})

export default router
