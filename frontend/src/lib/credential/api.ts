import api from '../api'
import type { ValidateCredentialRequest, TestResultVO, ApiResponse } from './types'

export async function validateCredential(data: ValidateCredentialRequest): Promise<ApiResponse<TestResultVO>> {
  const response = await api.post<ApiResponse<TestResultVO>>('/v1/credentials/validate', data)
  return response.data
}
