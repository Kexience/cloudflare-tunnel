package cloudflare

import (
	"context"
	"fmt"

	cf "github.com/cloudflare/cloudflare-go"
)

// TunnelClient Cloudflare Tunnel API 客户端封装
type TunnelClient interface {
	// ListTunnels 查询隧道列表
	ListTunnels(apiToken, accountID string) ([]cf.Tunnel, error)
	// GetTunnel 查询隧道详情
	GetTunnel(apiToken, accountID, tunnelID string) (cf.Tunnel, error)
	// CreateTunnel 创建隧道
	CreateTunnel(apiToken, accountID, name string) (cf.Tunnel, error)
	// DeleteTunnel 删除隧道
	DeleteTunnel(apiToken, accountID, tunnelID string) error
	// GetTunnelToken 获取隧道连接 Token
	GetTunnelToken(apiToken, accountID, tunnelID string) (string, error)
	// ListTunnelConnections 查询隧道连接列表
	ListTunnelConnections(apiToken, accountID, tunnelID string) ([]cf.Connection, error)
	// GetTunnelConfig 获取隧道配置（Ingress 规则等）
	GetTunnelConfig(apiToken, accountID, tunnelID string) (cf.TunnelConfigurationResult, error)
	// UpdateTunnelConfig 更新隧道配置
	UpdateTunnelConfig(apiToken, accountID, tunnelID string, config cf.TunnelConfiguration) (cf.TunnelConfigurationResult, error)
}

type tunnelClient struct{}

func NewTunnelClient() TunnelClient {
	return &tunnelClient{}
}

func (c *tunnelClient) ListTunnels(apiToken, accountID string) ([]cf.Tunnel, error) {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.AccountIdentifier(accountID)

	tunnels, _, err := client.ListTunnels(ctx, rc, cf.TunnelListParams{})
	if err != nil {
		return nil, fmt.Errorf("查询隧道列表失败: %w", err)
	}

	return tunnels, nil
}

func (c *tunnelClient) GetTunnel(apiToken, accountID, tunnelID string) (cf.Tunnel, error) {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return cf.Tunnel{}, fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.AccountIdentifier(accountID)

	tunnel, err := client.GetTunnel(ctx, rc, tunnelID)
	if err != nil {
		return cf.Tunnel{}, fmt.Errorf("查询隧道详情失败: %w", err)
	}

	return tunnel, nil
}

func (c *tunnelClient) CreateTunnel(apiToken, accountID, name string) (cf.Tunnel, error) {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return cf.Tunnel{}, fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.AccountIdentifier(accountID)

	tunnel, err := client.CreateTunnel(ctx, rc, cf.TunnelCreateParams{
		Name: name,
	})
	if err != nil {
		return cf.Tunnel{}, fmt.Errorf("创建隧道失败: %w", err)
	}

	return tunnel, nil
}

func (c *tunnelClient) DeleteTunnel(apiToken, accountID, tunnelID string) error {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.AccountIdentifier(accountID)

	if err := client.DeleteTunnel(ctx, rc, tunnelID); err != nil {
		return fmt.Errorf("删除隧道失败: %w", err)
	}

	return nil
}

func (c *tunnelClient) GetTunnelToken(apiToken, accountID, tunnelID string) (string, error) {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return "", fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.AccountIdentifier(accountID)

	token, err := client.GetTunnelToken(ctx, rc, tunnelID)
	if err != nil {
		return "", fmt.Errorf("获取隧道 Token 失败: %w", err)
	}

	return token, nil
}

func (c *tunnelClient) ListTunnelConnections(apiToken, accountID, tunnelID string) ([]cf.Connection, error) {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.AccountIdentifier(accountID)

	connections, err := client.ListTunnelConnections(ctx, rc, tunnelID)
	if err != nil {
		return nil, fmt.Errorf("查询隧道连接失败: %w", err)
	}

	return connections, nil
}

func (c *tunnelClient) GetTunnelConfig(apiToken, accountID, tunnelID string) (cf.TunnelConfigurationResult, error) {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return cf.TunnelConfigurationResult{}, fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.AccountIdentifier(accountID)

	result, err := client.GetTunnelConfiguration(ctx, rc, tunnelID)
	if err != nil {
		return cf.TunnelConfigurationResult{}, fmt.Errorf("获取隧道配置失败: %w", err)
	}

	return result, nil
}

func (c *tunnelClient) UpdateTunnelConfig(apiToken, accountID, tunnelID string, config cf.TunnelConfiguration) (cf.TunnelConfigurationResult, error) {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return cf.TunnelConfigurationResult{}, fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.AccountIdentifier(accountID)

	result, err := client.UpdateTunnelConfiguration(ctx, rc, cf.TunnelConfigurationParams{
		TunnelID: tunnelID,
		Config:   config,
	})
	if err != nil {
		return cf.TunnelConfigurationResult{}, fmt.Errorf("更新隧道配置失败: %w", err)
	}

	return result, nil
}
