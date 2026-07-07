package ctrl_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloudflared-tunnel/internal/middleware"
	"cloudflared-tunnel/internal/module/credential/repo"
	"cloudflared-tunnel/internal/module/credential/svc"
	"cloudflared-tunnel/internal/module/credential/ui/api/ctrl"
	"cloudflared-tunnel/pkg/cloudflare"
	"cloudflared-tunnel/pkg/core"
	"cloudflared-tunnel/pkg/crypto"
	"cloudflared-tunnel/pkg/errno"
	"cloudflared-tunnel/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "test-secret-key-for-jwt-signing"

func setupRouter(svc svc.CredentialSvc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	jwt := core.NewJWT(testJWTSecret, 24)
	c := ctrl.NewCtrl(svc)

	v1Group := r.Group("/v1")
	authorized := v1Group.Group("/credentials")
	authorized.Use(middleware.Auth(jwt))
	{
		authorized.POST("/validate", c.ValidateCredential)
		authorized.POST("", c.CreateCredential)
		authorized.GET("", c.GetCredentials)
		authorized.GET("/:id", c.GetCredential)
		authorized.PUT("/:id", c.UpdateCredential)
		authorized.DELETE("/:id", c.DeleteCredential)
		authorized.PUT("/:id/default", c.SetDefaultCredential)
	}

	return r
}

func generateToken(t *testing.T) string {
	t.Helper()
	jwt := core.NewJWT(testJWTSecret, 24)
	token, err := jwt.GenerateToken(1, "testuser")
	require.NoError(t, err)
	return token
}

func authHeader(token string) string {
	return "Bearer " + token
}

func TestValidateCredential(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	secret := []byte("test-secret-key-32bytes-long!!!!")

	credentialRepo := repo.NewCredentialRepo(container.Client, log)
	validator := cloudflare.NewValidator()
	credentialSvc := svc.NewCredentialSvc(credentialRepo, log, secret, validator)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("验证失败并记录日志", func(t *testing.T) {
		apiToken := "invalid-token"
		encryptedToken, err := crypto.Encrypt(apiToken, secret)
		require.NoError(t, err)

		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "account-123", false)

		r := setupRouter(credentialSvc)
		token := generateToken(t)

		body := `{"credential_id":` + string(rune('0'+cred.ID)) + `}`
		req := httptest.NewRequest(http.MethodPost, "/v1/credentials/validate", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		logs, err := credentialSvc.GetTestLogs(user.ID, cred.ID)
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, "failed", logs[0].Status)
	})

	t.Run("缺少必填参数", func(t *testing.T) {
		r := setupRouter(credentialSvc)
		token := generateToken(t)

		body := `{}`
		req := httptest.NewRequest(http.MethodPost, "/v1/credentials/validate", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrParam.Code, resp.Code)
	})

	t.Run("未携带 Token", func(t *testing.T) {
		r := setupRouter(credentialSvc)

		body := `{"credential_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/v1/credentials/validate", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrUnauthorized.Code, resp.Code)
	})
}

func parseResponse(t *testing.T, body *bytes.Buffer) core.Response {
	t.Helper()
	var resp core.Response
	err := json.Unmarshal(body.Bytes(), &resp)
	require.NoError(t, err)
	return resp
}
