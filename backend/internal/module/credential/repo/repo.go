package repo

import (
	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/internal/infra/logger"
)

type Repo struct {
	Client *ent.Client
	Log    logger.Logger
}

type CredentialRepo interface {
	CreateCredential(userID int64, name, encryptedToken, accountID string, isDefault bool) (*ent.Credential, error)
	GetCredentialByID(id int64) (*ent.Credential, error)
	GetCredentialByIDAndUserID(id, userID int64) (*ent.Credential, error)
	GetCredentialsByUserID(userID int64) ([]*ent.Credential, error)
	GetDefaultCredentialByUserID(userID int64) (*ent.Credential, error)
	UpdateCredential(id int64, name, encryptedToken, accountID string, isDefault bool) (*ent.Credential, error)
	DeleteCredential(id, userID int64) error
	ClearDefaultByUserID(userID int64) error
	CreateTestLog(credentialID int64, status string, errMsg *string) (*ent.CredentialTestLog, error)
	GetTestLogsByCredentialID(credentialID int64, limit int) ([]*ent.CredentialTestLog, error)
}

func NewCredentialRepo(client *ent.Client, log logger.Logger) CredentialRepo {
	return &Repo{
		Client: client,
		Log:    log,
	}
}
