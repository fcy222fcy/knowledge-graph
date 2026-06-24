// 登录请求参数
export interface LoginParams {
  username: string
  password: string
}

// 注册请求参数
export interface RegisterParams {
  username: string
  password: string
  email: string
  nickname?: string
}

// 登录响应
export interface LoginResponse {
  token: string
  user: UserInfo
}

// 用户信息
export interface UserInfo {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  role: 'admin' | 'teacher' | 'student'
  status: number
  created_at: string
  updated_at: string
}
