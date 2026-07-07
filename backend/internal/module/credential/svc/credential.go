package svc

import (
	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/internal/infra/logger"
	"cloudflared-tunnel/internal/module/credential/repo"
	v1 "cloudflared-tunnel/internal/module/credential/ui/api/req/v1"
	"cloudflared-tunnel/pkg/cloudflare"
	"cloudflared-tunnel/pkg/crypto"
	"cloudflared-tunnel/pkg/errno"
)

type svc struct {
	repo      repo.CredentialRepo
	log       logger.Logger
	secret    []byte
	validator cloudflare.Validator
}

func NewCredentialSvc(repo repo.CredentialRepo, log logger.Logger, secret []byte, validator cloudflare.Validator) CredentialSvc {
	return &svc{
		repo:      repo,
		log:       log,
		secret:    secret,
		validator: validator,
	}
}

func (s *svc) CreateCredential(userID int64, req *v1.CreateCredentialRequest) (*v1.CredentialVO, error) {
	if err := s.validator.Validate(req.ApiToken, req.AccountID); err != nil {
		s.log.Error("凭证验证失败", "error", err)
		return nil, errno.ErrCredentialInvalid.WithMessage(err.Error())
	}

	encryptedToken, err := crypto.Encrypt(req.ApiToken, s.secret)
	if err != nil {
		s.log.Error("加密 API Token 失败", "error", err)
		return nil, errno.ErrCredentialEncrypt
	}

	if req.IsDefault {
		if err := s.repo.ClearDefaultByUserID(userID); err != nil {
			s.log.Error("清除默认凭证失败", "error", err)
		}
	}

	cred, err := s.repo.CreateCredential(userID, req.Name, encryptedToken, req.AccountID, req.IsDefault)
	if err != nil {
		return nil, errno.ErrDB
	}

	return s.toVO(cred, req.ApiToken), nil
}

func (s *svc) GetCredential(userID, id int64) (*v1.CredentialVO, error) {
	cred, err := s.repo.GetCredentialByIDAndUserID(id, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errno.ErrCredentialNotFound
		}
		return nil, errno.ErrDB
	}

	decryptedToken, err := crypto.Decrypt(cred.APIToken, s.secret)
	if err != nil {
		s.log.Error("解密 API Token 失败", "id", id, "error", err)
		return nil, errno.ErrCredentialDecrypt
	}

	return s.toVO(cred, decryptedToken), nil
}

func (s *svc) GetCredentials(userID int64) ([]*v1.CredentialVO, error) {
	credentials, err := s.repo.GetCredentialsByUserID(userID)
	if err != nil {
		return nil, errno.ErrDB
	}

	vos := make([]*v1.CredentialVO, len(credentials))
	for i, cred := range credentials {
		decryptedToken, err := crypto.Decrypt(cred.APIToken, s.secret)
		if err != nil {
			s.log.Error("解密 API Token 失败", "id", cred.ID, "error", err)
			decryptedToken = "****"
		}
		vos[i] = s.toVO(cred, decryptedToken)
	}

	return vos, nil
}

func (s *svc) GetDefaultCredential(userID int64) (*v1.CredentialVO, error) {
	cred, err := s.repo.GetDefaultCredentialByUserID(userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errno.ErrCredentialNotFound
		}
		return nil, errno.ErrDB
	}

	decryptedToken, err := crypto.Decrypt(cred.APIToken, s.secret)
	if err != nil {
		s.log.Error("解密 API Token 失败", "id", cred.ID, "error", err)
		return nil, errno.ErrCredentialDecrypt
	}

	return s.toVO(cred, decryptedToken), nil
}

func (s *svc) UpdateCredential(userID, id int64, req *v1.UpdateCredentialRequest) (*v1.CredentialVO, error) {
	_, err := s.repo.GetCredentialByIDAndUserID(id, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errno.ErrCredentialNotFound
		}
		return nil, errno.ErrDB
	}

	if err := s.validator.Validate(req.ApiToken, req.AccountID); err != nil {
		s.log.Error("凭证验证失败", "error", err)
		return nil, errno.ErrCredentialInvalid.WithMessage(err.Error())
	}

	encryptedToken, err := crypto.Encrypt(req.ApiToken, s.secret)
	if err != nil {
		s.log.Error("加密 API Token 失败", "error", err)
		return nil, errno.ErrCredentialEncrypt
	}

	if req.IsDefault {
		if err := s.repo.ClearDefaultByUserID(userID); err != nil {
			s.log.Error("清除默认凭证失败", "error", err)
		}
	}

	cred, err := s.repo.UpdateCredential(id, req.Name, encryptedToken, req.AccountID, req.IsDefault)
	if err != nil {
		return nil, errno.ErrDB
	}

	return s.toVO(cred, req.ApiToken), nil
}

func (s *svc) DeleteCredential(userID, id int64) error {
	_, err := s.repo.GetCredentialByIDAndUserID(id, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return errno.ErrCredentialNotFound
		}
		return errno.ErrDB
	}

	return s.repo.DeleteCredential(id, userID)
}

func (s *svc) SetDefaultCredential(userID, id int64) (*v1.CredentialVO, error) {
	cred, err := s.repo.GetCredentialByIDAndUserID(id, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errno.ErrCredentialNotFound
		}
		return nil, errno.ErrDB
	}

	if err := s.repo.ClearDefaultByUserID(userID); err != nil {
		s.log.Error("清除默认凭证失败", "error", err)
	}

	updated, err := s.repo.UpdateCredential(id, cred.Name, cred.APIToken, cred.AccountID, true)
	if err != nil {
		return nil, errno.ErrDB
	}

	decryptedToken, err := crypto.Decrypt(updated.APIToken, s.secret)
	if err != nil {
		s.log.Error("解密 API Token 失败", "id", id, "error", err)
		return nil, errno.ErrCredentialDecrypt
	}

	return s.toVO(updated, decryptedToken), nil
}

func (s *svc) toVO(cred *ent.Credential, apiToken string) *v1.CredentialVO {
	return &v1.CredentialVO{
		ID:        cred.ID,
		Name:      cred.Name,
		ApiToken:  apiToken,
		AccountID: cred.AccountID,
		IsDefault: cred.IsDefault,
		CreatedAt: cred.CreatedAt,
		UpdatedAt: cred.UpdatedAt,
	}
}
