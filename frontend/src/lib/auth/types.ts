export interface RegisterRequest {
  nickname: string
  username: string
  email: string
  password: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface User {
  id: number
  nickname: string
  username: string
  email: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T | null
}

export interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}