import api from '../api'
import type { RegisterRequest, LoginRequest, LoginResponse, User, ApiResponse } from './types'

export async function register(data: RegisterRequest): Promise<ApiResponse<null>> {
  const response = await api.post<ApiResponse<null>>('/v1/user/register', data)
  return response.data
}

export async function login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
  const response = await api.post<ApiResponse<LoginResponse>>('/v1/user/login', data)
  return response.data
}

export async function getCurrentUser(): Promise<ApiResponse<User>> {
  const response = await api.get<ApiResponse<User>>('/v1/user/me')
  return response.data
}