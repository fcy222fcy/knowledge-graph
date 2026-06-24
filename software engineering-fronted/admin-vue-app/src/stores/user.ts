import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { UserInfo } from '@/types/auth'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  // 从 localStorage 恢复状态
  const token = ref<string>(localStorage.getItem('admin_token') || '')
  const userInfo = ref<UserInfo | null>(
    JSON.parse(localStorage.getItem('admin_user') || 'null')
  )

  // 设置 token
  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('admin_token', newToken)
  }

  // 清除 token
  function clearToken() {
    token.value = ''
    localStorage.removeItem('admin_token')
  }

  // 设置用户信息
  function setUserInfo(info: UserInfo) {
    userInfo.value = info
    localStorage.setItem('admin_user', JSON.stringify(info))
  }

  // 清除用户信息
  function clearUserInfo() {
    userInfo.value = null
    localStorage.removeItem('admin_user')
  }

  // 退出登录
  function logout() {
    clearToken()
    clearUserInfo()
    router.push('/login')
  }

  return {
    token,
    userInfo,
    setToken,
    clearToken,
    setUserInfo,
    clearUserInfo,
    logout,
  }
})
