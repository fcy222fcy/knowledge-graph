import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  // 登录页
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '教师登录', requiresAuth: false },
  },
  // 管理后台路由（需要认证）
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘' },
      },
      {
        path: 'documents',
        name: 'Documents',
        component: () => import('@/views/documents/index.vue'),
        meta: { title: '资料审核' },
      },
      {
        path: 'knowledge',
        name: 'Knowledge',
        component: () => import('@/views/knowledge/index.vue'),
        meta: { title: '知识点管理' },
      },
      {
        path: 'questions',
        name: 'Questions',
        component: () => import('@/views/questions/index.vue'),
        meta: { title: '题目管理' },
      },
      {
        path: 'students',
        name: 'Students',
        component: () => import('@/views/students/index.vue'),
        meta: { title: '学生管理' },
      },
    ],
  },
  // 默认重定向到登录页
  {
    path: '/',
    redirect: '/login',
  },
  // 404 页面
  {
    path: '/:pathMatch(.*)*',
    redirect: '/login',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 全局前置守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  const title = to.meta.title as string
  if (title) {
    document.title = `${title} - SE智图问答 教师端`
  }

  const token = localStorage.getItem('admin_token')

  // 需要认证的页面
  if (to.meta.requiresAuth && !token) {
    next('/login')
    return
  }

  // 已登录用户访问登录页，跳转到管理后台
  if (to.path === '/login' && token) {
    next('/admin')
    return
  }

  next()
})

export default router
