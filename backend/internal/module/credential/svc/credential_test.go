package svc

import (
	"testing"

	"cloudflared-tunnel/internal/infra/logger"
	"cloudflared-tunnel/internal/module/credential/repo"
	v1 "cloudflared-tunnel/internal/module/credential/ui/api/req/v1"
	"cloudflared-tunnel/pkg/crypto"
	"cloudflared-tunnel/pkg/errno"
	"cloudflared-tunnel/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger(t *testing.T) logger.Logger {
	t.Helper()
	log, err := logger.NewLoggerForTest()
	if err != nil {
		t.Fatalf("创建测试 logger 失败: %v", err)
	}
	return log
}

func newTestSecret() []byte {
	return []byte("0123456789abcdef0123456789abcdef")
}

func TestCreateCredential(t *testing.T) {
	t.Run("成功创建凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		req := &v1.CreateCredentialRequest{
			Name:      "生产环境",
			ApiToken:  "cf_api_token_123456",
			AccountID: "account-123",
			IsDefault: true,
		}

		vo, err := svc.CreateCredential(user.ID, req)
		require.NoError(t, err)
		assert.Equal(t, "生产环境", vo.Name)
		assert.Equal(t, "cf_api_token_123456", vo.ApiToken)
		assert.Equal(t, "account-123", vo.AccountID)
		assert.True(t, vo.IsDefault)
		assert.NotZero(t, vo.ID)
		assert.NotZero(t, vo.CreatedAt)
	})

	t.Run("创建第二个凭证时设置默认会清除之前的默认", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		req1 := &v1.CreateCredentialRequest{
			Name:      "生产环境",
			ApiToken:  "token1",
			AccountID: "account-1",
			IsDefault: true,
		}
		_, err := svc.CreateCredential(user.ID, req1)
		require.NoError(t, err)

		req2 := &v1.CreateCredentialRequest{
			Name:      "测试环境",
			ApiToken:  "token2",
			AccountID: "account-2",
			IsDefault: true,
		}
		vo2, err := svc.CreateCredential(user.ID, req2)
		require.NoError(t, err)
		assert.True(t, vo2.IsDefault)

		credentials, err := svc.GetCredentials(user.ID)
		require.NoError(t, err)
		assert.Len(t, credentials, 2)

		defaultCount := 0
		for _, cred := range credentials {
			if cred.IsDefault {
				defaultCount++
			}
		}
		assert.Equal(t, 1, defaultCount)
	})
}

func TestGetCredential(t *testing.T) {
	t.Run("成功获取凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		encryptedToken, _ := crypto.Encrypt("cf_token_123", secret)
		db.InsertCredential(t, user.ID, "生产环境", encryptedToken, "account-123", true)
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		vo, err := svc.GetCredential(user.ID, 1)
		require.NoError(t, err)
		assert.Equal(t, "生产环境", vo.Name)
		assert.Equal(t, "cf_token_123", vo.ApiToken)
		assert.Equal(t, "account-123", vo.AccountID)
		assert.True(t, vo.IsDefault)
	})

	t.Run("凭证不存在", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		_, err := svc.GetCredential(user.ID, 999)
		assert.ErrorIs(t, err, errno.ErrCredentialNotFound)
	})

	t.Run("无权访问其他用户的凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user1 := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		user2 := db.InsertUser(t, "用户2", "user2", "hashed", "user2@test.com")
		encryptedToken, _ := crypto.Encrypt("cf_token", secret)
		db.InsertCredential(t, user1.ID, "生产环境", encryptedToken, "account-123", true)
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		_, err := svc.GetCredential(user2.ID, 1)
		assert.ErrorIs(t, err, errno.ErrCredentialNotFound)
	})
}

func TestGetCredentials(t *testing.T) {
	t.Run("成功获取凭证列表", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		encryptedToken1, _ := crypto.Encrypt("token1", secret)
		encryptedToken2, _ := crypto.Encrypt("token2", secret)
		db.InsertCredential(t, user.ID, "生产环境", encryptedToken1, "account-1", true)
		db.InsertCredential(t, user.ID, "测试环境", encryptedToken2, "account-2", false)
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		vos, err := svc.GetCredentials(user.ID)
		require.NoError(t, err)
		assert.Len(t, vos, 2)
		assert.Equal(t, "token1", vos[0].ApiToken)
		assert.Equal(t, "token2", vos[1].ApiToken)
	})

	t.Run("用户没有凭证时返回空列表", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		vos, err := svc.GetCredentials(user.ID)
		require.NoError(t, err)
		assert.Empty(t, vos)
	})
}

