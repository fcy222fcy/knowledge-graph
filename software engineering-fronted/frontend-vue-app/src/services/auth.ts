import request, { USE_MOCK } from '@/utils/request'
import type { LoginParams, RegisterParams, LoginResponse, UserInfo } from '@/types/auth'
import { mockAuth } from './mock'

// 用户登录
export async function login(data: LoginParams) {
  if (USE_MOCK) {
    return mockAuth.login(data) as Promise<any>
  }
  return request.post<LoginResponse>('/student/auth/login', data)
}

// 用户注册
export async function register(data: RegisterParams) {
  if (USE_MOCK) {
    return mockAuth.register(data) as Promise<any>
  }
  return request.post('/student/auth/register', data)
}

// 刷新 Token
export function refreshToken(token: string) {
  return request.post<{ token: string }>('/student/auth/refresh', { token })
}

// 获取当前用户信息
export function getUserProfile() {
  return request.get<UserInfo>('/users/profile')
}

// 更新用户信息
export function updateUserProfile(data: { nickname?: string; avatar?: string }) {
  return request.put('/users/profile', data)
}

// 修改密码
export function changePassword(data: { old_password: string; new_password: string }) {
  return request.post('/users/password', data)
}
