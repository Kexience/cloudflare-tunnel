package svc

import (
	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/internal/infra/logger"
	"cloudflared-tunnel/internal/module/credential/repo"
	v1 "cloudflared-tunnel/internal/module/tunnel/ui/api/req/v1"
	"cloudflared-tunnel/pkg/cloudflare"
	"cloudflared-tunnel/pkg/crypto"
	"cloudflared-tunnel/pkg/errno"

	cf "github.com/cloudflare/cloudflare-go"
)

type svc struct {
	credentialRepo repo.CredentialRepo
	tunnelClient   cloudflare.TunnelClient
	dnsClient      cloudflare.DNSClient
	log            logger.Logger
	secret         []byte
}

// NewSvc 创建隧道管理服务（返回具体类型，由 FX 绑定接口）
func NewSvc(
	credentialRepo repo.CredentialRepo,
	tunnelClient cloudflare.TunnelClient,
	dnsClient cloudflare.DNSClient,
	log logger.Logger,
	secret []byte,
) *svc {
	return &svc{
		credentialRepo: credentialRepo,
		tunnelClient:   tunnelClient,
		dnsClient:      dnsClient,
		log:            log,
		secret:         secret,
	}
}

// getCredentialAndToken 获取凭证并解密 Token
func (s *svc) getCredentialAndToken(userID, credentialID int64) (*ent.Credential, string, error) {
	cred, err := s.credentialRepo.GetCredentialByIDAndUserID(credentialID, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, "", errno.ErrCredentialNotFound
		}
		return nil, "", errno.ErrDB
	}

	token, err := crypto.Decrypt(cred.APIToken, s.secret)
	if err != nil {
		s.log.Error("解密 API Token 失败", "credentialID", credentialID, "error", err)
		return nil, "", errno.ErrCredentialDecrypt
	}

	return cred, token, nil
}

// toTunnelVO 转换隧道 VO
func (s *svc) toTunnelVO(tunnel cf.Tunnel) *v1.TunnelVO {
	vo := &v1.TunnelVO{
		ID:        tunnel.ID,
		Name:      tunnel.Name,
		Status:    tunnel.Status,
		CreatedAt: tunnel.CreatedAt,
	}

	if len(tunnel.Connections) > 0 {
		conns := make([]v1.TunnelConnectionVO, len(tunnel.Connections))
		for i, conn := range tunnel.Connections {
			conns[i] = v1.TunnelConnectionVO{
				ID:            conn.ID,
				ColoName:      conn.ColoName,
				ClientID:      conn.ClientID,
				ClientVersion: conn.ClientVersion,
				OpenedAt:      conn.OpenedAt,
				OriginIP:      conn.OriginIP,
			}
		}
		vo.Connections = conns
	}

	return vo
}

// toDNSRecordVO 转换 DNS 记录 VO
func (s *svc) toDNSRecordVO(record cf.DNSRecord) *v1.DNSRecordVO {
	return &v1.DNSRecordVO{
		ID:      record.ID,
		Type:    record.Type,
		Name:    record.Name,
		Content: record.Content,
		Proxied: record.Proxied,
		TTL:     record.TTL,
	}
}
