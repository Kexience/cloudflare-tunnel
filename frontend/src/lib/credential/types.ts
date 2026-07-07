export interface ValidateCredentialRequest {
  api_token: string
  account_id: string
}

export interface CredentialVO {
  id: number
  name: string
  api_token: string
  account_id: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T | null
}
