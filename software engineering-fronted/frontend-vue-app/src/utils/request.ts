import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// Mock 模式开关 - 设为 true 启用 mock 数据，设为 false 使用真实后端
const USE_MOCK = false

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
    // 业务错误
    ElMessage.error(data.message || '请求失败')
    return Promise.reject(new Error(data.message))
  },
  (error) => {
    if (error.response) {
      const { status } = error.response
      if (status === 401) {
        ElMessage.error('登录已过期，请重新登录')
        localStorage.removeItem('token')
        router.push('/login')
      } else if (status === 403) {
        ElMessage.error('没有权限')
      } else if (status === 500) {
        ElMessage.error('服务器错误')
      } else {
        ElMessage.error(error.message || '请求失败')
      }
    } else {
      // 网络错误 - 如果是 mock 模式，静默处理
      if (!USE_MOCK) {
        ElMessage.error('网络错误，请检查网络连接')
      }
    }
    return Promise.reject(error)
  }
)

// 导出 mock 模式状态
export { USE_MOCK }

export default request
