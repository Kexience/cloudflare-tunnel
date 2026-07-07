export interface Credential {
  id: number
  name: string
  api_token: string
  account_id: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface CreateCredentialRequest {
  name: string
  api_token: string
  account_id: string
  is_default: boolean
}

export interface UpdateCredentialRequest {
  name: string
  api_token: string
  account_id: string
  is_default: boolean
}