func TestGetDefaultCredential(t *testing.T) {
	t.Run("成功获取默认凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		encryptedToken, _ := crypto.Encrypt("default_token", secret)
		db.InsertCredential(t, user.ID, "生产环境", encryptedToken, "account-123", true)
		encryptedToken2, _ := crypto.Encrypt("other_token", secret)
		db.InsertCredential(t, user.ID, "测试环境", encryptedToken2, "account-456", false)
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		vo, err := svc.GetDefaultCredential(user.ID)
		require.NoError(t, err)
		assert.Equal(t, "生产环境", vo.Name)
		assert.Equal(t, "default_token", vo.ApiToken)
		assert.True(t, vo.IsDefault)
	})

	t.Run("没有默认凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		encryptedToken, _ := crypto.Encrypt("token", secret)
		db.InsertCredential(t, user.ID, "测试环境", encryptedToken, "account-123", false)
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		_, err := svc.GetDefaultCredential(user.ID)
		assert.ErrorIs(t, err, errno.ErrCredentialNotFound)
	})
}

func TestUpdateCredential(t *testing.T) {
	t.Run("成功更新凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		encryptedToken, _ := crypto.Encrypt("old_token", secret)
		db.InsertCredential(t, user.ID, "生产环境", encryptedToken, "account-123", true)
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		req := &v1.UpdateCredentialRequest{
			Name:      "生产环境-v2",
			ApiToken:  "new_token",
			AccountID: "account-456",
			IsDefault: true,
		}

		vo, err := svc.UpdateCredential(user.ID, 1, req)
		require.NoError(t, err)
		assert.Equal(t, "生产环境-v2", vo.Name)
		assert.Equal(t, "new_token", vo.ApiToken)
		assert.Equal(t, "account-456", vo.AccountID)
		assert.True(t, vo.IsDefault)
	})

	t.Run("更新不存在的凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		req := &v1.UpdateCredentialRequest{
			Name:      "test",
			ApiToken:  "token",
			AccountID: "account",
			IsDefault: false,
		}

		_, err := svc.UpdateCredential(user.ID, 999, req)
		assert.ErrorIs(t, err, errno.ErrCredentialNotFound)
	})
}

func TestDeleteCredential(t *testing.T) {
	t.Run("成功删除凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		encryptedToken, _ := crypto.Encrypt("token", secret)
		db.InsertCredential(t, user.ID, "生产环境", encryptedToken, "account-123", true)
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		err := svc.DeleteCredential(user.ID, 1)
		require.NoError(t, err)

		_, err = svc.GetCredential(user.ID, 1)
		assert.ErrorIs(t, err, errno.ErrCredentialNotFound)
	})

	t.Run("删除不存在的凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		err := svc.DeleteCredential(user.ID, 999)
		assert.ErrorIs(t, err, errno.ErrCredentialNotFound)
	})

	t.Run("无权删除其他用户的凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user1 := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		user2 := db.InsertUser(t, "用户2", "user2", "hashed", "user2@test.com")
		encryptedToken, _ := crypto.Encrypt("token", secret)
		db.InsertCredential(t, user1.ID, "生产环境", encryptedToken, "account-123", true)
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		err := svc.DeleteCredential(user2.ID, 1)
		assert.ErrorIs(t, err, errno.ErrCredentialNotFound)
	})
}

func TestSetDefaultCredential(t *testing.T) {
	t.Run("成功设置默认凭证", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		encryptedToken1, _ := crypto.Encrypt("token1", secret)
		encryptedToken2, _ := crypto.Encrypt("token2", secret)
		db.InsertCredential(t, user.ID, "生产环境", encryptedToken1, "account-1", true)
		db.InsertCredential(t, user.ID, "测试环境", encryptedToken2, "account-2", false)
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		vo, err := svc.SetDefaultCredential(user.ID, 2)
		require.NoError(t, err)
		assert.Equal(t, "测试环境", vo.Name)
		assert.True(t, vo.IsDefault)

		defaultCred, err := svc.GetDefaultCredential(user.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(2), defaultCred.ID)
	})

	t.Run("设置不存在的凭证为默认", func(t *testing.T) {
		db := testutil.NewSQLiteContainer(t)
		log := newTestLogger(t)
		secret := newTestSecret()
		user := db.InsertUser(t, "管理员", "admin", "hashed", "admin@test.com")
		r := repo.NewCredentialRepo(db.Client, log)
		svc := NewCredentialSvc(r, log, secret)

		_, err := svc.SetDefaultCredential(user.ID, 999)
		assert.ErrorIs(t, err, errno.ErrCredentialNotFound)
	})
}
