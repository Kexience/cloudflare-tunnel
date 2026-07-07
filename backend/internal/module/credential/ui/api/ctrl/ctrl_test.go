package ctrl

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloudflared-tunnel/internal/middleware"
	v1 "cloudflared-tunnel/internal/module/credential/ui/api/req/v1"
	"cloudflared-tunnel/pkg/core"
	"cloudflared-tunnel/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCredentialSvc struct {
	CreateCredentialFn     func(userID int64, req *v1.CreateCredentialRequest) (*v1.CredentialVO, error)
	GetCredentialFn        func(userID, id int64) (*v1.CredentialVO, error)
	GetCredentialsFn       func(userID int64) ([]*v1.CredentialVO, error)
	GetDefaultCredentialFn func(userID int64) (*v1.CredentialVO, error)
	UpdateCredentialFn     func(userID, id int64, req *v1.UpdateCredentialRequest) (*v1.CredentialVO, error)
	DeleteCredentialFn     func(userID, id int64) error
	SetDefaultCredentialFn func(userID, id int64) (*v1.CredentialVO, error)
}

func (m *mockCredentialSvc) CreateCredential(userID int64, req *v1.CreateCredentialRequest) (*v1.CredentialVO, error) {
	return m.CreateCredentialFn(userID, req)
}

func (m *mockCredentialSvc) GetCredential(userID, id int64) (*v1.CredentialVO, error) {
	return m.GetCredentialFn(userID, id)
}

func (m *mockCredentialSvc) GetCredentials(userID int64) ([]*v1.CredentialVO, error) {
	return m.GetCredentialsFn(userID)
}

func (m *mockCredentialSvc) GetDefaultCredential(userID int64) (*v1.CredentialVO, error) {
	return m.GetDefaultCredentialFn(userID)
}

func (m *mockCredentialSvc) UpdateCredential(userID, id int64, req *v1.UpdateCredentialRequest) (*v1.CredentialVO, error) {
	return m.UpdateCredentialFn(userID, id, req)
}

func (m *mockCredentialSvc) DeleteCredential(userID, id int64) error {
	return m.DeleteCredentialFn(userID, id)
}

func (m *mockCredentialSvc) SetDefaultCredential(userID, id int64) (*v1.CredentialVO, error) {
	return m.SetDefaultCredentialFn(userID, id)
}

const testJWTSecret = "test-secret-key-for-jwt-signing"
const testUserID int64 = 42
const testUsername = "testuser"

func setupRouter(mock *mockCredentialSvc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	jwt := core.NewJWT(testJWTSecret, 24)
	c := NewCtrl(mock)
	router := &struct {
		ctrl *Ctrl
		jwt  *core.JWT
	}{ctrl: c, jwt: jwt}

	v1Group := r.Group("/v1")
	authorized := v1Group.Group("/credentials")
	authorized.Use(middleware.Auth(router.jwt))
	{
		authorized.POST("", router.ctrl.CreateCredential)
		authorized.GET("", router.ctrl.GetCredentials)
		authorized.GET("/:id", router.ctrl.GetCredential)
		authorized.PUT("/:id", router.ctrl.UpdateCredential)
		authorized.DELETE("/:id", router.ctrl.DeleteCredential)
		authorized.PUT("/:id/default", router.ctrl.SetDefaultCredential)
	}

	return r
}

func generateToken(t *testing.T) string {
	t.Helper()
	jwt := core.NewJWT(testJWTSecret, 24)
	token, err := jwt.GenerateToken(testUserID, testUsername)
	require.NoError(t, err)
	return token
}

func authHeader(token string) string {
	return "Bearer " + token
}

func parseResponse(t *testing.T, body *bytes.Buffer) core.Response {
	t.Helper()
	var resp core.Response
	err := json.Unmarshal(body.Bytes(), &resp)
	require.NoError(t, err)
	return resp
}

