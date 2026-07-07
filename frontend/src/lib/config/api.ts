import api from '../api'
import type { ApiResponse } from '../auth/types'
import type { Credential, CreateCredentialRequest, UpdateCredentialRequest } from './types'

export async function getCredentials(): Promise<ApiResponse<Credential[]>> {
  const response = await api.get<ApiResponse<Credential[]>>('/v1/credentials')
  return response.data
}

export async function getCredential(id: number): Promise<ApiResponse<Credential>> {
  const response = await api.get<ApiResponse<Credential>>(`/v1/credentials/${id}`)
  return response.data
}

export async function createCredential(data: CreateCredentialRequest): Promise<ApiResponse<Credential>> {
  const response = await api.post<ApiResponse<Credential>>('/v1/credentials', data)
  return response.data
}

export async function updateCredential(id: number, data: UpdateCredentialRequest): Promise<ApiResponse<Credential>> {
  const response = await api.put<ApiResponse<Credential>>(`/v1/credentials/${id}`, data)
  return response.data
}

export async function deleteCredential(id: number): Promise<ApiResponse<null>> {
  const response = await api.delete<ApiResponse<null>>(`/v1/credentials/${id}`)
  return response.data
}

export async function setDefaultCredential(id: number): Promise<ApiResponse<null>> {
  const response = await api.put<ApiResponse<null>>(`/v1/credentials/${id}/default`)
  return response.data
}
