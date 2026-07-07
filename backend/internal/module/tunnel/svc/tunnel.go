package svc

import (
	"os/exec"

	v1 "cloudflared-tunnel/internal/module/tunnel/ui/api/req/v1"
	"cloudflared-tunnel/pkg/errno"

	cf "github.com/cloudflare/cloudflare-go"
)

// TunnelSvc 隧道管理服务接口
type TunnelSvc interface {
	// ListTunnels 查询隧道列表
	ListTunnels(userID, credentialID int64) ([]*v1.TunnelVO, error)
	// GetTunnel 查询隧道详情
	GetTunnel(userID, credentialID int64, tunnelID string) (*v1.TunnelVO, error)
	// CreateTunnel 创建隧道
	CreateTunnel(userID, credentialID int64, name string) (*v1.TunnelVO, error)
	// DeleteTunnel 删除隧道
	DeleteTunnel(userID, credentialID int64, tunnelID string) error
	// GetTunnelToken 获取隧道连接 Token
	GetTunnelToken(userID, credentialID int64, tunnelID string) (*v1.TunnelTokenVO, error)
	// ListTunnelConnections 查询隧道连接列表
	ListTunnelConnections(userID, credentialID int64, tunnelID string) (*v1.TunnelVO, error)
	// GetTunnelConfig 获取隧道配置
	GetTunnelConfig(userID, credentialID int64, tunnelID string) (*v1.TunnelConfigVO, error)
	// UpdateTunnelConfig 更新隧道配置
	UpdateTunnelConfig(userID, credentialID int64, tunnelID string, config cf.TunnelConfiguration) (*v1.TunnelConfigVO, error)
	// StartTunnel 启动隧道
	StartTunnel(userID, credentialID int64, tunnelID string) error
	// StopTunnel 停止隧道
	StopTunnel(userID, credentialID int64, tunnelID string) error
	// GetTunnelStatus 获取隧道运行状态
	GetTunnelStatus(userID, credentialID int64, tunnelID string) (*v1.TunnelStatusVO, error)
}

func (s *svc) ListTunnels(userID, credentialID int64) ([]*v1.TunnelVO, error) {
	cred, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return nil, err
	}

	tunnels, err := s.tunnelClient.ListTunnels(token, cred.AccountID)
	if err != nil {
		s.log.Error("查询隧道列表失败", "error", err)
		return nil, errno.ErrTunnelNotFound
	}

	vos := make([]*v1.TunnelVO, len(tunnels))
	for i, tunnel := range tunnels {
		vos[i] = s.toTunnelVO(tunnel)
	}

	return vos, nil
}

func (s *svc) GetTunnel(userID, credentialID int64, tunnelID string) (*v1.TunnelVO, error) {
	cred, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return nil, err
	}

	tunnel, err := s.tunnelClient.GetTunnel(token, cred.AccountID, tunnelID)
	if err != nil {
		s.log.Error("查询隧道详情失败", "tunnelID", tunnelID, "error", err)
		return nil, errno.ErrTunnelNotFound
	}

	return s.toTunnelVO(tunnel), nil
}

func (s *svc) CreateTunnel(userID, credentialID int64, name string) (*v1.TunnelVO, error) {
	cred, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return nil, err
	}

	tunnel, err := s.tunnelClient.CreateTunnel(token, cred.AccountID, name)
	if err != nil {
		s.log.Error("创建隧道失败", "name", name, "error", err)
		return nil, errno.ErrTunnelCreateFailed
	}

	return s.toTunnelVO(tunnel), nil
}

func (s *svc) DeleteTunnel(userID, credentialID int64, tunnelID string) error {
	cred, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return err
	}

	if err := s.tunnelClient.DeleteTunnel(token, cred.AccountID, tunnelID); err != nil {
		s.log.Error("删除隧道失败", "tunnelID", tunnelID, "error", err)
		return errno.ErrTunnelDeleteFailed
	}

	return nil
}

func (s *svc) GetTunnelToken(userID, credentialID int64, tunnelID string) (*v1.TunnelTokenVO, error) {
	cred, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return nil, err
	}

	tunnelToken, err := s.tunnelClient.GetTunnelToken(token, cred.AccountID, tunnelID)
	if err != nil {
		s.log.Error("获取隧道 Token 失败", "tunnelID", tunnelID, "error", err)
		return nil, errno.ErrTunnelNotFound
	}

	return &v1.TunnelTokenVO{Token: tunnelToken}, nil
}

