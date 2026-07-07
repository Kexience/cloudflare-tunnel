package svc_test

import (
	"context"
	"testing"

	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/internal/module/credential/repo"
	"cloudflared-tunnel/internal/module/credential/svc"
	v1 "cloudflared-tunnel/internal/module/credential/ui/api/req/v1"
	"cloudflared-tunnel/pkg/cloudflare"
	"cloudflared-tunnel/pkg/crypto"
	"cloudflared-tunnel/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCredential(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	secret := []byte("test-secret-key-32bytes-long!!!!")

	credentialRepo := repo.NewCredentialRepo(container.Client, log)
	validator := cloudflare.NewValidator()
	svc := svc.NewCredentialSvc(credentialRepo, log, secret, validator)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("验证凭证并记录日志", func(t *testing.T) {
		apiToken := "invalid-token"
		encryptedToken, err := crypto.Encrypt(apiToken, secret)
		assert.NoError(t, err)

		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "account-123", false)

		req := &v1.ValidateCredentialRequest{
			CredentialID: cred.ID,
		}
		result, err := svc.ValidateCredential(user.ID, req)
		assert.NoError(t, err)
		assert.False(t, result.Success)
		assert.Equal(t, "API Token 无效", result.Message)

		log, err := svc.GetTestLogs(user.ID, cred.ID)
		assert.NoError(t, err)
		assert.NotNil(t, log)
		assert.Equal(t, "failed", log.Status)
		assert.NotNil(t, log.ErrorMessage)
		assert.Equal(t, "API Token 无效", *log.ErrorMessage)
	})
}

func TestDeleteCredential(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	secret := []byte("test-secret-key-32bytes-long!!!!")

	credentialRepo := repo.NewCredentialRepo(container.Client, log)
	validator := cloudflare.NewValidator()
	credentialSvc := svc.NewCredentialSvc(credentialRepo, log, secret, validator)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("删除无测试日志的凭证", func(t *testing.T) {
		apiToken := "test-token"
		encryptedToken, err := crypto.Encrypt(apiToken, secret)
		require.NoError(t, err)

		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "account-123", false)

		err = credentialSvc.DeleteCredential(user.ID, cred.ID)
		assert.NoError(t, err)

		_, err = container.Client.Credential.Get(context.Background(), cred.ID)
		assert.True(t, ent.IsNotFound(err))
	})

	t.Run("删除有测试日志的凭证", func(t *testing.T) {
		apiToken := "test-token-2"
		encryptedToken, err := crypto.Encrypt(apiToken, secret)
		require.NoError(t, err)

		cred := container.InsertCredential(t, user.ID, "测试凭证2", encryptedToken, "account-456", false)

		_, err = credentialRepo.CreateTestLog(cred.ID, "success", nil)
		require.NoError(t, err)

		_, err = credentialRepo.CreateTestLog(cred.ID, "failed", strPtr("测试错误"))
		require.NoError(t, err)

		err = credentialSvc.DeleteCredential(user.ID, cred.ID)
		assert.NoError(t, err)

		_, err = container.Client.Credential.Get(context.Background(), cred.ID)
		assert.True(t, ent.IsNotFound(err))

		logs, err := credentialRepo.GetTestLogsByCredentialID(cred.ID, 10)
		assert.NoError(t, err)
		assert.Empty(t, logs)
	})

	t.Run("删除不存在的凭证", func(t *testing.T) {
		err := credentialSvc.DeleteCredential(user.ID, 99999)
		assert.NoError(t, err)
	})
}

func strPtr(s string) *string {
	return &s
}
