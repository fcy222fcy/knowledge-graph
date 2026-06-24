import request from '@/utils/request'
import type { LoginParams, RegisterParams, LoginResponse } from '@/types/auth'

// 教师登录
export function login(data: LoginParams) {
  return request.post<any, LoginResponse>('/teacher/auth/login', data)
}

// 教师注册
export function register(data: RegisterParams) {
  return request.post('/teacher/auth/register', data)
}

// 获取当前用户信息
export function getUserProfile() {
  return request.get('/users/profile')
}

// 更新用户信息
export function updateUserProfile(data: { nickname?: string; email?: string }) {
  return request.put('/users/profile', data)
}

// 修改密码
export function changePassword(data: { old_password: string; new_password: string }) {
  return request.post('/users/password', data)
}