func parseDataToVO(t *testing.T, data any) *v1.CredentialVO {
	t.Helper()
	b, err := json.Marshal(data)
	require.NoError(t, err)
	var vo v1.CredentialVO
	err = json.Unmarshal(b, &vo)
	require.NoError(t, err)
	return &vo
}

func parseDataToVOList(t *testing.T, data any) []*v1.CredentialVO {
	t.Helper()
	b, err := json.Marshal(data)
	require.NoError(t, err)
	var vos []*v1.CredentialVO
	err = json.Unmarshal(b, &vos)
	require.NoError(t, err)
	return vos
}

func TestCreateCredential(t *testing.T) {
	t.Run("成功创建凭证", func(t *testing.T) {
		now := time.Now()
		mock := &mockCredentialSvc{
			CreateCredentialFn: func(userID int64, req *v1.CreateCredentialRequest) (*v1.CredentialVO, error) {
				assert.Equal(t, testUserID, userID)
				assert.Equal(t, "生产环境", req.Name)
				assert.Equal(t, "cf_token_123", req.ApiToken)
				assert.Equal(t, "account-123", req.AccountID)
				assert.True(t, req.IsDefault)
				return &v1.CredentialVO{
					ID:        1,
					Name:      req.Name,
					ApiToken:  req.ApiToken,
					AccountID: req.AccountID,
					IsDefault: req.IsDefault,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
		}

		r := setupRouter(mock)
		token := generateToken(t)

		body := `{"name":"生产环境","api_token":"cf_token_123","account_id":"account-123","is_default":true}`
		req := httptest.NewRequest(http.MethodPost, "/v1/credentials", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, 0, resp.Code)

		vo := parseDataToVO(t, resp.Data)
		assert.Equal(t, int64(1), vo.ID)
		assert.Equal(t, "生产环境", vo.Name)
		assert.Equal(t, "cf_token_123", vo.ApiToken)
		assert.Equal(t, "account-123", vo.AccountID)
		assert.True(t, vo.IsDefault)
	})

	t.Run("缺少必填参数", func(t *testing.T) {
		mock := &mockCredentialSvc{}
		r := setupRouter(mock)
		token := generateToken(t)

		body := `{"name":"生产环境"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/credentials", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrParam.Code, resp.Code)
	})

	t.Run("未携带 Token", func(t *testing.T) {
		mock := &mockCredentialSvc{}
		r := setupRouter(mock)

		body := `{"name":"test","api_token":"token","account_id":"acc"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/credentials", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrUnauthorized.Code, resp.Code)
	})

	t.Run("业务错误返回", func(t *testing.T) {
		mock := &mockCredentialSvc{
			CreateCredentialFn: func(userID int64, req *v1.CreateCredentialRequest) (*v1.CredentialVO, error) {
				return nil, errno.ErrCredentialEncrypt
			},
		}
		r := setupRouter(mock)
		token := generateToken(t)

		body := `{"name":"test","api_token":"token","account_id":"acc"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/credentials", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrCredentialEncrypt.Code, resp.Code)
	})

	t.Run("凭证验证失败", func(t *testing.T) {
		mock := &mockCredentialSvc{
			CreateCredentialFn: func(userID int64, req *v1.CreateCredentialRequest) (*v1.CredentialVO, error) {
				return nil, errno.ErrCredentialInvalid.WithMessage("API Token 无效")
			},
		}
		r := setupRouter(mock)
		token := generateToken(t)

		body := `{"name":"test","api_token":"bad_token","account_id":"acc"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/credentials", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrCredentialInvalid.Code, resp.Code)
		assert.Equal(t, "API Token 无效", resp.Message)
	})
}

func TestGetCredential(t *testing.T) {
	t.Run("成功获取凭证", func(t *testing.T) {
		now := time.Now()
		mock := &mockCredentialSvc{
			GetCredentialFn: func(userID, id int64) (*v1.CredentialVO, error) {
				assert.Equal(t, testUserID, userID)
				assert.Equal(t, int64(1), id)
				return &v1.CredentialVO{
					ID:        1,
					Name:      "生产环境",
					ApiToken:  "cf_token_123",
					AccountID: "account-123",
					IsDefault: true,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
		}

		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodGet, "/v1/credentials/1", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, 0, resp.Code)

		vo := parseDataToVO(t, resp.Data)
		assert.Equal(t, int64(1), vo.ID)
		assert.Equal(t, "生产环境", vo.Name)
	})

	t.Run("无效的凭证ID", func(t *testing.T) {
		mock := &mockCredentialSvc{}
		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodGet, "/v1/credentials/abc", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrParam.Code, resp.Code)
	})

	t.Run("凭证不存在", func(t *testing.T) {
		mock := &mockCredentialSvc{
			GetCredentialFn: func(userID, id int64) (*v1.CredentialVO, error) {
				return nil, errno.ErrCredentialNotFound
			},
		}
		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodGet, "/v1/credentials/999", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrCredentialNotFound.Code, resp.Code)
	})
}

func TestGetCredentials(t *testing.T) {
	t.Run("成功获取凭证列表", func(t *testing.T) {
		now := time.Now()
		mock := &mockCredentialSvc{
			GetCredentialsFn: func(userID int64) ([]*v1.CredentialVO, error) {
				assert.Equal(t, testUserID, userID)
				return []*v1.CredentialVO{
					{ID: 1, Name: "生产环境", ApiToken: "token1", AccountID: "acc-1", IsDefault: true, CreatedAt: now, UpdatedAt: now},
					{ID: 2, Name: "测试环境", ApiToken: "token2", AccountID: "acc-2", IsDefault: false, CreatedAt: now, UpdatedAt: now},
				}, nil
			},
		}

		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodGet, "/v1/credentials", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, 0, resp.Code)

		vos := parseDataToVOList(t, resp.Data)
		assert.Len(t, vos, 2)
		assert.Equal(t, "生产环境", vos[0].Name)
		assert.Equal(t, "测试环境", vos[1].Name)
	})

	t.Run("空列表", func(t *testing.T) {
		mock := &mockCredentialSvc{
			GetCredentialsFn: func(userID int64) ([]*v1.CredentialVO, error) {
				return []*v1.CredentialVO{}, nil
			},
		}

		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodGet, "/v1/credentials", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, 0, resp.Code)

		vos := parseDataToVOList(t, resp.Data)
		assert.Empty(t, vos)
	})
}

func TestUpdateCredential(t *testing.T) {
	t.Run("成功更新凭证", func(t *testing.T) {
		now := time.Now()
		mock := &mockCredentialSvc{
			UpdateCredentialFn: func(userID, id int64, req *v1.UpdateCredentialRequest) (*v1.CredentialVO, error) {
				assert.Equal(t, testUserID, userID)
				assert.Equal(t, int64(1), id)
				assert.Equal(t, "生产环境-v2", req.Name)
				return &v1.CredentialVO{
					ID:        1,
					Name:      req.Name,
					ApiToken:  req.ApiToken,
					AccountID: req.AccountID,
					IsDefault: req.IsDefault,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
		}

		r := setupRouter(mock)
		token := generateToken(t)

		body := `{"name":"生产环境-v2","api_token":"new_token","account_id":"acc-456","is_default":true}`
		req := httptest.NewRequest(http.MethodPut, "/v1/credentials/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, 0, resp.Code)

		vo := parseDataToVO(t, resp.Data)
		assert.Equal(t, "生产环境-v2", vo.Name)
		assert.Equal(t, "new_token", vo.ApiToken)
	})

	t.Run("无效的凭证ID", func(t *testing.T) {
		mock := &mockCredentialSvc{}
		r := setupRouter(mock)
		token := generateToken(t)

		body := `{"name":"test","api_token":"token","account_id":"acc"}`
		req := httptest.NewRequest(http.MethodPut, "/v1/credentials/abc", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrParam.Code, resp.Code)
	})

	t.Run("更新不存在的凭证", func(t *testing.T) {
		mock := &mockCredentialSvc{
			UpdateCredentialFn: func(userID, id int64, req *v1.UpdateCredentialRequest) (*v1.CredentialVO, error) {
				return nil, errno.ErrCredentialNotFound
			},
		}
		r := setupRouter(mock)
		token := generateToken(t)

		body := `{"name":"test","api_token":"token","account_id":"acc"}`
		req := httptest.NewRequest(http.MethodPut, "/v1/credentials/999", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrCredentialNotFound.Code, resp.Code)
	})

	t.Run("凭证验证失败", func(t *testing.T) {
		mock := &mockCredentialSvc{
			UpdateCredentialFn: func(userID, id int64, req *v1.UpdateCredentialRequest) (*v1.CredentialVO, error) {
				return nil, errno.ErrCredentialInvalid.WithMessage("账号 ID 无效或无权访问")
			},
		}
		r := setupRouter(mock)
		token := generateToken(t)

		body := `{"name":"test","api_token":"token","account_id":"bad_account"}`
		req := httptest.NewRequest(http.MethodPut, "/v1/credentials/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrCredentialInvalid.Code, resp.Code)
		assert.Equal(t, "账号 ID 无效或无权访问", resp.Message)
	})
}

func TestDeleteCredential(t *testing.T) {
	t.Run("成功删除凭证", func(t *testing.T) {
		mock := &mockCredentialSvc{
			DeleteCredentialFn: func(userID, id int64) error {
				assert.Equal(t, testUserID, userID)
				assert.Equal(t, int64(1), id)
				return nil
			},
		}

		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodDelete, "/v1/credentials/1", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, 0, resp.Code)
	})

	t.Run("无效的凭证ID", func(t *testing.T) {
		mock := &mockCredentialSvc{}
		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodDelete, "/v1/credentials/abc", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrParam.Code, resp.Code)
	})

	t.Run("删除不存在的凭证", func(t *testing.T) {
		mock := &mockCredentialSvc{
			DeleteCredentialFn: func(userID, id int64) error {
				return errno.ErrCredentialNotFound
			},
		}
		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodDelete, "/v1/credentials/999", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrCredentialNotFound.Code, resp.Code)
	})
}

func TestSetDefaultCredential(t *testing.T) {
	t.Run("成功设置默认凭证", func(t *testing.T) {
		now := time.Now()
		mock := &mockCredentialSvc{
			SetDefaultCredentialFn: func(userID, id int64) (*v1.CredentialVO, error) {
				assert.Equal(t, testUserID, userID)
				assert.Equal(t, int64(2), id)
				return &v1.CredentialVO{
					ID:        2,
					Name:      "测试环境",
					ApiToken:  "token2",
					AccountID: "acc-2",
					IsDefault: true,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
		}

		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodPut, "/v1/credentials/2/default", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, 0, resp.Code)

		vo := parseDataToVO(t, resp.Data)
		assert.Equal(t, int64(2), vo.ID)
		assert.True(t, vo.IsDefault)
	})

	t.Run("无效的凭证ID", func(t *testing.T) {
		mock := &mockCredentialSvc{}
		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodPut, "/v1/credentials/abc/default", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrParam.Code, resp.Code)
	})

	t.Run("设置不存在的凭证为默认", func(t *testing.T) {
		mock := &mockCredentialSvc{
			SetDefaultCredentialFn: func(userID, id int64) (*v1.CredentialVO, error) {
				return nil, errno.ErrCredentialNotFound
			},
		}
		r := setupRouter(mock)
		token := generateToken(t)

		req := httptest.NewRequest(http.MethodPut, "/v1/credentials/999/default", nil)
		req.Header.Set("Authorization", authHeader(token))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := parseResponse(t, w.Body)
		assert.Equal(t, errno.ErrCredentialNotFound.Code, resp.Code)
	})
}