func (s *svc) ListTunnelConnections(userID, credentialID int64, tunnelID string) (*v1.TunnelVO, error) {
	cred, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return nil, err
	}

	tunnel, err := s.tunnelClient.GetTunnel(token, cred.AccountID, tunnelID)
	if err != nil {
		s.log.Error("查询隧道详情失败", "tunnelID", tunnelID, "error", err)
		return nil, errno.ErrTunnelNotFound
	}

	connections, err := s.tunnelClient.ListTunnelConnections(token, cred.AccountID, tunnelID)
	if err != nil {
		s.log.Error("查询隧道连接失败", "tunnelID", tunnelID, "error", err)
		return nil, errno.ErrTunnelNotFound
	}

	vo := s.toTunnelVO(tunnel)
	if len(connections) > 0 {
		conns := make([]v1.TunnelConnectionVO, len(connections))
		for i, conn := range connections {
			conns[i] = v1.TunnelConnectionVO{
				ID:            conn.ID,
				ColoName:      conn.Connections[0].ColoName,
				ClientID:      conn.Connections[0].ClientID,
				ClientVersion: conn.Version,
				OpenedAt:      conn.Connections[0].OpenedAt,
				OriginIP:      conn.Connections[0].OriginIP,
			}
		}
		vo.Connections = conns
	}

	return vo, nil
}

func (s *svc) GetTunnelConfig(userID, credentialID int64, tunnelID string) (*v1.TunnelConfigVO, error) {
	cred, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return nil, err
	}

	result, err := s.tunnelClient.GetTunnelConfig(token, cred.AccountID, tunnelID)
	if err != nil {
		s.log.Error("获取隧道配置失败", "tunnelID", tunnelID, "error", err)
		return nil, errno.ErrTunnelConfigFailed
	}

	return &v1.TunnelConfigVO{
		TunnelID: result.TunnelID,
		Config:   result.Config,
		Version:  result.Version,
	}, nil
}

func (s *svc) UpdateTunnelConfig(userID, credentialID int64, tunnelID string, config cf.TunnelConfiguration) (*v1.TunnelConfigVO, error) {
	cred, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return nil, err
	}

	result, err := s.tunnelClient.UpdateTunnelConfig(token, cred.AccountID, tunnelID, config)
	if err != nil {
		s.log.Error("更新隧道配置失败", "tunnelID", tunnelID, "error", err)
		return nil, errno.ErrTunnelConfigFailed
	}

	return &v1.TunnelConfigVO{
		TunnelID: result.TunnelID,
		Config:   result.Config,
		Version:  result.Version,
	}, nil
}

func (s *svc) StartTunnel(userID, credentialID int64, tunnelID string) error {
	cred, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查隧道是否已在运行
	if cmd, exists := s.processes[tunnelID]; exists && cmd.Process != nil {
		if err := cmd.Process.Signal(nil); err == nil {
			return errno.ErrTunnelStartFailed.WithMessage("隧道已在运行中")
		}
	}

	// 获取隧道 Token
	tunnelToken, err := s.tunnelClient.GetTunnelToken(token, cred.AccountID, tunnelID)
	if err != nil {
		s.log.Error("获取隧道 Token 失败", "tunnelID", tunnelID, "error", err)
		return errno.ErrTunnelStartFailed
	}

	// 启动 cloudflared 进程
	cmd := exec.Command("cloudflared", "tunnel", "run", "--token", tunnelToken, tunnelID)
	if err := cmd.Start(); err != nil {
		s.log.Error("启动隧道进程失败", "tunnelID", tunnelID, "error", err)
		return errno.ErrTunnelStartFailed
	}

	s.processes[tunnelID] = cmd
	return nil
}

func (s *svc) StopTunnel(userID, credentialID int64, tunnelID string) error {
	// 验证凭证存在
	if _, _, err := s.getCredentialAndToken(userID, credentialID); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cmd, exists := s.processes[tunnelID]
	if !exists || cmd.Process == nil {
		return errno.ErrTunnelStopFailed.WithMessage("隧道未在运行")
	}

	// 终止进程
	if err := cmd.Process.Kill(); err != nil {
		s.log.Error("停止隧道进程失败", "tunnelID", tunnelID, "error", err)
		return errno.ErrTunnelStopFailed
	}

	// 等待进程退出
	go func() {
		cmd.Wait()
		s.mu.Lock()
		delete(s.processes, tunnelID)
		s.mu.Unlock()
	}()

	return nil
}

func (s *svc) GetTunnelStatus(userID, credentialID int64, tunnelID string) (*v1.TunnelStatusVO, error) {
	// 验证凭证存在
	if _, _, err := s.getCredentialAndToken(userID, credentialID); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	status := "stopped"
	if cmd, exists := s.processes[tunnelID]; exists && cmd.Process != nil {
		if err := cmd.Process.Signal(nil); err == nil {
			status = "running"
		}
	}

	return &v1.TunnelStatusVO{
		TunnelID: tunnelID,
		Status:   status,
	}, nil
}
