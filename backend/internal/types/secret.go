package types

// CredentialSecret 凭证加密密钥的具名类型，避免 FX 注入歧义
type CredentialSecret []byte
