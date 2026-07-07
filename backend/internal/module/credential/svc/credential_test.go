package svc_test

import (
	"testing"

	"cloudflared-tunnel/internal/module/credential/repo"
	"cloudflared-tunnel/internal/module/credential/svc"
	v1 "cloudflared-tunnel/internal/module/credential/ui/api/req/v1"
	"cloudflared-tunnel/pkg/cloudflare"
	"cloudflared-tunnel/pkg/crypto"
	"cloudflared-tunnel/pkg/testutil"

	"github.com/stretchr/testify/assert"
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

		logs, err := svc.GetTestLogs(user.ID, cred.ID)
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, "failed", logs[0].Status)
		assert.NotNil(t, logs[0].ErrorMessage)
		assert.Equal(t, "API Token 无效", *logs[0].ErrorMessage)
	})
}
