package svc_test

import (
	"testing"
	"time"

	"cloudflared-tunnel/internal/module/credential/repo"
	"cloudflared-tunnel/internal/module/tunnel/svc"
	"cloudflared-tunnel/pkg/cloudflare"
	"cloudflared-tunnel/pkg/crypto"
	"cloudflared-tunnel/pkg/errno"
	"cloudflared-tunnel/pkg/testutil"

	cf "github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret-key-32bytes-long!!!!"

func TestListTunnels(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功查询隧道列表", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		now := time.Now()
		tunnelClient := &mockTunnelClient{
			listTunnelsFn: func(apiToken, accountID string) ([]cf.Tunnel, error) {
				return []cf.Tunnel{
					{ID: "tunnel-1", Name: "测试隧道1", Status: "healthy", CreatedAt: &now},
					{ID: "tunnel-2", Name: "测试隧道2", Status: "healthy", CreatedAt: &now},
				}, nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, tunnelClient, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		vos, err := tunnelSvc.ListTunnels(user.ID, cred.ID)

		assert.NoError(t, err)
		assert.Len(t, vos, 2)
		assert.Equal(t, "tunnel-1", vos[0].ID)
		assert.Equal(t, "测试隧道1", vos[0].Name)
	})

	t.Run("凭证不存在", func(t *testing.T) {
		tunnelSvc := svc.NewSvc(credentialRepo, nil, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		_, err := tunnelSvc.ListTunnels(user.ID, 999)

		assert.ErrorIs(t, err, errno.ErrCredentialNotFound)
	})
}

func TestGetTunnel(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功查询隧道详情", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		now := time.Now()
		tunnelClient := &mockTunnelClient{
			getTunnelFn: func(apiToken, accountID, tunnelID string) (cf.Tunnel, error) {
				return cf.Tunnel{ID: tunnelID, Name: "测试隧道", Status: "healthy", CreatedAt: &now}, nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, tunnelClient, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		vo, err := tunnelSvc.GetTunnel(user.ID, cred.ID, "tunnel-1")

		assert.NoError(t, err)
		assert.Equal(t, "tunnel-1", vo.ID)
		assert.Equal(t, "测试隧道", vo.Name)
	})

	t.Run("隧道不存在", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证2", encryptedToken, "test-account-id", false)

		tunnelClient := &mockTunnelClient{
			getTunnelFn: func(apiToken, accountID, tunnelID string) (cf.Tunnel, error) {
				return cf.Tunnel{}, assert.AnError
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, tunnelClient, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		_, err := tunnelSvc.GetTunnel(user.ID, cred.ID, "tunnel-999")

		assert.ErrorIs(t, err, errno.ErrTunnelNotFound)
	})
}

func TestCreateTunnel(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功创建隧道", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		now := time.Now()
		tunnelClient := &mockTunnelClient{
			createTunnelFn: func(apiToken, accountID, name string) (cf.Tunnel, error) {
				return cf.Tunnel{ID: "new-tunnel", Name: name, Status: "healthy", CreatedAt: &now}, nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, tunnelClient, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		vo, err := tunnelSvc.CreateTunnel(user.ID, cred.ID, "新隧道")

		assert.NoError(t, err)
		assert.Equal(t, "new-tunnel", vo.ID)
		assert.Equal(t, "新隧道", vo.Name)
	})

	t.Run("创建失败", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证2", encryptedToken, "test-account-id", false)

		tunnelClient := &mockTunnelClient{
			createTunnelFn: func(apiToken, accountID, name string) (cf.Tunnel, error) {
				return cf.Tunnel{}, assert.AnError
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, tunnelClient, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		_, err := tunnelSvc.CreateTunnel(user.ID, cred.ID, "失败隧道")

		assert.ErrorIs(t, err, errno.ErrTunnelCreateFailed)
	})
}

func TestDeleteTunnel(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功删除隧道", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		tunnelClient := &mockTunnelClient{
			deleteTunnelFn: func(apiToken, accountID, tunnelID string) error {
				return nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, tunnelClient, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		err := tunnelSvc.DeleteTunnel(user.ID, cred.ID, "tunnel-1")

		assert.NoError(t, err)
	})

	t.Run("删除失败", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证2", encryptedToken, "test-account-id", false)

		tunnelClient := &mockTunnelClient{
			deleteTunnelFn: func(apiToken, accountID, tunnelID string) error {
				return assert.AnError
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, tunnelClient, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		err := tunnelSvc.DeleteTunnel(user.ID, cred.ID, "tunnel-1")

		assert.ErrorIs(t, err, errno.ErrTunnelDeleteFailed)
	})
}

func TestGetTunnelToken(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功获取 Token", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		tunnelClient := &mockTunnelClient{
			getTunnelTokenFn: func(apiToken, accountID, tunnelID string) (string, error) {
				return "test-tunnel-token-123", nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, tunnelClient, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		vo, err := tunnelSvc.GetTunnelToken(user.ID, cred.ID, "tunnel-1")

		assert.NoError(t, err)
		assert.Equal(t, "test-tunnel-token-123", vo.Token)
	})
}

func TestGetTunnelConfig(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功获取配置", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		tunnelClient := &mockTunnelClient{
			getTunnelConfigFn: func(apiToken, accountID, tunnelID string) (cf.TunnelConfigurationResult, error) {
				return cf.TunnelConfigurationResult{
					TunnelID: tunnelID,
					Config: cf.TunnelConfiguration{
						Ingress: []cf.UnvalidatedIngressRule{
							{Hostname: "app.example.com", Service: "http://localhost:8080"},
						},
					},
					Version: 1,
				}, nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, tunnelClient, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		vo, err := tunnelSvc.GetTunnelConfig(user.ID, cred.ID, "tunnel-1")

		assert.NoError(t, err)
		assert.Equal(t, "tunnel-1", vo.TunnelID)
		assert.Len(t, vo.Config.Ingress, 1)
		assert.Equal(t, "app.example.com", vo.Config.Ingress[0].Hostname)
	})
}

func TestUpdateTunnelConfig(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功更新配置", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		tunnelClient := &mockTunnelClient{
			updateTunnelConfigFn: func(apiToken, accountID, tunnelID string, config cf.TunnelConfiguration) (cf.TunnelConfigurationResult, error) {
				return cf.TunnelConfigurationResult{
					TunnelID: tunnelID,
					Config:   config,
					Version:  2,
				}, nil
			},
		}

		newConfig := cf.TunnelConfiguration{
			Ingress: []cf.UnvalidatedIngressRule{
				{Hostname: "new.example.com", Service: "http://localhost:9090"},
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, tunnelClient, nil, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		vo, err := tunnelSvc.UpdateTunnelConfig(user.ID, cred.ID, "tunnel-1", newConfig)

		assert.NoError(t, err)
		assert.Equal(t, 2, vo.Version)
		assert.Equal(t, "new.example.com", vo.Config.Ingress[0].Hostname)
	})
}

func TestListDNSRecords(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功查询 DNS 记录", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		proxied := true
		dnsClient := &mockDNSClient{
			listDNSRecordsFn: func(apiToken, zoneID string, params cf.ListDNSRecordsParams) ([]cf.DNSRecord, error) {
				return []cf.DNSRecord{
					{ID: "dns-1", Type: "CNAME", Name: "app.example.com", Content: "tunnel-1.cfargotunnel.com", Proxied: &proxied, TTL: 1},
				}, nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, nil, dnsClient, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		vos, err := tunnelSvc.ListDNSRecords(user.ID, cred.ID, "zone-1", "", "")

		assert.NoError(t, err)
		assert.Len(t, vos, 1)
		assert.Equal(t, "dns-1", vos[0].ID)
		assert.Equal(t, "CNAME", vos[0].Type)
	})

	t.Run("带过滤条件查询", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证2", encryptedToken, "test-account-id", false)

		var capturedParams cf.ListDNSRecordsParams
		dnsClient := &mockDNSClient{
			listDNSRecordsFn: func(apiToken, zoneID string, params cf.ListDNSRecordsParams) ([]cf.DNSRecord, error) {
				capturedParams = params
				return []cf.DNSRecord{}, nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, nil, dnsClient, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		_, err := tunnelSvc.ListDNSRecords(user.ID, cred.ID, "zone-1", "app.example.com", "CNAME")

		assert.NoError(t, err)
		assert.Equal(t, "app.example.com", capturedParams.Name)
		assert.Equal(t, "CNAME", capturedParams.Type)
	})
}

func TestCreateDNSRecord(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功创建 DNS 记录", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		proxied := true
		dnsClient := &mockDNSClient{
			createDNSRecordFn: func(apiToken, zoneID string, params cf.CreateDNSRecordParams) (cf.DNSRecord, error) {
				return cf.DNSRecord{
					ID:      "new-dns",
					Type:    params.Type,
					Name:    params.Name,
					Content: params.Content,
					Proxied: params.Proxied,
					TTL:     params.TTL,
				}, nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, nil, dnsClient, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		vo, err := tunnelSvc.CreateDNSRecord(user.ID, cred.ID, "zone-1", "app.example.com", "tunnel-1.cfargotunnel.com", &proxied, 1)

		assert.NoError(t, err)
		assert.Equal(t, "new-dns", vo.ID)
		assert.Equal(t, "CNAME", vo.Type)
		assert.Equal(t, "app.example.com", vo.Name)
	})

	t.Run("创建失败", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证2", encryptedToken, "test-account-id", false)

		dnsClient := &mockDNSClient{
			createDNSRecordFn: func(apiToken, zoneID string, params cf.CreateDNSRecordParams) (cf.DNSRecord, error) {
				return cf.DNSRecord{}, assert.AnError
			},
		}

		proxied := true
		tunnelSvc := svc.NewSvc(credentialRepo, nil, dnsClient, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		_, err := tunnelSvc.CreateDNSRecord(user.ID, cred.ID, "zone-1", "app.example.com", "tunnel-1.cfargotunnel.com", &proxied, 1)

		assert.ErrorIs(t, err, errno.ErrDNSRecordCreate)
	})
}

func TestUpdateDNSRecord(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功更新 DNS 记录", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		proxied := true
		dnsClient := &mockDNSClient{
			updateDNSRecordFn: func(apiToken, zoneID string, params cf.UpdateDNSRecordParams) (cf.DNSRecord, error) {
				return cf.DNSRecord{
					ID:      params.ID,
					Type:    params.Type,
					Name:    params.Name,
					Content: params.Content,
					Proxied: params.Proxied,
					TTL:     params.TTL,
				}, nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, nil, dnsClient, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		vo, err := tunnelSvc.UpdateDNSRecord(user.ID, cred.ID, "zone-1", "dns-1", "new.example.com", "tunnel-2.cfargotunnel.com", &proxied, 1)

		assert.NoError(t, err)
		assert.Equal(t, "dns-1", vo.ID)
		assert.Equal(t, "new.example.com", vo.Name)
	})
}

func TestDeleteDNSRecord(t *testing.T) {
	container := testutil.NewPostgresContainer(t)
	log := testutil.NewLogger(t)
	credentialRepo := repo.NewCredentialRepo(container.Client, log)

	user := container.InsertUser(t, "测试用户", "testuser", "password123", "test@example.com")

	t.Run("成功删除 DNS 记录", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证", encryptedToken, "test-account-id", true)

		dnsClient := &mockDNSClient{
			deleteDNSRecordFn: func(apiToken, zoneID, recordID string) error {
				return nil
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, nil, dnsClient, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		err := tunnelSvc.DeleteDNSRecord(user.ID, cred.ID, "zone-1", "dns-1")

		assert.NoError(t, err)
	})

	t.Run("删除失败", func(t *testing.T) {
		encryptedToken := encryptToken(t, "test-token")
		cred := container.InsertCredential(t, user.ID, "测试凭证2", encryptedToken, "test-account-id", false)

		dnsClient := &mockDNSClient{
			deleteDNSRecordFn: func(apiToken, zoneID, recordID string) error {
				return assert.AnError
			},
		}

		tunnelSvc := svc.NewSvc(credentialRepo, nil, dnsClient, log, []byte(testSecret), cloudflare.NewManager(), container.Client)
		err := tunnelSvc.DeleteDNSRecord(user.ID, cred.ID, "zone-1", "dns-1")

		assert.ErrorIs(t, err, errno.ErrDNSRecordDelete)
	})
}

func encryptToken(t *testing.T, token string) string {
	t.Helper()
	encrypted, err := crypto.Encrypt(token, []byte(testSecret))
	if err != nil {
		t.Fatalf("加密 Token 失败: %v", err)
	}
	return encrypted
}
