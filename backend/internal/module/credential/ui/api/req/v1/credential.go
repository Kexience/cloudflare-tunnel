package v1

import "time"

type CreateCredentialRequest struct {
	Name      string `json:"name" binding:"required"`
	ApiToken  string `json:"api_token" binding:"required"`
	AccountID string `json:"account_id" binding:"required"`
	IsDefault bool   `json:"is_default"`
}

type UpdateCredentialRequest struct {
	Name      string `json:"name" binding:"required"`
	ApiToken  string `json:"api_token" binding:"required"`
	AccountID string `json:"account_id" binding:"required"`
	IsDefault bool   `json:"is_default"`
}

type ValidateCredentialRequest struct {
	CredentialID int64 `json:"credential_id" binding:"required"`
}

type CredentialVO struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ApiToken  string    `json:"api_token"`
	AccountID string    `json:"account_id"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TestResultVO struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type TestLogVO struct {
	ID           int64     `json:"id"`
	Status       string    `json:"status"`
	ErrorMessage *string   `json:"error_message,omitempty"`
	TestedAt     time.Time `json:"tested_at"`
}
