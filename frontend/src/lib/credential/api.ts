import api from '../api'
import type { ValidateCredentialRequest, ApiResponse } from './types'

export async function validateCredential(data: ValidateCredentialRequest): Promise<ApiResponse<null>> {
  const response = await api.post<ApiResponse<null>>('/v1/credentials/validate', data)
  return response.data
}
