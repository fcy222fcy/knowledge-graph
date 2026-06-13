import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// Mock 模式开关 - 设为 true 启用 mock 数据，设为 false 使用真实后端
const USE_MOCK = false

// 防止 401 重定向重复触发
let isRedirectingToLogin = false

// 创建 axios 实例
const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器：自动携带 Token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：统一错误处理
request.interceptors.response.use(
  (response) => {
    const { data } = response
    // 假设后端返回格式为 { code: number, data: any, message: string }
    if (data.code === 200) {
      return data
    }
    // 业务错误 - 直接弹窗后 reject，不走 error 拦截器
    ElMessage.error(data.message || '请求失败')
    const err = new Error(data.message || '请求失败')
    ;(err as any).isBusinessError = true
    return Promise.reject(err)
  },
  (error) => {
    // 如果是业务错误（已被 success 拦截器处理），不再重复弹窗
    if (error?.isBusinessError) {
      return Promise.reject(error)
    }

    if (error.response) {
      const { status } = error.response
      if (status === 401) {
        if (!isRedirectingToLogin) {
          isRedirectingToLogin = true
          ElMessage.error('登录已过期，请重新登录')
          localStorage.removeItem('token')
          router.push('/login').finally(() => {
            isRedirectingToLogin = false
          })
        }
      } else if (status === 403) {
        ElMessage.error('没有权限')
      } else if (status === 500) {
        ElMessage.error('服务器错误')
      } else {
        ElMessage.error(error.response.data?.message || error.message || '请求失败')
      }
    } else {
      // 网络错误 - 如果是 mock 模式，静默处理
      if (!USE_MOCK) {
        // 超时和网络错误只弹一次
        if (error.code !== 'ERR_NETWORK' || !error.config?.__retryShown) {
          if (error.config) error.config.__retryShown = true
          ElMessage.error('网络错误，请检查网络连接')
        }
      }
    }
    return Promise.reject(error)
  }
)

// 导出 mock 模式状态
export { USE_MOCK }

export default request
