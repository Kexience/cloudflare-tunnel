package svc

import (
	v1 "cloudflared-tunnel/internal/module/credential/ui/api/req/v1"
)

type CredentialSvc interface {
	ValidateCredential(userID int64, req *v1.ValidateCredentialRequest) (*v1.TestResultVO, error)
	CreateCredential(userID int64, req *v1.CreateCredentialRequest) (*v1.CredentialVO, error)
	GetCredential(userID, id int64) (*v1.CredentialVO, error)
	GetCredentials(userID int64) ([]*v1.CredentialVO, error)
	GetDefaultCredential(userID int64) (*v1.CredentialVO, error)
	UpdateCredential(userID, id int64, req *v1.UpdateCredentialRequest) (*v1.CredentialVO, error)
	DeleteCredential(userID, id int64) error
	SetDefaultCredential(userID, id int64) (*v1.CredentialVO, error)
	GetTestLogs(userID, credentialID int64) (*v1.TestLogVO, error)
}
