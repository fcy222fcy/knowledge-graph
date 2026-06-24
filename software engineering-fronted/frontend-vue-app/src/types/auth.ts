export interface LoginParams {
  username: string
  password: string
}

export interface RegisterParams {
  username: string
  password: string
  email: string
  nickname?: string
}

export interface LoginResponse {
  token: string
  user: {
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
}

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
