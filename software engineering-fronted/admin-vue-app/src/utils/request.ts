import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 是否正在跳转登录页，防止重复跳转
let isRedirectingToLogin = false

// 创建 axios 实例
const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 自动添加 token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('admin_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 统一处理响应和错误
request.interceptors.response.use(
  (response) => {
    const res = response.data
    // 后端返回格式: { code: 200, data: ..., message: "success" }
    if (res.code === 200) {
      return res.data
    }
    // 业务错误
    ElMessage.error(res.message || '请求失败')
    return Promise.reject(new Error(res.message || '请求失败'))
  },
  (error) => {
    if (error.response) {
      const { status } = error.response
      switch (status) {
        case 401:
          // token 过期或无效，跳转登录页
          if (!isRedirectingToLogin) {
            isRedirectingToLogin = true
            localStorage.removeItem('admin_token')
            localStorage.removeItem('admin_user')
            ElMessage.error('登录已过期，请重新登录')
            router.push('/login').then(() => {
              isRedirectingToLogin = false
            })
          }
          break
        case 403:
          ElMessage.error('没有权限访问')
          break
        case 500:
          ElMessage.error('服务器错误')
          break
        default:
          ElMessage.error(error.response.data?.message || '请求失败')
      }
    } else {
      ElMessage.error('网络错误')
    }
    return Promise.reject(error)
  }
)

export default request
