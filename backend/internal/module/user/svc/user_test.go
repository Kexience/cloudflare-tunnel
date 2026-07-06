package svc

import (
	"context"
	"testing"

	"cloudflared-tunnel/internal/infra/logger"
	"cloudflared-tunnel/internal/module/user/repo"
	v1 "cloudflared-tunnel/internal/module/user/ui/api/req/v1"
	"cloudflared-tunnel/pkg/core"
	"cloudflared-tunnel/pkg/errno"
	"cloudflared-tunnel/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func newTestLogger(t *testing.T) logger.Logger {
	t.Helper()
	log, err := logger.NewLoggerForTest()
	if err != nil {
		t.Fatalf("创建测试 logger 失败: %v", err)
	}
	return log
}

func newTestJWT() *core.JWT {
	return core.NewJWT("test-secret-key", 24)
}

func insertUser(t *testing.T, pg *testutil.PostgresContainer, nickname, username, password, email string) {
	t.Helper()
	ctx := context.Background()
	err := pg.Client.User.Create().
		SetNickname(nickname).
		SetUsername(username).
		SetPassword(password).
		SetEmail(email).
		Exec(ctx)
	require.NoError(t, err, "插入用户数据失败")
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err, "密码加密失败")
	return string(hashed)
}

func TestRegister(t *testing.T) {
	t.Run("成功注册第一个用户", func(t *testing.T) {
		pg := testutil.NewPostgresContainer(t)
		log := newTestLogger(t)
		r := repo.NewUserRepo(pg.Client, log)
		svc := NewUserSvc(r, log, newTestJWT())

		vo, err := svc.Register("管理员", "admin", "password123", "admin@test.com")
		require.NoError(t, err)
		assert.Equal(t, "管理员", vo.Nickname)
		assert.Equal(t, "admin", vo.Username)
		assert.Equal(t, "admin@test.com", vo.Email)
		assert.NotZero(t, vo.ID)
	})

	t.Run("已有用户时注册失败", func(t *testing.T) {
		pg := testutil.NewPostgresContainer(t)
		log := newTestLogger(t)
		insertUser(t, pg, "管理员", "admin", hashPassword(t, "password123"), "admin@test.com")
		r := repo.NewUserRepo(pg.Client, log)
		svc := NewUserSvc(r, log, newTestJWT())

		_, err := svc.Register("新用户", "newuser", "password123", "new@test.com")
		assert.ErrorIs(t, err, errno.ErrUserExists)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		pg := testutil.NewPostgresContainer(t)
		log := newTestLogger(t)
		r := repo.NewUserRepo(pg.Client, log)
		svc := NewUserSvc(r, log, newTestJWT())

		_, err := svc.Register("管理员", "admin", "password123", "admin@test.com")
		require.NoError(t, err)

		_, err = svc.Register("管理员2", "admin", "password456", "admin2@test.com")
		assert.ErrorIs(t, err, errno.ErrUserExists)
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		pg := testutil.NewPostgresContainer(t)
		log := newTestLogger(t)
		r := repo.NewUserRepo(pg.Client, log)
		svc := NewUserSvc(r, log, newTestJWT())

		_, err := svc.Register("管理员", "admin", "password123", "admin@test.com")
		require.NoError(t, err)

		_, err = svc.Register("管理员2", "admin2", "password456", "admin@test.com")
		assert.ErrorIs(t, err, errno.ErrUserExists)
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("成功获取用户", func(t *testing.T) {
		pg := testutil.NewPostgresContainer(t)
		log := newTestLogger(t)
		insertUser(t, pg, "管理员", "admin", hashPassword(t, "password123"), "admin@test.com")
		r := repo.NewUserRepo(pg.Client, log)
		svc := NewUserSvc(r, log, newTestJWT())

		vo, err := svc.GetUserByID("1")
		require.NoError(t, err)
		assert.Equal(t, "管理员", vo.Nickname)
		assert.Equal(t, "admin", vo.Username)
		assert.Equal(t, "admin@test.com", vo.Email)
	})

	t.Run("用户不存在", func(t *testing.T) {
		pg := testutil.NewPostgresContainer(t)
		log := newTestLogger(t)
		r := repo.NewUserRepo(pg.Client, log)
		svc := NewUserSvc(r, log, newTestJWT())

		_, err := svc.GetUserByID("999")
		assert.ErrorIs(t, err, errno.ErrUserNotFound)
	})
}

func TestLogin(t *testing.T) {
	t.Run("成功登录", func(t *testing.T) {
		pg := testutil.NewPostgresContainer(t)
		log := newTestLogger(t)
		jwt := newTestJWT()
		insertUser(t, pg, "管理员", "admin", hashPassword(t, "password123"), "admin@test.com")
		r := repo.NewUserRepo(pg.Client, log)
		svc := NewUserSvc(r, log, jwt)

		loginVO, err := svc.Login("admin", "password123")
		require.NoError(t, err)
		assert.NotEmpty(t, loginVO.Token)
		assert.Equal(t, v1.UserVO{
			ID:       loginVO.User.ID,
			Nickname: "管理员",
			Username: "admin",
			Email:    "admin@test.com",
		}, loginVO.User)

		claims, err := jwt.ParseToken(loginVO.Token)
		require.NoError(t, err)
		assert.Equal(t, loginVO.User.ID, claims.UserID)
		assert.Equal(t, "admin", claims.Username)
	})

	t.Run("用户不存在", func(t *testing.T) {
		pg := testutil.NewPostgresContainer(t)
		log := newTestLogger(t)
		r := repo.NewUserRepo(pg.Client, log)
		svc := NewUserSvc(r, log, newTestJWT())

		_, err := svc.Login("nonexistent", "password123")
		assert.ErrorIs(t, err, errno.ErrUserNotFound)
	})

	t.Run("密码错误", func(t *testing.T) {
		pg := testutil.NewPostgresContainer(t)
		log := newTestLogger(t)
		insertUser(t, pg, "管理员", "admin", hashPassword(t, "password123"), "admin@test.com")
		r := repo.NewUserRepo(pg.Client, log)
		svc := NewUserSvc(r, log, newTestJWT())

		_, err := svc.Login("admin", "wrongpassword")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "密码错误")
	})
}
