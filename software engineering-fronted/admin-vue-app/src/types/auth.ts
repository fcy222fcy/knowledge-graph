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

// 登录响应（后端返回 teacher 字段）
export interface LoginResponse {
  token: string
  teacher: UserInfo
}

// 用户信息
export interface UserInfo {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  status: number
  created_at: string
  updated_at: string
}
